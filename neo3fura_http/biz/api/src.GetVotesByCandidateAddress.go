package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetVotesByCandidateAddress(args struct {
	CandidateAddress h160.T
	Filter           map[string]interface{}
}, ret *json.RawMessage) error {
	if args.CandidateAddress.Valid() == false {
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
		Index:      "GetVotesByCandidateAddress",
		Sort:       bson.M{},
		Filter:     bson.M{"candidate": args.CandidateAddress.TransferredVal()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
