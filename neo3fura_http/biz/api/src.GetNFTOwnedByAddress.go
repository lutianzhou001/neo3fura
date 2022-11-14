package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTOwnedByAddress(args struct {
	Address      h160.T
	MarketHash   h160.T
	ContractHash h160.T   //  asset
	AssetHash    h160.T   // auctionType
	NFTState     strval.T //state:aution  sale  notlisted  unclaimed
	Sort         strval.T //listedTime  price
	Order        int64    //-1:降序  +1：升序
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6

	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{}

	if len(args.ContractHash) > 0 {
		if args.ContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"asset": args.ContractHash}}
			pipeline = append(pipeline, a)
		}
	}

	if len(args.AssetHash) > 0 {
		if args.AssetHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"auctionAsset": args.AssetHash}}
			pipeline = append(pipeline, a)
		}
	}

	//按截止时间排序
	var deadlineCond bson.M
	if args.Sort == "deadline" { //按截止时间排序
		//deadlineCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": bson.M{"$subtract": []interface{}{"$deadline", currentTime}}, "else": bson.M{"$subtract": []interface{}{currentTime, "$deadline"}}}}
		deadlineCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": bson.M{"$subtract": []interface{}{"$deadline", currentTime}}, "else": currentTime}}
	}
	//按照时间价格排序
	var auctionAmountCond bson.M
	if args.Sort == "price" { // 将过期和未领取的放在后面

		if args.AssetHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			if args.Order == -1 { //降序
				auctionAmountCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": "$auctionAmount", "else": 0}}
			} else { //升序（默认）
				auctionAmountCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": "$auctionAmount", "else": 1e16}}
			}
		}

	}
	////按上架时间排序
	//var listedtimeCond bson.M
	//if args.Sort == "timestamp" { //将未上架和未领取、过期放在后面
	//	listedtimeCond = bson.M{"$cond": bson.M{"if": bson.M{"$and": []interface{}{bson.M{"$eq": []interface{}{"$owner", "$market"}}, bson.M{"$gt": []interface{}{"$deadline", currentTime}}}},
	//		"then": "$timestamp",
	//		"else": 0}}
	//
	//}
	if len(args.MarketHash) > 0 && args.MarketHash != "" {
		if args.MarketHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			//白名单
			raw1 := make(map[string]interface{})
			err1 := me.GetMarketWhiteList(struct {
				MarketHash h160.T
				Filter     map[string]interface{}
				Raw        *map[string]interface{}
			}{MarketHash: args.MarketHash, Raw: &raw1}, ret) //nonce 分组，并按时间排序
			if err1 != nil {
				return err1
			}

			whiteList := raw1["whiteList"]
			if whiteList == nil || whiteList == "" {
				return stderr.ErrWhiteList
			}
			s := whiteList.([]string)
			var wl []interface{}
			for _, w := range s {
				wl = append(wl, w)
			}
			if len(wl) > 0 {
				white := bson.M{"$match": bson.M{"asset": bson.M{"$in": wl}}}
				pipeline = append(pipeline, white)
			} else {
				return stderr.ErrWhiteList
			}

		}

	}

	if args.NFTState.Val() == NFTstate.Auction.Val() || args.NFTState.Val() == NFTstate.Sale.Val() || args.NFTState.Val() == NFTstate.Unclaimed.Val() {
		if args.MarketHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"market": args.MarketHash}}
			pipeline = append(pipeline, a)
		}
	}

	if args.NFTState.Val() == NFTstate.Auction.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"auctor": args.Address.Val()}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 2}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "properties": 1, "_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "auction"}},
			bson.M{"$match": bson.M{"difference": true}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"auctor": args.Address.Val()}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "properties": 1, "_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "sale"}},
			bson.M{"$match": bson.M{"difference": true}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.NotListed.Val() { //未上架  accont >0 && owner != market
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": args.Address.Val()}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "properties": 1, "_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "notlisted"}},
			bson.M{"$match": bson.M{"difference": false}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Unclaimed.Val() { //未领取 (买家和卖家) accont >0 &&  runtime > deadline && owner== market && bidAccount >0
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"$or": []interface{}{
				bson.M{"$and": []interface{}{ //卖家未领取(过期)
					bson.M{"auctor": args.Address.Val()},
					bson.M{"bidAmount": 0},
				}},
				bson.M{"$and": []interface{}{ //买家未领取
					bson.M{"bidder": args.Address.Val()},
					bson.M{"bidAmount": bson.M{"$gt": 0}},
				}},
			}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$lt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"_id": 1, "deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "unclaimed"}},
			bson.M{"$match": bson.M{"difference": true}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else { //默认  account > 0
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"$or": []interface{}{
				bson.M{"owner": args.Address.Val()}, //	未上架  owner
				bson.M{"$and": []interface{}{ //上架 售卖中、过期（无人出价）auctor
					bson.M{"auctor": args.Address.Val()},
					bson.M{"deadline": bson.M{"$gte": currentTime}},
				}},
				bson.M{"$and": []interface{}{ //未领取   竞价成功bidder
					bson.M{"bidder": args.Address.Val()},
					bson.M{"bidAmount": bson.M{"$gt": 0}},
					bson.M{"deadline": bson.M{"$lte": currentTime}},
				}},
			}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"_id": 1, "deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": ""}},
		}
		pipeline = append(pipeline, pipeline1...)
	}

	//按上架时间排序
	if args.Sort == "timestamp" {
		lookup := bson.M{"$lookup": bson.M{
			"from": "MarketNotification",
			"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"eventname": "Auction"}},
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
					bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					bson.M{"$eq": []interface{}{"$market", "$$market"}},
				}}}},

				bson.M{"$project": bson.M{"asset": 1, "nonce": 1, "tokenid": 1, "timestamp": 1}},
				bson.M{"$sort": bson.M{"nonce": -1}},
				bson.M{"$limit": 1},
			},
			"as": "marketnotification"},
		}

		pipeline2 := []bson.M{}
		if args.NFTState.Val() == NFTstate.Auction.Val() {
			pipeline2 = []bson.M{
				bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "marketnotification": 1, "_id": 1, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "auction"}},
				bson.M{"$match": bson.M{"difference": true}},
				bson.M{"$sort": bson.M{"deadline": -1}},
			}

		} else if args.NFTState.Val() == NFTstate.Sale.Val() {
			pipeline2 = []bson.M{
				bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "marketnotification": 1, "_id": 1, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "sale"}},
				bson.M{"$match": bson.M{"difference": true}},
				bson.M{"$sort": bson.M{"deadline": -1}},
			}
		} else if args.NFTState.Val() == NFTstate.Unclaimed.Val() {
			pipeline2 = []bson.M{
				bson.M{"$project": bson.M{"_id": 1, "deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "marketnotification": 1, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "unclaimed"}},
				bson.M{"$match": bson.M{"difference": true}},
				bson.M{"$sort": bson.M{"deadline": -1}},
			}
		} else {
			pipeline2 = []bson.M{
				bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "marketnotification": 1, "_id": 1, "properties": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": ""}},
				bson.M{"$sort": bson.M{"deadline": -1}},
			}
		}
		pipeline = append(pipeline, lookup)
		pipeline = append(pipeline, pipeline2...)

	}

	//按截止时间排序
	if args.Sort == "deadline" {
		sort := bson.M{"$sort": bson.M{"deadlineCond": 1}}
		pipeline = append(pipeline, sort)
	}
	//按价格排序
	if args.Sort == "price" { //按币种价格排序
		if args.AssetHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			s := bson.M{"$sort": bson.M{"auctionAmountCond": args.Order}}
			pipeline = append(pipeline, s)
		}
	}

	skip := bson.M{"$skip": args.Skip}
	limit := bson.M{"$limit": args.Limit}
	pipeline = append(pipeline, skip)
	pipeline = append(pipeline, limit)

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
			Index:      "GetNFTMarket",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{"_id", "date", "marketnotification", "deadlineCond", "auctionAmountCond", "properties", "asset", "tokenid", "amount", "owner", "market", "auctionType", "auctor", "auctionAsset", "auctionAmount", "deadline", "bidder", "bidAmount", "timestamp", "state"},
		}, ret)

	if err != nil {
		return err
	}

	//获取nft的属性
	for _, item := range r1 {
		if args.NFTState.Val() != NFTstate.Auction.Val() && args.NFTState.Val() != NFTstate.Sale.Val() && args.NFTState.Val() != NFTstate.NotListed.Val() && args.NFTState.Val() != NFTstate.Unclaimed.Val() {
			a := item["amount"].(primitive.Decimal128).String()
			amount, err1 := strconv.Atoi(a)
			if err1 != nil {
				return err1
			}

			bidAmount := item["bidAmount"].(primitive.Decimal128).String()

			deadline, _ := item["deadline"].(int64)
			auctionType, _ := item["auctionType"].(int32)

			if amount > 0 && auctionType == 2 && item["owner"] == item["market"] && deadline > currentTime {
				item["state"] = NFTstate.Auction.Val()
			} else if amount > 0 && auctionType == 1 && item["owner"] == item["market"] && deadline > currentTime {
				item["state"] = NFTstate.Sale.Val()
			} else if amount > 0 && item["owner"] != item["market"] {
				item["state"] = NFTstate.NotListed.Val()
			} else if amount > 0 && bidAmount == "0" && deadline < currentTime && item["owner"] == item["market"] {
				item["state"] = NFTstate.Unclaimed.Val()
			} else {
				item["state"] = ""
			}
		}
		//获得上架时间

		if item["marketnotification"] != nil && item["marketnotification"] != "" {
			switch item["marketnotification"].(type) {
			case string:
				item["listedTimestamp"] = int64(0)
			case primitive.A:
				marketnotification := item["marketnotification"].(primitive.A)
				if len(marketnotification) > 0 {
					mn := []interface{}(marketnotification)[0].(map[string]interface{})
					if item["deadline"].(int64) > currentTime {
						item["listedTimestamp"] = mn["timestamp"]
					} else {
						item["listedTimestamp"] = mn["timestamp"].(int64) - 1640966400000
					}
				} else {
					item["listedTimestamp"] = int64(0)
				}
			}
		} else {
			item["listedTimestamp"] = int64(0)
		}
		delete(item, "marketnotification")

		//获取nft 属性
		nftproperties := item["properties"]
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		if nftproperties != nil && nftproperties != "" {
			pp := nftproperties.(primitive.A)
			if len(pp) > 0 {
				it := pp[0].(map[string]interface{})
				extendData := it["properties"].(string)
				if extendData != "" {
					properties := make(map[string]interface{})
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
						image, ok := data["image"]
						if ok {
							properties["image"] = image
							//item["image"] = image
							item["image"] = ImagUrl(asset, image.(string), "images")
						} else {
							item["image"] = ""
						}
						tokenuri, ok := data["tokenURI"]
						if ok {
							ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)
							if err != nil {
								return err
							}
							for key, value := range ppjson {
								item[key] = value
								properties[key] = value
								if key == "image" {
									img := value.(string)
									item["thumbnail"] = ImagUrl(asset, img, "thumbnail")
									item["image"] = ImagUrl(asset, img, "images")
								}
							}
						}
						if item["name"] == "" || item["name"] == nil {
							name, ok1 := data["name"]
							if ok1 {
								item["name"] = name
							}
						}

						strArray := strings.Split(item["name"].(string), "#")
						if len(strArray) >= 2 {
							number := strArray[1]
							n, err22 := strconv.ParseInt(number, 10, 64)
							if err22 != nil {
								item["number"] = int64(-1)
							}
							item["number"] = n
							properties["number"] = n
						} else {
							item["number"] = int64(-1)
						}

						series, ok2 := data["series"]
						if ok2 {
							properties["series"] = series
						}
						supply, ok3 := data["supply"]
						if ok3 {
							properties["supply"] = supply
						}
						number, ok4 := data["number"]
						if ok4 {
							n, err22 := strconv.ParseInt(number.(string), 10, 64)
							if err22 != nil {
								item["number"] = int64(-1)
							}
							properties["number"] = n
							item["number"] = n
						}

						thumbnail, ok6 := data["thumbnail"]
						if ok6 {
							tb, err22 := base64.URLEncoding.DecodeString(thumbnail.(string))
							if err22 != nil {
								return err22
							}
							//item["image"] = string(tb[:])
							item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
						} else {
							if item["thumbnail"] == nil {
								if image != nil && image != "" {
									if image == nil {
										item["thumbnail"] = item["image"]
									} else {
										item["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
									}
								}
							}
						}

					} else {
						return err
					}
					if item["name"].(string) == "Video" {
						item["video"] = item["image"]
						delete(item, "image")
						properties["video"] = properties["image"]
						delete(properties, "image")
					}

					item["properties"] = properties
				} else {
					item["image"] = ""
					item["name"] = ""
					item["number"] = int64(-1)
					item["properties"] = ""
				}
			}
		}
	}

	// 按上架时间排序
	if args.Sort == "timestamp" {
		if args.Order == 1 {
			mapsort.MapSort2(r1, "listedTimestamp")
		} else {
			mapsort.MapSort(r1, "listedTimestamp")
		}

	}
	//获取查询总量
	pipeline = append(pipeline[:len(pipeline)-2], pipeline[len(pipeline):]...)
	var group = bson.M{"$group": bson.M{"_id": "$_id"}}

	pipeline = append(pipeline, group)
	var countKey = bson.M{"$count": "total counts"}

	pipeline = append(pipeline, countKey)

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
			Index:      "GetNFTMarket",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}

	var count interface{}
	if len(r2) != 0 {
		count = r2[0]["total counts"]
	} else {
		count = 0
	}
	if err != nil {
		return err
	}

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
