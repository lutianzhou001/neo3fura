package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetTransferEventByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
	Raw             *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, err := me.Client.QueryOneJob(struct {
		Collection string
		Filter     bson.M
	}{Collection: "TransferEvent", Filter: bson.M{"txid": args.TransactionHash.Val()}})
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
	}
	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
