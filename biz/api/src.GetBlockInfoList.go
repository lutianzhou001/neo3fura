package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBlockInfoList(args struct {
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
			Collection: "Block",
			Index:      "someIndex",
			Sort:       bson.M{"index":-1},
			Filter: bson.M{},
			Query: []string{"_id","index","size","timestamp","hash"},
			Limit: args.Limit,
			Skip: args.Skip,
		}, ret)
	if err != nil {
		return err
	}
	r2 := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		r3, err := me.Data.Client.QueryDocument(
			struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M

		}{  Collection: "[Block~Transaction(Transactions)]",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{"ParentID":item["_id"],
				}}, ret)
		if err != nil {
			return err
		}
		if (r3["total counts"] == nil){
			item["transactionNumber"] = 0
		}else {
			item["transactionNumber"] = r3["total counts"]
		}

		//delete(item,"_id")
		r2 = append(r2, item)

	}
	r4, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}


