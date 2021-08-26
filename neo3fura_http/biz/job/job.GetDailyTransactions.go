package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetDailyTransactions() error {
	message := make(json.RawMessage, 0)
	ret := &message

	r0, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Transaction", Index: "GetDailyTransactions", Sort: bson.M{"_id": -1}}, ret)
	if err != nil {
		return err
	}
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Transaction",
		Index:      "GetDailyTransactions",
		Sort:       bson.M{},
		Filter:     bson.M{"blocktime": bson.M{"$gt": r0["blocktime"].(int64) - 3600*24*1000}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	data := bson.M{"DailyTransactions": len(r1)}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "DailyTransactions", Data: data})
	if err != nil {
		return err
	}
	return nil
}
