package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura/lib/type/h160"
	"neo3fura/var/stderr"
)

func (me *T) GetAssetInfos(args struct {
	Filter    map[string]interface{}
	Addresses []h160.T
	Limit     int64
	Skip      int64
}, ret *json.RawMessage) error {
	var f bson.M
	if args.Addresses == nil {
		f = bson.M{}
	} else {
		addresses := make([]interface{}, 0)
		for _, address := range args.Addresses {
			if address.Valid() == false {
				return stderr.ErrInvalidArgs
			} else {
				addresses = append(addresses, bson.M{"hash": address.TransferredVal()})
			}
		}
		if len(addresses) == 0 {
			f = bson.M{}
		} else {
			f = bson.M{"$or": addresses}
		}
	}

	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "GetAssetInfos",
		Sort:       bson.M{},
		Filter:     f,
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	// retrieve all tokens
	r2, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil {
		return err
	}
	for _, item := range r1 {
		populars := r2["Populars"].(primitive.A)
		for _, v := range populars {
			if item["hash"] == v {
				item["ispopular"] = true
			}
		}
		item["ispopular"] = false
	}
	r3, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
