package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetCandidateCount(args struct{
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, err := me.Data.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Candidate",
		Index:      "someIndex",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{},
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

