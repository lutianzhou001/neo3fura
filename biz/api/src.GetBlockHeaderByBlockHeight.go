package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/uintval"
	"neo3fura/var/stderr"
)

func (me *T) GetBlockHeaderByBlockHeight(args struct {
	BlockHeight uintval.T
	Filter      map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHeight.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Header",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"index": args.BlockHeight},
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
