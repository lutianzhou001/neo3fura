package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetWhiteListByMarketHash(args struct {
	MarketHash h160.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash.Val()}},
		bson.M{"$match": bson.M{"$or": []interface{}{
			bson.M{"eventname": "AddAsset"},
			bson.M{"eventname": "RemoveAsset"},
		}}},
		bson.M{"$sort": bson.M{"timestamp": 1}},
		bson.M{"$group": bson.M{"_id": "$asset", "asset": bson.M{"$last": "$asset"}, "eventname": bson.M{"$last": "$eventname"}, "extendData": bson.M{"$last": "$extendData"}}},

		bson.M{"$lookup": bson.M{
			"from": "Asset",
			"let":  bson.M{"asset": "$asset"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$hash", "$$asset"}},
				}}}},
			},
			"as": "properties"},
		},
		bson.M{"$project": bson.M{"properties": 1, "asset": 1, "eventname": 1, "extendData": 1}},
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

	//result := make([]map[string]interface{},0 )

	result := make([]map[string]interface{}, 0)
	var assetArr []string
	for _, item := range r1 {
		if item["eventname"].(string) == "AddAsset" {
			rr := make(map[string]interface{})
			assetArr = append(assetArr, item["asset"].(string))
			rr["asset"] = item["asset"]
			p := item["properties"].(primitive.A)
			if len(p) > 0 {
				it := p[0].(map[string]interface{})
				rr["type"] = it["type"]
				rr["tokenname"] = it["tokenname"]
				rr["decimal"] = it["decimals"]
				rr["symbol"] = it["symbol"]
				rr["totalsupply"] = it["totalsupply"]
			}

			extendData := item["extendData"].(string)
			if extendData != "" {
				var dat map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {
					if err1 != nil {
						return err
					}
					rr["feeRate"] = dat["feeRate"]
					rr["rewardRate"] = dat["rewardRate"]
					rr["rewardReceiveAddress"] = dat["rewardReceiveAddress"]
				} else {
					return err
				}
			}
			result = append(result, rr)
		}

	}

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
