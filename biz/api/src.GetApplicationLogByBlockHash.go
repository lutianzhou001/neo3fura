package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetApplicationLogByBlockHash(args struct {
	BlockHash h256.T
	Limit     int64
	Skip      int64
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.BlockHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.BlockHash.IsZero() == true {
		return stderr.ErrZero
	}
	r1, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Execution",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"blockhash": args.BlockHash.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	for _, item2 := range r1 {
		r2, _, err := me.Data.Client.QueryAll(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{Collection: "[Execution~Notification(Notifications)]", Index: "someIndex", Sort: bson.M{}, Filter: bson.M{"ParentID": item2["_id"]}}, ret)
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
			item2["notifications"] = notifications
		} else {
			item2["notifications"] = []map[string]interface{}{}
		}

	}
	r4, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return nil
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
