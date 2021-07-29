package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

// this function may be not supported any more, we only support address in the formart of script hash
func (me *T) GetTransactionList(args struct {
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Transaction",
		Index:      "someIndex",
		Sort:       bson.M{"blocktime":-1},
		Filter:     bson.M{},
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

