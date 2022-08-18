package job

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/consts"
	"os"
)

//每小时更新获取当天的交易数据
func (me T) GetMarketHourlyVolume() {
	message := make(json.RawMessage, 0)
	ret := &message
	t := time.Now()
	newtime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	today := newtime.UnixMilli() //当天0点时间戳

	rt := os.ExpandEnv("${RUNTIME}")
	var secondMarketHash string
	if rt == "staging" {
		secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	} else if rt == "test2" {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	} else {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	}

	assetList, err2 := me.GetNep11Asset()
	if err2 != nil {
		log2.Fatal("GetMarketNep11Asset err")
	}

	for _, it := range assetList {

		r2, err := me.Client.QueryAggregate(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
				Pipeline   []bson.M
				Query      []string
			}{
				Collection: "MarketNotification",
				Index:      "someindex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"asset": it, "market": secondMarketHash,
						"eventname":  bson.M{"$in": []interface{}{"Claim", "CompleteOffer"}},
						"$timestamp": bson.M{"$gt": today},
					}}, //获取前一天的交易

				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("Get Market transaction err: ", err)
		}
		assetResult := make(map[string]interface{})
		date := time.UnixMilli(today).Format(consts.ShortForm)
		dateTime, _ := time.Parse(consts.ShortForm, date)
		assetResult["date"] = dateTime
		assetResult["asset"] = it
		dayVolume := big.NewFloat(0)
		dayAmount := len(r2)
		for _, item := range r2 {
			eventname := item["eventname"].(string)
			extendData := item["extendData"].(string)
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &data); err == nil {
				var toAmount *big.Float
				if eventname == "Claim" {
					auctionAsset := data["auctionAsset"].(string)
					auctionAmount := data["bidAmount"].(string)
					amount, _ := new(big.Int).SetString(auctionAmount, 10)
					if rt == "staging" {
						toAmount, err = TokenConversion(auctionAsset, amount, consts.BNEO_Main)
					} else {
						toAmount, err = TokenConversion(auctionAsset, amount, consts.BNEO_Test)
					}
					if err != nil {
						log2.Fatal("tokenConversion err:", err)
					}

				} else if eventname == "CompleteOffer" {
					offerAsset := data["offerAsset"].(string)
					offerAmount := data["offerAmount"].(string)
					amount, _ := new(big.Int).SetString(offerAmount, 10)
					if rt == "staging" {
						toAmount, err = TokenConversion(offerAsset, amount, consts.BNEO_Main)
					} else {
						toAmount, err = TokenConversion(offerAsset, amount, consts.BNEO_Test)
					}
					fmt.Println("completeOffer :", toAmount)
					if err != nil {
						log2.Fatal("tokenCOnversion err:", err)
					}
				}
				dayVolume = dayVolume.Add(dayVolume, toAmount)
			}
		}

		dv, err := primitive.ParseDecimal128(dayVolume.String())
		assetResult["dayVolume"] = dv

		p := dayVolume.Quo(dayVolume, big.NewFloat(float64(dayAmount)))
		ap, err := primitive.ParseDecimal128(p.String())
		assetResult["avgPrice"] = ap
		//存到本地数据库中
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketDayVolume", Data: assetResult, Filter: bson.M{"asset": it, "date": dateTime}})
		if err != nil {
			log2.Fatal("MarketDayVolume update err")
		}

	}

}

// 获取当前区块时间
func (me T) GetBlockTime() int64 {
	message := make(json.RawMessage, 0)
	ret := &message
	r0, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Transaction", Index: "GetDailyTransactions", Sort: bson.M{"_id": -1}}, ret)
	if err != nil {
		return int64(0)
	}
	if len(r0) > 0 {
		return r0["blocktime"].(int64)
	}
	return time.Now().UnixNano() / 1e6
}
