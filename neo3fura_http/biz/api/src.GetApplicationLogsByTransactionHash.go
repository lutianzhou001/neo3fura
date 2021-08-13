package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetApplicationLogsByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.TransactionHash.IsZero() == true {
		return stderr.ErrZero
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Execution",
		Index:      "GetApplicationLogByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
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
			Collection: "Notification",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"txid": r1["txid"].(string), "blockhash": r1["blockhash"].(string)}},
				bson.M{"$lookup": bson.M{
					"from":         "Contract",
					"localField":   "contract",
					"foreignField": "hash",
					"as":           "Contract"}},

				bson.M{"$project": bson.M{
					"_id":               0,
					"Contract.manifest": 1,
					"contract":          1,
					"eventname":         1,
					"state":             1,
					"timestamp":         1,
					"Vmstate":           1,
				},
				}},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	for _, item := range r2 {
		p := make(map[string]interface{})
		a := item["Contract"].(primitive.A)
		m := make(map[string]interface{})
		for _, v := range a {
			m = v.(map[string]interface{})
			err := json.Unmarshal([]byte(m["manifest"].(string)), &p)
			if err != nil {
				return err
			}
			m["manifest"] = p
		}
		item["Contract"] = m
	}

	r1["notification"] = r2

	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return nil
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
