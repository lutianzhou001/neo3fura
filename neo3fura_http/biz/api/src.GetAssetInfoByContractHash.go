package api

import (
	"encoding/json"
	"fmt"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetAssetInfoByContractHash(args struct {
	ContractHash h160.T
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
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
		Index:      "GetAssetInfoByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash.Val(), "totalsupply": bson.M{"$gt": 0}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
	}

	raw1 := make(map[string]interface{})
	if r1["type"] == "Unknown" {
		err := me.GetContractByContractHash(struct {
			ContractHash h160.T
			Filter       map[string]interface{}
			Raw          *map[string]interface{}
		}{ContractHash: h160.T(fmt.Sprint(r1["hash"])), Filter: nil, Raw: &raw1}, ret)
		if err != nil {
			return nil
		}
		m := make(map[string]interface{})
		json.Unmarshal([]byte(raw1["manifest"].(string)), &m)
		methods := m["abi"].(map[string]interface{})["methods"].([]interface{})
		i := 0
		for _, method := range methods {
			if method.(map[string]interface{})["name"].(string) == "transfer" {
				i = i + 1
			}
			if (method.(map[string]interface{})["name"].(string) == "transfer") && len(method.(map[string]interface{})["parameters"].([]interface{})) == 4 {
				i = i + 1
			}
			if (method.(map[string]interface{})["name"].(string) == "transfer") && len(method.(map[string]interface{})["parameters"].([]interface{})) == 3 {
				i = i + 2
			}
			if method.(map[string]interface{})["name"].(string) == "balanceOf" {
				i = i + 1
			}
			if method.(map[string]interface{})["name"].(string) == "totalSupply" {
				i = i + 1
			}
			if method.(map[string]interface{})["name"].(string) == "decimals" {
				i = i + 1
			}
		}

		if i == 5 {
			r1["type"] = "NEP17"
		}
		if i == 6 {
			r1["type"] = "NEP11"
		}
	}

	r2, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}

	r1["ispopular"] = false
	populars := r2["Populars"].(primitive.A)
	for _, v := range populars {
		if r1["hash"] == v {
			r1["ispopular"] = true
		}
	}

	_, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Address-Asset",
		Index:      "GetAssetInfos",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
	}
	r1["holders"] = count
	totalsuply, _, err := r1["totalsupply"].(primitive.Decimal128).BigInt()
	if err != nil {
		return err
	}
	r1["totalsupply"] = totalsuply
	if r1["type"].(string) == "NEP11" {
		r3, err1 := me.Client.QueryAggregate(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
				Pipeline   []bson.M
				Query      []string
			}{
				Collection: "Address-Asset",
				Index:      "GetContractList",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"asset": args.ContractHash.Val(), "balance": bson.M{"$gt": 0}}},
					bson.M{"$group": bson.M{"_id": "$address"}},
					bson.M{"$count": "addressCounts"},
				},
				Query: []string{},
			}, ret)
		if err1 != nil {
			return err1
		}

		r1["holders"] = r3[0]["addressCounts"]
	}

	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}

	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
