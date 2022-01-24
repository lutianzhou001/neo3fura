package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTByWords(args struct {
	Words  strval.T
	Filter map[string]interface{}
	Raw    *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.Words == "" {
		return stderr.ErrInvalidArgs
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
			Collection: "Nep11Properties",
			Index:      "GetNFTByWords",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"properties": bson.M{"$regex": args.Words, "$options": "$i"}}},

				bson.M{"$lookup": bson.M{
					"from": "Market",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						bson.M{"$project": bson.M{"id": 1, "tokenid": 1, "asset": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": 1}}},
					"as": "market"},
				},
			},

			Query: []string{},
		}, ret)

	if err != nil {
		return err
	}

	//result := make([]map[string]interface{}, 0)
	if len(r1) > 0 {
	}

	for _, item := range r1 {
		m := item["market"]
		if m != nil {
			market := m.(primitive.A)[0].(map[string]interface{})

			a := market["amount"].(primitive.Decimal128).String()
			amount, err1 := strconv.Atoi(a)
			if err1 != nil {
				return err1
			}

			bidAmount := market["bidAmount"].(primitive.Decimal128).String()

			deadline, _ := market["deadline"].(int64)
			auctionType, _ := market["auctionType"].(int32)

			if amount > 0 && auctionType == 2 && market["owner"] == market["market"] && deadline > currentTime {
				market["state"] = NFTstate.Auction.Val()
			} else if amount > 0 && auctionType == 1 && market["owner"] == market["market"] && deadline > currentTime {
				market["state"] = NFTstate.Sale.Val()
			} else if amount > 0 && market["owner"] != market["market"] {
				market["state"] = NFTstate.NotListed.Val()
			} else if amount > 0 && bidAmount != "0" && deadline < currentTime && market["owner"] == market["market"] {
				market["state"] = NFTstate.Unclaimed.Val()
			} else if amount > 0 && deadline < currentTime && bidAmount == "0" && market["owner"] == market["market"] {
				market["state"] = NFTstate.Expired.Val()
			} else {
				market["state"] = ""
			}

			item["amount"] = market["amount"]
			item["auctionAmount"] = market["auctionAmount"]
			item["auctionAsset"] = market["auctionAsset"]
			item["auctionType"] = market["auctionType"]
			item["auctor"] = market["auctor"]
			item["bidAmount"] = market["bidAmount"]
			item["bidder"] = market["bidder"]
			item["deadline"] = market["deadline"]

			item["market"] = market["market"]
			item["owner"] = market["owner"]
			item["state"] = market["state"]
			item["timestamp"] = market["timestamp"]
		}

		delete(item, "market")

		extendData := item["properties"].(string)
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
				name, ok1 := data["name"]
				if ok1 {
					item["name"] = name
					strArray := strings.Split(name.(string), "#")
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

				} else {
					item["name"] = ""
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
					item["image"] = string(tb[:])
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

	count := len(r1)

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
