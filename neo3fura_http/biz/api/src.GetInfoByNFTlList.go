package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"neo3fura_http/lib/type/Contract"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/lib/utils"
	"neo3fura_http/var/stderr"
	"os"
	"strconv"
	"time"
)

func (me *T) GetInfoByNFTList(args struct {
	NFT []struct {
		Asset   h160.T
		TokenId strval.T
	}
	Filter map[string]interface{}
	Raw    *map[string]interface{}
}, ret *json.RawMessage) error {

	nftlist := make([]map[string]interface{}, 0)

	for _, item := range args.NFT {
		if item.Asset.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		it := make(map[string]interface{})
		it["asset"] = item.Asset.Val()
		it["tokenid"] = item.TokenId.Val()
		nftlist = append(nftlist, it)
	}
	list := utils.GroupByAsset(nftlist)

	or := []interface{}{}
	for k, v := range list {
		f := bson.M{"asset": k, "tokenid": bson.M{"$in": v}}
		or = append(or, f)
	}
	filter := bson.M{"amount": 1, "$or": or}
	//获取上架以及Owner信息
	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Market",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": filter},
				bson.M{"$lookup": bson.M{
					"from": "MarketNotification",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{ //
						bson.M{"$match": bson.M{"eventname": bson.M{"$in": []interface{}{"CompleteOfferCollection", "CompleteOffer", "Claim"}}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						}}}},
						bson.M{"$sort": bson.M{"nonce": 1}},
						bson.M{"$group": bson.M{"_id": bson.M{"eventname": "$eventname"}, "eventname": bson.M{"$last": "$eventname"}, "market": bson.M{"$last": "$market"}, "timestamp": bson.M{"$last": "$timestamp"}, "extendData": bson.M{"$last": "$extendData"}}},
						bson.M{"$project": bson.M{"eventname": 1, "market": 1, "extendData": 1, "timestamp": 1}},
					},
					"as": "eventlist"}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	currentTime := time.Now().UnixNano() / 1e6

	for _, item := range r1 {
		//NFT状态   上架 （售卖中  成交未领取）  未上架
		ddl := item["deadline"].(int64)

		bidAmount := item["bidAmount"].(primitive.Decimal128).String()
		if item["market"] != item["owner"] || ddl < currentTime {
			item["state"] = "notlist"
		} else {
			item["state"] = "list"
		}
		if (item["market"] == item["owner"] && ddl > currentTime) || (item["market"] == item["owner"] && ddl < currentTime && bidAmount == "0") { //上架过期（无人出价）   上架未过期
			item["owner"] = item["auctor"]
		}
		if item["market"] == item["owner"] && ddl < currentTime && bidAmount != "0" { // 未领取
			item["owner"] = item["bidder"]
		}
		item["buyNowAsset"] = ""
		item["buyNowAmount"] = "0"
		item["lastSoldAsset"] = ""
		item["lastSoldAmount"] = "0"
		item["currentBidAsset"] = ""
		item["currentBidAmount"] = "0"
		item["offerAsset"] = ""
		item["offerAmount"] = "0"

		auctionType := item["auctionType"].(int32)
		if ddl > currentTime {
			if auctionType == 1 {
				item["buyNowAsset"] = item["auctionAsset"]
				item["buyNowAmount"] = item["auctionAmount"].(primitive.Decimal128).String()
			} else if auctionType == 2 {
				if bidAmount != "0" {
					item["currentBidAsset"] = item["auctionAsset"]
					item["currentBidAmount"] = item["bidAmount"].(primitive.Decimal128).String()
				} else {
					item["currentBidAsset"] = item["auctionAsset"]
					item["currentBidAmount"] = item["auctionAmount"].(primitive.Decimal128).String()
				}
			}
		} else {
			if auctionType == 2 && bidAmount != "0" {
				item["lastSoldAsset"] = item["auctionAsset"]
				item["lastSoldAmount"] = item["bidAmount"].(primitive.Decimal128).String()
			}
		}

		if item["eventlist"] != nil {
			eventlist := item["eventlist"].(primitive.A)
			for _, it := range eventlist {
				eventItem := it.(map[string]interface{})
				eventname := eventItem["eventname"]
				extendData := eventItem["extendData"]
				//market := eventItem["market"].(string)

				var finishTime int64
				data := make(map[string]interface{})
				if err := json.Unmarshal([]byte(extendData.(string)), &data); err == nil {
					if eventname == "Claim" {
						finishTime = eventItem["timestamp"].(int64)
						item["lastSoldAsset"] = data["auctionAsset"]
						item["lastSoldAmount"] = data["bidAmount"]
					} else if eventname == "CompleteOffer" || eventname == "CompleteOfferCollection" {
						time := eventItem["timestamp"].(int64)
						if time > finishTime {
							finishTime = time
							item["lastSoldAsset"] = data["offerAsset"]
							if err != nil {
								return err
							}
							item["lastSoldAmount"] = data["offerAmount"]

						}
					}
				}

			}
		}

		//获取Owner 设置的nns信息
		owner := item["owner"].(string)
		nns := ""
		if owner != "" {
			nns, err = GetNNSByAddress(owner)
			if err != nil {
				return err
			}
		}

		item["nns"] = nns
		delete(item, "eventlist")

	}
	//
	rt := os.ExpandEnv("${RUNTIME}")
	var market string
	if rt == "staging" {
		market = Contract.Main_SecondaryMarket.Val()
	} else if rt == "test2" {
		market = Contract.Test_SecondaryMarket.Val()
	} else {
		market = Contract.Test_SecondaryMarket.Val()
	}

	raw := make(map[string]interface{})
	err = me.GetHighestOfferByNFTList(struct {
		NFT []struct {
			Asset   h160.T
			TokenId strval.T
		}
		MarketHash h160.T
		Limit      int64
		Skip       int64
		Filter     map[string]interface{}
		Raw        *map[string]interface{}
	}{
		NFT:        args.NFT,
		MarketHash: h160.T(market),
		Raw:        &raw}, ret)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	for _, item := range r1 {
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		key := asset + tokenid
		if raw[key] != nil {
			value := raw[key].(map[string]interface{})
			if len(value) > 0 {
				deadline := value["deadline"].(int64)
				if deadline > currentTime {
					offerAmount := value["offerAmount"].(int64)
					guarantee := value["guarantee"].(*big.Int)
					amount := big.NewInt(offerAmount)
					if guarantee.Cmp(amount) != -1 {
						item["offerAmount"] = amount.String()
						item["offerAsset"] = value["offerAsset"]
					} else {
						item["offerAmount"] = "0"
						item["offerAsset"] = ""
					}
				} else {
					item["offerAmount"] = "0"
					item["offerAsset"] = ""
				}

			} else {
				item["offerAmount"] = "0"
				item["offerAsset"] = ""
			}
		} else {
			item["offerAmount"] = "0"
			item["offerAsset"] = ""
		}

		var order string
		if item["buyNowAmount"] != "0" {
			order = item["buyNowAmount"].(string)
		} else if item["currentBidAmount"] != "0" {
			order = item["currentBidAmount"].(string)
		} else if item["lastSoldAmount"] != "0" {
			order = item["lastSoldAmount"].(string)
		} else {
			order = item["offerAmount"].(string)
		}
		number, err := strconv.ParseInt(order, 10, 64)
		if err != nil {
			return err
		}
		item["order"] = number
		result[key] = item

	}

	r3, err := me.Filter(result, args.Filter)
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
