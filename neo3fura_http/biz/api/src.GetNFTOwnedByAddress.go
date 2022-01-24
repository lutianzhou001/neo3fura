package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
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
	var cond bson.M
	if args.Sort == "deadline" { //按截止时间排序
		cond = bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$deadline", currentTime}}, "then": bson.M{"$subtract": []interface{}{"$deadline", currentTime}}, "else": bson.M{"$subtract": []interface{}{currentTime, "$deadline"}}}}
	}

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
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"user": args.Address.Val()}},
					bson.M{"$match": bson.M{"eventname": "Auction"}},
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						bson.M{"$eq": []interface{}{"$market", "$$market"}},
					}}}},
					bson.M{"$sort": bson.M{"nonce": -1}},
					bson.M{"$limit": 1},
					bson.M{"$project": bson.M{"asset": 1, "nonce": 1, "tokenid": 1, "timestamp": 1}}},
				"as": "marketnotification"},
			},
			bson.M{"$project": bson.M{"date": cond, "_id": 1, "asset": 1, "marketnotification": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "auction"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$sort": bson.M{"marketnotification": -1}},
			bson.M{"$sort": bson.M{"date": 1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"auctor": args.Address.Val()}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"user": args.Address.Val()}},
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
			},
			bson.M{"$project": bson.M{"date": cond, "_id": 1, "asset": 1, "marketnotification": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "sale"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$sort": bson.M{"marketnotification": -1}},
			bson.M{"$sort": bson.M{"date": 1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.NotListed.Val() { //未上架  accont >0 && owner != market
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": args.Address.Val()}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},

			bson.M{"$project": bson.M{"date": cond, "_id": 1, "asset": 1, "marketnotification": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "notlisted"}},
			bson.M{"$match": bson.M{"difference": false}},
			bson.M{"$sort": bson.M{"date": 1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Unclaimed.Val() { //未领取 (买家和卖家) accont >0 &&  runtime > deadline && owner== market && bidAccount >0
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"$or": []interface{}{
				bson.M{"$and": []interface{}{ //卖家未领取(过期)
					bson.M{"auctor": args.Address.Val()},
					bson.M{"bidAmount": bson.M{"$eq": 0}},
				}},
				bson.M{"$and": []interface{}{ //买家未领取
					bson.M{"bidder": args.Address.Val()},
					bson.M{"bidAmount": bson.M{"$gt": 0}},
				}},
			}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$lt": currentTime}}},

			bson.M{"$lookup": bson.M{
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"user": args.Address.Val()}},
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
			},
			bson.M{"$project": bson.M{"_id": 1, "date": cond, "asset": 1, "marketnotification": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "unclaimed"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$sort": bson.M{"marketnotification": -1}},
			bson.M{"$sort": bson.M{"date": 1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else { //默认  account > 0
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"$or": []interface{}{
				bson.M{"auctor": args.Address.Val()},
				bson.M{"owner": args.Address.Val()}}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},

			bson.M{"$lookup": bson.M{
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"user": args.Address.Val()}},
					bson.M{"$match": bson.M{"eventname": "Auction"}},
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						//bson.M{"$eq": []interface{}{"$market", "$$market"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "nonce": 1, "tokenid": 1, "timestamp": 1}},
					bson.M{"$sort": bson.M{"nonce": -1}},
					bson.M{"$limit": 1},
				},
				"as": "marketnotification"},
			},

			bson.M{"$project": bson.M{"_id": 1, "date": cond, "asset": 1, "marketnotification": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": ""}},
			bson.M{"$sort": bson.M{"marketnotification": -1}},
			bson.M{"$sort": bson.M{"date": 1}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)
	}

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
			Query:      []string{"_id", "date", "marketnotification", "asset", "tokenid", "amount", "owner", "market", "auctionType", "auctor", "auctionAsset", "auctionAmount", "deadline", "bidder", "bidAmount", "timestamp", "state"},
		}, ret)

	if err != nil {
		return err
	}

	//获取nft的属性
	for _, item := range r1 {
		if args.NFTState.Val() != NFTstate.Auction.Val() && args.NFTState.Val() != NFTstate.Sale.Val() && args.NFTState.Val() != NFTstate.Expired.Val() && args.NFTState.Val() != NFTstate.Unclaimed.Val() {
			a := item["amount"].(primitive.Decimal128).String()
			amount, err := strconv.Atoi(a)
			if err != nil {
				return err
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

		//价格转换
		auctionAsset := item["auctionAsset"]
		auctionAmount, _, err2 := item["auctionAmount"].(primitive.Decimal128).BigInt()
		if err2 != nil {
			return err2
		}
		if auctionAsset != nil {
			dd, _ := OpenAssetHashFile()
			decimal := dd[auctionAsset.(string)] //获取精度
			if decimal == 0 {
				decimal = 1
			}
			price, err3 := GetPrice(auctionAsset.(string)) //  获取价格
			if err3 != nil {
				return err3
			}

			bfauctionAmount := new(big.Float).SetInt(auctionAmount)
			flag := auctionAmount.Cmp(big.NewInt(0))
			if flag == 1 {
				bfprice := big.NewFloat(price)
				ffprice := big.NewFloat(1).Mul(bfprice, bfauctionAmount)
				de := math.Pow(10, float64(decimal))
				usdAuctionAmount := new(big.Float).Quo(ffprice, big.NewFloat(de))
				item["usdAuctionAmount"] = usdAuctionAmount
			} else {
				item["usdAuctionAmount"] = 0
			}
		}

		//获得上架时间
		if item["marketnotification"] != nil {
			switch item["marketnotification"].(type) {
			case string:
				item["listedTimestamp"] = int64(0)
			case primitive.A:
				marketnotification := item["marketnotification"].(primitive.A)
				if len(marketnotification) > 0 {
					mn := []interface{}(marketnotification)[0].(map[string]interface{})
					item["listedTimestamp"] = mn["timestamp"]
				} else {
					item["listedTimestamp"] = int64(0)
				}
			}
		} else {
			item["listedTimestamp"] = int64(0)
		}
		delete(item, "marketnotification")

		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)

		var raw3 map[string]interface{}
		err1 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
		if err1 != nil {
			item["image"] = ""
			item["name"] = ""
			item["number"] = int64(-1)
			item["properties"] = ""

		}

		item["image"] = raw3["image"]
		item["name"] = raw3["name"]
		item["number"] = raw3["number"]
		item["properties"] = raw3["properties"]

	}

	//  按上架时间排序
	if args.Sort == "timestamp" {
		if args.Order == 1 {
			mapsort.MapSort(r1, "listedTimestamp")
		} else {
			mapsort.MapSort2(r1, "listedTimestamp")
		}

	}

	//按价格排序
	if args.Sort == "price" {
		if args.Order == 1 {
			mapsort.MapSort(r1, "usdAuctionAmount")
		}
	}

	// 按上架时间排序
	if args.Sort == "deadline" {
		count := 0
		var arr []int
		for _, i := range r1 {
			count++
			if i["state"] == "unclaimed" {
				arr = append(arr, count)
			}
		}
		//删除未领取的nft
		for i := len(arr) - 1; i >= 0; i-- {
			num := arr[i]
			r1 = append(r1[:num-1], r1[num:]...)

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
