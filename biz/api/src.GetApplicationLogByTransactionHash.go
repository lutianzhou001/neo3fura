package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetApplicationLogByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.TransactionHash.IsZero() == true {
		return stderr.ErrZero
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Execution",
		Index:      "GetApplicationLogByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	r2, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{Collection: "Notification", Index: "GetApplicationLogByBlockHash", Sort: bson.M{}, Filter: bson.M{"txid": r1["txid"].(string), "blockhash": r1["blockhash"].(string)}}, ret)
	if err != nil {
		return err
	}
	r1["notifications"] = r2

	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return nil
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
