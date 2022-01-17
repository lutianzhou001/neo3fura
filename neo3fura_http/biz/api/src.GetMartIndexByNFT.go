package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/big"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/utils"
	"neo3fura_http/var/stderr"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func (me *T) GetMarketIndexByAsset(args struct {
	MarketHash h160.T
	AssetHash  h160.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {

	currentTime := time.Now().UnixNano() / 1e6
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	result := make(map[string]interface{})
	var r1, err = me.Client.QueryAggregate(
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
				bson.M{"$match": bson.M{"asset": args.AssetHash.Val(), "amount": bson.M{"$gt": 0}}},
				bson.M{"$group": bson.M{"_id": "$tokenid"}},
				bson.M{"$count": "count"},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}
	if len(r1) > 0 {
		result["totalsupply"] = r1[0]["count"]
	} else {
		result["totalsupply"] = 0
	}

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
				bson.M{"$match": bson.M{"asset": args.AssetHash.Val(), "market": args.MarketHash.Val(), "amount": bson.M{"$gt": 0}}}, //上架（正常状态、过期）:auctor，未领取：bidder
				bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1}},
				bson.M{"$match": bson.M{"difference": true}},
			},
			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}
	owner := make([]map[string]interface{}, 0)
	for _, item := range r2 {
		bidAmount, _, err2 := item["bidAmount"].(primitive.Decimal128).BigInt()
		bidAmountFlag := bidAmount.Cmp(big.NewInt(0))
		//bidAmount, err2 := strconv.ParseInt(ba, 10, 64)
		if err2 != nil {
			return err
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
		owner = append(owner, item)
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
				bson.M{"$match": bson.M{"asset": args.AssetHash.Val(), "market": primitive.Null{}}}, //上架（正常状态、过期）:auctor，未领取：biddernu
				bson.M{"$group": bson.M{"_id": "$owner",
					"owner":    bson.M{"$last": "$owner"},
					"auctor":   bson.M{"$last": "$auctor"},
					"bidder":   bson.M{"$last": "$bidder"},
					"deadline": bson.M{"$last": "$deadline"},
					"market":   bson.M{"$last": "$market"},
				}},
			},
			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	if len(r3) > 0 {
		for _, item := range r3 {
			item["account"] = item["owner"]
			owner = append(owner, item)
		}

	}
	ownerGroup := utils.GroupBy(owner, "account") // owner 分组
	ownerCount := len(ownerGroup)
	result["totalowner"] = ownerCount

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
				bson.M{"$match": bson.M{"asset": args.AssetHash.Val(), "market": args.MarketHash, "eventname": "Claim"}},
			},
			Query: []string{"extendData"},
		}, ret)

	if err != nil {
		return err
	}

	var txAmount float64
	if len(r4) > 0 {
		for _, item := range r4 {
			extendData := item["extendData"].(string)
			if extendData != "" {
				var data map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
					auctionAsset := data["auctionAsset"].(string)
					dd, _ := OpenAssetHashFile()
					decimal := dd[auctionAsset]
					if decimal == 0 {
						decimal = 1
					}
					bidAmount, err2 := strconv.ParseInt(data["bidAmount"].(string), 10, 64)
					if err2 != nil {
						return err2
					}

					price, err3 := GetPrice(auctionAsset) //
					if err3 != nil {
						return err3
					}

					if price == 0 {
						price = 1
					}

					txprice := float64(bidAmount) * price / float64(decimal)
					txAmount += txprice

				} else {
					return err1
				}
			}
		}
	} else {
		txAmount = 0
	}

	result["totaltxamount"] = txAmount
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
				bson.M{"$match": bson.M{"asset": args.AssetHash.Val(), "market": args.MarketHash.Val(), "auctionType": bson.M{"$eq": 1}}},
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
			return err
		}

		//价格转换
		dd, _ := OpenAssetHashFile()
		decimal := dd[auctionAsset] //获取精度
		if decimal == 0 {
			decimal = 1
		}
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
			usdAuctionAmount := new(big.Float).Quo(ffprice, big.NewFloat(float64(decimal)))
			item["usdAmount"] = usdAuctionAmount
		} else {
			item["usdAmount"] = 0
		}

	}
	mapsort.MapSort7(r5, "usdAmount")

	if len(r5) > 0 {
		result["auctionAsset"] = r5[0]["auctionAsset"]
		result["auctionAmount"] = r5[0]["auctionAmount"]
		result["usdAmount"] = r5[0]["usdAmount"]
	} else {
		result["auctionAsset"] = "——"
		result["auctionAmount"] = "——"
		result["usdAmount"] = "——"
	}

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}

func GetPrice(asset string) (float64, error) {

	client := &http.Client{}
	reqBody := []byte(`["` + asset + `"]`)
	url := "https://onegate.space/api/quote?convert=usd"
	//str :=[]string{asset}
	req, _ :=
		http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, stderr.ErrPrice
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, stderr.ErrPrice
	}
	response := string(body)
	re := response[1 : len(response)-1]
	price, err1 := strconv.ParseFloat(re, 64)
	if err1 != nil {
		return 0, stderr.ErrPrice
	}
	return price, nil
}

func OpenAssetHashFile() (map[string]int64, error) {
	absPath, _ := filepath.Abs("./assethash.json")

	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Print(err)
	}
	whitelist := map[string]int64{}
	err = json.Unmarshal([]byte(string(b)), &whitelist)
	if err != nil {
		panic(err)
	}

	return whitelist, err
}
