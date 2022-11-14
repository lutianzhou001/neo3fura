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
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market, "eventname": bson.M{"$in": []interface{}{"Claim", "CompleteOffer", "CompleteOfferCollection"}}}},
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
			bson.M{"$match": bson.M{"asset": args.Asset, "market": args.Market, "eventname": bson.M{"$in": []interface{}{"Offer", "OfferCollection"}}}},
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
		r2["image"] = ""
		r2["name"] = item["name"]
		r2["txid"] = item["txid"]
		properties := item["properties"].(primitive.A)[0].(map[string]interface{})
		if properties["properties"] != nil {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(properties["properties"].(string)), &data); err == nil {

				image, ok := data["image"]
				if ok {
					r2["image"] = ImagUrl(item["asset"].(string), image.(string), "images")
				} else {
					r2["image"] = ""
				}

				thumbnail, ok1 := data["thumbnail"]
				if ok1 {
					//r1["image"] = thumbnail
					tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
					if err2 != nil {
						return err2
					}
					r2["thumbnail"] = ImagUrl(item["asset"].(string), string(tb[:]), "thumbnail")
				} else {
					if r2["thumbnail"] == nil {
						if r2["image"] != nil && r2["image"] != "" {
							if image == nil {
								r2["thumbnail"] = item["image"]
							} else {
								r2["thumbnail"] = ImagUrl(item["asset"].(string), image.(string), "thumbnail")
							}
						}
					}
				}

				tokenuri, ok := data["tokenURI"]
				if ok {
					ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), item["asset"].(string), item["tokenid"].(string))
					if err != nil {
						return err
					}
					for key, value := range ppjson {
						r2[key] = value
						if key == "image" {
							img := value.(string)
							r2["thumbnail"] = ImagUrl(item["asset"].(string), img, "thumbnail")
							r2["image"] = ImagUrl(item["asset"].(string), img, "images")
						}
					}
				}
				if r2["name"] == "" || r2["name"] == nil {
					name, ok := data["name"]
					if ok {
						r2["name"] = name
					}
				}

				if r2["name"].(string) == "Video" {
					r2["video"] = r2["image"]
					delete(r2, "image")
				}

			} else {
				return err
			}
		} else {

			assetInfo, err := me.Client.QueryOne(struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
				Query      []string
			}{Collection: "Asset",
				Index:  "GetAssetInfo",
				Sort:   bson.M{},
				Filter: bson.M{"hash": item["asset"]},
				Query:  []string{},
			}, ret)
			if err != nil {
				return err
			}
			r2["image"] = ""
			r2["name"] = assetInfo["tokenname"]
			r2["thumbnail"] = ""
			r2["video"] = ""

		}

		eventname := item["eventname"].(string)
		if eventname == "Claim" || eventname == "CompleteOffer" || eventname == "CompleteOfferCollection" {
			extendData := item["extendData"].(string)
			nonce := item["nonce"].(int64)
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
					auctionInfo, err := me.GetAuction(nonce, args.Market.Val())
					if err != nil {
						return err
					}
					var auctor string
					if auctionInfo != nil {
						auctor = auctionInfo["user"].(string)
					}

					r2["from"] = auctor
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

		} else if eventname == "Offer" || eventname == "OfferCollection" {
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

//获取上架的数据 （）
func (me T) GetAuction(nonce int64, market string) (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	r0, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "MarketNotification", Index: "GetAuction", Filter: bson.M{"eventname": "Auction", "market": market, "nonce": nonce}}, ret)
	if err != nil {
		return nil, err
	}

	return r0, nil
}
