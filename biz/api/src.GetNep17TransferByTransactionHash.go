package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetNep17TransferByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
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
		Collection: "TransferNotification",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
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
		TransactionHash: h256.T(fmt.Sprint(r1["txid"])),
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
	switch raw2["timestamp"].(type) {
	case int64:
		r1["timestamp"] = raw2["timestamp"].(int64)
	case float64:
		r1["timestamp"] = raw2["timestamp"].(float64)
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
