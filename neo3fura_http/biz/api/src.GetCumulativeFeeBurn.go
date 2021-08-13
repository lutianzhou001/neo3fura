package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetCumulativeFeeBurn(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Block",
			Index:      "GetCumulativeFeeBurn",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   []bson.M{bson.M{"$group": bson.M{"_id": "", "feeburn": bson.M{"$sum": "$systemFee"}}}},
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}
	r2, _, err := me.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Block",
			Index:      "GetCumulativeFeeBurn",
			Sort:       bson.M{"_id": -1},
			Filter:     bson.M{},
			Query:      []string{"index", "systemFee"},
			Limit:      10,
			Skip:       0,
		}, ret)
	if err != nil {
		return err
	}
	r1[0]["result"] = r2
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
