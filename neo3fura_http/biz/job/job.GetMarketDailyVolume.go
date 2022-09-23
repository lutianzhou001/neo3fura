package job

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/biz/api"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/consts"
	"neo3fura_http/var/stderr"
	"os"
)

func (me T) GetMarketDailyVolume() {
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

	assetList, err2 := me.GetNep11Asset()
	if err2 != nil {
		log2.Fatal("GetMarketNep11Asset err")
	}

	for _, it := range assetList {
		//获取上架记录
		fmt.Println(it)
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
						"eventname": bson.M{"$in": []interface{}{"Claim", "CompleteOffer"}},
						"timestamp": bson.M{"$gte": currentTime - 24*60*60*1000},
					}}, //获取前一天的交易

				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("Get Market transaction err: ", err)
		}
		assetResult := make(map[string]interface{})
		lastday := currentTime - 60*60*1000
		date := time.UnixMilli(lastday).Format(consts.ShortForm)
		assetResult["date"] = date
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
		assetResult["dayAmount"] = dayAmount
		if dayAmount == 0 {
			ap, _ := primitive.ParseDecimal128("0")
			assetResult["avePrice"] = ap
		} else {
			p := dayVolume.Quo(dayVolume, big.NewFloat(float64(dayAmount)))
			ap, _ := primitive.ParseDecimal128(p.String())
			assetResult["avePrice"] = ap
		}

		//存到本地数据库中
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketDayVolume", Data: assetResult, Filter: bson.M{"asset": it, "date": date}})
		if err != nil {
			log2.Fatal("MarketDayVolume update err")
		}

	}

}

//token
func TokenConversion(from string, amount *big.Int, to string) (*big.Float, error) {
	if from == to {
		return new(big.Float).SetInt(amount), nil
	}

	fromPrice, err := api.GetPrice(from)

	if err != nil {
		return big.NewFloat(0), err
	}
	toPrice, err := api.GetPrice(to)

	dd, _ := api.OpenAssetHashFile()
	fromDecimal := dd[from]
	toDecimal := dd[to]

	if amount.String() != "0" {
		fromPrice2 := big.NewFloat(fromPrice)
		ffprice := big.NewFloat(1).Mul(fromPrice2, new(big.Float).SetInt(amount))
		fromDeci := math.Pow(10, float64(fromDecimal))
		if fromDeci == 0 {
			fromDeci = 1
		}
		usdtAmount := new(big.Float).Quo(ffprice, big.NewFloat(fromDeci))

		toDeci := math.Pow(10, float64(toDecimal))
		mind := big.NewFloat(1).Mul(usdtAmount, big.NewFloat(toDeci))

		if toPrice == 0 {
			toPrice = 1
		}
		toAmount := big.NewFloat(1).Quo(mind, big.NewFloat(toPrice))

		return toAmount, nil

	} else {
		return big.NewFloat(0), stderr.ErrInvalidArgs
	}

}
