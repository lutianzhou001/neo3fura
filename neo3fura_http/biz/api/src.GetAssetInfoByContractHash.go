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
	r3, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "Holders"})
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

	holders := r3["Holders"].(primitive.A)
	for _, h := range holders {
		m := h.(map[string]interface{})
		for k, v := range m {
			if r1["hash"] == k {
				r1["holders"] = v
			}
		}
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
