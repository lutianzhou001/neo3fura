package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetActiveAddresses(args struct {
	Days   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Days == 0 {
		args.Days = 7
	}
	r1, err := me.Client.QueryLastJobs(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{Collection: "ActiveAddresses", Index: "GetActiveAddresses", Sort: bson.M{"id": -1}, Filter: bson.M{}, Query: []string{}, Limit: args.Days})
	if err != nil {
		return err
	}
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
