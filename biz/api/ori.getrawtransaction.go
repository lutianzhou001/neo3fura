package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) Getrawtransaction(args []interface{}, ret *json.RawMessage) error {
	if args[1] != true {
		return stderr.ErrInvalidArgs
	}
	switch args[0].(type) {
	case string:
		transactionHash := h256.T(args[0].(string))
		if transactionHash.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		r1, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Transaction",
			Index:      "GetRawTransactionByTransactionHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": transactionHash},
			Query:      []string{},
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
