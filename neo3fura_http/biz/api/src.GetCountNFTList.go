package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/NFTstate"
	"neo3fura_http/lib/type/h160"
	address "neo3fura_http/var/const"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetCountNFTList(args struct {
	SecondaryMarket h160.T //
	PrimaryMarket   h160.T
	ContractHash    h160.T
	Filter          map[string]interface{}
	Raw             *map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	pipeline := []bson.M{}
	pipeline2 := []bson.M{}

	if len(args.PrimaryMarket) > 0 && args.PrimaryMarket != "" {
		if args.PrimaryMarket.Valid() == false {
			return stderr.ErrInvalidArgs
		}
	}
	if len(args.ContractHash) > 0 && args.ContractHash != "" {
		if args.ContractHash.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		a := bson.M{"$match": bson.M{"asset": args.ContractHash.Val()}}
		pipeline = append(pipeline, a)
		pipeline2 = append(pipeline2, a)
	}

	if len(args.SecondaryMarket) > 0 && args.SecondaryMarket != "" {
		a := bson.M{"$match": bson.M{"market": args.SecondaryMarket}}
		pipeline = append(pipeline, a)
		pipeline2 = append(pipeline2, a)
	}
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
	var wl []interface{}
	for _, w := range s {
		wl = append(wl, w)
	}
	if len(wl) > 0 {
		white := bson.M{"$match": bson.M{"asset": bson.M{"$in": wl}}}
		//white := bson.M{"$match": bson.M{"asset": bson.M{"$in": []interface{}{"0x6c91e9997b8e74dcfa5ebb56fe5672dedd724b8f","0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f"}}}}
		pipeline = append(pipeline, white)
		pipeline2 = append(pipeline2, white)
	} else {
		return stderr.ErrWhiteList
	}

	if "auction" == NFTstate.Auction.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": bson.M{"$ne": address.NullAddress}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 2}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
		}
		pipeline = append(pipeline, pipeline1...)

	}
	if "sale" == NFTstate.Sale.Val() { //出售中 accont >0 && auctionType =1 && owner=market && runtime <deadline
		pipeline1 := []bson.M{
			bson.M{"$match": bson.M{"owner": bson.M{"$ne": address.NullAddress}}},
			bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}}},
			bson.M{"$match": bson.M{"auctionType": bson.M{"$eq": 1}}},
			bson.M{"$match": bson.M{"deadline": bson.M{"$gt": currentTime}}},
		}
		pipeline2 = append(pipeline2, pipeline1...)

	}

	//group
	setAndGroup := []bson.M{
		bson.M{"$lookup": bson.M{
			"from": "SelfControlNep11Properties",
			"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
			"pipeline": []bson.M{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
					bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
					bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
				}}}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x50ac1c37690cc2cfc594472833cf57505d5f46de"}}, "then": "$asset",
					"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x6a2893f97401e2b58b757f59d71238d91339856a"}}, "then": "$image",
						"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x9f344fe24c963d70f5dcf0cfdeb536dc9c0acb3a"}}, "then": "$tokenid",
							"else": "$name"}}}}}}}},
			},
			"as": "properties"},
		},
		bson.M{"$group": bson.M{"_id": bson.M{"asset": "$asset", "class": "$properties.class"}, "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "deadline": bson.M{"$last": "$deadline"}, "auctionAmount": bson.M{"$last": "$auctionAmount"}, "timestamp": bson.M{"$last": "$timestamp"}, "propertiesArr": bson.M{"$push": "$$ROOT"}}},
		bson.M{"$project": bson.M{"_id": 1, "properties": 1, "asset": 1, "tokenid": 1, "propertiesArr": 1, "auctionAmount": 1, "deadline": 1, "timestamp": 1}},

		//bson.M{"$count":"count"},
	}

	pipeline = append(pipeline, setAndGroup...)
	pipeline2 = append(pipeline2, setAndGroup...)

	var r1, err3 = me.Client.QueryAggregate(
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

	if err3 != nil {
		return err3
	}

	var r2, err4 = me.Client.QueryAggregate(
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
			Pipeline:   pipeline2,
			Query:      []string{},
		}, ret)

	if err4 != nil {
		return err4
	}
	fmt.Println(r2)
	var count int64
	var count2 int64
	result := make(map[string]interface{})
	for _, item := range r1 {
		if item["propertiesArr"] != nil {
			groupInfo := item["propertiesArr"].(primitive.A)
			count += int64(len(groupInfo))
		}

	}
	result["auction"] = count
	for _, item := range r2 {
		if item["propertiesArr"] != nil {
			groupInfo := item["propertiesArr"].(primitive.A)
			count2 += int64(len(groupInfo))
		}
	}
	result["sale"] = count2

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
