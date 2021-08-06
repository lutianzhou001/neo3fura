package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetTransferByBlockHash(args struct {
	BlockHash h256.T
	Limit       int64
	Skip        int64
	Filter      map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 500
	}

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Block",
		Index:      "GetBlockByBlockHeight",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.BlockHash},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	r2, _, err2 := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11TransferNotification",
		Index:      "GetNep11TransferByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": r1["timestamp"],"txid":"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Query:      []string{},
	}, ret)
	if err2 != nil {
		return err2
	}

	r3, _, err3 := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"timestamp": r1["timestamp"],"txid":"0x0000000000000000000000000000000000000000000000000000000000000000"},
		Query:      []string{},
	}, ret)

	if err3 != nil {
		return err3
	}
	r4 := append(r2, r3...)
	r5 := make([]map[string]interface{}, 0)
	for i, item := range r4 {
		if int64(i) < args.Skip {
			continue
		} else if int64(i) > args.Skip+args.Limit-1 {
			continue
		} else {
			r5 = append(r5, item)
		}
	}
	r6, err := me.FilterArrayAndAppendCount(r5, int64(len(r4)), args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r6)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
