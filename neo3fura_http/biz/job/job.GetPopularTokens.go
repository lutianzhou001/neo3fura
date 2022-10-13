package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me T) GetPopularTokens() error {
	message := make(json.RawMessage, 0)
	ret := &message
	// timeUnix := time.Now().Unix()*1000 - 24*86400*1000
	r0, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Block", Index: "GetPopularTokens", Sort: bson.M{"_id": -1}}, ret)
	if err != nil {
		return err
	}
	//分组查询
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"timestamp": bson.M{"$gt": r0["timestamp"].(int64) - 3*3600*24*1000}}},
		bson.M{"$lookup": bson.M{
			"from": "Asset",
			"let":  bson.M{"hash": "$contract"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$hash", "$$hash"}}}},
				bson.M{"$project": bson.M{"type": 1, "_id": -1}},
			},
			"as": "type"},
		},
		bson.M{"$group": bson.M{"_id": "$contract", "contract": bson.M{"$last": "$contract"}, "type": bson.M{"$last": "$type"}, "count": bson.M{"$sum": 1}}},
		bson.M{"$sort": bson.M{"count": -1}},
	}
	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Notification",
		Index:      "GetPopularTokens",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline:   pipeline,
		Query:      []string{},
	}, ret)

	if err != nil {
		return err
	}
	arr := MapArrGroup(r1)

	nep11Arr := arr["NEP11"]
	nep17Arr := arr["NEP17"]
	var values []string
	if len(nep17Arr) > 0 {
		for i, item := range nep17Arr {
			if i < 5 {
				contract := item["contract"].(string)
				values = append(values, contract)
			}
		}
	}
	if len(nep11Arr) > 0 {
		for i, item := range nep11Arr {
			if i < 5 {
				contract := item["contract"].(string)
				values = append(values, contract)
			}
		}
	}

	data := bson.M{"Populars": values}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "PopularTokens", Data: data})
	if err != nil {
		return err
	}
	return nil
}

// 根据Asset 的type 分类
func MapArrGroup(infos []map[string]interface{}) map[string][]map[string]interface{} {
	res := make(map[string][]map[string]interface{})
	nep11 := make([]map[string]interface{}, 0)
	nep17 := make([]map[string]interface{}, 0)
	for _, item := range infos {
		neptype := item["type"].(primitive.A)
		if len(neptype) > 0 {
			t := neptype[0].(map[string]interface{})
			if t["type"].(string) == "NEP11" {
				nep11 = append(nep11, item)
			} else if t["type"].(string) == "NEP17" {
				nep17 = append(nep17, item)
			}
		}

	}
	res["NEP11"] = nep11
	res["NEP17"] = nep17

	return res
}
