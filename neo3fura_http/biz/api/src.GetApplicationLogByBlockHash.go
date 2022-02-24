package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetApplicationLogByBlockHash(args struct {
	BlockHash h256.T
	Limit     int64
	Skip      int64
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.BlockHash.IsZero() == true {
		return stderr.ErrZero
	}

	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Execution",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"blockhash": args.BlockHash.Val()}},
				bson.M{"$lookup": bson.M{
					"from": "Notification",
					"let":  bson.M{"execution_txid": "$txid", "execution_blockhash": "$blockhash"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$txid", "$$execution_txid"}},
							bson.M{"$eq": []interface{}{"$blockhash", "$$execution_blockhash"}},
						}}}},
						bson.M{"$project": bson.M{"txid": 1, "contract": 1}}},
					"as": "notifications"},
				},
			},

			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(strconv.Itoa(len(r1)), 10, 64)
	if err != nil {
		return err
	}

	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return nil
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
