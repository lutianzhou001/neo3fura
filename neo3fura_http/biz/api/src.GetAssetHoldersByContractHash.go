package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/utils"
	"neo3fura_http/var/stderr"
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

		//it := big.NewInt(raw1["totalsupply"].(int64))
		it := raw1["totalsupply"].(*big.Int)
		if err != nil {
			return err
		}
		itf := new(big.Float).SetInt(it)

		var groups = utils.GroupBy(r1, "address")
		holders := make([]map[string]interface{}, 0)
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

			holder := make(map[string]interface{})
			holder["address"] = items[0]["address"].(string)
			holder["balance"] = bal
			holder["tokenid"] = tokenid
			holder["percentage"] = dv
			holder["asset"] = args.ContractHash.Val()

			holders = append(holders, holder)
			mapsort.MapSort(holders, "balance")
		}
		if args.Raw != nil {
			*args.Raw = holders
		}
		if args.Limit == 0 {
			args.Limit = int64(math.Inf(1))
		}
		pagedHolders := make([]map[string]interface{}, 0)
		for i, item := range holders {
			if int64(i) < args.Skip {
				continue
			} else if int64(i) > args.Skip+args.Limit-1 {
				continue
			} else {
				pagedHolders = append(pagedHolders, item)
			}
		}
		r2, err := me.FilterArrayAndAppendCount(pagedHolders, int64(len(holders)), args.Filter)
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
		var raw1 map[string]interface{}
		err = me.GetAssetInfoByContractHash(struct {
			ContractHash h160.T
			Filter       map[string]interface{}
			Raw          *map[string]interface{}
		}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
		if err != nil {
			return err
		}
		//it, _, err := raw1["totalsupply"].(primitive.Decimal128).BigInt()
		it := raw1["totalsupply"].(*big.Int)
		itf := new(big.Float).SetInt(it)

		for _, item := range r1 {
			ib, _, err := item["balance"].(primitive.Decimal128).BigInt()
			if err != nil {
				return err
			}

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

		if args.Limit == 0 {
			args.Limit = int64(math.Inf(1))
		}

		pagedHolders := make([]map[string]interface{}, 0)
		for i, item := range r1 {
			if int64(i) < args.Skip {
				continue
			} else if int64(i) > args.Skip+args.Limit-1 {
				continue
			} else {
				pagedHolders = append(pagedHolders, item)
			}
		}

		r2, err := me.FilterArrayAndAppendCount(pagedHolders, count, args.Filter)
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
