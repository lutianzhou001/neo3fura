package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	//"neo3fura/lib/type/h256"
	//"neo3fura/var/stderr"
)

func (me *T) GetAccountsList(args struct {
	Limit     int64
	Skip      int64
	Filter    map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Limit == 0 {
		args.Limit = 20
	}
	r1, count, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Address",
		Index:      "someIndex",
		Sort:       bson.M{"firstusetime": -1},
		Filter:     bson.M{},
		Query:      []string{"_id", "address", "firstusetime"},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	//gasContractHash := "0xd2a4cff31913016155e38e474a2c06d08be276cf"
	//fmt.Println("get accounts list", len(r1))
	//balances := make(chan string, 20)

	//for _, item := range r1 {
	//	go func() error {
	//		rTemp, err := me.Data.Client.QueryOne(struct {
	//			Collection string
	//			Index      string
	//			Sort       bson.M
	//			Filter     bson.M
	//			Query      []string
	//		}{
	//			Collection: "TransferNotification",
	//			Index:      "someIndex",
	//			Sort:       bson.M{"_id": -1},
	//			Filter: bson.M{"contract": gasContractHash, "$or": []interface{}{
	//				bson.M{"from": item["address"]},
	//				bson.M{"to": item["to"]},
	//			}},
	//			Query: []string{},
	//		}, ret)
	//		if err != nil {
	//			return err
	//		}
	//		fmt.Println(item["address"], rTemp["frombalance"].(string))

	//		var balance string
	//		if rTemp["from"] != nil{
	//			if rTemp["from"].(string) == item["address"] {
	//				balance, _ = rTemp["frombalance"].(string)
	//			} else {
	//				balance, _ = rTemp["tobalance"].(string)
	//			}
	//		} else {
	//			balance, _ = rTemp["tobalance"].(string)
	//		}

	//		balances <- balance
	//		return nil
	//	}()
	//	item["gasBalance"] = <- balances
	//}

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
