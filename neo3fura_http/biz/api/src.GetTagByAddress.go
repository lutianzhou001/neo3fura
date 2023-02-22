package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/consts"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetTagByAddress(args struct {
	Address h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	currentTime := time.Now().UnixNano() / 1e6
	before7days := currentTime - 7*24*60*60*1000

	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "TransferNotification",
		Index:      "GetAssetsHeldByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{
				"contract":  consts.NEO,
				"timestamp": bson.M{"$gt": before7days},
				"$or": []interface{}{
					bson.M{"from": args.Address},
					bson.M{"to": args.Address},
				},
			}},
			bson.M{"$group": bson.M{"_id": "$_id", "sum": bson.M{"$sum": "$value", "address": args.Address}}},

			bson.M{"$lookup": bson.M{
				"from": "Address-Asset",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
		},
		Query: []string{},
	}, ret)

	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, 1, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
