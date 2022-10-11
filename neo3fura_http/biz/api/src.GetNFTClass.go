package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
	"time"
)

func (me *T) GetNFTClass(args struct {
	MarketHash h160.T
	AssetHash  h160.T
	SubClass   [][]strval.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	currentTime := time.Now().UnixNano() / 1e6
	length := 0
	cond := bson.M{}
	var tokenidClassList [][]interface{}
	if len(args.SubClass) > 0 {
		for _, i := range args.SubClass {
			if len(i) != 2 {
				return stderr.ErrInvalidArgs
			} else {
				_tokenid, _ := base64.StdEncoding.DecodeString(i[1].Val())
				tokenid := string(_tokenid)

				category := strings.Split(tokenid, "#")
				str := category[0]
				var number = 1
				num := category[1]
				number, _ = strconv.Atoi(num)
				str = str + "#"

				//if len(tokenid) == 17 {
				//	series := tokenid[13:14]
				//	num := tokenid[15:17]
				//	numberstr = append(numberstr, num)
				//	number, _ = strconv.Atoi(num)
				//	str = "MetaPanacea #" + series + "-"
				//} else if len(tokenid) == 18 {
				//	series := tokenid[13:15]
				//	num := tokenid[16:18]
				//	numberstr = append(numberstr, num)
				//	number, _ = strconv.Atoi(num)
				//	str = "MetaPanacea #" + series + "-"
				//
				//}

				var tokenidList []interface{}
				for j := 1; j <= number; j++ {
					var s string
					if j < 10 {
						s = str + "0" + strconv.Itoa(j)
					} else {
						s = str + strconv.Itoa(j)
					}

					token := base64.StdEncoding.EncodeToString([]byte(s))
					tokenidList = append(tokenidList, token)
				}
				tokenidClassList = append(tokenidClassList, tokenidList)
			}
		}

		//fmt.Println("tokenList: ",tokenidClassList)
		for _, i := range tokenidClassList {
			//classSort[]
			b := bson.M{}
			//a := bson.M{"$and": []interface{}{bson.M{"$gte": []interface{}{"$tokenid", i[0].Val()}}, bson.M{"$lte": []interface{}{"$tokenid", i[1].Val()}}}}
			a := bson.M{"$and": []interface{}{bson.M{"$in": []interface{}{"$tokenid", i}}}}
			//a :=bson.M{"tokenid":bson.M{"$in":i}}
			if length == 0 {
				b = bson.M{"if": a, "then": length, "else": length - 1}
			} else {
				b = bson.M{"if": a, "then": length, "else": cond}
			}
			length++

			cond = bson.M{"$cond": b}
		}
	} else {
		return stderr.ErrInvalidArgs
	}
	result := make([]map[string]interface{}, 0)

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash}},
		bson.M{"$match": bson.M{"asset": args.AssetHash}},
		bson.M{"$match": bson.M{"eventname": "Auction"}},
		bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "tokenidArr": bson.M{"$push": "$$ROOT"}, "extendData": bson.M{"$last": "$extendData"}}},
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
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	//  获取claimed 的值
	pipeline2 := []bson.M{
		bson.M{"$match": bson.M{"market": args.MarketHash}},
		bson.M{"$match": bson.M{"asset": args.AssetHash}},
		bson.M{"$match": bson.M{"eventname": "Claim"}},
		bson.M{"$project": bson.M{"class": cond, "asset": 1, "tokenid": 1, "extendData": 1}},
		bson.M{"$group": bson.M{"_id": "$class", "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "claimedInfo": bson.M{"$push": "$$ROOT"}, "extendData": bson.M{"$last": "$extendData"}, "claimed": bson.M{"$sum": 1}}},
	}

	r2, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline2,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {
		deadline := int64(0)
		if item["_id"].(int32) != -1 {
			asset := item["asset"].(string)
			tokenid := item["tokenid"].(string)
			extendData := item["extendData"].(string)
			//
			tokenidList := make(map[string]interface{})
			if item["tokenidArr"] != nil {
				tokenidArr := item["tokenidArr"].(primitive.A)

				for _, it := range tokenidArr {
					nft := it.(map[string]interface{})
					tokenidList[nft["tokenid"].(string)] = 1
				}
			}

			var dat map[string]interface{}
			if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {
				item["deadline"] = dat["deadline"]
				ddl := dat["deadline"].(string)
				deadline, err = strconv.ParseInt(ddl, 10, 64)
				if err != nil {
					return err
				}
				auctionAsset := dat["auctionAsset"]
				auctionAmount := dat["auctionAmount"]
				item["price"] = auctionAmount
				item["sellAsset"] = auctionAsset
			} else {
				return err1
			}

			var raw3 map[string]interface{}
			err1 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
			if err1 != nil {
				item["image"] = ""
				item["name"] = ""
				item["number"] = int64(-1)
				item["properties"] = ""
			}
			item["image"] = raw3["image"]
			name := strings.Split(raw3["name"].(string), "#")
			item["name"] = strings.TrimSpace(name[0])
			item["number"] = raw3["number"]
			item["properties"] = raw3["properties"]
			p := raw3["properties"].(map[string]interface{})
			item["thumbnail"] = raw3["thumbnail"]

			supply, err2 := base64.URLEncoding.DecodeString(p["supply"].(string))

			if err2 != nil {
				return err2
			}
			//获取claimed
			if deadline > currentTime {
				if len(r2) > 0 {
					for _, item1 := range r2 {
						if item["_id"] == item1["_id"] {
							item["claimed"] = item1["claimed"]
							break
						} else {
							item["claimed"] = 0
						}
					}
				} else {
					item["claimed"] = 0
				}
			} else {
				claimed, err3 := strconv.Atoi(string(supply))
				if err3 != nil {
					return err3
				}
				item["claimed"] = claimed
			}

			//delete(item, "_id")
			delete(item, "extendData")
			delete(item, "tokenid")
			delete(item, "tokenidArr")
			result = append(result, item)
		}
	}

	mapsort.MapSort5(result, "_id")

	count := len(result)

	r3, err := me.FilterAggragateAndAppendCount(result, count, args.Filter)

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
