package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetBlockInfoByBlockHash(args struct {
	BlockHash    h256.T
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Block",
			Index:      "someIndex",
			Sort:       bson.M{"index":-1},
			Filter: bson.M{"hash": args.BlockHash},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	r2, err := me.Data.Client.QueryDocument(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M

		}{  Collection: "Transaction",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{"blockhash":args.BlockHash,
			}}, ret)
	if err != nil {
		return err
	}
	if (r2["total counts"] == nil){
		r1["transactionNumber"] = 0
		r1["transfersNumber"] = 0
	}else {
		r1["transactionNumber"] = r2["total counts"]
		r3, err := me.Data.Client.QueryDocument(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M

		}{  Collection: "TransferNotification",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{"blockhash":args.BlockHash,
			}}, ret)
		if err != nil {
			return err
		}
		if (r3["total counts"] == nil){
			r1["transfersNumber"] = 0
		}else {
			r1["transfersNumber"] = r3["total counts"]
		}
	}

		r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
