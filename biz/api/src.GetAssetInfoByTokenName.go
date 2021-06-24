package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/strval"
	"neo3fura/var/stderr"
)

func (me *T) GetAssetInfoByTokenName(args struct {
	TokenName strval.T
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TokenName.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Asset",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"tokenname": args.TokenName.Val()},
		Query:      []string{},
	}, ret)
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
