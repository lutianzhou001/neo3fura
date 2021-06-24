package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetNep11BalanceByContractHashAddress(args struct {
	ContractHash h160.T
	Address      h160.T
	Filter       map[string]interface{}
	Raw          *[]map[string]interface{}
	Limit        int64
	Skip         int64
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var r1 []map[string]interface{}
	var err error
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
		Filter: bson.M{"contract": args.ContractHash.Val(), "$or": []interface{}{
			bson.M{"from": args.Address.Val()},
			bson.M{"to": args.Address.Val()},
		}},
		Query: []string{},
		Limit: args.Limit,
		Skip:  args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2 := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		temp := make(map[string]interface{})
		if item["from"].(string) == args.Address.Val() {
			temp["balance"] = item["frombalance"]
		} else {
			temp["balance"] = item["tobalance"]
		}
		temp["latesttx"] = item
		r2 = append(r2, temp)
	}
	r3, err := me.FilterArrayAndAppendCount(r2, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r2
	}
	*ret = json.RawMessage(r)
	return nil
}
