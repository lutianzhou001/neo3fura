package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"
)

func (me *T) Getblock(args []interface{}, ret *json.RawMessage) error {
	if len(args) <= 1 {
		return stderr.ErrInvalidArgs
	}
	if args[1] != true {
		return stderr.ErrInvalidArgs
	}
	switch args[0].(type) {
	case string:
		blockHash := h256.T(args[0].(string))
		if blockHash.Valid() == false {
			return stderr.ErrInvalidArgs
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
				Collection: "Block",
				Index:      "someIndex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"hash": blockHash}},
					bson.M{"$lookup": bson.M{
						"from":         "Transaction",
						"localField":   "hash",
						"foreignField": "blockhash",
						"as":           "tx"}},
				},
				Query: []string{},
			}, ret)
		if err != nil {
			return err
		}

		r, err := json.Marshal(r1)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
		return nil
	case float64:
		blockHeight := uintval.T(uint64(args[0].(float64)))
		if blockHeight.Valid() == false {
			return stderr.ErrInvalidArgs
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
				Collection: "Block",
				Index:      "someIndex",
				Sort:       bson.M{},
				Filter:     bson.M{},
				Pipeline: []bson.M{
					bson.M{"$match": bson.M{"index": blockHeight}},
					bson.M{"$lookup": bson.M{
						"from":         "Transaction",
						"localField":   "index",
						"foreignField": "blockIndex",
						"as":           "tx"}},
				},
				Query: []string{},
			}, ret)

		if err != nil {
			return err
		}

		r, err := json.Marshal(r1)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}
