package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	log2 "neo3fura_http/lib/log"
	limit "neo3fura_http/var/const"
)

func (me T) GetMarketSupply() {
	message := make(json.RawMessage, 0)
	ret := &message

	assetList, err2 := me.GetNep11Asset()
	if err2 != nil {
		log2.Fatal("GetMarketNep11Asset err")
	}

	result := make(map[string]interface{})
	for _, it := range assetList {

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
				Index:      "someindex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"owner": bson.M{"$ne": limit.NullAddress}, "asset": it, "amount": bson.M{"$gt": 0}}},
					bson.M{"$group": bson.M{"_id": "$tokenid"}},
					bson.M{"$count": "count"},
				},
				Query: []string{},
			}, ret)

		if err != nil {
			log2.Fatal("GetMarketNftCount err")
		}

		if len(r1) > 0 {
			result["asset"] = it
			result["totalsupply"] = r1[0]["count"]
		}

		// 存到本地数据库中
		//data := bson.M{"MarketIndex": result}
		_, err = me.Client.UpdateJob(struct {
			Collection string
			Data       bson.M
			Filter     bson.M
		}{Collection: "MarketIndex", Data: result, Filter: bson.M{"asset": it}})
		if err != nil {
			log2.Fatal("MarketIndex update err")
		}

	}

}
