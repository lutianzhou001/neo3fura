package api

import (
	"encoding/base64"
	"encoding/json"
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
	Address    h160.T
	MarketHash h160.T // market合约地址
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	f := bson.M{"user": args.Address.Val()}
	var wl []interface{}
	if len(args.MarketHash) > 0 {
		if args.MarketHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			//f["market"] = args.MarketHash.Val()
			//白名单
			raw1 := make(map[string]interface{})
			err1 := me.GetMarketWhiteList(struct {
				MarketHash h160.T
				Filter     map[string]interface{}
				Raw        *map[string]interface{}
			}{MarketHash: args.MarketHash, Raw: &raw1}, ret) //nonce 分组，并按时间排序
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
		if err2 != nil {
			rr["image"] = ""
			rr["name"] = ""
			rr["number"] = int64(-1)
			rr["properties"] = ""

		}

		rr["image"] = raw2["image"]
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
			// 上架过期 （卖家事件）
			//if nowNFTState == NFTstate.Expired.Val() {
			//	rr1 := make(map[string]interface{})
			//	rr1["event"] = ""
			//	rr1["user"] = ""
			//	rr1["asset"] = raw1[0]["asset"]
			//	rr1["tokenid"] = raw1[0]["tokenid"]
			//	rr1["timestamp"] = raw1[0]["timestamp"]
			//	rr1["auctionAsset"] = raw1[0]["auctionAsset"]
			//	rr1["auctionAmount"] = raw1[0]["auctionAmount"]
			//	rr1["from"] = raw1[0]["auctor"]
			//	rr1["to"] = ""
			//	auctionType, _ := raw1[0]["auctionType"].(int32)
			//	if auctionType == 1 {
			//		rr1["state"] = NFTevent.Sell_Expired.Val() //上架过期
			//	} else if auctionType == 2 {
			//		rr1["state"] = NFTevent.Auction_Expired.Val() //拍卖过期
			//	}
			//
			//	rr1["image"] = raw1[0]["image"]
			//	rr1["name"] = raw1[0]["name"]
			//	rr1["number"] = raw1[0]["number"]
			//	rr1["properties"] = raw1[0]["properties"]
			//
			//	result = append(result, rr1)
			//}

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
						Filter:     bson.M{"nonce": nonce1, "eventname": "Claim", "asset": asset1, "tokenid": tokenid1, "market": args.MarketHash},
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
					MarketHash h160.T
					Filter     map[string]interface{}
					Raw        *[]map[string]interface{}
				}{AssetHash: h160.T(asset), TokenId: strval.T(tokenid), MarketHash: args.MarketHash, Raw: &raw3}, ret) //nonce 分组，并按时间排序
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
	if args.MarketHash.Valid() == true {
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
		Limit:      args.Limit,
		Skip:       args.Skip,
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

		if from != args.MarketHash.Val() && to != args.MarketHash.Val() {
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
			if err3 != nil {
				rr["image"] = ""
				rr["name"] = ""
				rr["number"] = int64(-1)
				rr["properties"] = ""
			}

			rr["image"] = raw3["image"]
			rr["name"] = raw3["name"]
			rr["number"] = raw3["number"]
			rr["properties"] = raw3["properties"]

			if from == args.Address.Val() && to != args.MarketHash.Val() {
				rr["user"] = from
				rr["state"] = NFTevent.Send.Val()
			} else if to == args.Address.Val() && from != args.MarketHash.Val() {
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

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Nep11Properties",
		Index:      "getNFTProperties",
		Sort:       bson.M{"balance": -1},
		Filter:     bson.M{"asset": contractHash.TransferredVal(), "tokenid": tokenId},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	extendData := r1["properties"].(string)
	if extendData != "" {
		properties := make(map[string]interface{})
		var data map[string]interface{}
		if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
			image, ok := data["image"]
			if ok {
				properties["image"] = image
				r1["image"] = image
			} else {
				r1["image"] = ""
			}
			name, ok1 := data["name"]
			if ok1 {
				r1["name"] = name
				strArray := strings.Split(name.(string), "#")
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

			} else {
				r1["name"] = ""
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
				r1["image"] = string(tb[:])
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
