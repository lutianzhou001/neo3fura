package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
)

func (me *T) GetNFTClass(args struct {
	MarketHash h160.T
	AssetHash  h160.T
	SubClass   [][]strval.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
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
	result := make([]map[string]interface{}, 0)

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash}},
		bson.M{"$match": bson.M{"asset": args.AssetHash}},
		bson.M{"$match": bson.M{"eventname": "Auction"}},
		bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "tokenidArr": bson.M{"$push": "$$ROOT"}, "extendData": bson.M{"$last": "$extendData"}}},
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
			Collection: "MarketNotification",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	//  获取claimed 的值
	pipeline2 := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash}},
		bson.M{"$match": bson.M{"asset": args.AssetHash}},
		bson.M{"$match": bson.M{"eventname": "Claim"}},
		bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "claimedInfo": bson.M{"$push": "$$ROOT"}, "extendData": bson.M{"$last": "$extendData"}, "claimed": bson.M{"$sum": 1}}},
	}

	r2, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline2,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {
		if item["_id"].(int32) != -1 {
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)
			extendData := item["extendData"].(string)
			//
			tokenidList := make(map[string]interface{})
			if item["tokenidArr"] != nil {
				tokenidArr := item["tokenidArr"].(primitive.A)

				for _, it := range tokenidArr {
					nft := it.(map[string]interface{})
					tokenidList[nft["tokenid"].(string)] = 1
				}
			}

			var dat map[string]interface{}
			if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {

				auctionAsset := dat["auctionAsset"]
				auctionAmount := dat["auctionAmount"]
				item["price"] = auctionAmount
				item["sellAsset"] = auctionAsset
			} else {
				return err1
			}

			var raw3 map[string]interface{}
			err1 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
			if err1 != nil {
				item["image"] = ""
				item["name"] = ""
				item["series"] = ""
				item["supply"] = ""
				item["number"] = ""
			}
			properties := raw3["properties"].(string)
			if properties != "" {
				var data map[string]interface{}
				if err11 := json.Unmarshal([]byte(properties), &data); err11 == nil {
					image, ok := data["image"]
					if ok {
						item["image"] = image
					} else {
						item["image"] = ""
					}
					name, ok1 := data["name"]
					if ok1 {
						item["name"] = name
						num := strings.Split(name.(string), "#")[1]
						number, err12 := strconv.ParseInt(num, 10, 64)
						if err12 != nil {
							return err12
						}
						item["number"] = number
					} else {
						item["name"] = ""
					}
					series, ok2 := data["series"]
					if ok2 {
						item["series"] = series
					} else {
						item["series"] = ""
					}
					supply, ok3 := data["supply"]
					if ok3 {
						item["supply"] = supply
					} else {
						item["supply"] = ""
					}

				} else {
					return err
				}

			} else {
				item["image"] = ""
				item["name"] = ""
				item["series"] = ""
				item["supply"] = ""
				item["number"] = ""

			}

			//获取climedi
			if len(r2) > 0 {
				for _, item1 := range r2 {
					if item["_id"] == item1["_id"] {
						item["claimed"] = item1["claimed"]
						break
					} else {
						item["claimed"] = 0
					}

				}
			} else {
				item["claimed"] = 0
			}
			delete(item, "_id")
			delete(item, "extendData")
			delete(item, "tokenid")
			delete(item, "tokenidArr")
			result = append(result, item)
		}
	}
	mapsort.MapSort2(result, "number")

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
