package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetCandidateByVoterAddress(args struct {
	VoterAddress h160.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.VoterAddress.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Vote",
		Index:      "GetCandidateByVoterAddress",
		Sort:       bson.M{"blockNumber": -1},
		Filter:     bson.M{"voter": args.VoterAddress.TransferredVal()},
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
