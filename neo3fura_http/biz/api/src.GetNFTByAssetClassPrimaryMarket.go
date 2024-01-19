package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/Contract"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"os"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTByAssetClassPrimaryMarket(args struct {
	Asset         h160.T
	PrimaryMarket h160.T
	Class         string
	ClassName     string
	Limit         int64
	Skip          int64
	Filter        map[string]interface{}
	Raw           *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 50

	}
	rt := os.ExpandEnv("${RUNTIME}")
	var nns, polemen, genesis string
	if rt == "staging" {
		nns = Contract.Main_NNS.Val()
		//metapanacea = Contract.Main_MetaPanacea.Val()
		genesis = Contract.Main_ILEXGENESIS.Val()
		polemen = Contract.Main_ILEXPOLEMEN.Val()

	} else if rt == "test2" {
		nns = Contract.Test_NNS.Val()
		//metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Test_ILEXGENESIS.Val()
		polemen = Contract.Test_ILEXPOLEMEN.Val()
	} else {
		nns = Contract.Main_NNS.Val()
		//metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Main_ILEXGENESIS.Val()
		polemen = Contract.Main_ILEXPOLEMEN.Val()
	}

	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "SelfControlNep11Properties",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": args.Asset}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", nns}}, "then": "$asset",
					"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", genesis}}, "then": "$image",
						"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", polemen}}, "then": "$tokenid",
							"else": "$name"}}}}}}}},
				bson.M{"$match": bson.M{"class": args.ClassName}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		tokenidArr := []string{tokenid}

		//if item["image"] == nil {
		if item["properties"] != nil { //
			jsonData := make(map[string]interface{})
			properties := item["properties"].(string)
			if properties != "" {
				err := json.Unmarshal([]byte(properties), &jsonData)
				if err != nil {
					return err
				}

				tokenURI, ok := jsonData["tokenURI"]
				if ok {
					ppjson, err := GetImgFromTokenURL(tokenurl(tokenURI.(string)), asset, tokenid)
					if err != nil {
						return err
					}
					for key, value := range ppjson {
						item[key] = value
						if key == "image" {
							img := value.(string)
							thumbnail := ImagUrl(asset, img, "thumbnail")
							flag := strings.HasSuffix(thumbnail, ".mp4")
							if flag {
								thumbnail = strings.Replace(thumbnail, ".mp4", "mp4", -1)
							}
							item["thumbnail"] = thumbnail
							item["image"] = ImagUrl(asset, img, "images")
						}
						if key == "name" {
							item["name"] = value
						}

					}

				}

				image, ok := jsonData["image"]
				if ok {
					item["image"] = ImagUrl(item["asset"].(string), image.(string), "images")
				} else {
					item["image"] = ""
				}

				thumbnail, ok1 := jsonData["thumbnail"]
				if ok1 {
					tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
					if err2 != nil {
						return err2
					}
					ss := string(tb[:])
					if ss == "" {
						item["thumbnail"] = ImagUrl(item["asset"].(string), item["image"].(string), "thumbnail")
					} else {
						item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
					}

				} else {
					if item["thumbnail"] == nil {
						if item["image"] != nil && item["image"] != "" {
							if image == nil {
								item["thumbnail"] = item["image"]
							} else {
								item["thumbnail"] = ImagUrl(item["asset"].(string), image.(string), "thumbnail")
							}
						}
					}
				}
			}

		}
		if item["tokenURI"] != nil {
			tokenUrl := item["tokenURI"].(string)
			ppjson, err := GetImgFromTokenURL(tokenurl(tokenUrl), asset, tokenid)
			if err != nil {
				return err
			}
			for key, value := range ppjson {
				//item[key] = value
				if key == "image" {
					img := value.(string)
					thumbnail := ImagUrl(asset, img, "thumbnail")
					flag := strings.HasSuffix(thumbnail, ".mp4")
					if flag {
						thumbnail = strings.Replace(thumbnail, ".mp4", "mp4", -1)
					}
					item["thumbnail"] = thumbnail
					item["image"] = ImagUrl(asset, img, "images")
				}
				if key == "name" {
					item["name"] = value
				}

			}
		}

		if item["name"] != nil && item["name"].(string) == "Nuanced Floral Symphony" {
			item["video"] = item["image"]
			delete(item, "image")
		}
		//}
		if item["name"] != nil && item["name"].(string) == "Virtual Visions #1" {
			item["video"] = item["image"]
			delete(item, "image")
		}

		re := map[string]interface{}{}
		err := me.GetInfoByNFTPrimaryMarket(struct {
			Asset         h160.T
			PrimaryMarket h160.T
			Tokenid       []string
			Filter        map[string]interface{}
			Raw           *map[string]interface{}
		}{Asset: h160.T(asset), PrimaryMarket: args.PrimaryMarket, Tokenid: tokenidArr, Filter: args.Filter, Raw: &re}, ret)

		if err != nil {
			return stderr.ErrGetNFTInfo
		}

		marketInfo := re["result"]
		if marketInfo != nil {
			marketItem := marketInfo.([]map[string]interface{})
			info := marketItem[0]
			for key, value := range info {
				item[key] = value
			}
			delete(item, "properties")
		}

	}
	// totalcount
	r2, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "SelfControlNep11Properties",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", nns}}, "then": "$asset",
					"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", genesis}}, "then": "$image",
						"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", polemen}}, "then": "$tokenid",
							"else": "$name"}}}}}}}},
				bson.M{"$match": bson.M{"class": args.ClassName}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	count := len(r2)
	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)

	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r3
	}
	*ret = json.RawMessage(r)
	return nil
}

func (me *T) GetNFTInfoPrimaryMarket(Market string, Asset string, Tokenid string) ([]map[string]interface{}, error) {

	var ret *json.RawMessage
	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Market",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": Asset}},
				bson.M{"$lookup": bson.M{
					"from": "MarketNotification",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{ //
						bson.M{"$match": bson.M{"$expr": bson.M{"$or": []interface{}{
							bson.M{"$and": []interface{}{
								bson.M{"$in": []interface{}{"$eventname", []interface{}{"CompleteOfferCollection", "Offer", "CompleteOffer", "Claim"}}},
								bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
								bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							}},
							bson.M{"$and": []interface{}{
								bson.M{"$eq": []interface{}{"$eventname", "OfferCollection"}},
								bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
							}},
						}}}},
						bson.M{"$sort": bson.M{"timestamp": 1}},
						bson.M{"$group": bson.M{"_id": "$eventname", "eventArr": bson.M{"$push": "$$ROOT"}, "eventname": bson.M{"$last": "$eventname"}, "market": bson.M{"$last": "$market"}, "timestamp": bson.M{"$last": "$timestamp"}, "extendData": bson.M{"$last": "$extendData"}}},
						//bson.M{"$project": bson.M{"eventname": 1,"eventArr" :1"market": 1, "extendData": 1, "timestamp": 1}},
					},
					"as": "eventlist"}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UnixNano() / 1e6

	for _, item := range r1 {
		//NFT状态   上架 （售卖中  成交未领取）  未上架
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		ddl := item["deadline"].(int64)

		bidAmount := item["bidAmount"].(primitive.Decimal128).String()
		if item["market"] == Market {
			if item["market"] != item["owner"] || ddl < currentTime {
				item["state"] = "notlist"
			} else {
				item["state"] = "list"
			}
		} else {
			item["state"] = "notlist"
		}
		if (item["market"] == item["owner"] && ddl > currentTime) || (item["market"] == item["owner"] && ddl < currentTime && bidAmount == "0") { //上架
			item["owner"] = item["auctor"]
		}
		if item["market"] == item["owner"] && ddl < currentTime && bidAmount != "0" { // 未领取
			item["owner"] = item["bidder"]
		}

		item["buyNowAsset"] = ""
		item["buyNowAmount"] = "0"
		item["lastSoldAsset"] = ""
		item["lastSoldAmount"] = "0"
		item["currentBidAsset"] = ""
		item["currentBidAmount"] = "0"
		item["offerAsset"] = ""
		item["offerAmount"] = "0"
		item["nonce"] = 0
		item["eventname"] = ""

		auctionType := item["auctionType"].(int32)
		if item["market"] == Market {
			if ddl > currentTime {
				if auctionType == 1 {
					item["buyNowAsset"] = item["auctionAsset"]
					item["buyNowAmount"] = item["auctionAmount"]
				} else if auctionType == 2 {
					if bidAmount != "0" {
						item["currentBidAsset"] = item["auctionAsset"]
						item["currentBidAmount"] = item["bidAmount"]
					} else {
						item["currentBidAsset"] = item["auctionAsset"]
						item["currentBidAmount"] = item["auctionAmount"]
					}
				}
			} else {
				if auctionType == 2 && bidAmount != "0" {
					item["lastSoldAsset"] = item["auctionAsset"]
					item["lastSoldAmount"] = item["bidAmount"]
				}
			}
		}

		var finishTime int64
		if item["eventlist"] != nil && len(item["eventlist"].(primitive.A)) > 0 {
			eventlist := item["eventlist"].(primitive.A)
			for _, it := range eventlist {
				eventItem := it.(map[string]interface{})
				eventname := eventItem["eventname"]
				extendData := eventItem["extendData"]
				market := eventItem["market"].(string)

				data := make(map[string]interface{})
				if err := json.Unmarshal([]byte(extendData.(string)), &data); err == nil {
					if eventname == "Claim" {
						time := eventItem["timestamp"].(int64)
						if time > finishTime {
							finishTime = time
							item["lastSoldAsset"] = data["auctionAsset"]
							item["lastSoldAmount"] = data["bidAmount"]
						}

					} else if eventname == "Offer" || eventname == "OfferCollection" {
						//判断offer 有效期以及是否有足够的保证金
						deadline := data["deadline"].(string)
						offerddl, _ := strconv.ParseInt(deadline, 10, 64)

						highestOffer := make(map[string]interface{})
						if offerddl > currentTime {
							err := me.GetHighestOfferByNFT(struct {
								Asset      h160.T
								TokenId    strval.T
								MarketHash h160.T
								Limit      int64
								Skip       int64
								Filter     map[string]interface{}
								Raw        *map[string]interface{}
							}{Asset: h160.T(asset), TokenId: strval.T(tokenid), MarketHash: h160.T(market), Raw: &highestOffer}, ret)
							if err != nil {
								return nil, stderr.ErrGetHighestOffer
							}
							if len(highestOffer) > 0 {
								offerAmount := highestOffer["offerAmount"].(int64)
								guarantee := highestOffer["guarantee"].(*big.Int)
								amount := big.NewInt(offerAmount)
								if guarantee.Cmp(amount) == 1 {
									item["offerAsset"] = highestOffer["offerAsset"]
									item["offerAmount"] = amount.String()
									item["nonce"] = highestOffer["nonce"]
									item["eventname"] = highestOffer["eventname"]
								}
							}
						}
					} else if eventname == "CompleteOffer" || eventname == "CompleteOfferCollection" {
						time := eventItem["timestamp"].(int64)
						if time > finishTime {
							finishTime = time
							item["lastSoldAsset"] = data["offerAsset"]
							if err != nil {
								return nil, err
							}
							item["lastSoldAmount"] = data["offerAmount"]

						}
					}
				}

			}
		}

		//获取Owner 地址的nns信息
		owner := item["owner"].(string)
		var nns, userName string
		if owner != "" {
			nns, userName, err = GetNNSByAddress(owner)
			if err != nil {
				return nil, err
			}
		}

		item["nns"] = nns
		item["userName"] = userName
		delete(item, "eventlist")
	}

	return r1, nil
}
