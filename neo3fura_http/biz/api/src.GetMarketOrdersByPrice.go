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
)

func (me *T) GetMarketOrdersByPrice(args struct {
	AssetHash  h160.T
	MarketHash h160.T
	Token      h160.T
	MinAmount  *big.Int
	MaxAmount  *big.Int
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {

	isbig := args.MinAmount.Cmp(args.MaxAmount)

	if !args.AssetHash.Valid() || !args.Token.Valid() || isbig == 1 || !args.MarketHash.Valid() {
		return stderr.ErrInvalidArgs
	}

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
				bson.M{"$match": bson.M{"asset": args.AssetHash, "market": args.MarketHash, "deadline": bson.M{"$gt": currenr()}, "amount": bson.M{"$gt": 0}, "auctionType": bson.M{"$eq": 1}}},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	//获取指定最大最小的价格
	dd, _ := OpenAssetHashFile()
	tokenPrice := float64(1)
	if len(r5) > 0 {
		tokenPrice, err = GetPrice(args.Token.Val()) //  获取价格
		//防止获取价格的api limit 获取价格为0
		if tokenPrice == 0 {
			tokenPrice = float64(1)
		}
		if err != nil {
			return err
		}
	}

	tokenDecimal := dd[args.Token.Val()]

	minPrice := GetAsset2Price(args.MinAmount, tokenDecimal, tokenPrice)
	maxPrice := GetAsset2Price(args.MaxAmount, tokenDecimal, tokenPrice)

	for _, item := range r5 {
		auctionAsset := item["auctionAsset"].(string)
		auctionAmount, _, err2 := item["auctionAmount"].(primitive.Decimal128).BigInt()
		if err2 != nil {
			return err2
		}
		//amount, _ := new(big.Float).SetString(auctionAmount.String())
		//item["tokenAmount"] = amount //价格转换

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
	//minAmount,_,_ := new(big.Float).Parse(args.MinAmount.String(),10)
	//maxAmount,_,_ := new(big.Float).Parse(args.MaxAmount.String(),10)
	//找出价格区间的订单
	//minindex := FindIndexLeft(r5,"tokenAmount",minAmount)   // ->
	//maxindex := FindIndexRight(r5,"tokenAmount",maxAmount)    // <-

	minindex := FindIndexLeft(r5, "usdAmount", minPrice.Sub(minPrice, big.NewFloat(1))) // ->
	maxindex := FindIndexRight(r5, "usdAmount", maxPrice)                               // <-

	result := r5[minindex:maxindex]

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}

func GetAsset2Price(amount *big.Int, assetDecimal int64, unitprice float64) *big.Float {
	bfAmount := new(big.Float).SetInt(amount)
	flag := amount.Cmp(big.NewInt(0))
	if flag == 1 {
		bfprice := big.NewFloat(unitprice)
		ffprice := big.NewFloat(1).Mul(bfprice, bfAmount)
		de := math.Pow(10, float64(assetDecimal))
		usdAuctionAmount := new(big.Float).Quo(ffprice, big.NewFloat(float64(de)))

		return usdAuctionAmount
	} else {
		return big.NewFloat(0)
	}

}

func FindIndexLeft(arr []map[string]interface{}, key string, target *big.Float) int {
	for i := 0; i < len(arr); i++ {
		if target.Cmp(arr[0][key].(*big.Float)) != 1 {
			return 0
		}
		if i == len(arr)-1 {
			return i
		}
		if arr[i][key].(*big.Float).Cmp(target) != 1 && target.Cmp(arr[i+1][key].(*big.Float)) == -1 { //target >=(*arr)[i] && target<= (*arr)[i+1]
			return i + 1
		}
	}
	return -1
}

func FindIndexRight(arr []map[string]interface{}, key string, target *big.Float) int {
	for i := 0; i < len(arr); i++ {
		if target.Cmp(arr[0][key].(*big.Float)) == -1 {
			return 0
		}
		if target.Cmp(arr[0][key].(*big.Float)) == 0 {
			return 1
		}
		if i == len(arr)-1 {
			return len(arr)
		}

		if target.Cmp(arr[i][key].(*big.Float)) == 1 && arr[i+1][key].(*big.Float).Cmp(target) != -1 { //target >(*arr)[i] && target< =(*arr)[i+1]
			if target.Cmp(arr[i][key].(*big.Float)) == 1 {
				return i + 2
			}
			return i + 1
		}
	}
	return -1

}
