package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTByWords(args struct {
	SecondaryMarket h160.T //
	PrimaryMarket   h160.T
	Words           strval.T
	Limit           int64
	Skip            int64
	Filter          map[string]interface{}
	Raw             *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.Words == "" {
		return stderr.ErrInvalidArgs
	}

	if len(args.PrimaryMarket) > 0 && args.PrimaryMarket != "" {
		if args.PrimaryMarket.Valid() == false {
			return stderr.ErrInvalidArgs
		}
	}

	if args.SecondaryMarket.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	var white []interface{}
	if len(args.SecondaryMarket) > 0 && args.SecondaryMarket != "" {
		//白名单
		raw1 := make(map[string]interface{})
		err1 := me.GetMarketWhiteList(struct {
			MarketHash h160.T
			Filter     map[string]interface{}
			Raw        *map[string]interface{}
		}{MarketHash: args.SecondaryMarket, Raw: &raw1}, ret) //nonce 分组，并按时间排序
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
			white = wl
		} else {
			return stderr.ErrWhiteList
		}
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
			Collection: "Market",
			Index:      "GetNFTByWords",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": bson.M{"$in": white}}},
				bson.M{"$match": bson.M{"market": bson.M{"$ne": args.PrimaryMarket.Val()}}},
				bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
				bson.M{"$lookup": bson.M{
					"from": "Nep11Properties",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$or": []interface{}{
							bson.M{"properties": bson.M{"$regex": "name\":\"" + args.Words, "$options": "$i"}},
							bson.M{"properties": bson.M{"$regex": "name\": \"" + args.Words, "$options": "$i"}},
						}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						bson.M{"$project": bson.M{"id": 1, "tokenid": 1, "asset": 1, "properties": 1}}},
					"as": "properties"},
				},
				bson.M{"$match": bson.M{"properties": bson.M{"$ne": []interface{}{}}}}, //过滤空的集合
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {

		//判断nft 状态
		a := item["amount"].(primitive.Decimal128).String()
		amount, err1 := strconv.Atoi(a)
		if err1 != nil {
			return err
		}

		bidAmount, _, err2 := item["bidAmount"].(primitive.Decimal128).BigInt()
		bidAmountFlag := bidAmount.Cmp(big.NewInt(0))
		if err2 != nil {
			return err
		}

		deadline, _ := item["deadline"].(int64)
		auctionType, _ := item["auctionType"].(int32)

		if amount > 0 && auctionType == 2 && item["owner"] == item["market"] && deadline > currentTime {
			item["state"] = NFTstate.Auction.Val()
		} else if amount > 0 && auctionType == 1 && item["owner"] == item["market"] && deadline > currentTime {
			item["state"] = NFTstate.Sale.Val()
		} else if amount > 0 && item["owner"] != item["market"] {
			item["state"] = NFTstate.NotListed.Val()
		} else if amount > 0 && bidAmountFlag == 1 && deadline < currentTime && item["owner"] == item["market"] {
			item["state"] = NFTstate.Unclaimed.Val()
		} else if amount > 0 && deadline < currentTime && bidAmountFlag == 0 && item["owner"] == item["market"] {
			item["state"] = NFTstate.Expired.Val()
		} else {
			item["state"] = ""
		}

		//获取nft 属性
		nftproperties := item["properties"]
		if nftproperties != nil && nftproperties != "" {
			pp := nftproperties.(primitive.A)
			if len(pp) > 0 {
				it := pp[0].(map[string]interface{})
				extendData := it["properties"].(string)
				asset := it["asset"].(string)
				tokenid := it["tokenid"].(string)
				if extendData != "" {
					properties := make(map[string]interface{})
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
						image, ok := data["image"]
						if ok {
							properties["image"] = image
							//item["image"] = image
							item["image"] = ImagUrl(asset, image.(string), "images")
						} else {
							item["image"] = ""
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
									item["thumbnail"] = ImagUrl(asset, img, "thumbnail")
									item["image"] = ImagUrl(asset, img, "images")
								}

							}
						}
						if item["name"] == "" || item["name"] == nil {
							name, ok1 := data["name"]
							if ok1 {
								item["name"] = name
							}
						}

						strArray := strings.Split(item["name"].(string), "#")
						if len(strArray) >= 2 {
							number := strArray[1]
							n, err22 := strconv.ParseInt(number, 10, 64)
							if err22 != nil {
								item["number"] = int64(-1)
							}
							item["number"] = n
							properties["number"] = n
						} else {
							item["number"] = int64(-1)
						}

						series, ok2 := data["series"]
						if ok2 {
							properties["series"] = series
						}
						supply, ok3 := data["supply"]
						if ok3 {
							properties["supply"] = supply
						}
						number, ok4 := data["number"]
						if ok4 {
							n, err22 := strconv.ParseInt(number.(string), 10, 64)
							if err22 != nil {
								item["number"] = int64(-1)
							}
							properties["number"] = n
							item["number"] = n
						}
						video, ok5 := data["video"]
						if ok5 {
							properties["video"] = video
						}
						thumbnail, ok6 := data["thumbnail"]
						if ok6 {
							//r1["image"] = thumbnail
							tb, err22 := base64.URLEncoding.DecodeString(thumbnail.(string))
							if err22 != nil {
								return err22
							}
							//item["image"] = string(tb[:])
							item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
						} else {
							if item["image"] != nil && item["image"] != "" {
								if image == nil {
									item["thumbnail"] = item["image"]
								} else {
									item["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
								}

							}

						}

					} else {
						return err
					}

					item["properties"] = properties
				} else {
					item["image"] = ""
					item["name"] = ""
					item["number"] = int64(-1)
					item["properties"] = ""
				}

			}
		} else {

		}

	}

	var r2, err2 = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Market",
			Index:      "GetNFTByWords",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": bson.M{"$in": white}}},
				bson.M{"$match": bson.M{"market": bson.M{"$ne": args.PrimaryMarket.Val()}}},
				bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
				bson.M{"$lookup": bson.M{
					"from": "Nep11Properties",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$or": []interface{}{
							bson.M{"properties": bson.M{"$regex": "name\":\"" + args.Words, "$options": "$i"}},
							bson.M{"properties": bson.M{"$regex": "name\": \"" + args.Words, "$options": "$i"}},
						}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						bson.M{"$project": bson.M{"id": 1, "tokenid": 1, "asset": 1, "properties": 1}}},
					"as": "properties"},
				},
				bson.M{"$match": bson.M{"properties": bson.M{"$ne": []interface{}{}}}}, //过滤空的集合

			},

			Query: []string{},
		}, ret)
	if err2 != nil {
		return err2
	}
	count := len(r2)

	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)

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
