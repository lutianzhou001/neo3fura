package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
)

func (me *T) GetContractInfoByContractHash(args struct {
	Hash   h256.T
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	r1, err := me.Data.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string

		}{
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter: bson.M{"hash": args.Hash},
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
}