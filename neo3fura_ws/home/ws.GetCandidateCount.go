package home

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Address
func (me *T) GetCandidateCount(ch *chan map[string]interface{}) error {
	candidateCount, err := me.getCandidateCount()
	if err != nil {
		return err
	}
	*ch <- candidateCount

	c, err := me.Client.GetCollection(struct{ Collection string }{Collection: "Candidate"})
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
		newCandidateCount, err := me.getCandidateCount()
		if err != nil {
			return err
		}
		if candidateCount["CandidateCount"].(map[string]interface{})["total counts"] != newCandidateCount["CandidateCount"].(map[string]interface{})["total counts"] {
			*ch <- newCandidateCount
			candidateCount = newCandidateCount
		}
	}
	return nil
}

func (me T) getCandidateCount() (map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message
	res := make(map[string]interface{})
	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Candidate",
		Index:      "GetCandidateCount",
		Sort:       bson.M{},
		Filter:     bson.M{"state": true},
	}, ret)
	if err != nil {
		return nil, err
	}
	res["CandidateCount"] = r1
	return res, nil
}
