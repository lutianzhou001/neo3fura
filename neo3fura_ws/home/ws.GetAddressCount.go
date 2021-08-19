package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Address
func (me *T) GetAddressCount(ch *chan map[string]interface{}) error {
	addressCount, err := me.getAddressCount()
	if err != nil {
		return err
	}
	*ch <- addressCount

	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Address"})
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
		newAddressCount, err := me.getAddressCount()
		if err != nil {
			return err
		}
		if addressCount["AddressCount"].(map[string]interface{})["total counts"] != newAddressCount["AddressCount"].(map[string]interface{})["total counts"] {
			*ch <- newAddressCount
			addressCount = newAddressCount
		}
	}
	return nil
}

func (me T) getAddressCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})
	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Address",
		Index:      "GetAddressCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["AddressCount"] = r1
	return res, nil
}
