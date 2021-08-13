package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetBlockRewardByBlockHash(args struct {
	BlockHash h256.T
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Block",
			Index:      "GetBlockRewardByBlockHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": args.BlockHash},
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}
	r2, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "TransferNotification",
			Index:      "GetBlockRewardByBlockHash",
			Sort:       bson.M{},
			Filter:     bson.M{"timestamp": r1["timestamp"], "from": nil},
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
