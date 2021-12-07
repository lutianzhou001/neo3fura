package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetNep11PropertiesByContractHashTokenId(args struct {
	ContractHash h160.T
	TokenIds     []strval.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if len(args.TokenIds) == 0 {
		var raw1 []map[string]interface{}
		err := me.GetAssetHoldersListByContractHash(struct {
			ContractHash h160.T
			Limit        int64
			Skip         int64
			Filter       map[string]interface{}
			Raw          *[]map[string]interface{}
		}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
		if err != nil {
			return err
		}
		var tokenIds []strval.T
		for _, raw := range raw1 {
			if len(raw["tokenid"].(string)) == 0 {
				continue
			} else {
				tokenIds = append(tokenIds, strval.T(raw["tokenid"].(string)))
			}
		}
		err = getNep11Properties(tokenIds, me, args.ContractHash, ret, args.Filter)
		if err != nil {
			return err
		}
	} else {
		err := getNep11Properties(args.TokenIds, me, args.ContractHash, ret, args.Filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func getNep11Properties(tokenIds []strval.T, me *T, contractHash h160.T, ret *json.RawMessage, filter map[string]interface{}) error {
	r4 := make([]map[string]interface{}, 0)
	for _, tokenId := range tokenIds {
		if len(tokenId) <= 0 {
			return stderr.ErrInvalidArgs
		}

		r1, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Nep11Properties",
			Index:      "GetNep11PropertiesByContractHashTokenId",
			Sort:       bson.M{"balance": -1},
			Filter:     bson.M{"asset": contractHash.TransferredVal(), "tokenid": tokenId},
			Query:      []string{},
		}, ret)
		if err != nil {
			return err
		}
		filter, err := me.Filter(r1, filter)
		if err != nil {
			return err
		}
		r4 = append(r4, filter)
	}
	r5, err := me.FilterArrayAndAppendCount(r4, int64(len(r4)), filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r5)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
