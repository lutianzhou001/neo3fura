package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetMarketFloorPriceByAsset(args struct {
	AssetHash h160.T
	Market    h160.T
	Filter    map[string]interface{}
	Raw       *map[string]interface{}
}, ret *json.RawMessage) error {
	//currentTime := time.Now().UnixMilli()
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Market.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	result := make(map[string]interface{})
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
				bson.M{"$match": bson.M{"asset": args.AssetHash, "market": args.Market, "deadline": bson.M{"$gt": currenr()}, "amount": bson.M{"$gt": 0}, "auctionType": bson.M{"$eq": 1}}},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r5 {
		auctionAsset := item["auctionAsset"].(string)
		auctionAmount, _, err2 := item["auctionAmount"].(primitive.Decimal128).BigInt()
		if err2 != nil {
			return err2
		}

		//价格转换
		dd, _ := OpenAssetHashFile()
		decimal := dd[auctionAsset]           //获取精度
		price, err3 := GetPrice(auctionAsset) //  获取价格
		if err3 != nil {
			return err3
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
			item["usdAmount"] = usdAuctionAmount
		} else {
			item["usdAmount"] = big.NewFloat(0)
		}

	}
	mapsort.MapSort7(r5, "usdAmount")

	if len(r5) > 0 {
		result = r5[0]

	} else {
		result = nil

	}

	if err != nil {
		return err
	}

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}

func currenr() int64 {
	return time.Now().UnixMilli()
}
