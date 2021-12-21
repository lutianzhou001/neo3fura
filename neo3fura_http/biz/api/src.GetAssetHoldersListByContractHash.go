package api

import (
	"encoding/json"
	"math/big"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetAssetHoldersListByContractHash(args struct {
	ContractHash h160.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
	Raw          *[]map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
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
		Collection: "Address-Asset",
		Index:      "GetAssetHoldersListByContractHash",
		Sort:       bson.M{"balance": -1},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)

	// 获取资产的totaluspply
	var raw1 map[string]interface{}
	err = me.GetAssetInfoByContractHash(struct {
		ContractHash h160.T
		Filter       map[string]interface{}
		Raw          *map[string]interface{}
	}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
	if err != nil {
		return err
	}
	// it, _ := new(big.Int).SetString(raw1["totalsupply"].(string), 10)
	it, _, err := raw1["totalsupply"].(primitive.Decimal128).BigInt()
	if err != nil {
		return err
	}
	itf := new(big.Float).SetInt(it)

	for _, item := range r1 {

		ib, _, err := item["balance"].(primitive.Decimal128).BigInt()
		if err != nil {
			return err
		}
		ibf := new(big.Float).SetInt(ib)
		dv := new(big.Float).Quo(ibf, itf)
		item["percentage"] = dv
	}

	if args.Raw != nil {
		*args.Raw = r1
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
