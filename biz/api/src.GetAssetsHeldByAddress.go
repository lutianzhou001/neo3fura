package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetAssetsHeldByAddress(args struct {
	Address h160.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var r1 map[string]interface{}
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
		Query:      []string{"_id"},
	}, ret)
	if err != nil {
		return err
	}
	fmt.Println(r1["_id"])
	r2, count, err := me.Data.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Asset-Address",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{"_id": r1["_id"]},
			Query: []string{},
			Limit: args.Limit,
			Skip: args.Skip,
		}, ret)
	if err != nil {
		return err
	}
	fmt.Println(r2)
	r3 := make([]map[string]interface{}, 0)
	for _, item := range r2 {
		r, err := me.Data.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Asset", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"_id": item["ParentID"]}}, ret)
		if err != nil {
			return err
		}
		var raw map[string]interface{}
		var filter map[string]interface{}
		if args.Filter["balanceinfo"] == nil {
			filter = nil
		} else {
			filter = args.Filter["balanceinfo"].(map[string]interface{})
		}
		err = me.GetBalanceByContractHashAddress(struct {
			ContractHash h160.T
			Address      h160.T
			Filter       map[string]interface{}
			Raw          *map[string]interface{}
		}{
			ContractHash: h160.T(fmt.Sprint(r["hash"])),
			Address:      args.Address,
			Filter:       filter,
			Raw:          &raw,
		}, ret)
		if err != nil {
			return err
		}
		r["balanceinfo"] = raw
		r3 = append(r3, r)
	}
	r4, err := me.FilterArrayAndAppendCount(r3, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
