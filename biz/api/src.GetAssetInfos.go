package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetAssetInfos(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, count, err := me.Client.QueryAll(struct {
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
	// retrieve all tokens
	r2, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}
	for _, item := range r1 {
		populars := r2["Populars"].(primitive.A)
		for _, v := range populars {
			if item["hash"] == v {
				item["ispopular"] = true
			}
		}
		item["ispopular"] = false
	}
	r3, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
