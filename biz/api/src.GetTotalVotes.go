package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func (me *T) GetTotalVotes(args struct {
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	r1,_,err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64

	}{
		Collection: "Candidate",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{"votesOfCandidate"},
		Limit:      args.Limit,
		Skip:       args.Skip,

	}, ret)
	if err != nil {
		return err
	}
	r := map[string] uint64{}
	for _,item:= range r1{
		ib, _, err := item["votesOfCandidate"].(primitive.Decimal128).BigInt()
		if err != nil {
			return err
		}
		r["totalVotes"] = ib.Uint64() + r["totalVotes"]
	}
	r2, err := json.Marshal(r)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r2)
	return nil
}
