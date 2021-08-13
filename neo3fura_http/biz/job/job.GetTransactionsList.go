package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetTransactionList() error {
	message := make(json.RawMessage, 0)
	ret := &message
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
		Index:      "GetTransactionList",
		Sort:       bson.M{"blocktime": -1},
		Filter:     bson.M{},
		Query:      []string{},
		Limit:      10,
		Skip:       0,
	}, ret)
	if err != nil {
		return err
	}
	data := bson.M{"TransactionList": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "TransactionList", Data: data})
	if err != nil {
		return err
	}
	return nil
}
