package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNep11TransferByAddress(args struct {
	Address h160.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
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
		Collection: "Nep11TransferNotification",
		Index:      "GetNep11TransferByAddress",
		Sort:       bson.M{},
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal()},
			bson.M{"to": args.Address.TransferredVal()},
		}},
		Query: []string{},
		Limit: args.Limit,
		Skip:  args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	var raw1 map[string]interface{}
	var raw3 map[string]interface{}
	for _, item := range r1 {
		err = me.GetVmStateByTransactionHash(struct {
			TransactionHash h256.T
			Filter          map[string]interface{}
			Raw             *map[string]interface{}
		}{
			TransactionHash: h256.T(fmt.Sprint(item["txid"])),
			Filter:          nil,
			Raw:             &raw1,
		}, ret)
		if err != nil {
			return err
		}
		item["vmstate"] = raw1["vmstate"].(string)

		if fmt.Sprint(item["txid"]) != "0x0000000000000000000000000000000000000000000000000000000000000000" {
			err = me.GetRawTransactionByTransactionHash(struct {
				TransactionHash h256.T
				Filter          map[string]interface{}
				Raw             *map[string]interface{}
			}{TransactionHash: h256.T(fmt.Sprint(item["txid"])), Raw: &raw3}, ret)
			if err != nil {
				return err
			}
			item["netfee"] = raw3["netfee"]
			item["sysfee"] = raw3["sysfee"]
		}
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
