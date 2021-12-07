package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetNFTMarket(args struct {
	ContractHash h160.T   //  asset
	AssetHash    h160.T   // auctionType
	NFTState     strval.T //state:aution  sale  notlisted  unclaimed
	Sort         strval.T //listedTime  price
	Order        int64    //-1:降序  +1：升序
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
	Raw          *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixMilli()
	pipeline := []bson.M{}

	if args.ContractHash != "" {
		if args.ContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"asset": args.ContractHash}}
			pipeline = append(pipeline, a)
		}
	}

	if args.AssetHash != "" {
		if args.AssetHash.Valid() == false {
			return stderr.ErrInvalidArgs
		} else {
			a := bson.M{"$match": bson.M{"asset": args.ContractHash}}
			pipeline = append(pipeline, a)
		}
	}

	if args.Sort == "timestamp" || args.Sort == "price" {
		b := bson.M{}
		if args.Order == -1 || args.Order == 1 {
			b = bson.M{"$sort": bson.M{args.Sort.Val(): args.Order}}
		} else {
			b = bson.M{"$sort": bson.M{args.Sort.Val(): 1}}
		}
		pipeline = append(pipeline, b)
	}

	if args.NFTState.Val() == NFTstate.Auction.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline = []bson.M{

			bson.M{"$match": bson.M{"asset": args.ContractHash}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 2}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "auction"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}

	} else if args.NFTState.Val() == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline = []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "sale"}},
			bson.M{"$match": bson.M{"difference": true}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}

	} else if args.NFTState.Val() == NFTstate.NotListed.Val() { //未上架  accont >0 && owner != market
		pipeline = []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "difference": bson.M{"$eq": []string{"$owner", "$market"}}, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "notlisted"}},
			bson.M{"$match": bson.M{"difference": false}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}

	} else if args.NFTState.Val() == NFTstate.Unclaimed.Val() { //未领取  accont >0 &&  runtime > deadline && bidAccount >0
		pipeline = []bson.M{

			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"bidAmount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$lt": currentTime}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": "unclaimed"}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}

	} else { //默认  account > 0
		pipeline = []bson.M{
			//bson.M{"$sort": bson.M{"_id":-1}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$project": bson.M{"_id": 1, "asset": 1, "tokenid": 1, "amount": 1, "owner": 1, "market": 1, "auctionType": 1, "auctor": 1, "auctionAsset": 1, "auctionAmount": 1, "deadline": 1, "bidder": 1, "bidAmount": 1, "timestamp": 1, "state": ""}},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		}
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
			Collection: "Market",
			Index:      "GetNFTMarket",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{"_id", "asset", "tokenid", "amount", "owner", "market", "auctionType", "auctor", "auctionAsset", "auctionAmount", "deadline", "bidder", "bidAmount", "timestamp", "state"},
		}, ret)

	if err != nil {
		return err
	}
	pipeline = append(pipeline[:len(pipeline)-1], pipeline[len(pipeline):]...)
	var group = bson.M{"$group": bson.M{"_id": "$_id"}}

	pipeline = append(pipeline, group)
	var countKey = bson.M{"$count": "total counts"}

	pipeline = append(pipeline, countKey)

	r2, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Market",
			Index:      "GetNFTMarket",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}

	var count interface{}
	if len(r2) != 0 {
		count = r2[0]["total counts"]
	} else {
		count = 0
	}

	if err != nil {
		return err
	}
	//r2, err := me.Filter(r1, args.Filter)
	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)
	//	r3, err := me.FilterArrayAndAppendCount(r1,count ,args.Filter)

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
