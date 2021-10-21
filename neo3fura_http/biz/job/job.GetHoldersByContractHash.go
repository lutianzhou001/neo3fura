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
		Index:      "GetHoldersByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	for _, item := range r1 {
		holders, count, err := me.Client.QueryAll(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Address-Asset",
			Index:      "GetHoldersByContractHash",
			Sort:       bson.M{"balance": -1},
			Filter:     bson.M{"asset": item["hash"], "balance": bson.M{"$gt": 0}},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}
		f := make(map[string]interface{})
		count = 0
		for _, holder := range holders {
			f[holder["address"].(string)] = 1
		}
		for _, _ = range f {
			count = count + 1
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
