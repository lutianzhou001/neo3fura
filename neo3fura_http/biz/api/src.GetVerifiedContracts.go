package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetVerifiedContracts(args struct {
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
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
		Collection: "VerifyContractModel",
		Index:      "GetVerifiedContracts",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r2)
	return nil
}
