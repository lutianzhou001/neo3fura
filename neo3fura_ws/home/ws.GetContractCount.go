package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Address
func (me *T) GetContractCount(ch *chan map[string]interface{}) error {
	contractCount, err := me.getContractCount()
	if err != nil {
		return err
	}
	*ch <- contractCount

	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Contract"})
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
		newContractCount, err := me.getContractCount()
		if err != nil {
			return err
		}
		if contractCount["ContractCount"].(map[string]interface{})["total counts"] != newContractCount["ContractCount"].(map[string]interface{})["total counts"] {
			*ch <- newContractCount
			contractCount = newContractCount
		}
	}
	return nil
}

func (me T) getContractCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})
	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Contract",
		Index:      "GetContractCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["ContractCount"] = r1
	return res, nil
}
