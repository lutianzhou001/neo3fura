package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/h160"
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

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Asset",
		Index:      "GetAssetInfo",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash.Val()},
		Query:      []string{},
	}, ret)
	asset_type := r1["type"]
	asset_totalsupply, _, err := r1["totalsupply"].(primitive.Decimal128).BigInt()
	if err != nil {
		return err
	}

	var pipeline []bson.M

	if asset_type == "NEP11" {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}}},
			bson.M{"$group": bson.M{"_id": "$address", "count": bson.M{"$sum": 1}, "tokenidArr": bson.M{"$push": "$$ROOT"}}},
			bson.M{"$sort": bson.M{"count": -1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
	} else if asset_type == "NEP17" {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}}},
			bson.M{"$sort": bson.M{"balance": -1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
	} else {
		return stderr.ErrAssetType
	}

	r2, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Address-Asset",
		Index:      "GetAssetHoldersByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline:   pipeline,
		Query:      []string{},
	}, ret)

	if asset_type == "NEP11" {
		for _, item := range r2 {
			item["address"] = item["_id"]
			balance := item["count"].(int32)
			item["balance"] = balance
			item["asset"] = args.ContractHash
			arr := item["tokenidArr"].(primitive.A)
			var tokenidArr []string
			for _, it := range arr {
				tokenid := it.(map[string]interface{})
				tokenidArr = append(tokenidArr, tokenid["tokenid"].(string))
			}
			item["tokenid"] = tokenidArr

			itf := new(big.Float).SetInt(asset_totalsupply)
			var b2 *big.Float = big.NewFloat(float64(balance))
			dv := new(big.Float).Quo(b2, itf)
			item["percentage"] = dv

			delete(item, "_id")
			delete(item, "count")
			delete(item, "tokenidArr")

		}

	} else if asset_type == "NEP17" {
		for _, item := range r2 {
			balance, _, err := item["balance"].(primitive.Decimal128).BigInt()
			if err != nil {
				return err
			}
			itf := new(big.Float).SetInt(asset_totalsupply)
			var b2 *big.Float = new(big.Float).SetInt(balance)
			dv := new(big.Float).Quo(b2, itf)
			item["percentage"] = dv
		}

	}

	count, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Address-Asset",
		Index:      "GetAssetHoldersByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}},
	}, ret)
	if err != nil {
		return err
	}

	r3, err := me.FilterArrayAndAppendCount(r2, count["total counts"].(int64), args.Filter)
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
