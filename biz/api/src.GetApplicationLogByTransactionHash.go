package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetApplicationLogByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.TransactionHash.IsZero() == true {
		return stderr.ErrZero
	}
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Execution",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2, _, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{Collection: "[Execution~Notification(Notifications)]", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"ParentID": r1["_id"]}}, ret)
	if err != nil {
		return err
	}
	notifications := make([]map[string]interface{}, 0)
	for _, item3 := range r2 {
		r3, err := me.Data.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Notification", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"_id": item3["ChildID"]}}, ret)
		if err != nil {
			return err
		}
		notifications = append(notifications, r3)
	}
	if len(notifications) > 0 {
		r1["notifications"] = notifications
	} else {
		r1["notifications"] = []map[string]interface{}{}
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
