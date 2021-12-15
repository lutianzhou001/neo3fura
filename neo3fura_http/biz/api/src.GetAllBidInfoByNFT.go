package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/lib/utils"
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

	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "MarketNotification",
		Index:      "GetAllBidInfoByNFT",
		Sort:       bson.M{"timestamp": -1},
		Filter:     f,
		Query:      []string{},
	}, ret)

	groups := utils.GroupByString(r1, "nonce")

	result := make([]map[string]interface{}, 0)
	for _, items := range groups {
		items = mapsort.MapSort(items, "timestamp") //
		bidInfo := make(map[string]interface{})
		bidInfo["asset"] = items[0]["asset"]
		bidInfo["tokenid"] = items[0]["tokenid"]
		bidInfo["nonce"] = items[0]["nonce"]

		//var bal int64 = 0
		bidAmounts := []int64{}
		bidders := []string{}
		for _, item := range items {
			bidder := item["user"].(string)

			extendData := item["extendData"].(string)
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

			bidders = append(bidders, bidder)

		}
		bidInfo["bidder"] = bidders
		bidInfo["bidAmount"] = bidAmounts
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
