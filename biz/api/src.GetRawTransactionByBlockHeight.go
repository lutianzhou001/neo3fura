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
	Raw         *[]map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHeight.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.BlockHeight.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Transaction",
		Index:      "GetRawTransactionByBlockHeight",
		Sort:       bson.M{},
		Filter:     bson.M{"blockIndex": args.BlockHeight.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
	}
	r3, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
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
