package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/NFTevent"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"reflect"
	"strconv"
)

func (me *T) GetNFTRecordByAddress(args struct {
	Address            h160.T
	MarketContractHash h160.T // market合约地址
	Limit              int64
	Skip               int64
	Filter             map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	if args.MarketContractHash != "" {
		if args.MarketContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
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
		Filter:     bson.M{"user": args.Address.Val()},
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

		//获取Nft的属性
		var raw2 map[string]interface{}
		err := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw2)
		if err != nil {
			return err
		}

		extendData := raw2["properties"]
		if extendData != nil {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData.(string)), &dat); err == nil {
				image, ok := dat["image"]
				if ok {
					rr["image"] = image
				} else {
					rr["image"] = ""
				}
				name, ok1 := dat["name"]
				if ok1 {
					rr["name"] = name
				} else {
					rr["name"] = ""
				}

			} else {
				return err
			}

		} else {
			rr["image"] = ""
			rr["name"] = ""
		}

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
		nowNFTState := raw1[0]["state"]

		// 上架过期 （卖家事件）
		if nowNFTState == NFTstate.Expired.Val() {
			rr1 := make(map[string]interface{})
			rr1["event"] = ""
			rr1["user"] = ""
			rr1["asset"] = raw1[0]["asset"]
			rr1["tokenid"] = raw1[0]["tokenid"]
			rr1["timestamp"] = raw1[0]["timestamp"]
			rr1["auctionAsset"] = raw1[0]["auctionAsset"]
			rr1["auctionAmount"] = raw1[0]["auctionAmount"]
			rr1["from"] = raw1[0]["timestamp"]
			rr1["to"] = ""
			auctionType, _ := raw1[0]["auctionType"].(int32)
			if auctionType == 1 {
				rr1["state"] = NFTevent.Sell_Expired.Val() //上架过期
			} else if auctionType == 1 {
				rr1["state"] = NFTevent.Auction_Expired.Val() //拍卖过期
			}

			result = append(result, rr1)
		}

		if item["eventname"].(string) == "Auction" { //2种状态   正常   已过期  (卖家事件)
			extendData := item["extendData"].(string)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
				auctionType, err := strconv.ParseInt(dat["auctionType"].(string), 10, 64)
				if err != nil {
					return err
				}
				auctionAsset := dat["auctionAsset"]
				auctionAmount, err := strconv.ParseInt(dat["auctionAmount"].(string), 10, 64)

				if err != nil {
					return err
				}
				rr["auctionAsset"] = auctionAsset
				rr["auctionAmount"] = auctionAmount
				rr["from"] = item["user"]
				rr["to"] = ""
				if auctionType == 1 {
					rr["state"] = NFTevent.Sell_Listed.Val() //上架  正常状态
				} else if auctionType == 2 {
					rr["state"] = NFTevent.Auction_Listed.Val() // 拍卖  正常状态
				}
			} else {
				return err
			}

		} else if item["eventname"].(string) == "Cancel" { //下架  (卖家事件)

			rr["auctionAsset"] = ""
			rr["auctionAmount"] = ""
			rr["from"] = ""
			rr["to"] = item["user"]
			rr["state"] = NFTevent.Cancel.Val()

		} else if item["eventname"].(string) == "Bid" { //3种状态  正常 已过期  已成交 (买家事件)

			//获取nft所有的bid信息
			var raw2 []map[string]interface{}
			err4 := me.GetAllBidInfoByNFT(struct {
				AssetHash h160.T
				TokenId   strval.T
				Filter    map[string]interface{}
				Raw       *[]map[string]interface{}
			}{AssetHash: h160.T(asset), TokenId: strval.T(tokenid), Raw: &raw2}, ret) //nonce 分组，并按时间排序
			if err4 != nil {
				return err4
			}

			extendData := item["extendData"].(string)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &dat); err == nil {

				bidAmount, err := strconv.ParseInt(dat["bidAmount"].(string), 10, 64)
				if err != nil {
					return err
				}

				auctionAsset := dat["auctionAsset"]
				rr["auctionAsset"] = auctionAsset
				rr["auctionAmount"] = bidAmount
				rr["from"] = user
				rr["to"] = ""

				for _, it := range raw2 {
					ba := reflect.ValueOf(it["bidAmount"]) //获取竞价数组
					bd := reflect.ValueOf(it["bidder"])    //获取竞价数组
					println(" ")
					if nowNFTState == NFTstate.Auction.Val() && raw2[0]["nonce"] == it["nonce"] { //最新上架  拍卖中 2种状态：已退回  正常s
						if bidAmount == ba.Index(0).Int() && user == bd.Index(0).String() {
							rr["state"] = NFTevent.Auction_Bid.Val() //state :正常
						} else {
							rr["state"] = NFTevent.Auction_Return.Val() //state :已退回
						}
					} else {
						if bidAmount == ba.Index(0).Int() && user == bd.Index(0).String() { //上架 ：2种状态： 已成交  已退回
							rr["state"] = NFTevent.Auction_Bid_Deal.Val() //state :已成交
						} else {
							rr["state"] = NFTevent.Auction_Return.Val() //state :已退回
						}

					}
				}

			} else {
				return err
			}

		} else if item["eventname"].(string) == "Claim" { //  领取  （买家事件）
			extendData := item["extendData"].(string)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
				bidAmount, err := strconv.ParseInt(dat["bidAmount"].(string), 10, 64)
				auctionType, err := strconv.ParseInt(dat["auctionType"].(string), 10, 64)
				if err != nil {
					return err
				}
				auctionAsset := dat["auctionAsset"]
				user := item["user"]
				rr["auctionAsset"] = auctionAsset
				rr["auctionAmount"] = bidAmount
				rr["from"] = raw1[0]["auctor"]
				rr["to"] = user
				//卖家售出事件
				rr1 := make(map[string]interface{})
				rr1 = rr
				rr1["from"] = rr["to"]
				rr1["to"] = rr["from"]

				if auctionType == 1 {
					rr["state"] = NFTevent.Sell_Buy.Val()   // 直买直卖 购买(买家)
					rr1["state"] = NFTevent.Sell_Sold.Val() // 直买直卖  售出(卖家)

				} else if auctionType == 2 {
					rr["state"] = NFTevent.Auction_Withdraw.Val() //拍卖:领取（买家）
					rr1["state"] = NFTevent.Aucion_Deal.Val()     //拍卖:成交（卖家）
				}
				result = append(result, rr1)

			} else {
				return err
			}
		}

		result = append(result, rr)

	}

	//普通账户见的NFT转账 ,去掉和市场之间的转账
	// 获取NFT的Transfer
	var raw2 []map[string]interface{}
	err1 := me.GetNep11TransferByAddress(struct {
		Address h160.T
		Limit   int64
		Skip    int64
		Start   int64
		End     int64
		Filter  map[string]interface{}
		Raw     *[]map[string]interface{}
	}{Address: args.Address, Raw: &raw2}, ret)
	if err1 != nil {
		return err1
	}
	for _, item := range raw2 {
		from := ""
		to := ""
		if item["from"] != nil {
			from = item["from"].(string)
		}
		if item["from"] != nil {
			to = item["to"].(string)
		}

		if from != args.MarketContractHash.Val() && to != args.MarketContractHash.Val() {
			rr := make(map[string]interface{})

			asset := item["contract"].(string)
			tokenid := item["tokenId"].(string)

			rr["event"] = "transfer"
			rr["user"] = item["user"]
			rr["asset"] = asset
			rr["tokenid"] = tokenid
			rr["timestamp"] = item["timestamp"]
			rr["from"] = from
			rr["to"] = to
			rr["auctionAsset"] = ""
			rr["auctionAmount"] = ""

			//获取nft的属性
			var raw3 map[string]interface{}
			err := getNFTProperties(strval.T(tokenid), h160.T(asset), me, ret, args.Filter, &raw3)
			if err != nil {
				return err
			}

			extendData := raw3["properties"].(string)
			if extendData != "" {
				var dat map[string]interface{}
				if err := json.Unmarshal([]byte(extendData), &dat); err == nil {
					image, ok := dat["image"]
					if ok {
						rr["image"] = image
					} else {
						rr["image"] = ""
					}
					name, ok1 := dat["name"]
					if ok1 {
						rr["name"] = name
					} else {
						rr["name"] = ""
					}

				} else {
					return err
				}

			} else {
				rr["image"] = ""
				rr["name"] = ""
			}

			if from == args.Address.Val() && to != args.MarketContractHash.Val() {
				rr["state"] = NFTevent.Send.Val()
			} else if to == args.Address.Val() && from != args.MarketContractHash.Val() {
				rr["state"] = NFTevent.Receive.Val()
			}
			result = append(result, rr)
		}
	}
	num, err := strconv.ParseInt(strconv.Itoa(len(result)), 10, 64)
	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(result, num, args.Filter)
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

	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11Properties",
		Index:      "getNFTProperties",
		Sort:       bson.M{"balance": -1},
		Filter:     bson.M{"asset": contractHash.Val(), "tokenid": tokenId},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	r2, err := me.FilterArrayAndAppendCount(r1, count, filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}

	if Raw != nil {
		*Raw = r2
	}

	*ret = json.RawMessage(r)
	return nil
}
