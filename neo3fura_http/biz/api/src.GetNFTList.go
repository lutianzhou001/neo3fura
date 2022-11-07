package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	address "neo3fura_http/var/const"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTList(args struct {
	SecondaryMarket h160.T //
	PrimaryMarket   h160.T
	ContractHash    h160.T
	NFTState        strval.T //state:aution  sale  notlisted  unclaimed
	Sort            strval.T //listedTime  price  deadline
	Order           int64    //-1:降序  +1：升序
	Limit           int64
	Skip            int64
	Filter          map[string]interface{}
	Raw             *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	pipeline := []bson.M{}

	if len(args.PrimaryMarket) > 0 && args.PrimaryMarket != "" {
		if args.PrimaryMarket.Valid() == false {
			return stderr.ErrInvalidArgs
		}
	}
	if len(args.ContractHash) > 0 && args.ContractHash != "" {
		if args.ContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		a := bson.M{"$match": bson.M{"asset": args.ContractHash.Val()}}
		pipeline = append(pipeline, a)
	}

	if len(args.SecondaryMarket) > 0 && args.SecondaryMarket != "" {
		if args.NFTState.Val() == NFTstate.Auction.Val() || args.NFTState.Val() == NFTstate.Sale.Val() {
			if args.SecondaryMarket.Valid() == false {
				return stderr.ErrInvalidArgs
			} else {
				a := bson.M{"$match": bson.M{"market": args.SecondaryMarket}}
				pipeline = append(pipeline, a)
			}
		}
		//白名单
		raw1 := make(map[string]interface{})
		err1 := me.GetMarketWhiteList(struct {
			MarketHash h160.T
			Filter     map[string]interface{}
			Raw        *map[string]interface{}
		}{MarketHash: args.SecondaryMarket, Raw: &raw1}, ret) //nonce 分组，并按时间排序
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
			//white := bson.M{"$match": bson.M{"asset": bson.M{"$in": []interface{}{"0x6c91e9997b8e74dcfa5ebb56fe5672dedd724b8f","0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f"}}}}
			pipeline = append(pipeline, white)
		} else {
			return stderr.ErrWhiteList
		}
	}

	if args.NFTState.Val() == NFTstate.Auction.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": bson.M{"$ne": address.NullAddress}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 2}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": bson.M{"$ne": address.NullAddress}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
		}
		pipeline = append(pipeline, pipeline1...)

	} else { //默认  account > 0
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": bson.M{"$ne": address.NullAddress}}},
			bson.M{"$match": bson.M{"market": bson.M{"$ne": args.PrimaryMarket.Val()}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
		}
		pipeline = append(pipeline, pipeline1...)
	}
	var deadlineCond bson.M
	if args.Sort == "deadline" { //按截止时间排序
		deadlineCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": bson.M{"$subtract": []interface{}{"$deadline", currentTime}}, "else": currentTime}}
	}
	var auctionAmountCond bson.M
	if args.Sort == "price" { // 将过期和未领取的放在后面
		if args.Order == -1 { //降序
			auctionAmountCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": "$auctionAmount", "else": 0}}
		} else { //升序（默认）
			auctionAmountCond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": "$auctionAmount", "else": 1e16}}
		}
	}

	//group
	setAndGroup := []bson.M{

		bson.M{"$lookup": bson.M{
			"from": "SelfControlNep11Properties",
			"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
					bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
				}}}},
				bson.M{"$set": bson.M{"class": bson.M{"$ifNull": []interface{}{"$name", "$tokenid"}}}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x50ac1c37690cc2cfc594472833cf57505d5f46de"}}, "then": "$asset", "else": "$class"}}}},
			},
			"as": "properties"},
		},
		bson.M{"$group": bson.M{"_id": bson.M{"asset": "$asset", "class": "$properties.class"}, "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "deadline": bson.M{"$last": "$deadline"}, "auctionAmount": bson.M{"$last": "$auctionAmount"}, "timestamp": bson.M{"$last": "$timestamp"}, "propertiesArr": bson.M{"$push": "$$ROOT"}}},
		bson.M{"$project": bson.M{"deadlineCond": deadlineCond, "auctionAmountCond": auctionAmountCond, "_id": 1, "properties": 1, "asset": 1, "tokenid": 1, "propertiesArr": 1, "auctionAmount": 1, "deadline": 1, "timestamp": 1}},

		bson.M{"$sort": bson.M{"timestamp": -1}},
	}
	var sort bson.M
	if args.Sort == "timestamp" { //上架时间
		sort = bson.M{"$sort": bson.M{"timestamp": args.Order}}
	} else if args.Sort == "price" { //价格
		sort = bson.M{"$sort": bson.M{"auctionAmountCond": args.Order}}
	} else if args.Sort == "deadline" { //截止时间
		sort = bson.M{"$sort": bson.M{"deadlineCond": args.Order}}
	} else {
		sort = bson.M{"$sort": bson.M{"timestamp": -1}}
	}
	setAndGroup = append(setAndGroup, sort)
	pipeline = append(pipeline, setAndGroup...)

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
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		if item["propertiesArr"] != nil {
			groupInfo := item["propertiesArr"].(primitive.A)
			var asset = item["asset"].(string)
			var tokenidArr []string
			for _, pitem := range groupInfo {
				it := pitem.(map[string]interface{})
				tokenid := it["tokenid"].(string)
				tokenidArr = append(tokenidArr, tokenid)
			}
			raw := make(map[string]interface{})
			err := me.GetInfoByNFT(struct {
				Asset   h160.T
				Tokenid []string
				Filter  map[string]interface{}
				Raw     *map[string]interface{}
			}{Asset: h160.T(asset), Tokenid: tokenidArr, Raw: &raw}, ret)
			if err != nil {
				return stderr.ErrGetNFTInfo
			}
			//rawResult := raw["result"].([]map[string]interface{})[0]

			newRaw := raw["result"].([]map[string]interface{})

			if args.NFTState.Val() == NFTstate.Sale.Val() {
				mapsort.MapSort9(newRaw, "buyNowAmount")
			} else if args.NFTState.Val() == NFTstate.Sale.Val() {
				mapsort.MapSort9(newRaw, "currentBidAmount")
			} else {
				mapsort.MapSort9(newRaw, "lastSoldAmount")
			}

			rawResult := newRaw[0]

			rawTokenid := rawResult["tokenid"]
			for _, it := range groupInfo {
				newit := it.(map[string]interface{})
				if rawTokenid == newit["tokenid"] {
					dst := make(map[string]interface{})
					dst = CopyMap(dst, newit)
					if dst["properties"] != nil {
						properties := dst["properties"].(primitive.A)
						pp := properties[0].(map[string]interface{})
						newProperties, err1 := ReSetProperties(pp)
						if err1 != nil {
							continue
						}
						dst["buyNowAsset"] = rawResult["buyNowAsset"]
						dst["buyNowAmount"] = rawResult["buyNowAmount"]
						dst["lastSoldAsset"] = rawResult["lastSoldAsset"]
						dst["lastSoldAmount"] = rawResult["lastSoldAmount"]
						dst["currentBidAsset"] = rawResult["currentBidAsset"]
						dst["currentBidAmount"] = rawResult["currentBidAmount"]
						dst["offerAsset"] = rawResult["offerAsset"]
						dst["offerAmount"] = rawResult["offerAmount"]

						dst["image"] = newProperties["image"]
						dst["thumbnail"] = newProperties["thumbnail"]
						dst["number"] = newProperties["number"]
						dst["properties"] = newProperties
						dst["class"] = newProperties["class"]
						dst["count"] = len(groupInfo)
						state := rawResult["state"]
						if state == "list" {
							auctionType := rawResult["auctionType"].(int32)
							if auctionType == 1 {
								dst["state"] = NFTstate.Sale.Val()
							} else {
								dst["state"] = NFTstate.Auction.Val()
							}
						}
						result = append(result, dst)
					}
				}
			}
		}

	}

	//  分页
	if args.Limit == 0 {
		args.Limit = int64(math.Inf(1))
	}

	pageResult := make([]map[string]interface{}, 0)
	for i, item := range result {
		if int64(i) < args.Skip {
			continue
		} else if int64(i) > args.Skip+args.Limit-1 {
			continue
		} else {
			pageResult = append(pageResult, item)
		}
	}

	r3, err := me.FilterAggragateAndAppendCount(pageResult, len(result), args.Filter)

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

func ReSetProperties(p map[string]interface{}) (map[string]interface{}, error) {
	if p["image"] == nil {
		if p["properties"] != nil && p["tokenURI"] == nil {
			pp := p["properties"].(string)
			if pp != "" {
				var data map[string]interface{}
				if err1 := json.Unmarshal([]byte(pp), &data); err1 == nil {
					for key, value := range data {
						p[key] = value
					}
				} else {
					return nil, err1
				}
			}
		} else if p["tokenURI"] != nil {
			asset := p["asset"].(string)
			tokenid := p["tokenid"].(string)
			tokenuri := p["tokenURI"].(string)
			ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri), asset, tokenid)
			if err != nil {
				return nil, err
			}
			for key, value := range ppjson {
				p[key] = value
			}
		}
	}

	//处理number
	name := p["name"]
	if name != nil {
		strArray := strings.Split(name.(string), "#")
		if len(strArray) >= 2 {
			number := strArray[1]
			n, err22 := strconv.ParseInt(number, 10, 64)
			if err22 != nil {
				p["number"] = int64(-1)
			}
			p["number"] = n
		} else {
			p["number"] = int64(-1)
		}
	}

	delete(p, "properties")
	delete(p, "_id")
	delete(p, "tokenURI")
	return p, nil
}