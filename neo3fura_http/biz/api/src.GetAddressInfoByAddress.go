package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetAddressInfoByAddress(args struct {
	Address h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Address",
		Index:      "GetAddressInfoByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"address": args.Address.TransferredVal()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	_, count, err := me.Client.QueryAll(struct {
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
		Filter:     bson.M{"sender": args.Address.TransferAddress()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r1["transactionssent"] = count
	r3, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Transaction",
		Index:      "GetRawTransactionByAddress",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"sender": args.Address.TransferAddress()},
		Query:      []string{},
	}, ret)

	r1["lastusetime"] = r3["blocktime"]
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
