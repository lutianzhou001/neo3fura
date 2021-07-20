package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
)

func (me *T) GetTopTokens(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, _, err := me.Data.Client.QueryAll(struct {
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
	r3 := make(map[string]int)
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
	for i, kv := range kvs {
		if i < 10 {
			r3[kv.Key] = kv.Value
		}
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
