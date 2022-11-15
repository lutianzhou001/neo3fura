package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"
)

func (me *T) GetHighestOfferByNFT(args struct {
	Asset      h160.T
	TokenId    strval.T
	MarketHash h160.T
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
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
		Index:      "GetHighestOfferByNFT",
		Sort:       bson.M{"nonce": -1},
		Filter:     bson.M{"asset": args.Asset.Val(), "tokenid": args.TokenId.Val(), "market": args.MarketHash, "eventname": bson.M{"$in": []interface{}{"Offer", "OfferCollection"}}},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)

	if err != nil {
		return err
	}
	result := make([]map[string]interface{}, 0)

	for _, item := range r1 {
		//获取有效期内的offer
		offer := make(map[string]interface{})
		extendData := item["extendData"].(string)
		var dat map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &dat); err1 == nil {
			dl := dat["deadline"].(string)
			deadline, err := strconv.ParseInt(dl, 10, 64)
			if err != nil {
				return err
			}
			if currentTime < deadline {
				//查看offer 当前状态
				offer_nonce := item["nonce"]
				eventname := item["eventname"]
				var f bson.M
				if eventname == "Offer" {
					f = bson.M{
						"nonce":   offer_nonce,
						"asset":   item["asset"],
						"tokenid": item["tokenid"],
						"$or": []interface{}{
							bson.M{"eventname": "CompleteOffer"},
							bson.M{"eventname": "CancelOffer"},
						},
					}
				} else if eventname == "OfferCollection" {
					f = bson.M{
						"nonce":   offer_nonce,
						"asset":   item["asset"],
						"tokenid": item["tokenid"],
						"$or": []interface{}{
							bson.M{"eventname": "CompleteOfferCollection"},
							bson.M{"eventname": "CancelOfferCollection"},
						},
					}
				}

				offerstate, _ := me.Client.QueryOne(struct {
					Collection string
					Index      string
					Sort       bson.M
					Filter     bson.M
					Query      []string
				}{
					Collection: "MarketNotification",
					Index:      "getOfferSate",
					Sort:       bson.M{},
					Filter:     f,
					Query:      []string{},
				}, ret)

				if len(offerstate) > 0 {
					if eventname == "Offer" {
						continue
					} else if eventname == "OfferCollection" {
						count := dat["count"].(string)
						if count == "0" {
							continue
						}
					}

				} else {
					offer["user"] = item["user"]
					offer["asset"] = item["asset"]
					offer["tokenid"] = item["tokenid"]
					offer["originOwner"] = dat["originOwner"]
					offer["offerAsset"] = dat["offerAsset"]
					offerAmount, err := strconv.ParseInt(dat["offerAmount"].(string), 10, 64)
					if err != nil {
						return err
					}
					offer["offerAmount"] = offerAmount
					offer["deadline"] = deadline
					offer["nonce"] = item["nonce"]
					offer["eventname"] = item["eventname"]
					result = append(result, offer)

				}
			}
		} else {
			return err1
		}

	}

	//	排序
	result = mapsort.MapSort(result, "offerAmount")

	offerCount := len(result)

	skip := 5

	page := offerCount/skip + 1
	if offerCount%skip == 0 {
		page = offerCount / skip
	}
	hightestOffer := map[string]interface{}{}
	flag := true
	for i := 0; i < page; i++ { //skip
		addressArr := make([]map[string]interface{}, 0)
		var addressList []string
		if i < page-1 {
			for j := i * skip; j < (i+1)*skip; j++ {
				addressArr = append(addressArr, result[j])
				addressList = append(addressList, result[j]["user"].(string))
			}
		} else {
			for j := i * skip; j < offerCount; j++ {
				addressArr = append(addressArr, result[j])
				addressList = append(addressList, result[j]["user"].(string))
			}
		}

		re, err := GetSavings(args.MarketHash, "getUserSavingsAmount", addressList, result[0]["offerAsset"].(string))
		if err != nil {
			return err
		}

		for m := 0; m < len(re); m++ {

			offerAmount := result[i*skip+m]["offerAmount"].(int64)
			oa := new(big.Int).SetInt64(offerAmount)

			if !(re[m].Cmp(oa) == -1) {
				hightestOffer = result[i*skip+m]
				hightestOffer["guarantee"] = re[m]
				flag = false
				break
			}
		}

		if !flag {
			break
		}
	}

	if args.Raw != nil {
		*args.Raw = hightestOffer
	}
	r2, err := me.Filter(hightestOffer, args.Filter)
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
