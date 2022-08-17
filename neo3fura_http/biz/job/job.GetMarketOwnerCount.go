package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	log2 "neo3fura_http/lib/log"
	limit "neo3fura_http/var/const"
	"os"
	"time"
)

func (me T) GetMarketOwnerCount() {
	message := make(json.RawMessage, 0)
	ret := &message
	currentTime := time.Now().UnixNano() / 1e6

	rt := os.ExpandEnv("${RUNTIME}")
	var secondMarketHash string
	if rt == "staging" {
		secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	} else if rt == "test2" {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	} else {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	}

	result := make(map[string]interface{})

	assetList, err2 := me.GetNep11Asset()
	if err2 != nil {
		log2.Fatal("GetMarketNep11Asset err")
	}

	for _, it := range assetList {
		//获取上架记录
		r2, err := me.Client.QueryAggregate(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
				Pipeline   []bson.M
				Query      []string
			}{
				Collection: "Market",
				Index:      "someindex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"owner": bson.M{"$ne": limit.NullAddress}, "asset": it, "market": secondMarketHash, "amount": bson.M{"$gt": 0}}}, //上架（正常状态、过期）:auctor，未领取：bidder
					bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1}},
					bson.M{"$match": bson.M{"difference": true}},
				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("GetListedOwner err ")
		}
		owner := make(map[string]interface{})
		for _, item := range r2 {
			bidAmount, _, err2 := item["bidAmount"].(primitive.Decimal128).BigInt()
			bidAmountFlag := bidAmount.Cmp(big.NewInt(0))
			//bidAmount, err2 := strconv.ParseInt(ba, 10, 64)
			if err2 != nil {
				log2.Fatal("OwnerCount: Covert err")

			}
			deadline, _ := item["deadline"].(int64)
			if item["owner"] == item["market"] && deadline > currentTime { //在售
				item["account"] = item["auctor"]
			} else if bidAmountFlag == 1 && deadline < currentTime && item["owner"] == item["market"] { //未领取
				item["account"] = item["bidder"]
			} else if deadline < currentTime && bidAmountFlag == 0 && item["owner"] == item["market"] { //过期
				item["account"] = item["auctor"]
			} else {
				item["account"] = ""
			}
			owner[item["account"].(string)] = 1
		}

		//二级市场未上架数据
		r3, err := me.Client.QueryAggregate(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
				Pipeline   []bson.M
				Query      []string
			}{
				Collection: "Market",
				Index:      "someindex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"owner": bson.M{"$ne": limit.NullAddress}, "amount": bson.M{"$gt": 0}, "asset": it, "market": bson.M{"$in": []interface{}{limit.NullAddress, nil}}}}, //上架（正常状态、过期）:auctor，未领取：biddernu
					bson.M{"$group": bson.M{"_id": "$owner"}},
					//bson.M{"$count":"count"},
				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("GetNotListedOwner err ")
		}

		if len(r3) > 0 {
			for _, item := range r3 {
				owner[item["_id"].(string)] = 1
			}

		}

		result["totalowner"] = int32(len(owner))

		// 存到本地数据库中
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketIndex", Data: result, Filter: bson.M{"asset": it}})
		if err != nil {
			log2.Fatal("MarketIndex update err")
		}

	}

}
