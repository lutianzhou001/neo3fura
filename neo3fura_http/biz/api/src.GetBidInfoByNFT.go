package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetBidInfoByNFT(args struct {
	Address   h160.T
	AssetHash h160.T
	TokenId   strval.T
	Filter    map[string]interface{}
	Raw       *[]map[string]interface{}
}, ret *json.RawMessage) error {
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	//获取NFT 最新一轮上架的竞价信息
	var f bson.M
	if args.Address != "" {
		if args.AssetHash.Valid() == false || args.Address.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		if args.TokenId == "" {
			f = bson.M{"asset": args.AssetHash.Val(), "user": args.Address.Val(), "eventname": "Bid"}
		} else {
			f = bson.M{"asset": args.AssetHash.Val(), "tokenid": args.TokenId.Val(), "user": args.Address.Val(), "eventname": "Bid"}
		}

	} else {
		if args.TokenId == "" {
			f = bson.M{"asset": args.AssetHash.Val(), "eventname": "Bid"}
		} else {
			f = bson.M{"asset": args.AssetHash.Val(), "tokenid": args.TokenId.Val(), "eventname": "Bid"}
		}
	}
	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "MarketNotification",
		Index:      "GetBidInfoByNFT",
		Sort:       bson.M{"timestamp": -1},
		Filter:     f,
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		rr := make(map[string]interface{})
		rr["tokenid"] = item["tokenid"]
		rr["asset"] = item["asset"]
		rr["bidder"] = item["user"]
		rr["timestamp"] = item["timestamp"]

		extendData := item["extendData"].(string)
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
			if err != nil {
				return err
			}
			auctionAsset := dat["auctionAsset"]
			bidAmount := dat["bidAmount"]
			rr["bidAmount"] = bidAmount
			rr["auctionAsset"] = auctionAsset
		} else {
			return err
		}
		result = append(result, rr)
	}
	r2, err := me.FilterArrayAndAppendCount(result, count, args.Filter)
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
