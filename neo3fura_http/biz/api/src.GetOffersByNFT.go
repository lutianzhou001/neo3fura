package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetOffersByNFT(args struct {
	Asset      h160.T
	TokenId    strval.T
	MarketHash h160.T
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"asset": args.Asset.Val(), "tokenid": args.TokenId.Val(), "eventname": "Offer", "market": args.MarketHash.Val()}},

		bson.M{"$lookup": bson.M{
			"from": "Nep11Properties",
			"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
					bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
				}}}},
				bson.M{"$project": bson.M{"properties": 1}},
			},
			"as": "properties"},
		},

		bson.M{"$sort": bson.M{"timestamp": -1}},
		bson.M{"$limit": args.Limit},
		bson.M{"$skip": args.Skip},
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
			Index:      "GetOffersByNFT",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {
		extendData := item["extendData"].(string)
		var dat map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {
			item["originOwner"] = dat["originOwner"]
			item["offerAsset"] = dat["offerAsset"]
			item["offerAmount"] = dat["offerAmount"]
			item["deadline"] = dat["deadline"]

		} else {
			return err1
		}

		nftproperties := item["properties"]
		if nftproperties != nil && nftproperties != "" {
			pp := nftproperties.(primitive.A)
			if len(pp) > 0 {
				it := pp[0].(map[string]interface{})
				extendData := it["properties"].(string)
				if extendData != "" {
					properties := make(map[string]interface{})
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
						image, ok := data["image"]
						if ok {
							properties["image"] = image
							item["image"] = image
						} else {
							item["image"] = ""
						}

					} else {
						return err
					}

				} else {
					item["image"] = ""
				}
			}
		}
		delete(item, "extendData")
		delete(item, "properties")
	}
	count := int64(len(r1))
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil

}
