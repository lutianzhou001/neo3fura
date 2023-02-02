package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/consts"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetAddressList(args struct {
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Limit <= 0 {
		args.Limit = 20
	}

	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Address",
		Index:      "GetAddressInfo",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$sort": bson.M{"firstusetime": -1}},
			bson.M{"$lookup": bson.M{
				"from": "Address-Asset",
				"let":  bson.M{"address": "$address"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$address", "$$address"}},
						bson.M{"$or": []interface{}{
							bson.M{"$eq": []interface{}{"$asset", consts.NEO}},
							bson.M{"$eq": []interface{}{"$asset", consts.GAS}},
						}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "balance": 1}},
				},
				"as": "nep17balance"},
			},

			bson.M{"$lookup": bson.M{
				"from": "TransferNotification",
				"let":  bson.M{"address": "$address"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$or": []interface{}{
						bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$from", "$$address"}},
							bson.M{"$ne": []interface{}{"$to", nil}},
						}},
						bson.M{"$eq": []interface{}{"$to", "$$address"}},
					}}}},
					bson.M{"$group": bson.M{"_id": "$_id"}},
					bson.M{"$count": "count"},
				},
				"as": "nep17transfer"},
			},

			bson.M{"$lookup": bson.M{
				"from": "Nep11TransferNotification",
				"let":  bson.M{"address": "$address"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$or": []interface{}{
						bson.M{"$eq": []interface{}{"$from", "$$address"}},
						bson.M{"$eq": []interface{}{"$to", "$$address"}},
					}}}},
					bson.M{"$group": bson.M{"_id": "$_id"}},
					bson.M{"$count": "count"},
				},
				"as": "nep11transfer"},
			},

			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {
		nep17balance := item["nep17balance"].(primitive.A)
		nep17transfer := item["nep17transfer"].(primitive.A)
		nep11transfer := item["nep11transfer"].(primitive.A)
		item["neobalance"] = "0"
		item["gasbalance"] = "0"
		if len(nep17balance) > 0 {
			for _, it := range nep17balance {
				balance := it.(map[string]interface{})
				if balance["asset"].(string) == consts.GAS {
					item["gasbalance"] = balance["balance"]
				}
				if balance["asset"].(string) == consts.NEO {
					item["neobalance"] = balance["balance"]
				}

			}
		}
		item["nep17TransferCount"] = 0
		item["nep11TransferCount"] = 0
		if len(nep17transfer) > 0 {
			transfer := nep17transfer[0].(map[string]interface{})
			item["nep17TransferCount"] = transfer["count"]
		}
		if len(nep11transfer) > 0 {
			transfer := nep11transfer[0].(map[string]interface{})
			item["nep11TransferCount"] = transfer["count"]
		}

		delete(item, "nep17balance")
		delete(item, "nep17transfer")
		delete(item, "nep11transfer")

	}

	count, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Address",
		Index:      "someindex",
		Sort:       bson.M{},
		Filter:     bson.M{}}, ret)

	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count["total counts"].(int64), args.Filter)
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
