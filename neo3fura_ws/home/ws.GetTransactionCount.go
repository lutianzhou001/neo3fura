package home

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (me *T) GetTransactionCount(ch *chan map[string]interface{}) error {
	var transactionCount interface{}
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "TransactionCount"})
	if err != nil {
		return err
	}
	lastJob, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "TransactionCount"})
	if err != nil {
		return err
	}
	transactionCount = lastJob["TransactionCount"].(map[string]interface{})["total counts"]
	*ch <- lastJob

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
		if transactionCount != changeEvent["fullDocument"].(map[string]interface{})["TransactionCount"].(map[string]interface{})["total counts"] {
			*ch <- changeEvent["fullDocument"].(map[string]interface{})
			transactionCount = changeEvent["fullDocument"].(map[string]interface{})["TransactionCount"].(map[string]interface{})["total counts"]
		}
	}
	return nil
}
