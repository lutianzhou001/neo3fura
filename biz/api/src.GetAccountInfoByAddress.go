package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"

	//"neo3fura/lib/type/h256"
	//"neo3fura/var/stderr"
)

func (me *T) GetAccountInfoByAddress(args struct {
	Address h160.T
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Address",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"address": args.Address.Val()},
		Query:      []string{"_id", "firstusetime"},
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
