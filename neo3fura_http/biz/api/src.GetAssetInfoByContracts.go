package api

import (
	"encoding/json"
	"fmt"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetAssetInfoByContracts(args struct {
	ContractHash []h160.T
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
}, ret *json.RawMessage) error {
	if len(args.ContractHash) < 1 {
		return stderr.ErrInvalidArgs
	}

	var hashArr []interface{}
	for _, item := range args.ContractHash {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			hashArr = append(hashArr, item)
		}
	}

	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Asset",
		Index:      "GetAssetInfoByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"hash": bson.M{"$in": hashArr}, "totalsupply": bson.M{"$gt": 0}}},
			bson.M{"$lookup": bson.M{
				"from": "Address-Asset",
				"let":  bson.M{"asset": "$hash"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						bson.M{"$gt": []interface{}{"$balance", 0}},
					}}}},
					bson.M{"$group": bson.M{"_id": "$address"}},
					bson.M{"$count": "count"},
				},
				"as": "addressCount"},
			},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	r22, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}

	for _, item := range r1 {
		raw1 := make(map[string]interface{})
		if item["type"] == "Unknown" {
			err := me.GetContractByContractHash(struct {
				ContractHash h160.T
				Filter       map[string]interface{}
				Raw          *map[string]interface{}
			}{ContractHash: h160.T(fmt.Sprint(item["hash"])), Filter: nil, Raw: &raw1}, ret)
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
				item["type"] = "NEP17"
			}
			if i == 6 {
				item["type"] = "NEP11"
			}
		}

		// holders

		addressCount := item["addressCount"].(primitive.A)
		if len(addressCount) > 0 {
			count := addressCount[0].(map[string]interface{})["count"]
			item["holders"] = count
		}
		delete(item, "addressCount")

		item["ispopular"] = false
		if r22["Populars"] != nil {
			populars := r22["Populars"].(primitive.A)
			for _, v := range populars {
				if item["hash"] == v {
					item["ispopular"] = true
				}
			}
		}

	}

	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, int64(len(args.ContractHash)), args.Filter)
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
