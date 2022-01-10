package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"strconv"
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

	rr1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Market",
		Index:      "GetNFTByContractHashTokenId",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": args.ContractHash.Val(), "$or": orArray, "amount": bson.M{"$gt": 0}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range rr1 {
		a := item["amount"]
		tokenId := item["tokenid"].(string)
		var a1 string
		switch a.(type) {
		case string:
			a1 = item["amount"].(string)
		case primitive.Decimal128:
			a1 = item["amount"].(primitive.Decimal128).String()
		}

		amount, err := strconv.Atoi(a1)
		if err != nil {
			return err
		}

		b := item["bidAmount"]
		var ba string
		switch b.(type) {
		case string:
			ba = item["bidAmount"].(string)
		case primitive.Decimal128:
			ba = item["bidAmount"].(primitive.Decimal128).String()
		}

		//ba := r1["bidAmount"].(primitive.Decimal128).String()
		//ba := r1["bidAmount"].(string)
		bidAmount, err := strconv.ParseInt(ba, 10, 64)
		if err != nil {
			return err
		}

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

		if amount > 0 && auctionType == 2 && item["owner"] == item["market"] && deadline > currentTime {
			item["state"] = NFTstate.Auction.Val()
		} else if amount > 0 && auctionType == 1 && item["owner"] == item["market"] && deadline > currentTime {
			item["state"] = NFTstate.Sale.Val()
		} else if amount > 0 && item["owner"] != item["market"] {
			item["state"] = NFTstate.NotListed.Val()
		} else if amount > 0 && bidAmount > 0 && deadline < currentTime && item["owner"] == item["market"] {
			item["state"] = NFTstate.Unclaimed.Val()
		} else if amount > 0 && deadline < currentTime && bidAmount == 0 && item["owner"] == item["market"] {
			item["state"] = NFTstate.Expired.Val()
		} else {
			item["state"] = ""
		}

		var raw3 map[string]interface{}
		err1 := getNFTProperties(strval.T(tokenId), args.ContractHash, me, ret, args.Filter, &raw3)
		if err1 != nil {
			item["image"] = ""
			item["name"] = ""
		}
		extendData := raw3["properties"]
		if extendData != nil {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData.(string)), &dat); err == nil {
				value, ok := dat["image"]
				if ok {
					item["image"] = value
				} else {
					item["image"] = ""
				}
				name, ok1 := dat["name"]
				if ok1 {
					item["name"] = name
				} else {
					item["name"] = ""
				}

			} else {
				return err
			}
		} else {
			item["image"] = ""
			item["name"] = ""
		}

	}

	r5, err := me.FilterArrayAndAppendCount(rr1, count, args.Filter)
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
