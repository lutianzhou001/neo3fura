package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep17TransferByContractHash(args struct {
	ContractHash h160.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"contract": args.ContractHash.Val()}},
			bson.M{"$sort": bson.M{"timestamp0": -1, "_id": -1}},
			bson.M{"$lookup": bson.M{
				"from": "Execution",
				"let":  bson.M{"txid": "$txid", "blockhash": "$blockhash"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$txid", "$$txid"}},
						bson.M{"$eq": []interface{}{"$blockhash", "$$blockhash"}},
					}}}},
					bson.M{"$project": bson.M{"vmstate": 1}},
				},
				"as": "execution"},
			},
			bson.M{"$skip": args.Skip},
			bson.M{"$limit": args.Limit},
		},
		Query: []string{},
	}, ret)

	_, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByContractHash",
		Sort:       bson.M{"timestamp": -1},
		Filter:     bson.M{"contract": args.ContractHash.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	for _, item := range r1 {
		if item["execution"] != nil {
			execution := item["execution"].(primitive.A)
			item["vmstate"] = execution[0].(map[string]interface{})["vmstate"]

		} else {
			item["vmstate"] = "FAULT"
		}

		delete(item, "execution")
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
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
