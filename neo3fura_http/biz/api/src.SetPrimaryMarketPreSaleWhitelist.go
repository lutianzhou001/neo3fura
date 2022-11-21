package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) SetPrimaryMarketPreSaleWhitelist(args struct {
	Filter     map[string]interface{}
	MarketHash h160.T
	Address    []h160.T
}, ret *json.RawMessage) error {
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	key := args.MarketHash.Val()
	var hashArr []interface{}
	for _, item := range args.Address {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		hashArr = append(hashArr, item.Val())
	}

	success, err := me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "PrimaryMarketPreSaleWhitelist", Data: bson.M{"market": key, "PreSaleWhitelist": hashArr}})
	if err != nil {
		return err
	}
	result := make(map[string]interface{})
	if success {
		result["msg"] = "Insert document done!"

	} else {
		result["msg"] = "Insert document failed!"
	}
	r, err := json.Marshal(result)
	if err != nil {
		return stderr.ErrInsertDocument
	}
	*ret = json.RawMessage(r)
	return nil
}
