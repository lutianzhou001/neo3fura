package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/lib/type/strval"
	"neo3fura/var/stderr"
)

func (me *T) GetNep11BalanceByContractHashAddressTokenId(args struct {
	ContractHash h160.T
	Address      h160.T
	TokenId      strval.T
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if len(args.TokenId.Val()) <= 0 {
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
		Index:      "GetNep11BalanceByContractHashAddressTokenId",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "address": args.Address.TransferredVal(), "tokenid": args.TokenId},
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
	if args.Raw != nil {
		*args.Raw = r2
	}
	*ret = json.RawMessage(r)
	return nil
}
