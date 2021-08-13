package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetAssetCount() error {
	message := make(json.RawMessage, 0)
	ret := &message

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
		return err
	}
	data := bson.M{"AssetCount": r1}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "AssetCount", Data: data})
	if err != nil {
		return err
	}
	return nil
}
