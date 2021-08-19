package home

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Asset
func (me *T) GetAssetCount(ch *chan map[string]interface{}) error {
	assetCount, err := me.getAssetCount()
	if err != nil {
		return err
	}
	*ch <- assetCount

	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Asset"})
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
		newAssetCount, err := me.getAssetCount()
		if err != nil {
			return err
		}
		if assetCount["AssetCount"].(map[string]interface{})["total counts"] != newAssetCount["AssetCount"].(map[string]interface{})["total counts"] {
			*ch <- newAssetCount
			assetCount = newAssetCount
		}
	}
	return nil
}

func (me T) getAssetCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})

	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Asset",
		Index:      "GetAssetCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["AssetCount"] = r1
	return res, nil
}
