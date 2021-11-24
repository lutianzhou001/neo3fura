package api

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetContractCount(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, err := me.Client.GetDistinctCount(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "Contract",
		Index:      "GetContractCount",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline:   []bson.M{},
		Query:      []string{},
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
