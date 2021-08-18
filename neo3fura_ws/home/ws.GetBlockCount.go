package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Asset
func (me *T) GetBlockCount(ch *chan map[string]interface{}) error {
	blockCount, err := me.getBlockCount()
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
		newBlockCount, err := me.getBlockCount()
		if err != nil {
			return err
		}
		if blockCount["BlockCount"].(map[string]interface{})["total counts"] != newBlockCount["BlockCount"].(map[string]interface{})["total counts"] {
			*ch <- newBlockCount
			blockCount = newBlockCount
		}
	}
	return nil
}

func (me T) getBlockCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})

	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Block",
		Index:      "GetBlockCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["BlockCount"] = r1
	return res, nil
}
