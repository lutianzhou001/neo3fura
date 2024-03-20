package home

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// bridge
func (me *T) GetBridge(contract string, nonce int32, ch *chan map[string]interface{}) error {
	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Notification"})
	if err != nil {
		return err
	}

	//matchStage := bson.D{
	//	{"$match", bson.D{
	//		{"operationType", "insert"},
	//		//{"fullDocument.index", bson.D{
	//		//	{"index", 1},
	//		//}},
	//	}},
	//}
	//cs, err := c.Watch(context.TODO(), mongo.Pipeline{matchStage})   //需要开启副本集模式

	cs, err := c.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		return err
	}
	defer cs.Close(context.TODO())
	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent map[string]interface{}
		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}
		fullDocument := changeEvent["fullDocument"].(map[string]interface{})
		contractAdd := fullDocument["contract"]
		eventName := fullDocument["eventname"]
		if contractAdd == contract {
			if eventName == "Withdrawal" || eventName == "Claimable" {
				state := fullDocument["state"].(map[string]interface{})
				stateValue := state["value"].(primitive.A)
				eventNonce := stateValue[0].(int32)
				if nonce == eventNonce {
					*ch <- fullDocument
				}

			}
		}

	}
	return nil
}
