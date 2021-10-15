package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNep11PropertiesByContractHashTokenId(args struct {
	Address      h160.T
	ContractHash h160.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Address-Asset",
		Index:      "GetAssetsHeldByContractHashAddress",
		Sort:       bson.M{"balance": -1},
		Filter:     bson.M{"address": args.Address.TransferredVal(), "asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2, err := me.Filter(r1, args.Filter)
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
