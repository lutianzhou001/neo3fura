package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	currentTime := time.Now().UnixMilli()
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
			Sort:       bson.M{"auctionAsset": -1},
			Filter:     bson.M{"asset": args.ContractHash, "tokenid": tokenId, "amount": bson.M{"$gt": 0}},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}

		a := r1["amount"].(primitive.Decimal128).String()
		amount, err := strconv.Atoi(a)
		if err != nil {
			return err
		}

		ba := r1["bidAmount"].(primitive.Decimal128).String()
		bidAmount, err := strconv.ParseInt(ba, 10, 64)
		if err != nil {
			return err
		}

		deadline, _ := r1["deadline"].(int64)
		auctionType, _ := r1["auctionType"].(int32)

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
