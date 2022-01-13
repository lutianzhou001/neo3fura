package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"
)

func (me *T) GetNFTMarket(args struct {
	ContractHash h160.T   //  asset
	AssetHash    h160.T   // auctionType
	MarketHash   h160.T   //
	NFTState     strval.T //state:aution  sale  notlisted  unclaimed
	Sort         strval.T //listedTime  price  deadline
	Order        int64    //-1:降序  +1：升序
	SubClass     [][]strval.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	pipeline := []bson.M{}

	if len(args.ContractHash) > 0 && args.ContractHash != "" {
		if args.ContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"asset": args.ContractHash}}
			pipeline = append(pipeline, a)
		}
	}

	if len(args.AssetHash) > 0 && args.AssetHash != "" {
		if args.AssetHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"auctionAsset": args.AssetHash}}
			pipeline = append(pipeline, a)
		}
	}

	if len(args.MarketHash) > 0 && args.MarketHash != "" {
		if args.MarketHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"market": args.MarketHash}}
			pipeline = append(pipeline, a)
		}
	}

	if args.Sort == "deadline" { //按截止时间排序
		b := bson.M{}
		if args.Order == -1 || args.Order == 1 {
			b = bson.M{"$sort": bson.M{args.Sort.Val(): args.Order}}
		} else {
			b = bson.M{"$sort": bson.M{args.Sort.Val(): -1}}
		}
		pipeline = append(pipeline, b)
	} else if args.Sort == "price" { //按上架价格排序
		b := bson.M{}
		if args.Order == -1 || args.Order == 1 {
			b = bson.M{"$sort": bson.M{"auctionAmount": args.Order}}
		} else {
			b = bson.M{"$sort": bson.M{"auctionAmount": 1}}
		}
		pipeline = append(pipeline, b)
	}

	if args.NFTState.Val() == NFTstate.Auction.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline1 := []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 2}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"eventname": "Auction"}},
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						bson.M{"$eq": []interface{}{"$market", "$$market"}},
					}}}},
					bson.M{"$group": bson.M{"_id": bson.M{"tokenid": "$$tokenid", "asset": "$$asset", "market": "$$market"},
						"nonce":     bson.M{"$last": "$nonce"},
						"asset":     bson.M{"$last": "$asset"},
						"tokenid":   bson.M{"$last": "$tokenid"},
						"timestamp": bson.M{"$last": "$timestamp"},
					}},
					bson.M{"$project": bson.M{"asset": 1, "nonce": 1, "tokenid": 1, "timestamp": 1}},
				},
				"as": "marketnotification"},
			},

			bson.M{"$project": bson.M{"_id": 1, "marketnotification": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "auction"}},
			bson.M{"$match": bson.M{"difference": true}},

			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline1 := []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$lookup": bson.M{
				"from": "MarketNotification",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid", "market": "$market"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"eventname": "Auction"}},
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						bson.M{"$eq": []interface{}{"$market", "$$market"}},
					}}}},
					bson.M{"$group": bson.M{"_id": bson.M{"tokenid": "$$tokenid", "asset": "$$asset", "market": "$$market"},
						"nonce":     bson.M{"$last": "$nonce"},
						"asset":     bson.M{"$last": "$asset"},
						"tokenid":   bson.M{"$last": "$tokenid"},
						"timestamp": bson.M{"$last": "$timestamp"},
					}},
					bson.M{"$project": bson.M{"asset": 1, "nonce": 1, "tokenid": 1, "timestamp": 1}},
				},
				"as": "marketnotification"},
			},

			bson.M{"$project": bson.M{"_id": 1, "marketnotification": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "sale"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else if args.NFTState.Val() == NFTstate.NotListed.Val() { //未上架  accont >0 && owner != market  ||  owner == market && deadline < currentTime
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$project": bson.M{"_id": 1, "marketnotification": "", "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "notlisted"}},
			bson.M{"$match": bson.M{"$or": []interface{}{
				bson.M{"difference": false},
				bson.M{"$and": []interface{}{
					bson.M{"deadline": bson.M{"$lt": currentTime}},
					bson.M{"difference": true},
				}}}}},

			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
		pipeline = append(pipeline, pipeline1...)

	} else { //默认  account > 0
		pipeline1 := []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "notlisted"}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
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
			Query:      []string{"_id", "asset", "marketnotification", "tokenid", "amount", "owner", "market", "auctionType", "auctor", "auctionAsset", "auctionAmount", "deadline", "bidder", "bidAmount", "timestamp", "state"},
		}, ret)

	if err != nil {
		return err
	}

	//获取nft的属性
	for _, item := range r1 {
		if args.NFTState.Val() != NFTstate.Auction.Val() && args.NFTState.Val() != NFTstate.Sale.Val() {
			a := item["amount"].(primitive.Decimal128).String()
			amount, err1 := strconv.Atoi(a)
			if err1 != nil {
				return err
			}

			ba := item["bidAmount"].(primitive.Decimal128).String()
			bidAmount, err2 := strconv.ParseInt(ba, 10, 64)
			if err2 != nil {
				return err
			}

			deadline, _ := item["deadline"].(int64)
			auctionType, _ := item["auctionType"].(int32)

			if amount > 0 && auctionType == 2 && item["owner"] == item["market"] && deadline > currentTime {
				item["state"] = NFTstate.Auction.Val()
			} else if amount > 0 && auctionType == 1 && item["owner"] == item["market"] && deadline > currentTime {
				item["state"] = NFTstate.Sale.Val()
			} else if amount > 0 && item["owner"] != item["market"] {
				item["state"] = NFTstate.NotListed.Val()
			} else if amount > 0 && bidAmount > 0 && deadline < currentTime && item["owner"] == item["market"] {
				item["state"] = NFTstate.Unclaimed.Val()
			} else if amount > 0 && deadline < currentTime && bidAmount == 0 && item["owner"] == item["market"] {
				item["state"] = NFTstate.Expired.Val()
			} else {
				item["state"] = ""
			}

		}
		//获得上架时间

		if item["marketnotification"] != "" {
			marketnotification := item["marketnotification"].(primitive.A)
			if len(marketnotification) > 0 {
				mn := []interface{}(marketnotification)[0].(map[string]interface{})
				item["listedTimestamp"] = mn["timestamp"]
			}

		} else {
			item["listedTimestamp"] = 0
		}

		delete(item, "marketnotification")

		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)

		var raw3 map[string]interface{}
		err1 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
		if err1 != nil {
			item["image"] = ""
			item["name"] = ""
		}
		extendData := raw3["properties"].(string)
		if extendData != "" {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &data); err == nil {
				image, ok := data["image"]
				if ok {
					item["image"] = image
				} else {
					item["image"] = ""
				}
				name, ok1 := data["name"]
				if ok1 {
					item["name"] = name
				} else {
					item["name"] = ""
				}

			} else {
				return err
			}

		} else {
			item["image"] = ""
			item["name"] = ""
		}

	}
	//  按上架时间排序
	if args.Sort == "timestamp" {
		if args.Order == 1 {
			mapsort.MapSort(r1, "listedTimestamp")
		} else {
			mapsort.MapSort2(r1, "listedTimestamp")
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
