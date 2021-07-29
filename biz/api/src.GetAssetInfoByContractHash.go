package api

import (
	"encoding/json"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetAssetInfoByContractHash(args struct {
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
		Collection: "Asset",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash.Val()},
		Query:      []string{},
	}, ret)
	r2, err := me.Data.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Address-Asset",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"asset": r1["hash"]},
	}, ret)
	r1["total_holders"] = r2["total counts:"]
	_, err = me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "TransferNotification", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"contract": r1["hash"]}, Query: []string{},
	}, ret)
	if err != nil {
		r1["standard"] = "NEP11"
	} else {
		r1["standard"] = "NEP17"
	}
	delete(r1, "_id")
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
