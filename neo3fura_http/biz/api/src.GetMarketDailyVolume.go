package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"

	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/consts"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetMarketDailyVolume(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	//currentTime := time.Now().UnixNano() / 1e6

	rt := os.ExpandEnv("${RUNTIME}")
	fmt.Println(rt)
	//var secondMarketHash string
	//if rt == "staging" {
	//	secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	//} else if rt == "test2" {
	//	secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	//} else {
	//	secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	//}

	result := make(map[string]interface{})

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
					bson.M{"$match": bson.M{"asset": it, "eventname": bson.M{"$in": []interface{}{"Claim", "CompleteOffer"}}}}, //交易:售卖，竞拍，offer
					bson.M{"$group": bson.M{"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": bson.M{"$toDate": "$timestamp"}}},
						"count":  bson.M{"$sum": 1},
						"events": bson.M{"$push": "$$ROOT"},
					}},
				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("Get Market transaction err: ", err)
		}

		for _, item := range r2 {
			assertResult := make(map[string]interface{})
			date := item["_id"].(string)
			events := item["events"].(primitive.A)
			dayVolume := big.NewFloat(0) //日交易额   bneo
			dayAmount := len(events)     //日交易数量

			dateTime, _ := time.Parse(consts.ShortForm, date)
			assertResult["asset"] = it
			assertResult["date"] = primitive.NewDateTimeFromTime(dateTime)

			assertResult["dayAmount"] = dayAmount

			for _, event := range events {
				e := event.(map[string]interface{})
				eventname := e["eventname"].(string)
				extendData := e["extendData"].(string)
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
							log2.Fatal("tokenCOnversion err:", err)
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
			assertResult["dayVolume"] = dv

			p := dayVolume.Quo(dayVolume, big.NewFloat(float64(dayAmount)))
			ap, err := primitive.ParseDecimal128(p.String())
			assertResult["avgPrice"] = ap

			//存到本地数据库中
			_, err = me.Client.UpdateJob(struct {
				Collection string
				Data       bson.M
				Filter     bson.M
			}{Collection: "MarketDayVolume", Data: assertResult, Filter: bson.M{"asset": it, "date": dateTime}})
			if err != nil {
				log2.Fatal("MarketIndex update err")
			}
		}

		//填充历史交易量为空的数据

		//获取合约创建时间
		r0, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Contract",
			Index:      "GetAddressInfoByAsset",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": it},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}
		var createtime int64
		ct := r0["createtime"]
		switch ct.(type) {
		case float64:
			time := ct.(float64)
			createtime = int64(time)
		case int64:
			createtime = ct.(int64)
		case string:
			createtime, _ = strconv.ParseInt(ct.(string), 10, 64)
		}

		tt := time.UnixMilli(createtime)
		newtime := time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, tt.Location())
		createtime0 := newtime.UnixMilli() //当天0点时间戳
		today := time.Now()
		todaytime := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		todaytime0 := todaytime.UnixMilli() //当天0点时间戳

		for i := int64(createtime0); i <= todaytime0; i += 24 * 60 * 60 * 1000 {
			ct := time.UnixMilli(i).Format(consts.ShortForm)

			dateTime, _ := time.Parse(consts.ShortForm, ct)
			//查询有没有数据
			r, err := me.Client.QueryOneJob(struct {
				Collection string
				Filter     bson.M
			}{Collection: "MarketDayVolume", Filter: bson.M{"asset": it, "date": dateTime}})

			if r == nil && err.Error() == "mongo: no documents in result" {
				data := make(map[string]interface{})
				data["asset"] = it
				data["date"] = dateTime
				fmt.Println(dateTime)
				avg, _ := primitive.ParseDecimal128("0")
				data["avgPrice"] = avg
				data["dayAmount"] = int32(0)
				dv, _ := primitive.ParseDecimal128("0")
				data["dayVolume"] = dv

				_, err = me.Client.UpdateJob(struct {
					Collection string
					Data       bson.M
					Filter     bson.M
				}{Collection: "MarketDayVolume", Data: data, Filter: bson.M{"asset": it, "date": dateTime}})
				if err != nil {
					log2.Fatal("MarketDayVolume update err")
				}
			}

		}
	}

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}

//获取二级市场白名单中的所有NEP11资产
func (me T) GetNep11Asset() ([]string, error) {
	rt := os.ExpandEnv("${RUNTIME}")
	var secondMarketHash string
	if rt == "staging" {
		secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	} else if rt == "test2" {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	} else {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": secondMarketHash}},
		bson.M{"$match": bson.M{"$or": []interface{}{
			bson.M{"eventname": "AddAsset"},
			bson.M{"eventname": "RemoveAsset"},
		}}},
		bson.M{"$sort": bson.M{"timestamp": 1}},
		bson.M{"$group": bson.M{"_id": "$asset", "asset": bson.M{"$last": "$asset"}, "eventname": bson.M{"$last": "$eventname"}}},

		bson.M{"$lookup": bson.M{
			"from": "Asset",
			"let":  bson.M{"hash": "$asset"},
			"pipeline": []bson.M{

				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$hash", "$$hash"}},
					bson.M{"$eq": []interface{}{"$type", "NEP11"}},
				}}}},
			},
			"as": "type"}},
	}
	message := make(json.RawMessage, 0)
	ret := &message
	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return nil, err
	}

	var assetArr []string
	for _, item := range r1 {
		assetType := item["type"].(primitive.A)
		if item["eventname"].(string) == "AddAsset" && len(assetType) > 0 {
			assetArr = append(assetArr, item["asset"].(string))
		}
	}

	return assetArr, nil
}

//token
func TokenConversion(from string, amount *big.Int, to string) (*big.Float, error) {
	if from == to {
		return new(big.Float).SetInt(amount), nil
	}

	fromPrice, err := GetPrice(from)

	if err != nil {
		return big.NewFloat(0), err
	}
	toPrice, err := GetPrice(to)

	dd, _ := OpenAssetHashFile()
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
