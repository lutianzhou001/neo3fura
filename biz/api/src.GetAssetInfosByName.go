package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetAssetInfosByName(args struct {
	Name   string
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
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
		Filter:     bson.M{"tokenname": bson.M{"$regex": args.Name, "$options": "$i"}},
		Query:      []string{"hash", "tokenname", "symbol", "_id"},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	// retrieve all tokens
	r2, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}
	r3, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "Holders"})
	if err != nil {
		return err
	}
	for _, item := range r1 {
		populars := r2["Populars"].(primitive.A)

		item["ispopular"] = false
		for _, v := range populars {
			if item["hash"] == v {
				item["ispopular"] = true
			}
		}
		holders := r3["Holders"].(primitive.A)
		for _, h := range holders {
			m := h.(map[string]interface{})
			for k, v := range m {
				if item["hash"] == k {
					item["holders"] = v
				}
			}
		}
	}
	r4, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil

}
