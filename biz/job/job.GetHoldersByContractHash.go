package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetHoldersByContractHash() error {
	message := make(json.RawMessage, 0)
	ret := &message
	data := make([]bson.M, 0)
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "GetAssetInfos",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	for _, item := range r1 {
		_, count, err := me.Client.QueryAll(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Address-Asset",
			Index:      "GetAssetHoldersByContractHash",
			Sort:       bson.M{"balance": -1},
			Filter:     bson.M{"asset": item["hash"]},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}
		data = append(data, bson.M{item["hash"].(string): count})
	}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "Holders", Data: bson.M{"Holders": data}})
	if err != nil {
		return err
	}
	return nil
}
