package job

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetContractCount() error {
	message := make(json.RawMessage, 0)
	ret := &message

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

	data := bson.M{"ContractCount": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "ContractCount", Data: data})
	if err != nil {
		return err
	}
	return nil
}
