package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Address
func (me *T) GetTransactionCount(ch *chan map[string]interface{}) error {
	transactionCount, err := me.getTransactionCount()
	if err != nil {
		return err
	}
	*ch <- transactionCount

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
		newTransactionCount, err := me.getTransactionCount()
		if err != nil {
			return err
		}
		if transactionCount["TransactionCount"].(map[string]interface{})["total counts"] != newTransactionCount["TransactionCount"].(map[string]interface{})["total counts"] {
			*ch <- newTransactionCount
			transactionCount = newTransactionCount
		}
	}
	return nil
}

func (me T) getTransactionCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})
	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Transaction",
		Index:      "GetTransactionCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["TransactionCount"] = r1
	return res, nil
}
