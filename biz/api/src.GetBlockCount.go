package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBlockCount(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Block",
		Index:      "someIndex",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{},
		Query:      []string{"index"},
	}, ret)
	if err != nil {
		return err
	}
	r1, err = me.Filter(r1, args.Filter)
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
