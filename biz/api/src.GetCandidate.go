package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetCandidate(args struct {
	Limit            int64
	Skip             int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	r1, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Candidate",
		Index:      "someIndex",
		Sort:       bson.M{"votesOfCandidate":-1},
		Filter:     bson.M{},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}