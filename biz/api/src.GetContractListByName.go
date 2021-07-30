package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetContractListByName(args struct {
	Name   string
	Filter map[string]interface{}
	Limit        int64
	Skip         int64
}, ret *json.RawMessage) error {
	var r1, err = me.Data.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:  []bson.M{
				bson.M{"$match": bson.M{"name": bson.M{"$regex":args.Name, "$options": "$i"}}},
				bson.M{"$limit":args.Limit},
				bson.M{"$skip":args.Skip},
				bson.M{"$lookup": bson.M{
					"from": "Transaction",
					"localField": "createTxid",
					"foreignField": "hash",
					"as": "Transaction"}},

				bson.M{"$project": bson.M{
					"_id":0,
					"Transaction.sender":1,
					"hash":1,
					"createtime":1,
					"name":1,
					"id":1},
				},},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	 r2, err := me.Data.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:  []bson.M{
				bson.M{"$lookup": bson.M{
					"from": "Transaction",
					"localField": "createTxid",
					"foreignField": "hash",
					"as": "Transaction"}},
				bson.M{"$match": bson.M{"name": bson.M{"$regex":args.Name, "$options": "$i"}}},
				bson.M{"$count":"total counts"},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	r3, err := me.FilterAggragateAndAppendCount(r1, r2[0]["total counts"], args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}