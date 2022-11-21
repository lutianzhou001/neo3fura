package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"reflect"
)

func (me *T) SetMarketCollectionWhitelist(args struct {
	Filter       map[string]interface{}
	MarketHash   h160.T
	ContractHash []h160.T
}, ret *json.RawMessage) error {
	//var hashArr []interface{}
	var hashArr2 []interface{}
	for _, item := range args.ContractHash {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		//hashArr = append(hashArr, item)
		hashArr2 = append(hashArr2, item.Val())
	}
	raw := make(map[string]interface{})
	err := me.GetMarketWhiteList(struct {
		MarketHash h160.T
		Filter     map[string]interface{}
		Raw        *map[string]interface{}
	}{MarketHash: args.MarketHash, Raw: &raw}, ret)
	whitelist := raw["whiteList"].([]string)

	for _, item := range args.ContractHash {
		flag := InArray(whitelist, item.Val())
		if !flag {
			return stderr.ErrNotInMarketWhiteList
		}
	}

	success, err := me.Client.SaveJob(struct {
		Collection string
		Data       bson.M
	}{Collection: "MarketCollectionWhitelist", Data: bson.M{"CollectionWhitelist": hashArr2}})
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

func InArray(array []string, element string) bool {
	// 实现查找整形、string类型和bool类型是否在数组中
	if element == "" || array == nil {
		return false
	}
	for _, value := range array {
		// 首先判断类型是否一致
		if reflect.TypeOf(value).Kind() == reflect.TypeOf(element).Kind() {
			// 比较值是否一致
			if value == element {
				return true
			}
		}
	}
	return false
}
