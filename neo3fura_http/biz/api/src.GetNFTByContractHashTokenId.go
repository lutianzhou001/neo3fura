package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTByContractHashTokenId(args struct {
	ContractHash h160.T     //  asset
	TokenIds     []strval.T // tokenId
	Filter       map[string]interface{}
	Raw          *[]map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var tokenIds []strval.T
	if len(args.TokenIds) == 0 {
		var raw1 []map[string]interface{}
		err := me.GetAssetHoldersListByContractHash(struct {
			ContractHash h160.T
			Limit        int64
			Skip         int64
			Filter       map[string]interface{}
			Raw          *[]map[string]interface{}
		}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
		if err != nil {
			return err
		}
		for _, raw := range raw1 {
			if len(raw["tokenid"].(string)) == 0 {
				continue
			} else {
				tokenIds = append(tokenIds, strval.T(raw["tokenid"].(string)))
			}
		}

	} else {
		tokenIds = args.TokenIds
	}

	orArray := []interface{}{}
	for _, tokenId := range tokenIds {

		if len(tokenId) <= 0 {
			return stderr.ErrInvalidArgs
		}
		a := bson.M{"tokenid": tokenId.Val()}
		orArray = append(orArray, a)
	}

	//rr1, count, err := me.Client.QueryAll(struct {
	//	Collection string
	//	Index      string
	//	Sort       bson.M
	//	Filter     bson.M
	//	Query      []string
	//	Limit      int64
	//	Skip       int64
	//}{
	//	Collection: "Market",
	//	Index:      "GetNFTByContractHashTokenId",
	//	Sort:       bson.M{},
	//	Filter:     bson.M{"asset": args.ContractHash.Val(), "$or": orArray, "amount": bson.M{"$gt": 0}},
	//	Query:      []string{},
	//}, ret)
	//if err != nil {
	//	return err
	//}

	rr1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Market",
		Index:      "GetNFTByContractHashTokenId",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"asset": args.ContractHash.Val(), "$or": orArray, "amount": bson.M{"$gt": 0}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					//	bson.M{"$sort"}
				},
				"as": "properties"},
			},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range rr1 {
		a := item["amount"]
		//tokenId := item["tokenid"].(string)
		var a1 string
		switch a.(type) {
		case string:
			a1 = item["amount"].(string)
		case primitive.Decimal128:
			a1 = item["amount"].(primitive.Decimal128).String()
		}

		amount, err1 := strconv.Atoi(a1)
		if err1 != nil {
			return err
		}

		bidAmount := item["bidAmount"].(primitive.Decimal128).String()

		dl := item["deadline"]

		at := item["auctionType"]

		var deadline, auctionType int64
		switch dl.(type) {
		case float64:
			deadline = f2i(dl.(float64), 0)
		case int64:
			deadline = dl.(int64)
		case int32:
			deadline = int64(dl.(int32))

		}
		switch at.(type) {
		case float64:
			auctionType = f2i(at.(float64), 0)
		case int64:
			auctionType = at.(int64)
		case int32:
			auctionType = int64(at.(int32))
		}
		// NFT实际持有人
		if amount > 0 && auctionType == 2 && item["owner"] == item["market"] && deadline > currentTime {
			item["owner"] = item["auctor"]
			item["state"] = NFTstate.Auction.Val()
		} else if amount > 0 && auctionType == 1 && item["owner"] == item["market"] && deadline > currentTime {
			item["state"] = NFTstate.Sale.Val()
			item["owner"] = item["auctor"]
		} else if amount > 0 && item["owner"] != item["market"] {
			item["state"] = NFTstate.NotListed.Val()
		} else if amount > 0 && bidAmount != "0" && deadline < currentTime && item["owner"] == item["market"] {
			item["owner"] = item["bidder"]
			item["state"] = NFTstate.Unclaimed.Val()
		} else if amount > 0 && deadline < currentTime && bidAmount == "0" && item["owner"] == item["market"] {
			item["owner"] = item["auctor"]
			item["state"] = NFTstate.Expired.Val()
		} else {
			item["state"] = ""
		}

		if item["properties"] != nil {
			var pp map[string]interface{}
			switch item["properties"].(type) {
			case primitive.A:
				pp = item["properties"].(primitive.A)[0].(map[string]interface{})
			case map[string]interface{}:
				pp = item["properties"].(map[string]interface{})
			default:
				return stderr.ErrInvalidArgs
			}

			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)
			if pp["properties"] != nil && pp["properties"] != "" {
				extendData := pp["properties"].(string)

				properties := make(map[string]interface{})
				var data map[string]interface{}
				if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
					image, ok := data["image"]
					if ok {
						properties["image"] = image
						//r1["image"] = image
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

							if key == "name" || key == "number" {
								item[key] = value
							}
							if key != "name" {
								properties[key] = value
							}
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

					if item["name"] == "" || item["name"] == nil {
						name, ok1 := data["name"]
						if ok1 {
							item["name"] = name

						} else {
							item["name"] = ""
						}
					}
					strArray := strings.Split(item["name"].(string), "#")
					if len(strArray) >= 2 {
						number := strArray[1]
						n, err2 := strconv.ParseInt(number, 10, 64)
						if err2 != nil {
							item["number"] = int64(-1)
						}
						item["number"] = n
						properties["number"] = n
					} else {
						item["number"] = int64(-1)
					}

					series, ok2 := data["series"]
					if ok2 {
						decodeSeries, err2 := base64.URLEncoding.DecodeString(series.(string))
						if err2 != nil {
							properties["series"] = series
						}
						properties["series"] = string(decodeSeries)
					}
					supply, ok3 := data["supply"]
					if ok3 {
						decodeSupply, err2 := base64.URLEncoding.DecodeString(supply.(string))
						if err2 != nil {
							properties["supply"] = supply
						}
						properties["supply"] = string(decodeSupply)
					}
					number, ok4 := data["number"]
					if ok4 {
						n, err2 := strconv.ParseInt(number.(string), 10, 64)
						if err2 != nil {
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
						tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
						if err2 != nil {
							return err2
						}
						if len(tb) > 0 {
							item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
						} else {
							if image != nil && image != "" {
								item["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
							}
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

					if item["name"] != nil && item["name"].(string) == "Nuanced Floral Symphony" || item["name"].(string) == "Sunshine #1" {
						item["video"] = item["image"]
						delete(item, "image")
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

		owner := item["owner"].(string)
		nns := ""
		if owner != "" {
			nns, err = GetNNSByAddress(owner)
			if err != nil {
				return err
			}
		}
		item["nns"] = nns

	}

	r5, err := me.FilterArrayAndAppendCount(rr1, int64(len(rr1)), args.Filter)
	if err != nil {
		return err
	}

	if args.Raw != nil {
		*args.Raw = rr1
	}

	r, err := json.Marshal(r5)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}

func f2i(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}
