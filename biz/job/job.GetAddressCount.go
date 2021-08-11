package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetAddressCount() error {

	message := make(json.RawMessage, 0)
	ret := &message

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

	data := bson.M{"AddressCount": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "AddressCount", Data: data})
	if err != nil {
		return err
	}
	return nil
}
