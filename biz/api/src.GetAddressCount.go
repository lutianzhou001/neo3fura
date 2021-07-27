package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetAddressCount(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	_, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{Collection: "Address"}, ret)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	r1["count"] = count
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
