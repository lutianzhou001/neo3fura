package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetTransferByAddress(args struct {
	Address h160.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 512
	}
	r1, _, err1 := me.Client.QueryAll(struct {
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
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal()},
			bson.M{"to": args.Address.TransferredVal()},
		}},
		Query: []string{},
	}, ret)
	if err1 != nil {
		return err1
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
		Collection: "TransferNotification",
		Index:      "GetNep17TransferByAddress",
		Sort:       bson.M{},
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal()},
			bson.M{"to": args.Address.TransferredVal()},
		}},
		Query: []string{},
	}, ret)

	if err2 != nil {
		return err2
	}
	r3 := append(r1, r2...)
	r4 := make([]map[string]interface{}, 0)
	for i, item := range r3 {
		if int64(i) < args.Skip {
			continue
		} else if int64(i) > args.Skip+args.Limit-1 {
			continue
		} else {
			r4 = append(r4, item)
		}
	}
	r5, err := me.FilterArrayAndAppendCount(r4, int64(len(r3)), args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r5)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
