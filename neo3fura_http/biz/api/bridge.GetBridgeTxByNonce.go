package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strconv"
)

func (me *T) GetBridgeTxByNonce(args struct {
	ContractHash h160.T
	Nonce        int64
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	nonceStr := strconv.FormatInt(args.Nonce, 10)
	r1, _, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Notification",
		Index:      "GetBridgeTxByNonce",
		Sort:       bson.M{"_id": -1},
		Filter: bson.M{"contract": args.ContractHash.Val(),
			"$or":                 []interface{}{bson.M{"eventname": "Withdrawal"}, bson.M{"eventname": "Claimable"}},
			"state.value.0.value": nonceStr,
		},
		Query: []string{},
		Limit: args.Limit,
		Skip:  args.Skip,
	}, ret)

	var result map[string]interface{}
	if len(r1) > 0 {
		result = r1[0]
	}
	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}
