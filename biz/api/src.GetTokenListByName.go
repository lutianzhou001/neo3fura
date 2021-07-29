package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetTokenListByName(args struct {
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
			Collection: "Asset",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter: bson.M{"tokenname": bson.M{"$regex":args.Name, "$options": "$i"}},
			Query: []string{"hash", "tokenname", "symbol", "_id"},
			Limit: args.Limit,
			Skip: args.Skip,
		}, ret)
	if err != nil {
		return err
	}
	for _, item := range r1 {
		r, err := me.Data.Client.QueryDocument(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
		}{
			Collection: "Address-Asset",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{"asset": item["hash"]},
		}, ret)
		item["total_holders"] = r["total counts:"]
		_, err = me.Data.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "TransferNotification", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"contract": item["hash"]}, Query: []string{},
		}, ret)
		if err != nil {
			item["standard"] = "NEP11"
		} else {
			item["standard"] = "NEP17"
		}
		delete(item, "_id")
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