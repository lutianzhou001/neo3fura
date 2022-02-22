package api

import (
	"encoding/json"
	"fmt"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep17TransferByTransactionHash(args struct {
	TransactionHash h256.T
	Limit           int64
	Skip            int64
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
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
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByTransactionHash",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	var raw1 map[string]interface{}

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
		r, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Asset",
			Index:      "GetNep17TransferByTransactionHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": item["contract"]},
			Query:      []string{"tokenname", "decimals", "symbol"},
		}, ret)
		if err == nil {
			item["tokenname"] = r["tokenname"]
			item["decimals"] = r["decimals"]
			item["symbol"] = r["symbol"]

		} else if err.Error() == "NOT FOUND" {
			item["tokenname"] = ""
			item["decimals"] = ""
			item["symbol"] = ""
		} else {
			return err
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
