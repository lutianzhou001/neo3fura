package job

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetBlockCount() error {
	message := make(json.RawMessage, 0)
	ret := &message
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Block",
		Index:      "GetBlockCount",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{},
		Query:      []string{"index"},
	}, ret)
	if err != nil {
		return err
	}

	data := bson.M{"BlockCount": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "BlockCount", Data: data})
	if err != nil {
		return err
	}
	return nil
}
