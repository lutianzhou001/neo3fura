package home

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Block
func (me *T) GetBlockInfoList(ch *chan map[string]interface{}) error {
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "BlockInfoList"})
	if err != nil {
		return err
	}
	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		return err
	}
	var blockInfoList interface{}
	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		for i, item := range changeEvent["fullDocument"].(map[string]interface{})["BlockInfoList"].(primitive.A) {
			if i == 0 && blockInfoList != item.(map[string]interface{})["hash"] {
				*ch <- changeEvent["fullDocument"].(map[string]interface{})
				blockInfoList = item.(map[string]interface{})["hash"]
			} else {
				break
			}
		}
	}
	return nil
}
