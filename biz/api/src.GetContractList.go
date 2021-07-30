package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetContractList(args struct {
	Filter map[string]interface{}
	Limit        int64
	Skip         int64
}, ret *json.RawMessage) error {

	var r1, err = me.Data.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Contract",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:  []bson.M{
			bson.M{"$skip":args.Skip},
			bson.M{"$limit":args.Limit},
			bson.M{"$lookup": bson.M{
				"from": "Transaction",
				"localField": "createTxid",
				"foreignField": "hash",
				"as": "Transaction"}},
			bson.M{"$project":
				bson.M{"_id":0,"Transaction.sender":1,"hash":1,"createtime":1,"name":1,"id":1}}},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	r3, err := me.Data.Client.QueryDocument(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M

		}{  Collection: "Contract",
			Index: "someIndex",
			Sort: bson.M{},
			Filter: bson.M{},
		}, ret)
	if err != nil {
		return err
	}
	//r1 = append(r1, r3)
	r2, err := me.FilterAggragateAndAppendCount(r1,r3["total counts"], args.Filter)
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
