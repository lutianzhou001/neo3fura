package home

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (me *T) GetAddressCount(ch *chan map[string]interface{}) error {
	// var hash string
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "AddressCount"})
	if err != nil {
		return err
	}
	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		return err
	}
	var addressCount map[string]interface{}
	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		if addressCount["AddressCount"] != changeEvent["fullDocument"].(map[string]interface{})["AddressCount"] {
			*ch <- changeEvent["fullDocument"].(map[string]interface{})
			addressCount["AddressCount"] = changeEvent["fullDocument"].(map[string]interface{})["AddressCount"]
		}
	}
	return nil
}
