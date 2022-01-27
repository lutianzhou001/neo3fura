package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetMarketTokenidList(args struct {
	AssetHash  h160.T
	MarketHash h160.T
	SubClass   [][]strval.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	length := 0
	cond := bson.M{}
	if len(args.SubClass) > 0 {
		for _, i := range args.SubClass {
			b := bson.M{}
			if len(i) != 2 || i[0] > i[1] {
				return stderr.ErrInvalidArgs
			} else {
				a := bson.M{"$and": []interface{}{bson.M{"$gte": []interface{}{"$tokenid", i[0].Val()}}, bson.M{"$lte": []interface{}{"$tokenid", i[1].Val()}}}}
				if length == 0 {
					b = bson.M{"if": a, "then": length, "else": length - 1}
				} else {
					b = bson.M{"if": a, "then": length, "else": cond}
				}
				length++
			}
			cond = bson.M{"$cond": b}
		}
	} else {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
		bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
		bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
		bson.M{"$match": bson.M{"market": bson.M{"$eq": args.MarketHash.Val()}}},
		bson.M{"$match": bson.M{"asset": bson.M{"$eq": args.AssetHash.Val()}}},

		bson.M{"$project": bson.M{"class": cond, "_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1}},
		bson.M{"$match": bson.M{"difference": true}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "market": bson.M{"$last": "$market"}, "tokenid": bson.M{"$push": "$$ROOT"}}},
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
			Query:      []string{"asset", "tokenid", "market", "_id"},
		}, ret)

	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)

	for _, item := range r1 {
		res := make([]map[string]interface{}, 0)
		tokenids := item["tokenid"]
		if tokenids != nil {
			tokenidArr := tokenids.(primitive.A)
			for _, it := range tokenidArr {
				r := make(map[string]interface{})
				nft := it.(map[string]interface{})
				r["tokenid"] = nft["tokenid"]

				res = append(res, r)
			}
			mapsort.MapSort8(res, "tokenid")
			item["tokenidArr"] = res
		}

	}

	for _, item := range r1 {

		id := item["_id"].(int32)
		if id != -1 {
			re := make(map[string]interface{})
			tokenidArr := item["tokenidArr"].([]map[string]interface{})
			tokenids := []string{}
			for _, tk := range tokenidArr {
				tokenids = append(tokenids, tk["tokenid"].(string))
			}
			re["tokenid"] = tokenids
			re["id"] = id
			result = append(result, re)
		}

	}
	mapsort.MapSort5(result, "id")
	count := len(result)

	r3, err := me.FilterAggragateAndAppendCount(result, count, args.Filter)

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
