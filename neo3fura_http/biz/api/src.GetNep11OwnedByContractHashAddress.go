package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/jsonrpc2"
	"neo3fura_http/lib/type/h160"
)

// this function may be not supported any more, we only support address in the formart of script hash
func (me *T) GetNep11OwnedByContractHashAddress(args struct {
	ContractHash h160.T
	Address      h160.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false && args.ContractHash.Valid() == false {
		return jsonrpc2.NewError(-32001, "You should enter a parameter in ContractHash or Address")
	}
	var filter bson.M

	if args.Address.Valid() && args.ContractHash.Valid() {
		filter = bson.M{"tokenid": bson.M{"$ne": ""}, "balance": bson.M{"$gt": 0}, "address": args.Address.Val(), "asset": args.ContractHash.Val()}
	} else if args.Address.Valid() && !args.ContractHash.Valid() {
		filter = bson.M{"tokenid": bson.M{"$ne": ""}, "balance": bson.M{"$gt": 0}, "address": args.Address.Val()}
	} else if !args.Address.Valid() && args.ContractHash.Valid() {
		filter = bson.M{"tokenid": bson.M{"$ne": ""}, "balance": bson.M{"$gt": 0}, "asset": args.ContractHash.Val()}
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
		Collection: "Address-Asset",
		Index:      "GetNep11OwnedByContractHashAddress",
		Sort:       bson.M{},
		Filter:     filter,
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	//r2, err := me.Deduplicate(r1)
	r3, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
