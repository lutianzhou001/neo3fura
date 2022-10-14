package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
)

func (me *T) GetPopularToken(args struct {
	Filter   map[string]interface{}
	Limit    int64
	Skip     int64
	Standard strval.T
}, ret *json.RawMessage) error {

	if args.Standard != "NEP11" && args.Standard != "NEP11" {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 100
	}
	//获取前三天的排名token
	r1, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PopularTokens"})
	if err != nil && err != stderr.ErrNotFound {
		return err
	}
	popularTokens := r1["Populars"]
	r2, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "GetPopularAsset",
		Sort:       bson.M{},
		Filter:     bson.M{"type": args.Standard.Val(), "hash": bson.M{"$in": popularTokens}},
		Query:      []string{},
	}, ret)

	for _, item := range r2 {
		delete(item, "_id")
		delete(item, "firsttransfertime")
	}
	//获取在白名单的token
	r3, err := me.Client.QueryAggregateJob(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{
		Collection: "PopularTokenWhitelist",
		Index:      "GetPopularTokens",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"type": args.Standard.Val()}},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}
	for _, item := range r3 {
		delete(item, "_id")
		r2 = append(r2, item)
	}

	r4, err := me.FilterArrayAndAppendCount(r2, int64(len(r2)), args.Filter)
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
