package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/uintval"
	"neo3fura/var/stderr"
)

func (me *T) GetRawTransactionByBlockHeight(args struct {
	BlockHeight uintval.T
	Limit       int64
	Skip        int64
	Filter      map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Limit == 0 {
		args.Limit = 200
	}
	if args.BlockHeight.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Block",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"index": args.BlockHeight},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Transaction",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"blockhash": r1["hash"]},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r3, err := me.FilterArrayAndAppendCount(r2, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil

}
