package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) SetPopularTokenWhitelist(args struct {
	Filter       map[string]interface{}
	ContractHash []h160.T
}, ret *json.RawMessage) error {
	var hashArr []interface{}
	for _, item := range args.ContractHash {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		hashArr = append(hashArr, item)
	}
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "SetPopularTokens",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": bson.M{"$in": hashArr}},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}

	var data []interface{}
	for _, item := range r1 {
		asset := bson.M{
			"tokenname":   item["tokenname"],
			"symbol":      item["symbol"],
			"totalsupply": item["totalsupply"],
			"decimals":    item["decimals"],
			"hash":        item["hash"],
			"type":        item["type"],
		}
		data = append(data, asset)
	}
	success, err := me.Client.SaveManyJob(struct {
		Collection string
		Data       []interface{}
	}{Collection: "PopularTokenWhitelist", Data: data})
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
