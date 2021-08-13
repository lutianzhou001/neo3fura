package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetBlockHeaderByBlockHeight(args struct {
	BlockHeight uintval.T
	Filter      map[string]interface{}
	Raw         *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHeight.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Header",
		Index:      "GetBlockHeaderByBlockHeight",
		Sort:       bson.M{},
		Filter:     bson.M{"index": args.BlockHeight},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r1
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
