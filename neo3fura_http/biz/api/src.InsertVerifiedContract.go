package api

import (
	"encoding/json"
	"neo3fura_http/lib/cli"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) InsertVerifiedContract(args struct {
	ContractHash h160.T
	UpdateCounter int
	Id			  int
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, err := me.Client.InsertDocument(struct {
		Collection string
		Index      string
		Insert 	 *cli.Insert
	}{
		Collection: "VerifyContractModel",
		Index:      "InsertVerifiedContract",
		Insert:     &cli.Insert{
			Hash: args.ContractHash,
			Id:args.Id,
			UpdateCounter: args.UpdateCounter,
			},
	}, ret)

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
