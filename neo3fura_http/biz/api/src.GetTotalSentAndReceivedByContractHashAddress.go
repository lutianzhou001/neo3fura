package api

import (
	"encoding/json"
	"math/big"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetTotalSentAndReceivedByContractHashAddress(args struct {
	ContractHash h160.T
	Address      h160.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByAddress",
		Sort:       bson.M{"_id": -1},
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal()},
			bson.M{"to": args.Address.TransferredVal()},
		}, "contract": args.ContractHash.Val()},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2 := make(map[string]interface{}, 0)
	var totalSent = new(big.Int)
	var totalReceived = new(big.Int)
	for _, item := range r1 {
		if item["from"] == args.Address.TransferredVal() {
			it, _, err := item["value"].(primitive.Decimal128).BigInt()
			if err != nil {
				return err
			}
			totalSent.Add(totalSent, it)
		} else {
			ir, _, err := item["value"].(primitive.Decimal128).BigInt()
			if err != nil {
				return err
			}
			totalReceived.Add(totalReceived, ir)
		}
	}
	r2["ContractHash"] = args.ContractHash
	r2["Address"] = args.Address
	r2["sent"] = totalSent
	r2["received"] = totalReceived
	r3, err := me.Filter(r2, args.Filter)
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
