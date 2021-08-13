package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetBlockInfoByBlockHash(args struct {
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
			Index:      "GetBlockInfoByBlockHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": args.BlockHash},
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}
	r4, err := me.Client.QueryDocument(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
		}{Collection: "Transaction",
			Index:  "GetBlockInfoByBlockHash",
			Sort:   bson.M{},
			Filter: bson.M{"blockhash": args.BlockHash}}, ret)
	if err != nil {
		return err
	}
	if r4["total counts"] == nil {
		r1["transactioncount"] = 0
	} else {
		r1["transactioncount"] = r4["total counts"]
	}

	r2, err2 := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Nep11TransferNotification",
		Index:      "GetBlockInfoByBlockHash",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": r1["timestamp"]},
	}, ret)
	if err2 != nil {
		return err2
	}
	if r2["total counts"] == nil {
		r1["nep11count"] = 0
	} else {
		r1["nep11count"] = r2["total counts"]
	}
	r3, err3 := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "TransferNotification",
		Index:      "GetBlockInfoByBlockHash",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": r1["timestamp"]},
	}, ret)
	if err3 != nil {
		return err3
	}
	if r3["total counts"] == nil {
		r1["nep17count"] = 0
	} else {
		r1["nep17count"] = r3["total counts"]
	}

	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
