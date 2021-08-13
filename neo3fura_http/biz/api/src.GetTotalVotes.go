package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (me *T) GetTotalVotes(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Candidate",
		Index:      "GetTotalVotes",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{"votesOfCandidate"},
	}, ret)
	if err != nil {
		return err
	}
	r := map[string]uint64{}
	for _, item := range r1 {
		ib, _, err := item["votesOfCandidate"].(primitive.Decimal128).BigInt()
		if err != nil {
			return err
		}
		r["totalvotes"] = ib.Uint64() + r["totalvotes"]
	}
	r2, err := json.Marshal(r)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r2)
	return nil
}
