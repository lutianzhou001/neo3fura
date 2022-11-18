package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"
)

func (me *T) GetInfoByNFT(args struct {
	Asset   h160.T
	Tokenid []string
	Filter  map[string]interface{}
	Raw     *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	//获取上架以及Owner信息
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
				bson.M{"$match": bson.M{"amount": 1, "asset": args.Asset.Val(), "tokenid": bson.M{"$in": args.Tokenid}}},
				bson.M{"$lookup": bson.M{
					"from": "MarketNotification",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{ //
						bson.M{"$match": bson.M{"$expr": bson.M{"$or": []interface{}{
							bson.M{"$and": []interface{}{
								bson.M{"$eq": []interface{}{"$eventname", []interface{}{"CompleteOfferCollection", "Offer", "CompleteOffer", "Claim"}}},
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
		return err
	}

	currentTime := time.Now().UnixNano() / 1e6

	for _, item := range r1 {
		//NFT状态   上架 （售卖中  成交未领取）  未上架
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		ddl := item["deadline"].(int64)

		bidAmount := item["bidAmount"].(primitive.Decimal128).String()
		if item["market"] != item["owner"] || ddl < currentTime {
			item["state"] = "notlist"
		} else {
			item["state"] = "list"
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
		var finishTime int64
		if item["eventlist"] != nil {
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
								return stderr.ErrGetHighestOffer
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
								return err
							}
							item["lastSoldAmount"] = data["offerAmount"]

						}
					}
				}

			}
		}

		//获取Owner 地址的用户信息
		//TODO
		delete(item, "eventlist")
	}

	count := len(r1)
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
