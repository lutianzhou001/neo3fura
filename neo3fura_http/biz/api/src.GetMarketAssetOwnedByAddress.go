package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetMarketAssetOwnedByAddress(args struct {
	Address    h160.T
	MarketHash h160.T
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6

	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{}

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

	//默认  account > 0
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

		bson.M{"$group": bson.M{"_id": "asset", "asset": bson.M{"$last": "$asset"}, "marketAsset": bson.M{"$push": "$$ROOT"}}},
		bson.M{"$project": bson.M{"_id": 1, "asset": 1, "marketAsset": 1}},
	}
	pipeline = append(pipeline, pipeline1...)

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

	//获取nft的属性
	var assetList []interface{}
	for _, item := range r1 {
		asset := item["asset"]
		assetList = append(assetList, asset)
	}

	result := make(map[string]interface{})
	result["assetlist"] = assetList
	result["result"] = r1

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
