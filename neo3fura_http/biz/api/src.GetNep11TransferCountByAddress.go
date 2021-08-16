package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNep11TransferCountByAddress(args struct {
	Address h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Nep11TransferNotification",
		Index:      "GetNep11TransferCountByAddress",
		Sort:       bson.M{},
		Filter: bson.M{"$or": []interface{}{
			bson.M{"from": args.Address.TransferredVal()},
			bson.M{"to": args.Address.TransferredVal()},
		}},
	}, ret)
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
