package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/consts"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"time"
)

func (me *T) GetMarketDayVolumeByAsset(args struct {
	AssetHash h160.T
	LastDays  int64
	Filter    map[string]interface{}
	Raw       *map[string]interface{}
}, ret *json.RawMessage) error {

	lastday := args.LastDays * 24 * 60 * 60 * 1000
	today := time.Now()
	newtime := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	ld := newtime.UnixMilli() - lastday

	fmt.Println(ld, time.UnixMilli(ld).Format(consts.ShortForm))
	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	result, err := me.Client.QueryAggregateJob(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{Collection: "MarketDayVolume",
		Index:  "MarketDayVolume",
		Sort:   bson.M{},
		Filter: bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"asset": args.AssetHash, "date": bson.M{"$gt": primitive.DateTime(ld)}}},
			bson.M{"$sort": bson.M{"date": 1}},
		},
		Query: []string{},
	}, ret)

	if err != nil {
		return err
	}

	for _, item := range result {
		date := item["date"].(primitive.DateTime)
		item["date"] = date.Time().UnixMilli()
	}

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}
