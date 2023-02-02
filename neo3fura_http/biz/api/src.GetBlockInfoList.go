package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBlockInfoList(args struct {
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
}, ret *json.RawMessage) error {

	if args.Limit == 0 {
		args.Limit = 20
	}

	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{Collection: "Block",
			Index:  "GetBlockInfoList",
			Sort:   bson.M{},
			Filter: bson.M{},
			Pipeline: []bson.M{
				bson.M{"$lookup": bson.M{
					"from": "Transaction",
					"let":  bson.M{"blockhash": "$hash"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$blockhash", "$$blockhash"}},
						}}}},
						bson.M{"$group": bson.M{"_id": "$_id"}},
						bson.M{"$count": "count"},
					},
					"as": "info"},
				},
				//bson.M{"$match": bson.M{"info": bson.M{"$elemMatch": bson.M{"$ne": nil}}}},
				bson.M{"$project": bson.M{"transactioncount": bson.M{"$cond": bson.M{"if": bson.M{"$gt": []interface{}{"$info", []interface{}{}}}, "then": bson.M{"$arrayElemAt": []interface{}{"$info.count", 0}}, "else": 0}},
					"_id": 1, "index": 1, "size": 1, "timestamp": 1, "hash": 1}},
				bson.M{"$sort": -1},
				bson.M{"$limit": args.Limit},
				bson.M{"$skip": args.Skip},
			},
			Query: []string{}}, ret)

	count, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Block",
		Index:      "GetBlockInfoList",
		Sort:       bson.M{},
		Filter:     bson.M{}}, ret)
	if err != nil {
		return err
	}

	r4, err := me.FilterArrayAndAppendCount(r1, count["total counts"].(int64), args.Filter)
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
