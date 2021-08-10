package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetActiveAddresses() error {
	message := make(json.RawMessage, 0)
	ret := &message

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
		Collection: "Transaction",
		Index:      "GetActiveAddresses",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": bson.M{"$gt": r0["timestamp"].(int64) - 3600*24*1000}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2 := make(map[string]interface{})
	for _, item := range r1 {
		r2[item["sender"].(string)] = true
	}
	var i = 0
	for _, _ = range r2 {
		i++
	}
	data := bson.M{"ActiveAddresses": i}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "ActiveAddresses", Data: data})
	if err != nil {
		return err
	}
	return nil
}
