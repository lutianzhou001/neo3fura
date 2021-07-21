package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
)

func (me T) GetTopTokens() error {
	var ret *json.RawMessage
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
		Index:      "GetTopTokens",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": bson.M{"$gt": 100000000}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2 := make(map[string]int)
	for _, item := range r1 {
		r2[item["contracthash"].(string)] = r2[item["contracthash"].(string)] + 1
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
	var data []interface{}
	for i, kv := range kvs {
		if i < 10 {
			data = append(data, bson.M{kv.Key: kv.Value})
		}
	}
	_, err = me.Client.Save(struct {
		Collection string
		Data       []interface{}
	}{Collection: "TopTokens", Data: data})
	if err != nil {
		return err
	}
	return nil
}
