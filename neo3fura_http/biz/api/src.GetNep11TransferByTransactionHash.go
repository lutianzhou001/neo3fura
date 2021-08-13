package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep11TransferByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Nep11TransferNotification",
		Index:      "GetNep11TransferByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	r, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Asset",
		Index:      "GetNep11TransferByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": r1["contract"]},
		Query:      []string{"tokenname"},
	}, ret)
	if err == nil {
		r1["tokenname"] = r["tokenname"]
		r1["decimals"] = r["decimals"]
	} else if err.Error() == "NOT FOUND" {
		r1["tokenname"] = ""
		r1["decimals"] = ""
	} else {
		return err
	}

	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}
	r2, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r2)
	return nil
}