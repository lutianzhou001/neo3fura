package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetContractByContractHash(args struct {
	ContractHash h160.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Contract",
		Index:      "someIndex",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"hash": args.ContractHash},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r1, err = me.Filter(r1, args.Filter)
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
