package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetContractByContractHash(args struct {
	ContractHash h160.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Contract",
		Index:      "GetContractByContractHash",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{"hash": args.ContractHash},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2, err := me.Client.QueryDocument(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
		}{Collection: "ScCall",
			Index:  "GetContractByContractHash",
			Sort:   bson.M{},
			Filter: bson.M{"contractHash": args.ContractHash}}, ret)
	if err != nil {
		return err
	}
	if r2["total counts"] == nil {
		r1["totalsccall"] = 0

	} else {
		r1["totalsccall"] = r2["total counts"]
	}
	if r1["createTxid"] != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		r3, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Transaction",
			Index:      "GetContractByContractHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": r1["createTxid"]},
			Query:      []string{"sender"},
		}, ret)
		if err != nil {
			return err
		}
		r1["sender"] = r3["sender"]
	} else {
		r1["sender"] = nil
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
