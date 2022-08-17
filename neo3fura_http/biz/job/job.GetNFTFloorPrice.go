package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/biz/api"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/mapsort"
	"os"
	"time"
)

func (me T) GetNFTFloorPrice() {
	message := make(json.RawMessage, 0)
	ret := &message
	result := make(map[string]interface{})
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
		//地板价
		r5, err := me.Client.QueryAggregate(
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
					bson.M{"$match": bson.M{"asset": it, "market": secondMarketHash, "deadline": bson.M{"$gt": currentTime}, "amount": bson.M{"$gt": 0}, "auctionType": bson.M{"$eq": 1}}},
				},

				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("GetMarketNep11Asset err")
		}

		for _, item := range r5 {
			auctionAsset := item["auctionAsset"].(string)
			auctionAmount, _, err2 := item["auctionAmount"].(primitive.Decimal128).BigInt()
			if err2 != nil {
				log2.Fatal("FloorPrice:: data conversion err:", err)
			}

			//价格转换
			dd, _ := api.OpenAssetHashFile()
			decimal := dd[auctionAsset]               //获取精度
			price, err3 := api.GetPrice(auctionAsset) //  获取价格
			if err3 != nil {
				log2.Fatal("FloorPrice:: get price err:", err)
			}
			if price == 0 {
				price = 1
			}

			bfauctionAmount := new(big.Float).SetInt(auctionAmount)
			flag := auctionAmount.Cmp(big.NewInt(0))

			if flag == 1 {
				bfprice := big.NewFloat(price)
				ffprice := big.NewFloat(1).Mul(bfprice, bfauctionAmount)
				de := math.Pow(10, float64(decimal))
				usdAuctionAmount := new(big.Float).Quo(ffprice, big.NewFloat(float64(de)))
				item["usdAmount"] = usdAuctionAmount.String()
			} else {
				item["usdAmount"] = "0"
			}

		}
		mapsort.MapSort7(r5, "usdAmount")

		if len(r5) > 0 {
			result["auctionAsset"] = r5[0]["auctionAsset"].(string)
			result["auctionAmount"] = r5[0]["auctionAmount"].(primitive.Decimal128).String()
			result["usdAmount"] = r5[0]["usdAmount"]
		} else {
			result["auctionAsset"] = nil
			result["auctionAmount"] = 0
			result["usdAmount"] = 0
		}

		// 存到本地数据库中
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketIndex", Data: result, Filter: bson.M{"asset": it}})
		if err != nil {
			log2.Fatal("floorPrice: MarketIndex update err")
		}
	}

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
