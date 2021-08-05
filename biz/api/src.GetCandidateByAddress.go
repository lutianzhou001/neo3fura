package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetCandidateByAddress(args struct {
	Address h160.T
	Filter  map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Candidate",
		Index:      "GetCandidateByAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"candidate": args.Address.TransferredVal()},
		Query:      []string{},
	}, ret)

	if err != nil {
		return err
	}

	r2, err := me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
