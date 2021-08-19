package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// TransactionList
func (me *T) GetTransactionList(ch *chan map[string]interface{}) error {
	transactionList, err := me.getTransactionList()
	if err != nil {
		return err
	}
	*ch <- transactionList

	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Transaction"})
	if err != nil {
		return err
	}
	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		return err
	}

	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		newTransactionList, err := me.getTransactionList()
		if err != nil {
			return err
		}
		if transactionList["TransactionList"].([]map[string]interface{})[0]["hash"] == newTransactionList["TransactionList"].([]map[string]interface{})[0]["hash"] {
			*ch <- newTransactionList
			transactionList = newTransactionList
		}
	}
	return nil
}

func (me T) getTransactionList() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})
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
		return nil, err
	}
	res["TransactionList"] = r1
	return res, nil
}
