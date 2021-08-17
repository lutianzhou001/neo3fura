package home

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (me *T) GetContractCount(ch *chan map[string]interface{}) error {
	// var hash string
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "ContractCount"})
	if err != nil {
		return err
	}
	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		return err
	}
	var contractCount map[string]interface{}
	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		if contractCount["ContractCount"] != changeEvent["fullDocument"].(map[string]interface{})["ContractCount"] {
			*ch <- changeEvent["fullDocument"].(map[string]interface{})
			contractCount["ContractCount"] = changeEvent["fullDocument"].(map[string]interface{})["ContractCount"]
		}
	}
	return nil
}