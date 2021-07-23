package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetRawTransactionByAddress(args struct {
	Address h160.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	sender, err := args.Address.ScriptHashToAddress()
	if err != nil {
		return nil
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
		Index:      "GetRawTransactionByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"sender": sender},
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