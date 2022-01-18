package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetMarketWhiteList(args struct {
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
		bson.M{"$group": bson.M{"_id": "$asset", "info": bson.M{"$push": "$$ROOT"}, "asset": bson.M{"$last": "$asset"}, "eventname": bson.M{"$last": "$eventname"}, "extendData": bson.M{"$last": "$extendData"}}},
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
	result := map[string]interface{}{}
	var assetArr []string
	for _, item := range r1 {
		if item["eventname"].(string) == "AddAsset" {
			assetArr = append(assetArr, item["asset"].(string))
		}
	}
	result["market"] = args.MarketHash.Val()
	result["whiteList"] = assetArr

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = result
	}
	*ret = json.RawMessage(r)
	return nil
}
