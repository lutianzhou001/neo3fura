package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
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
	//if args.MarketHash.Valid() == false {
	//	return stderr.ErrInvalidArgs
	//}

	length := 0
	cond := bson.M{}
	if len(args.SubClass) > 0 {

		for _, i := range args.SubClass {
			b := bson.M{}
			if len(i) != 2 || i[0] > i[1] {
				return stderr.ErrInvalidArgs
			} else {
				//a:=bson.M{"$and":[]interface{}{bson.M{"tokenid":bson.M{"$gte":i[0].Val()}},bson.M{"tokenid":bson.M{"$lte":i[1].Val()}}}}
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
		//bson.M{"$match": bson.M{"market": args.AssetHash}},
		bson.M{"$match": bson.M{"asset": args.AssetHash}},
		bson.M{"$match": bson.M{"eventname": "Claim"}},
		bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "extendData": bson.M{"$last": "$extendData"}, "claimed": bson.M{"$sum": 1}}},

		//bson.M{"$project":bson.M{"class":bson.M{"if":bson.M{"$and":[]interface{}{bson.M{"$gte":}}}}}},
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

		if item["_id"].(int32) != -1 {
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)
			extendData := item["extendData"].(string)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
				if err != nil {
					return err
				}
				auctionAsset := dat["auctionAsset"]
				bidAmount := dat["bidAmount"]
				item["price"] = bidAmount
				item["sellAsset"] = auctionAsset
			} else {
				return err
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
