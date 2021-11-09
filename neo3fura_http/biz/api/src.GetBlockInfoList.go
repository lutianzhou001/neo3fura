package api

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBlockInfoList(args struct {
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
}, ret *json.RawMessage) error {

	if args.Limit == 0 {
		args.Limit = 20
	}

	r1, count, err := me.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Block",
			Index:      "GetBlockInfoList",
			Sort:       bson.M{"_id": -1},
			Filter:     bson.M{},
			Query:      []string{"_id", "index", "size", "timestamp", "hash"},
			Limit:      args.Limit,
			Skip:       args.Skip,
		}, ret)
	if err != nil {
		return err
	}
	r2 := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		r3, err := me.Client.QueryDocument(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
			}{Collection: "Transaction",
				Index:  "GetBlockInfoList",
				Sort:   bson.M{},
				Filter: bson.M{"blockhash": item["hash"]}}, ret)
		if err != nil {
			return err
		}
		if r3["total counts"] == nil {
			item["transactioncount"] = 0
		} else {
			item["transactioncount"] = r3["total counts"]
		}
		r2 = append(r2, item)
	}
	r4, err := me.FilterArrayAndAppendCount(r2, count, args.Filter)
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
