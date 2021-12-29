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
	r4 := make([]map[string]interface{}, 0)
	for _, tokenId := range tokenIds {
		if len(tokenId) <= 0 {
			return stderr.ErrInvalidArgs
		}

		r1, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Market",
			Index:      "GetNFTByContractHashTokenId",
			Sort:       bson.M{},
			Filter:     bson.M{"asset": args.ContractHash.Val(), "tokenid": tokenId, "amount": bson.M{"$gt": 0}},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}

		a := r1["amount"]
		var a1 string
		switch a.(type) {
		case string:
			a1 = r1["amount"].(string)
		case primitive.Decimal128:
			a1 = r1["amount"].(primitive.Decimal128).String()
		}

		amount, err := strconv.Atoi(a1)
		if err != nil {
			return err
		}

		b := r1["bidAmount"]
		var ba string
		switch b.(type) {
		case string:
			ba = r1["bidAmount"].(string)
		case primitive.Decimal128:
			ba = r1["bidAmount"].(primitive.Decimal128).String()
		}

		//ba := r1["bidAmount"].(primitive.Decimal128).String()
		//ba := r1["bidAmount"].(string)
		bidAmount, err := strconv.ParseInt(ba, 10, 64)
		if err != nil {
			return err
		}

		dl := r1["deadline"]

		at := r1["auctionType"]

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
			deadline = int64(at.(int32))
		}

		if amount > 0 && auctionType == 2 && r1["owner"] == r1["market"] && deadline > currentTime {
			r1["state"] = NFTstate.Auction.Val()
		} else if amount > 0 && auctionType == 1 && r1["owner"] == r1["market"] && deadline > currentTime {
			r1["state"] = NFTstate.Sale.Val()
		} else if amount > 0 && r1["owner"] != r1["market"] {
			r1["state"] = NFTstate.NotListed.Val()
		} else if amount > 0 && bidAmount > 0 && deadline < currentTime && r1["owner"] == r1["market"] {
			r1["state"] = NFTstate.Unclaimed.Val()
		} else if amount > 0 && deadline < currentTime && bidAmount == 0 && r1["owner"] == r1["market"] {
			r1["state"] = NFTstate.Expired.Val()
		} else {
			r1["state"] = ""
		}

		var raw3 map[string]interface{}
		err1 := getNFTProperties(tokenId, args.ContractHash, me, ret, args.Filter, &raw3)
		if err1 != nil {
			return err1
		}
		extendData := raw3["properties"]
		if extendData != nil {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(extendData.(string)), &dat); err == nil {
				value, ok := dat["image"]
				if ok {
					r1["image"] = value
				} else {
					r1["image"] = ""
				}
				name, ok1 := dat["name"]
				if ok1 {
					r1["name"] = name
				} else {
					r1["name"] = ""
				}

			} else {
				return err
			}
		} else {
			r1["image"] = ""
			r1["name"] = ""
		}

		filter, err := me.Filter(r1, args.Filter)
		if err != nil {
			return err
		}
		r4 = append(r4, filter)
	}
	r5, err := me.FilterArrayAndAppendCount(r4, int64(len(r4)), args.Filter)
	if err != nil {
		return err
	}

	if args.Raw != nil {
		*args.Raw = r4
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
