package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
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
	r3, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "Holders"})
	if err != nil {
		return err
	}
	for _, item := range r1 {
		populars := r2["Populars"].(primitive.A)

		item["ispopular"] = false
		for _, v := range populars {
			if item["hash"] == v {
				item["ispopular"] = true
			}
		}
		holders := r3["Holders"].(primitive.A)
		for _, h := range holders {
			m := h.(map[string]interface{})
			for k, v := range m {
				if item["hash"] == k {
					item["holders"] = v
				}
			}
		}
	}
	r4, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
