package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/addr"
	"neo3fura/var/stderr"
)

// this function may be not supported any more, we only support address in the formart of script hash
func (me *T) GetRawTransactionBySender(args struct {
	Sender addr.T
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Sender.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, count, err := me.Data.Client.QueryAll(struct {
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
		Filter:     bson.M{"sender": args.Sender.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
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
