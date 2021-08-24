package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetNewAddresses() error {
	message := make(json.RawMessage, 0)
	ret := &message

	r0, err := me.Client.QueryOne(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{Collection: "Transaction", Index: "GetNewAddresses", Sort: bson.M{"_id": -1}}, ret)
	if err != nil {
		return err
	}

	_, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{Collection: "Address", Index: "GetNewAddresses", Filter: bson.M{"firstusetime": bson.M{"$gt": r0["blocktime"].(int64) - 3600*24*1000}}}, ret)
	if err != nil {
		return err
	}

	data := bson.M{"NewAddresses": count}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "NewAddresses", Data: data})
	if err != nil {
		return err
	}
	return nil
}
