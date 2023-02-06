package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/lib/type/consts"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetTagByAddresses(args struct {
	Address []h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if len(args.Address) < 1 {
		return stderr.ErrInvalidArgs
	}

	var addressArr []interface{}
	for _, item := range args.Address {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			addressArr = append(addressArr, item)
		}
	}
	currentTime := time.Now().UnixNano() / 1e6
	before7days := currentTime - 7*24*60*60*1000

	///---------FT(NEO,bNEO)------------
	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "TransferNotification",
		Index:      "someindex",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"contract": bson.M{"$in": []interface{}{consts.NEO, consts.BNEO_Main, consts.BNEO_Test}},
				"from":      bson.M{"$in": addressArr},
				"timestamp": bson.M{"$gt": before7days}}},

			//  地址仅有一次出现在from/to间，一笔交易会被计算2次
			//bson.M{"$project": bson.M{"_id": 1, "txid": 1, "blockhash": 1, "from": 1, "to": 1, "timestamp": 1, "value": 1,
			//	"address": bson.M{"$map": bson.M{"input": bson.M{"$literal": []interface{}{"p1", "p2"}},
			//		"as": "p",
			//		"in": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$$p", "p1"}}, "$from", "$to"}},
			//	}},
			//}},
			//bson.M{"$unwind": "$address"},
			bson.M{"$group": bson.M{"_id": bson.M{"contract": "$contract", "address": "$from"}, "sum": bson.M{"$sum": "$value"}}},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	r2, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "TransferNotification",
		Index:      "someindex",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"contract": bson.M{"$in": []interface{}{consts.NEO, consts.BNEO_Main, consts.BNEO_Test}},
				"to":        bson.M{"$in": addressArr},
				"timestamp": bson.M{"$gt": before7days}}},
			bson.M{"$group": bson.M{"_id": bson.M{"contract": "$contract", "address": "$to"}, "sum": bson.M{"$sum": "$value"}}},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	///---------NFT-------------
	r3, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Address-Asset",
		Index:      "someindex",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"address": bson.M{"$in": addressArr}, "balance": bson.M{"$gt": 0}, "tokenid": bson.M{"$ne": ""}}},
			bson.M{"$group": bson.M{"_id": bson.M{"address": "$address", "count": bson.M{"$sum": 1}}}},
		},
		Query: []string{},
	}, ret)
	fmt.Println(r3)
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range args.Address {
		info := make(map[string]interface{})
		info["address"] = item
		info["neoSum"] = big.NewInt(0)
		info["bneoSum"] = big.NewInt(0)
		///-----------FT--------------
		//from
		for _, fromItem := range r1 {
			address_contract := fromItem["_id"].(map[string]interface{})
			if address_contract["address"].(string) == item.Val() {
				if address_contract["contract"].(string) == consts.NEO {
					neoSum, _, err := fromItem["sum"].(primitive.Decimal128).BigInt()
					if err != nil {
						return err
					}
					info["neoSum"] = neoSum
				} else if address_contract["contract"].(string) == consts.BNEO_Main || address_contract["contract"].(string) == consts.BNEO_Test {
					bneoSum, _, err := fromItem["sum"].(primitive.Decimal128).BigInt()
					if err != nil {
						return err
					}
					info["bneoSum"] = bneoSum
				}
			}
		}

		//to
		for _, toItem := range r2 {
			address_contract := toItem["_id"].(map[string]interface{})
			fmt.Println(address_contract["address"], ",", item, ",", address_contract["address"] == item)
			if address_contract["address"].(string) == item.Val() {
				if address_contract["contract"].(string) == consts.NEO {
					neoSum, _, err := toItem["sum"].(primitive.Decimal128).BigInt()
					if err != nil {
						return err
					}
					if info["neoSum"] != nil {
						info["neoSum"] = big.NewInt(0).Add(neoSum, info["neoSum"].(*big.Int))
					} else {
						info["neoSum"] = neoSum
					}
				} else if address_contract["contract"].(string) == consts.BNEO_Main || address_contract["contract"].(string) == consts.BNEO_Test {
					bneoSum, _, err := toItem["sum"].(primitive.Decimal128).BigInt()
					if err != nil {
						return err
					}
					if info["bneoSum"] != nil {
						info["bneoSum"] = big.NewInt(0).Add(bneoSum, info["bneoSum"].(*big.Int))
					} else {
						info["bneoSum"] = bneoSum
					}
				}
			}
		}
		bneo := new(big.Float).SetInt(info["bneoSum"].(*big.Int))
		bneo2neo := new(big.Float).Quo(bneo, big.NewFloat(math.Pow10(8)))
		total := new(big.Float).Add(new(big.Float).SetInt(info["neoSum"].(*big.Int)), bneo2neo)
		info["ft_total"] = total

		if total.Cmp(big.NewFloat(50)) == 1 {
			info["ft_tag"] = "Semi-Experienced Trader"
		} else if total.Cmp(big.NewFloat(1000)) == 1 {
			info["ft_tag"] = "Experienced Trader"
		} else {
			info["ft_tag"] = ""
		}

		///-----------NFT--------------
		info["nft_total"] = int32(0)
		for _, nftItem := range r3 {
			nftinfo := nftItem["_id"].(map[string]interface{})
			if nftinfo["address"].(string) == item.Val() {
				info["nft_total"] = nftinfo["count"].(int32)
			}
		}
		nft_total := info["nft_total"].(int32)
		if nft_total > 0 {
			info["nft_tag"] = "NFT Player"
		} else if nft_total > 9 {
			info["nft_tag"] = "NFT Killer"
		} else if nft_total > 99 {
			info["nft_tag"] = "NFT Whale"
		} else {
			info["nft_tag"] = ""
		}

		result = append(result, info)
	}

	if err != nil {
		return err
	}
	r4, err := me.FilterArrayAndAppendCount(result, int64(len(result)), args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
