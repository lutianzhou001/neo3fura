package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
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
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Notification",
		Index:      "GetPopularTokens",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": bson.M{"$gt": r0["timestamp"].(int64) - 3600*1*1000}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2 := make(map[string]int)
	for _, item := range r1 {
		r2[item["contract"].(string)] = r2[item["contract"].(string)] + 1
	}
	type kv struct {
		Key   string
		Value int
	}
	var kvs []kv
	for k, v := range r2 {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})
	var values []string
	for i, kv := range kvs {
		if i < 10 {
			values = append(values, kv.Key)
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
