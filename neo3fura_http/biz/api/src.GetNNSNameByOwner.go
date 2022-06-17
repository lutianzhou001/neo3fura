package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strconv"
)

func (me *T) GetNNSNameByOwner(args struct {
	Asset  h160.T
	Owner  h160.T
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Owner.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Address-Asset",
			Index:      "GetNNSNameByOwner",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": args.Asset, "address": args.Owner, "balance": bson.M{"$gt": 0}}},
				bson.M{"$sort": bson.M{"id": 1}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
				bson.M{"$lookup": bson.M{
					"from": "Nep11Properties",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{

						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						//bson.M{"$project": bson.M{"id": 1, "tokenid": 1, "asset": 1, "properties": 1}},
					},
					"as": "properties"},
				},
				bson.M{"$project": bson.M{"_id": 0, "properties.properties": 1, "asset": 1, "tokenid": 1, "address": 1}}},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {
		properties := item["properties"].(primitive.A)[0].(map[string]interface{})

		if properties["properties"] != nil {
			extendData := properties["properties"].(string)
			if extendData != "" {
				var data map[string]interface{}
				if err2 := json.Unmarshal([]byte(extendData), &data); err2 == nil {
					name, ok := data["name"]
					if ok {
						item["name"] = name
					}
					admin, ok2 := data["admin"]
					if ok2 {
						item["admin"] = admin
					}
					expiration, ok3 := data["expiration"]
					if ok3 {
						time, err := strconv.ParseInt(expiration.(string), 10, 64)
						if err != nil {
							return err
						}
						item["expiration"] = time
					}
				} else {
					return err2
				}

			}
		}

		delete(item, "properties")

	}

	r1 = mapsort.MapSort2(r1, "expiration") //升序
	if err != nil {
		return err
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}

	*ret = json.RawMessage(r)
	return nil
}
