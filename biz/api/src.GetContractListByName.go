package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetContractListByName(args struct {
	Name   string
	Filter map[string]interface{}
	Limit        int64
	Skip         int64
}, ret *json.RawMessage) error {

	r1, count, err := me.Data.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter: bson.M{"name": bson.M{"$regex":args.Name, "$options": "$i"}},
			Query: []string{"hash","id","name","createtime"},
			Limit: args.Limit,
			Skip: args.Skip,
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