package api

import (
	"encoding/json"
	"fmt"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetAssetHoldersByContractHash(args struct {
	ContractHash h160.T
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var r1 map[string]interface{}
	r1, err :=me.Client.QueryOne(struct {
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
		Query:      []string{"hash", "totalsupply"},
	}, ret)
	if err != nil {
		return err
	}
	supply, err := strconv.Atoi(r1["totalsupply"].(string))
	r2, count, err :=me.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Address-Asset",
			Index:      "someIndex",
			Sort:       bson.M{"balance": -1},
			Filter:     bson.M{"asset": r1["hash"]},
			Query:      []string{"address", "balance"},
			Limit:      args.Limit,
			Skip:       args.Skip},
		ret)
	if err != nil {
		return err
	}
	for _, item := range r2 {
		balance, err := strconv.Atoi(item["balance"].(string))
		if err != nil {
			return err
		}
		if supply != 0 {
			item["percentage"] = float64(balance) / float64(supply)
		} else {
			item["percentage"] = -1
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
			ContractHash: args.ContractHash,
			Address:      h160.T(fmt.Sprint(item["address"])),
			Filter:       filter,
			Raw:          &raw,
		}, ret)
		if err != nil {
			return err
		}
		item["lasttx"] = raw["latesttx"]
	}
	r4, err := me.FilterArrayAndAppendCount(r2, count, args.Filter)
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
