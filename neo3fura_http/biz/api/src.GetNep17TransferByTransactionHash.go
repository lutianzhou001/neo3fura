package api

import (
	"encoding/json"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNep17TransferByTransactionHash(args struct {
	TransactionHash h256.T
	Limit           int64
	Skip            int64
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
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
		Collection: "Notification",
		Index:      "GetNep17TransferByTransactionHash",
		Sort:       bson.M{"index": 1},
		Filter:     bson.M{"eventname": "Transfer", "txid": args.TransactionHash.Val()},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {

		item["vmstate"] = item["Vmstate"].(string)
		state := item["state"].(map[string]interface{})
		value := state["value"].(primitive.A)

		base64from := value[0].(map[string]interface{})["value"]
		base64to := value[1].(map[string]interface{})["value"]

		if base64from != nil {
			from, err1 := crypto.Base64Decode(base64from.(string))
			if err1 != nil {
				return err1
			}
			item["from"] = "0x" + helper.BytesToHex(helper.ReverseBytes(from))
		} else {
			item["from"] = nil
		}
		if base64to != nil {
			to, err1 := crypto.Base64Decode(base64to.(string))

			if err1 != nil {
				return err1
			}

			item["to"] = "0x" + helper.BytesToHex(helper.ReverseBytes(to))
		} else {
			item["to"] = nil
		}

		item["value"] = value[2].(map[string]interface{})["value"]

		delete(item, "state")
		delete(item, "Vmstate")
		delete(item, "eventname")

		r, err := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "Asset",
			Index:      "GetNep17TransferByTransactionHash",
			Sort:       bson.M{},
			Filter:     bson.M{"hash": item["contract"]},
			Query:      []string{"tokenname", "decimals", "symbol"},
		}, ret)
		if err == nil {
			item["tokenname"] = r["tokenname"]
			item["decimals"] = r["decimals"]
			item["symbol"] = r["symbol"]

		} else if err.Error() == "NOT FOUND" {
			item["tokenname"] = ""
			item["decimals"] = ""
			item["symbol"] = ""
		} else {
			return err
		}
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
