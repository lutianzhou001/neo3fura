package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v2"
	"math"
	"math/big"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/NFTevent"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (me *T) GetNFTRecordByAddress(args struct {
	Address         h160.T
	SecondaryMarket h160.T //
	PrimaryMarket   h160.T
	Limit           int64
	Skip            int64
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	f := bson.M{"user": args.Address.Val()}
	var wl []interface{}
	if len(args.SecondaryMarket) > 0 {
		if args.SecondaryMarket.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
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

			for _, w := range s {
				wl = append(wl, w)
			}
			if len(wl) > 0 {
				f["asset"] = bson.M{"$in": wl}
			} else {
				return stderr.ErrWhiteList
			}

		}
	}
	result := make([]map[string]interface{}, 0)

	//获取某个用户对NFT所有操作
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
		Index:      "GetNFTRecordByAddress",
		Sort:       bson.M{},
		Filter:     f,
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {
		rr := make(map[string]interface{})
		tokenids := []strval.T{}
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		tokenids = append(tokenids, strval.T(tokenid))
		rr["event"] = item["eventname"]
		user := item["user"].(string)
		rr["user"] = user
		rr["asset"] = asset
		rr["tokenid"] = tokenid
		rr["timestamp"] = item["timestamp"]
		nonce := item["nonce"].(int64)
		rr["nonce"] = nonce
		rr["market"] = item["market"]

		//获取Nft的属性
		var raw2 map[string]interface{}
		err2 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw2)
		fmt.Println("error", raw2)
		log2.Infof("error", raw2)
		if err2 != nil {
			rr["image"] = ""
			rr["name"] = ""
			rr["number"] = int64(-1)
			rr["properties"] = ""

		}

		rr["image"] = raw2["image"]
		rr["thumbnail"] = raw2["thumbnail"]
		rr["name"] = raw2["name"]
		rr["number"] = raw2["number"]
		rr["properties"] = raw2["properties"]

		//获取此时Nft的状态
		var raw1 []map[string]interface{}
		err1 := me.GetNFTByContractHashTokenId(struct {
			ContractHash h160.T
			TokenIds     []strval.T
			Filter       map[string]interface{}
			Raw          *[]map[string]interface{}
		}{ContractHash: h160.T(asset), TokenIds: tokenids, Raw: &raw1}, ret)
		if err1 != nil {
			return err1
		}
		if len(raw1) > 0 {
			nowNFTState := raw1[0]["state"]

			if item["eventname"].(string) == "Auction" { //2种状态   正常   已过期  (卖家事件)

				extendData1 := item["extendData"].(string)
				var dat map[string]interface{}
				if err11 := json.Unmarshal([]byte(extendData1), &dat); err11 == nil {
					auctionType, err12 := strconv.ParseInt(dat["auctionType"].(string), 10, 64)
					if err12 != nil {
						return err12
					}
					auctionAsset := dat["auctionAsset"]
					auctionAmount := dat["auctionAmount"]
					rr["auctionAsset"] = auctionAsset
					rr["auctionAmount"] = auctionAmount
					rr["from"] = item["user"]
					rr["to"] = ""
					if nowNFTState == NFTstate.Expired.Val() {
						if auctionType == 1 {
							rr["state"] = NFTevent.Sell_Expired.Val() //上架过期
						} else if auctionType == 2 {
							rr["state"] = NFTevent.Auction_Expired.Val() //拍卖过期
						}
					} else {
						if auctionType == 1 {
							rr["state"] = NFTevent.Sell_Listed.Val() //上架  正常状态
						} else if auctionType == 2 {
							rr["state"] = NFTevent.Auction_Listed.Val() // 拍卖  正常状态
						}
					}

					////卖家售出事件
					nonce1 := item["nonce"]
					tokenid1 := item["tokenid"]
					asset1 := item["asset"]

					rr1, count, err14 := me.Client.QueryAll(struct {
						Collection string
						Index      string
						Sort       bson.M
						Filter     bson.M
						Query      []string
						Limit      int64
						Skip       int64
					}{
						Collection: "MarketNotification",
						Index:      "someindex",
						Sort:       bson.M{},
						Filter:     bson.M{"nonce": nonce1, "eventname": "Claim", "asset": asset1, "tokenid": tokenid1, "market": bson.M{"$in": []interface{}{args.SecondaryMarket, args.PrimaryMarket}}},
						Query:      []string{},
					}, ret)
					if err14 != nil {
						return err14
					}
					if count > 0 {
						//卖家售出事件

						rr2 := make(map[string]interface{})
						for _, it := range rr1 {
							rr2["asset"] = it["asset"]
							rr2["tokenid"] = it["tokenid"]
							rr2["timestamp"] = it["timestamp"]
							rr2["event"] = it["eventname"]
							rr2["market"] = it["market"]
							rr2["nonce"] = it["nonce"]
							rr2["image"] = rr["image"]
							rr2["name"] = rr["name"]
							rr2["number"] = rr["number"]
							rr2["properties"] = rr["properties"]
							rr2["from"] = it["market"]
							rr2["to"] = it["user"]

							extendData2 := it["extendData"].(string)
							var dat2 map[string]interface{}
							if err3 := json.Unmarshal([]byte(extendData2), &dat2); err3 == nil {

								auctionAsset1 := dat2["auctionAsset"]
								auctionAmount1 := dat2["bidAmount"]
								rr2["auctionAsset"] = auctionAsset1
								rr2["auctionAmount"] = auctionAmount1
							} else {
								return err1
							}

							if auctionType == 1 {

								rr2["state"] = NFTevent.Sell_Sold.Val() // 直买直卖  售出(卖家)

							} else if auctionType == 2 {

								rr2["state"] = NFTevent.Aucion_Deal.Val() //拍卖:成交（卖家）
							}
							result = append(result, rr2)
						}

					}

				} else {
					return err11
				}

			} else if item["eventname"].(string) == "Cancel" { //下架  (卖家事件)

				rr["auctionAsset"] = ""
				rr["auctionAmount"] = ""
				rr["from"] = ""
				rr["to"] = item["user"]
				rr["state"] = NFTevent.Cancel.Val()

			} else if item["eventname"].(string) == "Bid" { //3种状态  正常 已过期  已成交 (买家事件)

				//获取nft所有的bid信息
				var raw3 []map[string]interface{}
				err4 := me.GetAllBidInfoByNFT(struct {
					AssetHash  h160.T
					TokenId    strval.T
					MarketHash []h160.T
					Filter     map[string]interface{}
					Raw        *[]map[string]interface{}
				}{AssetHash: h160.T(asset), TokenId: strval.T(tokenid), MarketHash: []h160.T{args.SecondaryMarket, args.PrimaryMarket}, Raw: &raw3}, ret) //nonce 分组，并按时间排序
				if err4 != nil {
					return err4
				}

				extendData2 := item["extendData"].(string)
				var dat map[string]interface{}
				if err21 := json.Unmarshal([]byte(extendData2), &dat); err21 == nil {

					//bidAmount, err22 := strconv.ParseInt(dat["bidAmount"].(string), 10, 64)
					//bidAmount, flag := new(big.Int).SetString(dat["bidAmount"].(string),10)
					bidAmount := dat["bidAmount"].(string)

					auctionAsset := dat["auctionAsset"]
					rr["auctionAsset"] = auctionAsset
					rr["auctionAmount"] = bidAmount
					rr["from"] = user
					rr["to"] = ""

					for _, it := range raw3 {
						ba := it["bidAmount"].([]*big.Int) //获取竞价数组
						bd := it["bidder"].([]string)      //获取竞价数组

						if nowNFTState == NFTstate.Auction.Val() && raw3[0]["nonce"] == it["nonce"] { //最新上架  拍卖中 2种状态：已退回  正常s

							if bidAmount == ba[0].String() && user == bd[0] { //最高竞价人
								rr["state"] = NFTevent.Auction_Bid.Val() //state :正常

							} else {
								rr["state"] = NFTevent.Auction_Return.Val() //state :已退回

							}
						} else { //历史上架 ：2种状态： 已成交  已退回
							if bidAmount == ba[0].String() && user == bd[0] {
								rr["state"] = NFTevent.Auction_Bid_Deal.Val() //state :已成交
							} else {
								rr["state"] = NFTevent.Auction_Return.Val() //state :已退回
							}

						}
					}

				} else {
					return err21
				}

			} else if item["eventname"].(string) == "Claim" { //  领取  （买家事件）
				extendData3 := item["extendData"].(string)
				var dat map[string]interface{}
				if err31 := json.Unmarshal([]byte(extendData3), &dat); err31 == nil {
					bidAmount := dat["bidAmount"].(string)
					auctionType, err33 := strconv.ParseInt(dat["auctionType"].(string), 10, 64)
					if err33 != nil {
						return err33
					}
					auctionAsset := dat["auctionAsset"]
					user1 := item["user"]
					rr["auctionAsset"] = auctionAsset
					rr["auctionAmount"] = bidAmount
					rr["from"] = raw1[0]["auctor"]
					rr["to"] = user1

					if auctionType == 1 {
						rr["state"] = NFTevent.Sell_Buy.Val() // 直买直卖 购买(买家)

					} else if auctionType == 2 {
						rr["state"] = NFTevent.Auction_Withdraw.Val() //拍卖:领取（买家）
					}

				} else {
					return err31
				}
			} else if item["eventname"].(string) == "Offer" {
				extendData3 := item["extendData"].(string)
				var dat map[string]interface{}
				if err31 := json.Unmarshal([]byte(extendData3), &dat); err31 == nil {
					bidAmount := dat["offerAmount"].(string)
					offerAsset := dat["offerAsset"]
					rr["auctionAsset"] = offerAsset
					rr["auctionAmount"] = bidAmount
					user1 := item["user"]

					//查看offer 当前状态
					offer_nonce := item["nonce"]
					offer, _ := me.Client.QueryOne(struct {
						Collection string
						Index      string
						Sort       bson.M
						Filter     bson.M
						Query      []string
					}{
						Collection: "MarketNotification",
						Index:      "getOfferSate",
						Sort:       bson.M{},
						Filter: bson.M{
							"nonce":   offer_nonce,
							"asset":   item["asset"],
							"tokenid": item["tokenid"],
							"$or": []interface{}{
								bson.M{"eventname": "CompleteOffer"},
								bson.M{"eventname": "CancelOffer"},
							}},
						Query: []string{},
					}, ret)

					if len(offer) > 0 {
						offer_event := offer["eventname"]
						if offer_event == "CompleteOffer" {
							rr["state"] = NFTevent.Offer_Accept.Val() //出价被卖家接受
							rr["from"] = user1
							rr["to"] = offer["user"]
						} else if offer_event == "CancelOffer" {
							rr["state"] = NFTevent.Offer_Cancel.Val() //出价被买家取消
							rr["from"] = ""
							rr["to"] = user1
						}
					} else {
						rr["state"] = NFTevent.Offer.Val() //拍卖:领取（买家）
						rr["from"] = user1
						rr["to"] = ""

					}

				} else {
					return err31
				}

			} else if item["eventname"].(string) == "CompleteOffer" {
				extendData3 := item["extendData"].(string)
				var dat map[string]interface{}
				if err31 := json.Unmarshal([]byte(extendData3), &dat); err31 == nil {
					bidAmount := dat["offerAmount"].(string)

					auctionAsset := dat["offerAsset"]
					user1 := item["user"]
					rr["auctionAsset"] = auctionAsset
					rr["auctionAmount"] = bidAmount
					rr["from"] = user1
					rr["to"] = dat["offerer"]
					rr["state"] = NFTevent.Offer_Complete.Val()

				} else {
					return err31
				}

			} else if item["eventname"].(string) == "CancelOffer" {
				extendData3 := item["extendData"].(string)
				var dat map[string]interface{}
				if err31 := json.Unmarshal([]byte(extendData3), &dat); err31 == nil {
					bidAmount := dat["offerAmount"].(string)
					auctionAsset := dat["offerAsset"]
					user1 := item["user"]
					rr["auctionAsset"] = auctionAsset
					rr["auctionAmount"] = bidAmount
					rr["from"] = ""
					rr["to"] = user1
					rr["state"] = NFTevent.Offer_Cancel.Val()

				} else {
					return err31
				}

			}
			result = append(result, rr)
		}

	}

	//普通账户见的NFT转账 ,去掉和市场之间的转账
	// 获取NFT的Transfer
	market, err := OpenMarketHashFile()
	marketArray := market.MarketHash
	if err != nil {
		return stderr.ErrAMarketConfig
	}
	andArrayto := []interface{}{}
	for _, item := range marketArray {
		andArrayto = append(andArrayto, item)
	}
	filter := bson.M{"$and": []interface{}{
		bson.M{"to": bson.M{"$nin": marketArray}},
		bson.M{"from": bson.M{"$nin": marketArray}},
	}}
	filter["$or"] = []interface{}{
		bson.M{"from": args.Address.TransferredVal()},
		bson.M{"to": args.Address.TransferredVal()},
	}
	//
	if args.SecondaryMarket.Valid() == true {
		if len(wl) > 0 {
			filter["asset"] = bson.M{"$in": wl}
		}
	}

	r3, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11TransferNotification",
		Index:      "GetNep11TransferByAddress",
		Sort:       bson.M{},
		Filter:     filter,
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range r3 {
		from := ""
		to := ""
		if item["from"] != nil && item["from"] != "" {
			from = item["from"].(string)
		}
		if item["to"] != nil && item["to"] != "" {
			to = item["to"].(string)
		}

		log2.Infof("TESTERROR:", item)
		log2.Infof("TESTERROR:", from != args.SecondaryMarket.Val() && to != args.SecondaryMarket.Val() && from != args.PrimaryMarket.Val() && to != args.PrimaryMarket.Val())
		if from != args.SecondaryMarket.Val() && to != args.SecondaryMarket.Val() && from != args.PrimaryMarket.Val() && to != args.PrimaryMarket.Val() {
			rr := make(map[string]interface{})

			asset := item["contract"].(string)
			tokenid := item["tokenId"].(string)

			rr["event"] = "transfer"

			rr["asset"] = asset
			rr["tokenid"] = tokenid
			rr["timestamp"] = item["timestamp"]
			rr["from"] = from
			rr["to"] = to
			rr["auctionAsset"] = ""
			rr["auctionAmount"] = ""

			//获取nft的属性
			var raw3 map[string]interface{}
			err3 := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
			fmt.Println("error", raw3)
			log2.Infof("error", raw3)
			if err3 != nil {
				rr["image"] = ""
				rr["name"] = ""
				rr["number"] = int64(-1)
				rr["properties"] = ""
			}
			fmt.Println("TESTERROR:", raw3)
			rr["image"] = raw3["image"]
			rr["name"] = raw3["name"]
			rr["number"] = raw3["number"]
			rr["properties"] = raw3["properties"]

			if from == args.Address.Val() && (to != args.SecondaryMarket.Val() || to != args.PrimaryMarket.Val()) {
				rr["user"] = from
				rr["state"] = NFTevent.Send.Val()
			} else if to == args.Address.Val() && (from != args.SecondaryMarket.Val() || from != args.PrimaryMarket.Val()) {
				rr["user"] = to
				rr["state"] = NFTevent.Receive.Val()
			}
			result = append(result, rr)
		}
	}

	mapsort.MapSort(result, "timestamp") //按时间排序
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

	num, err := strconv.ParseInt(strconv.Itoa(len(result)), 10, 64)
	if err != nil {
		return err
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
func getNFTProperties(tokenId strval.T, contractHash h160.T, me *T, ret *json.RawMessage, filter map[string]interface{}, Raw *map[string]interface{}) error {

	r4 := make([]map[string]interface{}, 0)
	fmt.Println("query nft properites....")
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Nep11Properties",
		Index:      "getNFTProperties",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": contractHash.Val(), "tokenid": tokenId},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	fmt.Println("error: ", r1)
	asset := r1["asset"].(string)
	tokenid := r1["tokenid"].(string)
	extendData := r1["properties"].(string)
	log2.Infof("error: ", r1)
	fmt.Println("TESTERROR: ", extendData)
	if extendData != "" {
		properties := make(map[string]interface{})
		var data map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
			image, ok := data["image"]
			if ok {
				properties["image"] = image
				//r1["image"] = image
				r1["image"] = ImagUrl(asset, image.(string), "images")
			} else {
				r1["image"] = ""
			}

			tokenuri, ok := data["tokenURI"]
			if ok {
				ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)
				fmt.Println(ppjson)
				if err != nil {
					return err
				}
				for key, value := range ppjson {
					r1[key] = value
					if key != "name" {
						properties[key] = value
					}
					if key == "image" {
						img := value.(string)
						r1["thumbnail"] = ImagUrl(asset, img, "thumbnail")
						r1["image"] = ImagUrl(asset, img, "images")
					}

				}
			}

			if r1["name"] == "" || r1["name"] == nil {
				name, ok1 := data["name"]
				if ok1 {
					r1["name"] = name

				} else {
					r1["name"] = ""
				}
			}
			strArray := strings.Split(r1["name"].(string), "#")
			if len(strArray) >= 2 {
				number := strArray[1]
				n, err2 := strconv.ParseInt(number, 10, 64)
				if err2 != nil {
					r1["number"] = int64(-1)
				}
				r1["number"] = n
				properties["number"] = n
			} else {
				r1["number"] = int64(-1)
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
					r1["number"] = int64(-1)
				}
				properties["number"] = n
				r1["number"] = n
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
				//r1["image"] = string(tb[:])
				r1["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
			} else {
				if image != nil && image != "" {
					if image == nil {
						r1["thumbnail"] = r1["image"]
					} else {
						r1["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
					}
				}
			}

		} else {
			return err
		}

		r1["properties"] = properties
	} else {
		r1["image"] = ""
		r1["name"] = ""
		r1["number"] = int64(-1)
		r1["properties"] = ""
	}

	filter1, err := me.Filter(r1, filter)
	if err != nil {
		return err
	}

	r4 = append(r4, filter1)

	if Raw != nil {
		*Raw = r1
	}
	return nil
}

func getNFTProperties1(tokenId strval.T, contractHash h160.T, me *T, ret *json.RawMessage, filter map[string]interface{}, Raw *map[string]interface{}) error {

	r4 := make([]map[string]interface{}, 0)

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Nep11Properties",
		Index:      "getNFTProperties",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": contractHash.Val(), "tokenid": tokenId},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	asset := r1["asset"].(string)
	tokenid := r1["tokenid"].(string)
	extendData := r1["properties"].(string)
	if extendData != "" {
		properties := make(map[string]interface{})
		var data map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
			image, ok := data["image"]
			if ok {
				properties["image"] = image
				//r1["image"] = image
				r1["image"] = ImagUrl(asset, image.(string), "images")
			} else {
				r1["image"] = ""
			}

			tokenuri, ok := data["tokenURI"]
			if ok {
				ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)

				if err != nil {
					return err
				}
				for key, value := range ppjson {
					r1[key] = value
					if key != "name" {
						properties[key] = value
					}
					if key == "image" {
						img := value.(string)
						r1["thumbnail"] = ImagUrl(asset, img, "thumbnail")
						r1["image"] = ImagUrl(asset, img, "images")
					}

				}
			}

			if r1["name"] == "" || r1["name"] == nil {
				name, ok1 := data["name"]
				if ok1 {
					r1["name"] = name

				} else {
					r1["name"] = ""
				}
			}
			strArray := strings.Split(r1["name"].(string), "#")
			if len(strArray) >= 2 {
				number := strArray[1]
				n, err2 := strconv.ParseInt(number, 10, 64)
				if err2 != nil {
					r1["number"] = int64(-1)
				}
				r1["number"] = n
				properties["number"] = n
			} else {
				r1["number"] = int64(-1)
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
					r1["number"] = int64(-1)
				}
				properties["number"] = n
				r1["number"] = n
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
				//r1["image"] = string(tb[:])
				r1["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
			} else {
				if image != nil && image != "" {
					if image == nil {
						r1["thumbnail"] = r1["image"]
					} else {
						r1["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
					}
				}
			}

		} else {
			return err
		}

		r1["properties"] = properties
	} else {
		r1["image"] = ""
		r1["name"] = ""
		r1["number"] = int64(-1)
		r1["properties"] = ""
	}

	filter1, err := me.Filter(r1, filter)
	if err != nil {
		return err
	}

	r4 = append(r4, filter1)

	if Raw != nil {
		*Raw = r1
	}
	return nil
}

func OpenMarketHashFile() (Config, error) {
	absPath, _ := filepath.Abs("./markethash.yml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log2.Fatalf("Closing file error: %v", err)
		}
	}(f)
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

type Config struct {
	MarketHash []string `yaml:"MarketHash"`
}
