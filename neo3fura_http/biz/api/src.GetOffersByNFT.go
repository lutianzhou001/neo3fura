package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
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
	currentTime := time.Now().UnixNano() / 1e6
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash.Val(),
			"$or": []interface{}{
				bson.M{"asset": args.Asset.Val(), "tokenid": args.TokenId.Val(), "eventname": "Offer"},
				bson.M{"asset": args.Asset.Val(), "eventname": "OfferCollection"},
			},
		}},

		bson.M{"$lookup": bson.M{
			"from": "Nep11Properties",
			"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$tokenid", args.TokenId}},
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

	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		//查看offer 当前状态
		offer_nonce := item["nonce"]
		eventname := item["eventname"].(string)

		var filter bson.M
		if eventname == "Offer" {
			filter = bson.M{
				"nonce":   offer_nonce,
				"asset":   item["asset"],
				"tokenid": item["tokenid"],
				//"eventname":"CancelOffer",
				"$or": []interface{}{
					bson.M{"eventname": "CompleteOffer"},
					bson.M{"eventname": "CancelOffer"},
				},
			}
		} else if eventname == "OfferCollection" {
			filter = bson.M{
				"nonce": offer_nonce,
				"asset": item["asset"],
				//"tokenid": item["tokenid"],
				//"eventname":"CancelOffer",
				"$or": []interface{}{
					bson.M{"eventname": "CompleteOfferCollection"},
					bson.M{"eventname": "CancelOfferCollection"},
				},
			}
		}
		offer, _ := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "getOfferSate",
			Sort:       bson.M{"timestamp": -1},
			Filter:     filter,
			Query:      []string{},
		}, ret)

		if len(offer) > 0 {
			offereventname := offer["eventname"].(string)
			if offereventname == "CancelOfferCollection" || offereventname == "CancelOffer" || offereventname == "CompleteOffer" {
				continue
			} else {
				extendData := item["extendData"].(string)
				var data map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
					count := data["count"].(string)
					if count == "0" {
						continue
					}

				}
			}

		}

		extendData := item["extendData"].(string)
		var dat map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {

			item["offerAsset"] = dat["offerAsset"]
			item["offerAmount"] = dat["offerAmount"]
			item["deadline"] = dat["deadline"]
			offerCollectionDDL := dat["deadline"].(string)

			DDL, err := strconv.ParseInt(offerCollectionDDL, 10, 64)
			if err != nil {
				return err
			}
			if DDL < currentTime {
				continue
			}

			if eventname == "Offer" {
				item["originOwner"] = dat["originOwner"]
			} else if eventname == "OfferCollection" {
				//item["originOwner"] = dat["originOwner"]

				item["tokenid"] = args.TokenId

				marketInfo, err := me.Client.QueryOne(struct {
					Collection string
					Index      string
					Sort       bson.M
					Filter     bson.M
					Query      []string
				}{Collection: "Market",
					Index:  "GetMarketInfo",
					Sort:   bson.M{},
					Filter: bson.M{"amount": bson.M{"$gt": 0}, "asset": item["asset"], "tokenid": args.TokenId},
					Query:  []string{},
				}, ret)

				if err != nil {
					return stderr.ErrGetNFTInfo
				}
				market := marketInfo["market"]
				owner := marketInfo["owner"]
				bidder := marketInfo["bidder"]
				bidAmount := marketInfo["bidAmount"].(primitive.Decimal128).String()
				deadline := marketInfo["deadline"].(int64)

				if market == owner && deadline > currentTime { // 上架未过期
					item["originOwner"] = marketInfo["auctor"]
				} else if market == owner && deadline < currentTime { //上架过期
					if bidAmount == "0" {
						item["originOwner"] = marketInfo["auctor"]
					} else {
						item["originOwner"] = bidder
					}
				} else { //未上架
					item["originOwner"] = owner
				}

			}

		} else {
			return err1
		}

		nftproperties := item["properties"]
		if nftproperties != nil && nftproperties != "" {
			pp := nftproperties.(primitive.A)
			if len(pp) > 0 {
				it := pp[0].(map[string]interface{})
				extendData1 := it["properties"].(string)
				asset := item["asset"].(string)
				var tokenid string
				switch item["tokenid"].(type) {
				case strval.T:
					tokenid = item["tokenid"].(strval.T).Val()
				default:
					tokenid = item["tokenid"].(string)
				}
				if extendData1 != "" {
					properties := make(map[string]interface{})
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData1), &data); err1 == nil {
						image, ok := data["image"]
						if ok {
							properties["image"] = image
							//item["image"] = image
							item["image"] = ImagUrl(asset, image.(string), "images")
						} else {
							item["image"] = ""
						}

						thumbnail, ok := data["thumbnail"]
						if ok {
							tb, err22 := base64.URLEncoding.DecodeString(thumbnail.(string))
							if err22 != nil {
								return err22
							}
							ss := string(tb[:])
							if ss == "" {
								item["thumbnail"] = ImagUrl(item["asset"].(string), item["image"].(string), "thumbnail")
							} else {
								item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
							}

						} else {
							if item["thumbnail"] == nil {
								if image != nil && image != "" {
									if image == nil {
										item["thumbnail"] = item["image"]
									} else {
										item["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
									}
								}
							}
						}

						tokenuri, ok := data["tokenURI"]
						if ok {
							ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)
							if err != nil {
								return err
							}
							for key, value := range ppjson {
								item[key] = value
								properties[key] = value
								if key == "image" {
									img := value.(string)
									tb := ImagUrl(asset, img, "thumbnail")
									flag := strings.HasSuffix(tb, ".mp4")
									if flag {
										tb = strings.Replace(tb, ".mp4", "mp4", -1)
									}
									item["thumbnail"] = tb
									item["image"] = ImagUrl(asset, img, "images")
								}
							}
						}

						if item["name"] != nil && item["name"].(string) == "Nuanced Floral Symphony" {
							item["video"] = item["image"]
							delete(item, "image")
						}
						if item["name"] != nil && item["name"].(string) == "Virtual Visions #1" {
							item["video"] = item["image"]
							delete(item, "image")
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

		result = append(result, item)
	}
	count := int64(len(result))
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
