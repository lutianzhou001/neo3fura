package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"neo3fura_http/lib/type/NFTevent"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"
)

func (me *T) GetNFTActivityByAsset(args struct {
	Asset  h160.T
	Market h160.T
	State  string
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Market.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	var pipeline []bson.M
	if args.State == "sales" {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market, "eventname": bson.M{"$in": []interface{}{"Claim", "CompleteOffer"}}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$sort": bson.M{"timestamp": -1}},
		}
	} else if args.State == "listings" {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market, "eventname": "Auction"}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$sort": bson.M{"timestamp": -1}},
		}
	} else if args.State == "offers" {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market, "eventname": "Offer"}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$sort": bson.M{"timestamp": -1}},
		}
	} else {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
			bson.M{"$sort": bson.M{"timestamp": -1}},
		}
	}

	result := make([]map[string]interface{}, 0)

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
			Index:      "GetNFTActivityByAsset",
			Sort:       bson.M{"timestamp": -1},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {
		r2 := make(map[string]interface{})
		r2["asset"] = item["asset"]
		r2["tokenid"] = item["tokenid"]
		r2["timestamp"] = item["timestamp"]
		r2["event"] = item["eventname"]
		r2["market"] = item["market"]
		r2["nonce"] = item["nonce"]
		r2["image"] = item["image"]
		r2["name"] = item["name"]
		r2["txid"] = item["txid"]
		properties := item["properties"].(primitive.A)[0].(map[string]interface{})
		if properties["properties"] != nil {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(properties["properties"].(string)), &data); err == nil {
				name, ok := data["name"]
				if ok {
					r2["name"] = name
				} else {
					r2["name"] = ""
				}
				thumbnail, ok1 := data["thumbnail"]
				if ok1 {
					//r1["image"] = thumbnail
					tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
					if err2 != nil {
						return err2
					}
					r2["image"] = string(tb[:])
				} else {
					r2["image"] = ""
				}
			} else {
				return err
			}
		} else {
			r2["image"] = ""
			r2["name"] = ""
		}

		eventname := item["eventname"].(string)
		if eventname == "Claim" || eventname == "CompleteOffer" {
			extendData := item["extendData"].(string)
			var data map[string]interface{}
			var auctionType = 0
			if err := json.Unmarshal([]byte(extendData), &data); err == nil {
				if item["eventname"] == "Claim" {
					auctionAsset := data["auctionAsset"]
					auctionAmount := data["bidAmount"]
					r2["auctionAsset"] = auctionAsset
					r2["auctionAmount"] = auctionAmount
					auctionType, err = strconv.Atoi(data["auctionType"].(string))
					if err != nil {
						return err
					}
					r2["from"] = item["market"]
					r2["to"] = item["user"]

					if auctionType == 1 {
						r2["state"] = NFTevent.Sale.Val()

					} else if auctionType == 2 {
						r2["state"] = NFTevent.Bid.Val()
					}

				} else if item["eventname"] == "CompleteOffer" {
					offerAsset := data["offerAsset"]
					offerAmount := data["offerAmount"]
					r2["auctionAsset"] = offerAsset
					r2["auctionAmount"] = offerAmount
					r2["from"] = data["offerer"]
					r2["to"] = item["user"]
					r2["state"] = NFTevent.CompleteOffer.Val()
				}

			} else {
				return err
			}
			result = append(result, r2)

		} else if eventname == "Auction" {
			extendData := item["extendData"].(string)
			var data map[string]interface{}
			var auctionType = 0
			if err := json.Unmarshal([]byte(extendData), &data); err == nil {

				deadline := data["deadline"].(string)
				ddl, err := strconv.ParseInt(deadline, 10, 64)
				if err != nil {
					return err
				}

				auctionType, err = strconv.Atoi(data["auctionType"].(string))
				if err != nil {
					return err
				}
				r2["to"] = item["market"]
				r2["from"] = item["user"]

				if auctionType == 1 {
					auctionAsset := data["auctionAsset"]
					auctionAmount := data["auctionAmount"]
					r2["auctionAsset"] = auctionAsset
					r2["auctionAmount"] = auctionAmount
					// 只展示直买直卖
					if ddl > currentTime {
						r2["state"] = NFTevent.List.Val()
					} else {
						r2["state"] = NFTevent.List_Expired.Val()
					}
					result = append(result, r2)
				}
			} else {
				return err
			}

		} else if eventname == "Offer" {
			extendData := item["extendData"].(string)
			var data map[string]interface{}

			if err := json.Unmarshal([]byte(extendData), &data); err == nil {
				offerAsset := data["offerAsset"]
				offerAmount := data["offerAmount"]
				originOwner := data["originOwner"]
				r2["from"] = originOwner
				r2["to"] = item["user"]
				r2["auctionAsset"] = offerAsset
				r2["auctionAmount"] = offerAmount
				deadline := data["deadline"].(string)
				ddl, err := strconv.ParseInt(deadline, 10, 64)
				if err != nil {
					return err
				}
				if currentTime < ddl {
					r2["state"] = NFTevent.Offers
				} else {
					r2["state"] = NFTevent.Offer_Expired
				}
				result = append(result, r2)
			}
		}
	}

	num, err := strconv.ParseInt(strconv.Itoa(len(result)), 10, 64)
	if err != nil {
		return err
	}

	if args.Limit == 0 {
		args.Limit = int64(math.Inf(1))
	}

	pagedNFT := make([]map[string]interface{}, 0)
	for i, item := range result {
		if int64(i) < args.Skip {
			continue
		} else if int64(i) > args.Skip+args.Limit-1 {
			continue
		} else {
			pagedNFT = append(pagedNFT, item)
		}
	}

	r2, err := me.FilterArrayAndAppendCount(pagedNFT, num, args.Filter)
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
