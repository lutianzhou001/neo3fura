package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetScCallByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
	Skip             int64
	Limit            int64
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, count,err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit int64
		Skip int64
	}{
		Collection: "ScCall",
		Index:      "GetScCallByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
		Limit :args.Limit,
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
