package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNFTClass(args struct {
	MarketHash h160.T
	AssetHash  h160.T
	SubClass   [][]strval.T
	State      strval.T
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
	if args.State.Val() == NFTstate.Mint.Val() {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"asset": args.AssetHash}},
			bson.M{"$group": bson.M{"_id": "$tokenid", "tokenid": bson.M{"$last": "$tokenid"}, "asset": bson.M{"$last": "$asset"}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"properties": 1, "_id": 0}}},
				"as": "properties"},
			},
			bson.M{"$project": bson.M{"class": cond, "asset": 1, "set": 1, "tokenid": 1, "properties": 1}},
			bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "properties": bson.M{"$last": "$properties"}}},
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
				Collection: "Address-Asset",
				Index:      "GetNFTClass",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline:   pipeline,
				Query:      []string{},
			}, ret)

		if err != nil {
			return err
		}
		for _, item := range r1 {
			if item["_id"].(int32) != -1 {
				p := item["properties"].(primitive.A)[0].(map[string]interface{})

				properties := p["properties"].(string)
				if properties != "" {
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(properties), &data); err1 == nil {
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
						//owner, ok2 := data["owner"]
						//if ok2{
						//	item["owner"] = owner
						//} else {
						//	item["owner"] = ""
						//}

					} else {
						return err1
					}

				} else {
					item["image"] = ""
					item["name"] = ""
				}
				delete(item, "_id")
				delete(item, "properties")
				item["price"] = "——"
				item["sellAsset"] = "——"
				result = append(result, item)
			}

		}

	} else if args.State.Val() == NFTstate.Listed.Val() {
		pipeline := []bson.M{
			bson.M{"$match": bson.M{"market": args.MarketHash}},
			bson.M{"$match": bson.M{"asset": args.AssetHash}},
			bson.M{"$match": bson.M{"eventname": "Auction"}},
			bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
			bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "extendData": bson.M{"$last": "$extendData"}, "claimed": bson.M{"$sum": 1}}},
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

		for _, item := range r1 {
			if item["_id"].(int32) != -1 {
				asset := item["asset"].(string)
				tokenid := item["tokenid"].(string)
				extendData := item["extendData"].(string)
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
				}
				properties := raw3["properties"].(string)
				if properties != "" {
					var data map[string]interface{}
					if err := json.Unmarshal([]byte(properties), &data); err == nil {
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
				delete(item, "_id")
				delete(item, "extendData")
				delete(item, "tokenid")

				result = append(result, item)
			}

		}

	} else if args.State.Val() == NFTstate.Selling.Val() {

		pipeline := []bson.M{
			bson.M{"$match": bson.M{"market": args.MarketHash}},
			bson.M{"$match": bson.M{"asset": args.AssetHash}},
			bson.M{"$match": bson.M{"eventname": "Claim"}},
			bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
			bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "extendData": bson.M{"$last": "$extendData"}, "claimed": bson.M{"$sum": 1}}},
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

		for _, item := range r1 {
			if item["_id"].(int32) != -1 {
				asset := item["asset"].(string)
				tokenid := item["tokenid"].(string)
				extendData := item["extendData"].(string)
				var dat map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {

					auctionAsset := dat["auctionAsset"]
					bidAmount := dat["bidAmount"]
					item["price"] = bidAmount
					item["sellAsset"] = auctionAsset
				} else {
					return err1
				}

				var raw3 map[string]interface{}
				err1 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
				if err1 != nil {
					item["image"] = ""
					item["name"] = ""
				}
				properties := raw3["properties"].(string)
				if properties != "" {
					var data map[string]interface{}
					if err2 := json.Unmarshal([]byte(properties), &data); err2 == nil {
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
						return err2
					}

				} else {
					item["image"] = ""
					item["name"] = ""
				}
				delete(item, "_id")
				delete(item, "extendData")
				delete(item, "tokenid")

				result = append(result, item)
			}

		}

	} else {
		return stderr.ErrInvalidArgs
	}

	count := len(args.SubClass)

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
