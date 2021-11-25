package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/utils"
	"neo3fura_http/var/stderr"
	"sort"
)

func (me *T) GetAssetHoldersByContractHash(args struct {
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
		Index:      "GetAssetHoldersByContractHash",
		Sort:       bson.M{"balance": -1},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)

	//Nep11
	if count != 0 && r1[0]["tokenid"] != "" {
		var raw1 map[string]interface{}
		err = me.GetAssetInfoByContractHash(struct {
			ContractHash h160.T
			Filter       map[string]interface{}
			Raw          *map[string]interface{}
		}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
		if err != nil {
			return err
		}

		it, _, err := raw1["totalsupply"].(primitive.Decimal128).BigInt()
		if err != nil {
			return err
		}
		itf := new(big.Float).SetInt(it)

		var groups = utils.GroupBy(r1, "address")
		var holders = []Nep11Holder{}
		for _, items := range groups {
			var bal int64 = 0
			tokenid := []string{}
			for _, item := range items {
				tid := item["tokenid"]
				tokenid = append(tokenid, tid.(string))
				bal++
			}

			var b2 *big.Float = big.NewFloat(float64(bal))

			dv := new(big.Float).Quo(b2, itf)

			holder := Nep11Holder{
				Address:    items[0]["address"].(string),
				Balance:    bal,
				TokenId:    tokenid,
				Percentage: dv,
			}
			holders = append(holders, holder)
		}
		sort.Sort(Nep11HolderByBalance(holders))

		result := make(map[string]interface{})
		result["asset"] = args.ContractHash
		result["holder"] = holders
		var results []map[string]interface{}
		results = append(results, result)

		if args.Raw != nil {
			*args.Raw = results
		}
		r2, err := me.FilterArrayAndAppendCount(results, count, args.Filter)
		if err != nil {
			return err
		}
		r, err := json.Marshal(r2)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
	} else {
		//Nep17
		for _, item := range r1 {
			var raw1 map[string]interface{}
			err = me.GetAssetInfoByContractHash(struct {
				ContractHash h160.T
				Filter       map[string]interface{}
				Raw          *map[string]interface{}
			}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
			if err != nil {
				return err
			}
			ib, _, err := item["balance"].(primitive.Decimal128).BigInt()
			if err != nil {
				return err
			}
			// it, _ := new(big.Int).SetString(raw1["totalsupply"].(string), 10)
			it, _, err := raw1["totalsupply"].(primitive.Decimal128).BigInt()

			if err != nil {
				return err
			}
			ibf := new(big.Float).SetInt(ib)
			itf := new(big.Float).SetInt(it)
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
	}
	return nil
}
