package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetVerifiedContractByContractHash(args struct {
	ContractHash  h160.T
	UpdateCounter uintval.T
	Filter        map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.UpdateCounter.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "VerifyContract",
		Index:      "GetVerifiedContractByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash, "updatecounter": args.UpdateCounter},
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
