package home

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (me *T) GetAddressCount(ch *chan map[string]interface{}) error {
	var addressCount interface{}
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "AddressCount"})
	if err != nil {
		return err
	}
	lastJob, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "AddressCount"})
	if err != nil {
		return err
	}
	addressCount = lastJob["AddressCount"].(map[string]interface{})["total counts"]
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
		if addressCount != changeEvent["fullDocument"].(map[string]interface{})["AddressCount"].(map[string]interface{})["total counts"] {
			*ch <- changeEvent["fullDocument"].(map[string]interface{})
			addressCount = changeEvent["fullDocument"].(map[string]interface{})["AddressCount"].(map[string]interface{})["total counts"]
		}
	}
	return nil
}
