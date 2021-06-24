package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/strval"
	"neo3fura/var/stderr"
)

func (me *T) GetExecutionByTrigger(args struct {
	Trigger strval.T
	Limit   int64
	Skip    int64
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Limit == 0 {
		args.Limit = 200
	}
	in := args.Trigger.In([]string{"OnPersist", "PostPersist", "Application", "Verification", "System", "All"})
	if in == false {
		return stderr.ErrInvalidArgs
	}
	var filter bson.M
	if args.Trigger.Val() == "All" {
		filter = bson.M{"$or": []interface{}{
			bson.M{"trigger": "OnPersist"},
			bson.M{"trigger": "PostPersist"},
			bson.M{"trigger": "Application"},
			bson.M{"trigger": "Verification"},
		}}
	} else if args.Trigger.Val() == "System" {
		filter = bson.M{"$or": []interface{}{
			bson.M{"trigger": "OnPersist"},
			bson.M{"trigger": "PostPersist"},
			bson.M{"trigger": "Application"},
			bson.M{"trigger": "Verification"},
		}}
	} else {
		filter = bson.M{
			"trigger": args.Trigger.Val(),
		}
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
		Filter:     filter,
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
