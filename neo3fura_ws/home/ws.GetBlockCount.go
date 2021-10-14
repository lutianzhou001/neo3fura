package home

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Asset
func (me *T) GetBlockCount(ch *chan map[string]interface{}) error {
	blockCount := make(map[string]interface{})
	totalCounts := make(map[string]interface{})
	lastestBlock,err := me.Client.QueryLastOne(struct{ Collection string }{Collection: "Block"})
	totalCounts["total counts"] =  lastestBlock["index"]
	blockCount["BlockCount"] = totalCounts
	if err != nil {
		return err
	}
	*ch <- blockCount
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Block"})
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
		newBlockCount := make(map[string]interface{})
		newTotalCounts := make(map[string]interface{})
		newTotalCounts["total counts"] = changeEvent["fullDocument"].(map[string]interface{})["index"]
		newBlockCount["BlockCount"] = newTotalCounts
		if blockCount["BlockCount"].(map[string]interface{})["total counts"] != newBlockCount["BlockCount"].(map[string]interface{})["total counts"] {
			*ch <- newBlockCount
			blockCount = newBlockCount
		}
	}
	return nil
}

