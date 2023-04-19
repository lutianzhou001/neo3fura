package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetInfoByNFTPrimaryMarket(args struct {
	Asset         h160.T
	PrimaryMarket h160.T
	Tokenid       []string
	Filter        map[string]interface{}
	Raw           *map[string]interface{}
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
						bson.M{"$match": bson.M{"market": args.PrimaryMarket, "eventname": bson.M{"$in": []interface{}{"Claim", "Auction"}}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						bson.M{"$sort": bson.M{"timestamp": 1}},
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
		//asset := item["asset"].(string)
		//tokenid := item["tokenid"].(string)
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

		auctionType := item["auctionType"].(int32)
		if item["market"] == args.PrimaryMarket && ddl > currentTime && item["market"] == item["owner"] { //一级市场上架状态
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
		} else if item["market"] == args.PrimaryMarket && ddl <= currentTime && item["market"] == item["owner"] { //一级市场未领取
			if auctionType == 2 && bidAmount != "0" {
				item["lastSoldAsset"] = item["auctionAsset"]
				item["lastSoldAmount"] = item["bidAmount"]
			} else {
				item["lastSoldAsset"] = item["auctionAsset"]
				item["lastSoldAmount"] = item["auctionAmount"]
			}
		} else { // 一级市场 售卖成功状态
			if item["eventlist"] != nil && len(item["eventlist"].(primitive.A)) > 0 {
				eventlist := item["eventlist"].(primitive.A)[0]
				eventItem := eventlist.(map[string]interface{})
				extendData := eventItem["extendData"]
				eventname := eventItem["eventname"]
				data := make(map[string]interface{})
				if err := json.Unmarshal([]byte(extendData.(string)), &data); err == nil {
					if eventname == "Clain" {
						item["lastSoldAsset"] = data["auctionAsset"]
						item["lastSoldAmount"] = data["bidAmount"]
					} else if eventname == "Auction" {
						item["lastSoldAsset"] = data["auctionAsset"]
						item["lastSoldAmount"] = data["auctionAmount"]
					}

				}
			}

		}
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
