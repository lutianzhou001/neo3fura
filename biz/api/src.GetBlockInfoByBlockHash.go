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

		}{  Collection: "[Block~Transaction(Transactions)]",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{"ParentID":r1["_id"],
			}}, ret)
	if err != nil {
		return err
	}
	r1["transactionNumber"] = r2["total counts"]
	r3, err := me.Data.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Transaction",
		Index:      "someIndex",
		Pipeline: []bson.M{bson.M{"$match":bson.M{"blockhash":args.BlockHash}},
			bson.M{"$group":bson.M{"_id":"$blockhash","systemFee":bson.M{"$sum":"$sysfee"},"networkFee":bson.M{"$sum":"$netfee"}}}},
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	if len(r3) != 0{
		r1["totalNetworkFee"] = r3[0]["systemFee"]
		r1["totalSystemFee"] = r3[0]["networkFee"]
	}else {
		r1["totalNetworkFee"] = 0
		r1["totalSystemFee"] = 0
	}
		r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
