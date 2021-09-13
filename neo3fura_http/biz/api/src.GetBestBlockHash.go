package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBestBlockHash(args struct {
	Filter map[string]interface{}
	Raw    *map[string]interface{}
}, ret *json.RawMessage) error {
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Block",
		Index:      "GetBestBlockHash",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{},
		Query:      []string{"hash"},
	}, ret)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
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
