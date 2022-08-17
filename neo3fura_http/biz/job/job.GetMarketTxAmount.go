package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"

	"math"
	"math/big"
	"neo3fura_http/biz/api"
	log2 "neo3fura_http/lib/log"
	"os"
)

func (me T) GetMarketTxAmount() {
	message := make(json.RawMessage, 0)
	ret := &message

	rt := os.ExpandEnv("${RUNTIME}")
	var primaryMarketHash string
	var secondMarketHash string
	if rt == "staging" {
		secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
		primaryMarketHash = "0xa41600dec34741b143c66f2d3448d15c7d79a0b7"
	} else if rt == "test2" {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
		primaryMarketHash = "0x6f1ef5147a00ebbb7de1cf82420485674c5c55bc"
	} else {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
		primaryMarketHash = "0x6f1ef5147a00ebbb7de1cf82420485674c5c55bc"
	}

	result := make(map[string]interface{})

	assetList, err2 := me.GetNep11Asset()
	if err2 != nil {
		log2.Fatal("GetMarketNep11Asset err")
	}

	for _, it := range assetList {
		//交易数额
		r4, err := me.Client.QueryAggregate(
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
					bson.M{"$match": bson.M{"asset": it, "market": bson.M{"$in": []interface{}{primaryMarketHash, secondMarketHash}}, "eventname": "Claim"}},
				},
				Query: []string{"extendData"},
			}, ret)

		if err != nil {
			log2.Errorf("GetTxCount ERROR", err)
		}

		var txAmount = big.NewFloat(0)

		if len(r4) > 0 {
			for _, item := range r4 {
				extendData := item["extendData"].(string)
				if extendData != "" {
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
						auctionAsset := data["auctionAsset"].(string)
						dd, _ := api.OpenAssetHashFile()
						decimal := dd[auctionAsset]

						ba := data["bidAmount"].(string)
						bidAmount, err2 := new(big.Int).SetString(ba, 10)
						if err2 == false {
							bidAmount = big.NewInt(0)
						}
						price, err3 := api.GetPrice(auctionAsset) //
						if err3 != nil {
							log2.Errorf("TxCount getPrice", err)
						}

						if price == 0 {
							price = 1
						}

						bfbidAmount := new(big.Float).SetInt(bidAmount)
						flag := bidAmount.Cmp(big.NewInt(0))

						if flag == 1 {
							bfprice := big.NewFloat(price)
							ffprice := big.NewFloat(1).Mul(bfprice, bfbidAmount)
							de := math.Pow(10, float64(decimal))
							usdbidAmount := new(big.Float).Quo(ffprice, big.NewFloat(float64(de)))
							txAmount = new(big.Float).Add(txAmount, usdbidAmount)

						} else {
							txAmount = new(big.Float).Add(txAmount, big.NewFloat(0))
						}
					} else {
						log2.Errorf("GetTxCount json err:", err)
					}
				}
			}
		} else {
			txAmount = new(big.Float).Add(txAmount, big.NewFloat(0))
		}

		result["totaltxamount"] = txAmount.String()

		// 存到本地数据库中
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketIndex", Data: result, Filter: bson.M{"asset": it}})
		if err != nil {
			log2.Fatal("totaltxamount: MarketIndex update err")
		}
	}

}
