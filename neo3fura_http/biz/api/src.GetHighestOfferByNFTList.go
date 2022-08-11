package api

import (
	"encoding/json"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/lib/utils"
	"neo3fura_http/var/stderr"
	"os"
	"strconv"
	"time"
)

func (me *T) GetHighestOfferByNFTList(args struct {
	NFT []struct {
		Asset   h160.T
		TokenId strval.T
	}
	MarketHash h160.T
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	nftlist := make([]map[string]interface{}, 0)

	for _, item := range args.NFT {
		if item.Asset.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		it := make(map[string]interface{})
		it["asset"] = item.Asset.Val()
		it["tokenid"] = item.TokenId.Val()

		nftlist = append(nftlist, it)

	}
	list := utils.GroupByAsset(nftlist)

	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	or := []interface{}{}
	for k, v := range list {
		f := bson.M{"asset": k, "tokenid": bson.M{"$in": v}}
		or = append(or, f)
	}
	filter := bson.M{"market": args.MarketHash, "eventname": "Offer", "$or": or}

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
			Index:      "someindex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": filter},

				bson.M{"$group": bson.M{"_id": bson.M{"asset": "$asset", "tokenid": "$tokenid"}, "info": bson.M{"$push": "$$ROOT"}}},
				//bson.M{"$count": "count"},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	hightestOfferList := make(map[string]interface{})

	for _, tokenidOffer := range r1 {
		id := tokenidOffer["_id"].(map[string]interface{})
		asset := id["asset"].(string)
		tokenid := id["tokenid"].(string)
		key := asset + tokenid

		result := make([]map[string]interface{}, 0)

		items := tokenidOffer["info"].(primitive.A)

		for _, itemOffer := range items {
			//获取有效期内的offer
			item := itemOffer.(map[string]interface{})
			offer := make(map[string]interface{})
			extendData := item["extendData"].(string)
			var dat map[string]interface{}
			if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {
				dl := dat["deadline"].(string)
				deadline, err := strconv.ParseInt(dl, 10, 64)
				if err != nil {
					return err
				}

				if currentTime < deadline {
					//查看offer 当前状态
					offer_nonce := item["nonce"]
					offerstate, _ := me.Client.QueryOne(struct {
						Collection string
						Index      string
						Sort       bson.M
						Filter     bson.M
						Query      []string
					}{
						Collection: "MarketNotification",
						Index:      "getOfferSate",
						Sort:       bson.M{},
						Filter: bson.M{
							"nonce":   offer_nonce,
							"asset":   item["asset"],
							"tokenid": item["tokenid"],
							"$or": []interface{}{
								bson.M{"eventname": "CompleteOffer"},
								bson.M{"eventname": "CancelOffer"},
							},
						},
						Query: []string{},
					}, ret)

					if len(offerstate) > 0 {
						continue
					} else {
						offer["user"] = item["user"]
						offer["asset"] = item["asset"]
						offer["tokenid"] = item["tokenid"]
						offer["originOwner"] = dat["originOwner"]
						offer["offerAsset"] = dat["offerAsset"]
						offerAmount, err := strconv.ParseInt(dat["offerAmount"].(string), 10, 64)
						if err != nil {
							return err
						}
						offer["offerAmount"] = offerAmount
						offer["deadline"] = deadline

						// 获取对应usd的价格
						dd, _ := OpenAssetHashFile()
						decimal := dd[offer["offerAsset"].(string)]
						price, err := GetPrice(offer["offerAsset"].(string))
						//price,err :=GetPrice("0xd2a4cff31913016155e38e474a2c06d08be276cf")
						if err != nil {
							return err
						}

						if price == 0 {
							price = 1
						}

						bfofferAmount, _ := new(big.Float).SetString(dat["offerAmount"].(string))

						flag := bfofferAmount.Cmp(big.NewFloat(0))

						if flag == 1 {
							bfprice := big.NewFloat(price)
							ffprice := big.NewFloat(1).Mul(bfprice, bfofferAmount)
							de := math.Pow(10, float64(decimal))
							usdAuctionAmount := new(big.Float).Quo(ffprice, big.NewFloat(float64(de)))
							offer["usdAmount"] = usdAuctionAmount
						} else {
							offer["usdAmount"] = 0
						}
						//offer["usdAmount"] = price

						result = append(result, offer)
					}
				}
			} else {
				return err1
			}

		}

		//	排序
		result = mapsort.MapSort(result, "offerAmount")

		offerCount := len(result)

		skip := 5

		page := offerCount/skip + 1
		if offerCount%skip == 0 {
			page = offerCount / skip
		}
		hightestOffer := map[string]interface{}{}
		flag := true
		for i := 0; i < page; i++ { //skip
			addressArr := make([]map[string]interface{}, 0)
			var addressList []string
			if i < page-1 {
				for j := i * skip; j < (i+1)*skip; j++ {
					addressArr = append(addressArr, result[j])
					addressList = append(addressList, result[j]["user"].(string))
				}
			} else {
				for j := i * skip; j < offerCount; j++ {
					addressArr = append(addressArr, result[j])
					addressList = append(addressList, result[j]["user"].(string))
				}
			}

			re, err := GetSavings(args.MarketHash, "getUserSavingsAmount", addressList, result[0]["offerAsset"].(string))
			if err != nil {
				return err
			}

			for m := 0; m < len(re); m++ {

				offerAmount := result[i*skip+m]["offerAmount"].(int64)
				oa := new(big.Int).SetInt64(offerAmount)

				if !(re[m].Cmp(oa) == -1) {
					hightestOffer = result[i*skip+m]
					hightestOffer["guarantee"] = re[m]
					flag = false
					break
				}
			}

			if !flag {
				break
			}
		}

		hightestOfferList[key] = hightestOffer
	}

	r2, err := me.Filter(hightestOfferList, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil

}

func GetSavings(scriptHash h160.T, operation string, address []string, assetStr string) ([]*big.Int, error) {

	rt := os.ExpandEnv("${RUNTIME}")
	testNetEndPoint := "http://seed2.neo.org:10332"
	switch rt {

	case "test":
		testNetEndPoint = "http://seed2t4.neo.org:20332"
	case "test2":
		testNetEndPoint = "http://seed2t5.neo.org:20332"
	case "staging":
		testNetEndPoint = "http://seed2.neo.org:10332"
	default:
		log2.Fatalf("runtime environment mismatch")
	}

	client := rpc.NewClient(testNetEndPoint)

	sb := sc.NewScriptBuilder()
	sh, err := helper.UInt160FromString(scriptHash.Val())
	if err != nil {
		return nil, err
	}

	for _, item := range address {
		print(item)
		user, err := helper.UInt160FromString(item)
		if err != nil {
			return nil, err
		}
		asset, err := helper.UInt160FromString(assetStr)
		if err != nil {
			return nil, err
		}
		var arg = []interface{}{user, asset}
		sb.EmitDynamicCall(sh, operation, arg)
	}
	script, err := sb.ToArray()
	if err != nil {
		return nil, err
	}

	response := client.InvokeScript(crypto.Base64Encode(script), nil)
	stack_len := len(response.Result.Stack)

	var result []*big.Int
	for i := 0; i < stack_len; i++ {
		stack := response.Result.Stack[i]
		p, err := stack.ToParameter()
		if err != nil {
			return nil, err
		}
		result = append(result, p.Value.(*big.Int))
	}

	return result, nil
}
