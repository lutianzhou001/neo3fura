package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura/lib/type/strval"
	"neo3fura/var/stderr"
)

func (me *T) GetAssetInfoByTokenName(args struct {
	TokenName strval.T
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TokenName.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Asset",
		Index:      "GetAssetInfoByTokenName",
		Sort:       bson.M{},
		Filter:     bson.M{"tokenname": args.TokenName.Val()},
		Query:      []string{},
	}, ret)
	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}

	r2, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}
	r3, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "Holders"})
	if err != nil {
		return err
	}

	r1["ispopular"] = false
	populars := r2["Populars"].(primitive.A)
	for _, v := range populars {
		if r1["hash"] == v {
			r1["ispopular"] = true
		}
	}


	holders := r3["Holders"].(primitive.A)
	for _, h := range holders {
		m := h.(map[string]interface{})
		for k, v := range m {
			if r1["hash"] == k {
				r1["holders"] = v
			}
		}
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
