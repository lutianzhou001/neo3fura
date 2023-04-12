package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strings"
	"time"
)

func (me *T) GetNFTClass(args struct {
	MarketHash h160.T
	AssetHash  h160.T
	NFTState   string
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	var filter bson.M
	if args.NFTState == NFTstate.Auction.Val() {
		filter = bson.M{"market": args.MarketHash, "amount": 1, "auctionType": 2}
	} else if args.NFTState == NFTstate.Sale.Val() {
		filter = bson.M{"market": args.MarketHash, "amount": 1, "auctionType": 1}
	} else {
		filter = bson.M{"amount": 1}
	}

	var r2, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "SelfControlNep11Properties",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": args.AssetHash}},
				bson.M{"$lookup": bson.M{
					"from": "Market",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": filter},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						//	bson.M{"$sort"}
					},
					"as": "marketInfo"},
				},
				bson.M{"$lookup": bson.M{
					"from": "MarketNotification",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"market": args.MarketHash, "eventname": "Claim"}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						//bson.M{"$sort":bson.M{"timestamp":1}},
					},
					"as": "marketNotification"},
				},
				bson.M{"$set": bson.M{"class": "$name"}},
				bson.M{"$group": bson.M{"_id": bson.M{"asset": "$asset", "class": "$class"}, "class": bson.M{"$last": "$class"}, "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"},
					"name": bson.M{"$last": "$name"}, "image": bson.M{"$last": "$image"}, "supply": bson.M{"$last": "$supply"}, "thumbnail": bson.M{"$last": "$thumbnail"}, "marketNotification": bson.M{"$last": "$marketNotification"},
					"properties": bson.M{"$last": "$properties"}, "marketArr": bson.M{"$last": "$marketInfo"}, "itemList": bson.M{"$push": "$$ROOT"}}},
			},
			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}
	currentTime := time.Now().UnixNano() / 1e6
	result := make([]map[string]interface{}, 0)
	for _, item := range r2 {
		Info := item["marketArr"].(primitive.A)
		if len(Info) == 0 {
			continue
		}
		marketInfo := Info[0].(map[string]interface{})
		marketInfo = GetNFTState(marketInfo, args.MarketHash)
		item["auctionType"] = marketInfo["auctionType"]
		item["currentBidAmount"] = marketInfo["bidAmount"]
		item["currentBidAsset"] = marketInfo["auctionAsset"]
		//item["state"] =marketInfo["state"]
		auctionType := marketInfo["auctionType"].(int32)
		if auctionType == 2 {
			item["currentBidAsset"] = marketInfo["auctionAsset"]
		} else {
			item["currentBidAsset"] = ""
		}
		item["auctionAsset"] = marketInfo["auctionAsset"]
		item["auctionAmount"] = marketInfo["auctionAmount"]
		item["deadline"] = marketInfo["deadline"]
		item["lastSoldAsset"] = ""
		item["lastSoldAmount"] = "0"

		notify := item["marketNotification"].(primitive.A)

		if len(notify) > 0 {
			notification := notify[0].(map[string]interface{})
			extendData := notification["extendData"]
			if extendData != nil && extendData != "" {
				data := make(map[string]interface{})
				if err := json.Unmarshal([]byte(extendData.(string)), &data); err == nil {
					item["lastSoldAsset"] = data["auctionAsset"]
					item["lastSoldAmount"] = data["bidAmount"]
				}
			}
		}

		deadline := marketInfo["deadline"].(int64)
		bidAmount := marketInfo["bidAmount"].(primitive.Decimal128).String()
		if deadline < currentTime && bidAmount != "0" {
			item["lastSoldAsset"] = marketInfo["auctionAsset"]
			item["lastSoldAmount"] = marketInfo["auctionAmount"]
		}

		count := 0
		groupinfo := item["itemList"].(primitive.A)
		for _, it := range groupinfo {
			pit := it.(map[string]interface{})
			market := pit["marketInfo"].(primitive.A)[0].(map[string]interface{})
			market = GetNFTState(market, args.MarketHash)
			if market["state"].(string) == "sale" || market["state"].(string) == "auction" {
				count++
			}
		}
		item["claimed"] = len(groupinfo) - count

		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)

		if item["properties"] != nil {
			properties := item["properties"].(string)
			if properties != "" {
				jsonData := make(map[string]interface{})
				err := json.Unmarshal([]byte(properties), &jsonData)
				if err != nil {
					return err
				}
				tokenURI, ok := jsonData["tokenURI"]
				if ok {
					//item["image"] = ImagUrl(item["asset"].(string), image.(string), "images")
					//tokenUrl := item["tokenURI"].(string)
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

				} else {
					item["image"] = ""
				}
			}

		}

		//if image != nil {
		//	item["image"] = ImagUrl(asset, item["image"].(string), "images")
		//} else {
		//	item["image"] = ""
		//}
		name := item["name"].(string)
		nameArr := strings.Split(name, "#")

		if len(nameArr) > 1 {
			item["series"] = strings.Trim(nameArr[0], " ")
			item["number"] = nameArr[1]

		}

		if item["name"] != nil && item["name"].(string) == "Nuanced Floral Symphony" {
			item["video"] = item["image"]
			delete(item, "image")

		}

		// 处理ilex genesis 以image 分类的特殊情况
		class := item["class"]
		item["classname"] = class
		if class != nil {
			if isHttp(class.(string)) {
				item["class"] = item["series"]

			}
		}

		delete(item, "_id")
		delete(item, "itemList")
		delete(item, "marketArr")
		delete(item, "properties")
		//delete(item, "class")
		delete(item, "marketNotification")

		item["count"] = len(groupinfo)
		//处理排序 video优先
		//tokenid := item["tokenid"].(string)
		if tokenid == "AQ==" {
			item["order"] = "00"
		} else {
			item["order"] = "11"
		}
		result = append(result, item)
	}

	mapsort.MapSort8(result, "order")

	r3, err := me.FilterAggragateAndAppendCount(result, len(result), args.Filter)

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

func GetNFTState(info map[string]interface{}, primarymarket interface{}) map[string]interface{} {
	if len(info) > 0 {
		market := info["market"]
		currentTime := time.Now().UnixNano() / 1e6
		if market == nil || market.(string) == primarymarket.(h160.T).Val() {
			deadline := info["deadline"].(int64)
			auctionType := info["auctionType"].(int32)
			bidAmount := info["bidAmount"].(primitive.Decimal128).String()

			info["currentBidAmount"] = info["bidAmount"]
			info["currentBidAmount"] = info["auctionAsset"]

			if deadline > currentTime && market != nil && market.(string) == primarymarket.(h160.T).Val() {
				if auctionType == 1 {
					info["state"] = "sale" //
				} else if auctionType == 2 {
					info["state"] = "auction"
				}
			} else if deadline <= currentTime && market != nil && market.(string) == primarymarket.(h160.T).Val() {
				if auctionType == 2 && bidAmount != "0" {
					info["state"] = "soldout" //竞拍有人出价
				} else {
					info["state"] = "expired"
				}
			} else {
				info["state"] = "soldout"
			}

		} else {
			info["state"] = ""
		}
	} else {
		info["state"] = "no"
	}

	//delete(info, "bidAmount")
	delete(info, "bidder")
	delete(info, "auctor")
	delete(info, "timestamp")
	return info
}

func GetNFTStateByNotication(info map[string]interface{}, primarymarket interface{}) map[string]interface{} {
	if len(info) > 0 {
		market := info["market"]
		if market == nil || market == primarymarket {
			deadline := info["deadline"].(int64)
			auctionType := info["auctionType"].(int32)
			bidAmount := info["bidAmount"].(primitive.Decimal128).String()

			info["currentBidAmount"] = info["bidAmount"]
			info["currentBidAmount"] = info["auctionAsset"]
			currentTime := time.Now().UnixNano() / 1e6
			if deadline > currentTime && market == primarymarket {
				if auctionType == 1 {
					info["state"] = "sale" //
				} else if auctionType == 2 {
					info["state"] = "auction"
				}
			} else if deadline <= currentTime && market == primarymarket {
				if auctionType == 2 && bidAmount != "0" {
					info["state"] = "soldout" //竞拍有人出价
				} else {
					info["state"] = "expired"
				}
			} else {
				info["state"] = "soldout"
			}

		} else {
			info["state"] = ""
		}
	} else {
		info["state"] = "no"
	}

	delete(info, "bidAmount")
	delete(info, "bidder")
	delete(info, "auctor")
	delete(info, "timestamp")
	return info
}
