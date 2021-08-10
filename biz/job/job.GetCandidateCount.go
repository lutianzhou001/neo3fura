package job

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetCandidateCount() error {

	message := make(json.RawMessage, 0)
	ret := &message

	r1, err := me.Client.QueryDocument(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
	}{
		Collection: "Candidate",
		Index:      "GetCandidateCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
	}, ret)

	data := bson.M{"CandidateCount": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "CandidateCount", Data: data})
	if err != nil {
		return err
	}
	return nil
}
