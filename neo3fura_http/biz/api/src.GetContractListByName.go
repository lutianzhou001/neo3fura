package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetContractListByName(args struct {
	Name   string
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
}, ret *json.RawMessage) error {

	if args.Limit == 0 {
		args.Limit = 512
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
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$group": bson.M{"_id": "$hash",
					"hash":          bson.M{"$last": "$hash"},
					"updatecounter": bson.M{"$last": "$updatecounter"},
					"createtime":    bson.M{"$last": "$createtime"},
					"name":          bson.M{"$last": "$name"},
					"id":            bson.M{"$last": "$id"},
					"createTxid":    bson.M{"$last": "$createTxid"},
				},
				},
				bson.M{"$match": bson.M{"name": bson.M{"$regex": args.Name, "$options": "$i"}}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
				bson.M{"$lookup": bson.M{
					"from":         "Transaction",
					"localField":   "createTxid",
					"foreignField": "hash",
					"as":           "Transaction"}},
				bson.M{"$project": bson.M{
					"_id":                0,
					"Transaction.sender": 1,
					"hash":               1,
					"createtime":         1,
					"name":               1,
					"id":                 1,
					"updatecounter":      1},
				}},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	r2, err := me.Client.QueryAggregate(
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
			Pipeline: []bson.M{
				bson.M{"$lookup": bson.M{
					"from":         "Transaction",
					"localField":   "createTxid",
					"foreignField": "hash",
					"as":           "Transaction"}},
				bson.M{"$match": bson.M{"name": bson.M{"$regex": args.Name, "$options": "$i"}}},
				bson.M{"$count": "total counts"},
			},
			Query: []string{},
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
	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)
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
