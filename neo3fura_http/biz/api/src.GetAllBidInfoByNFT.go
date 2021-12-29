package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
)

func (me *T) GetAllBidInfoByNFT(args struct {
	AssetHash h160.T
	TokenId   strval.T
	Filter    map[string]interface{}
	Raw       *[]map[string]interface{}
}, ret *json.RawMessage) error {
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	var f bson.M
	if args.TokenId == "" {
		f = bson.M{"asset": args.AssetHash.Val(), "eventname": "Bid"}
	} else {
		f = bson.M{"asset": args.AssetHash.Val(), "tokenid": args.TokenId.Val(), "eventname": "Bid"}
	}

	r11, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetAllBidInfoByNFT",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": f},
				bson.M{"$group": bson.M{"_id": "$nonce",
					"nonce":   bson.M{"$last": "$nonce"},
					"asset":   bson.M{"$last": "$asset"},
					"tokenid": bson.M{"$last": "$tokenid"},
					"bidInfo": bson.M{"$push": "$$ROOT"}}},
				bson.M{"$sort": bson.M{"_id": -1}},
			},
			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, items := range r11 {
		bidInfo := make(map[string]interface{})
		bidinfos := items["bidInfo"].(primitive.A)
		bidInfo["tokenid"] = items["tokenid"]
		bidInfo["nonce"] = items["nonce"]
		bidInfo["asset"] = items["asset"]
		bidAmounts := []int64{}
		bidders := []string{}
		for _, i := range bidinfos {
			info := i.(map[string]interface{})
			bidder := info["user"].(string)
			bidders = append(bidders, bidder)

			extendData := info["extendData"].(string)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
				bidAmount, err := strconv.ParseInt(dat["bidAmount"].(string), 10, 64)
				if err != nil {
					return err
				}
				bidAmounts = append(bidAmounts, bidAmount)
			} else {
				return err
			}
		}
		bidInfo["bidAmount"] = bidAmounts
		bidInfo["bidder"] = bidders
		result = append(result, bidInfo)

	}

	num, err := strconv.ParseInt(strconv.Itoa(len(result)), 10, 64)
	if err != nil {
		return err
	}

	if args.Raw != nil {
		*args.Raw = result
	}

	r2, err := me.FilterArrayAndAppendCount(result, num, args.Filter)
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
