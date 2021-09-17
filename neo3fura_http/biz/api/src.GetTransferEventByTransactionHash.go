package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/lib/type/strval"
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
	var reversedArray []string
	if len(r1["hexStringParams"].(primitive.A)) != 0 {
		for _, item := range r1["hexStringParams"].(primitive.A) {
			st := strval.T(item.(string))
			reversedItem := st.Reverse()
			reversedArray = append(reversedArray, reversedItem)
		}
		r1["hexStringParams"] = reversedArray
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
