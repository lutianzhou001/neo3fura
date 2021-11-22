package api

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetCandidateCount(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
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
		return err
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
