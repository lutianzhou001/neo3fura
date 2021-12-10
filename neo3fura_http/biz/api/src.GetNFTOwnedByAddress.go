package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strconv"
	"time"
)

func (me *T) GetNFTOwnedByAddress(args struct {
	Address h160.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
	Raw     *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixMilli()

	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Market",
		Index:      "GetNFTOwnedByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"auctor": args.Address.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {

		a := item["amount"].(primitive.Decimal128).String()
		amount, err := strconv.Atoi(a)
		if err != nil {
			return err
		}

		ba := item["bidAmount"].(primitive.Decimal128).String()
		bidAmount, err := strconv.ParseInt(ba, 10, 64)
		if err != nil {
			return err
		}

		deadline, _ := item["deadline"].(int64)
		auctionType, _ := item["auctionType"].(int32)

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
	}

	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)

	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r2
	}
	*ret = json.RawMessage(r)
	return nil
}
