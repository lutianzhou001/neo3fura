package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep17TransferCountByAddress(args struct {
	Address h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	_, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "TransferNotification",
		Index:      "GetNep17TransferCountByAddress",
		Sort:       bson.M{},
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal(), "to": bson.M{"$ne": nil}},
			bson.M{"to": args.Address.TransferredVal()},
		}},
	}, ret)

	if err != nil {
		return err
	}
	r, err := json.Marshal(count)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
