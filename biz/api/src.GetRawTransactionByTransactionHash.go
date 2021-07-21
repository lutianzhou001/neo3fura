package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetRawTransactionByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
	Raw             *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Transaction",
		Index:      "GetRawTransactionByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.TransactionHash.Val()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	var raw1 map[string]interface{}
	var raw2 map[string]interface{}

	err = me.GetVmStateByTransactionHash(struct {
		TransactionHash h256.T
		Filter          map[string]interface{}
		Raw             *map[string]interface{}
	}{
		TransactionHash: h256.T(fmt.Sprint(args.TransactionHash.Val())),
		Filter:          nil,
		Raw:             &raw1,
	}, ret)
	if err != nil {
		return err
	}
	r1["vmstate"] = raw1["vmstate"].(string)

	err = me.GetBlockByBlockHash(struct {
		BlockHash h256.T
		Filter    map[string]interface{}
		Raw       *map[string]interface{}
	}{
		BlockHash: h256.T(fmt.Sprint(r1["blockhash"])),
		Filter:    nil,
		Raw:       &raw2,
	}, ret)
	if err != nil {
		return err
	}
	r1["timestamp"] = raw2["timestamp"]
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
