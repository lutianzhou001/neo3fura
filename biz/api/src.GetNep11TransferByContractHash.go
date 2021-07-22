package api

import (
	"encoding/json"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep11TransferByContractHash(args struct {
	ContractHash h160.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Nep11TransferNotification",
		Index:      "someIndex",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"contract": args.ContractHash.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	for _, item := range r1 {
		r, err := me.Data.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Block",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": item["blockhash"]},
			Query:      []string{"timestamp"},
		}, ret)
		if err != nil {
			return err
		}
		item["time"] = r["timestamp"]
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
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
