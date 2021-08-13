package job

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (me T) GetBlockInfoList() error {
	message := make(json.RawMessage, 0)
	ret := &message

	r1, _, err := me.Client.QueryAll(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Block",
			Index:      "GetBlockInfoList",
			Sort:       bson.M{"_id": -1},
			Filter:     bson.M{},
			Query:      []string{"_id", "index", "size", "timestamp", "hash"},
			Limit:      10,
			Skip:       0,
		}, ret)
	if err != nil {
		return err
	}

	r2 := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		r3, err := me.Client.QueryDocument(
			struct {
				Collection string
				Index      string
				Sort       bson.M
				Filter     bson.M
			}{Collection: "Transaction",
				Index:  "GetBlockInfoList",
				Sort:   bson.M{},
				Filter: bson.M{"blockhash": item["hash"]}}, ret)
		if err != nil {
			return err
		}
		if r3["total counts"] == nil {
			item["transactioncount"] = 0
		} else {
			item["transactioncount"] = r3["total counts"]
		}
		r2 = append(r2, item)
	}

	data := bson.M{"BlockInfoList": r2}
	_, err = me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "BlockInfoList", Data: data})
	if err != nil {
		return err
	}
	return nil
}
