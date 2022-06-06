package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"io/ioutil"
	"neo3fura_http/lib/joh"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	//"go.mongodb.org/mongo-driver/bson"
	"net/http"
	//"strconv"
)

func (me *T) GetNep11PropertiesByContractHashTokenId(args struct {
	ContractHash h160.T
	TokenIds     []strval.T
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if len(args.TokenIds) == 0 {
		var raw1 []map[string]interface{}
		err := me.GetAssetHoldersListByContractHash(struct {
			ContractHash h160.T
			Limit        int64
			Skip         int64
			Filter       map[string]interface{}
			Raw          *[]map[string]interface{}
		}{ContractHash: args.ContractHash, Raw: &raw1}, ret)
		if err != nil {
			return err
		}
		var tokenIds []strval.T
		for _, raw := range raw1 {
			if len(raw["tokenid"].(string)) == 0 {
				continue
			} else {
				tokenIds = append(tokenIds, strval.T(raw["tokenid"].(string)))
			}
		}
		err = getNep11Properties(tokenIds, me, args.ContractHash, ret, args.Filter)
		if err != nil {
			return err
		}
	} else {
		err := getNep11Properties(args.TokenIds, me, args.ContractHash, ret, args.Filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func getNep11Properties(tokenIds []strval.T, me *T, contractHash h160.T, ret *json.RawMessage, filter map[string]interface{}) error {
	r4 := make([]map[string]interface{}, 0)
	for _, tokenId := range tokenIds {
		if len(tokenId) <= 0 {
			return stderr.ErrInvalidArgs
		}

		//r1, err := me.Client.QueryOne(struct {
		//	Collection string
		//	Index      string
		//	Sort       bson.M
		//	Filter     bson.M
		//	Query      []string
		//}{
		//	Collection: "Nep11Properties",
		//	Index:      "GetNep11PropertiesByContractHashTokenId",
		//	Sort:       bson.M{"balance": -1},
		//	Filter:     bson.M{"asset": contractHash.TransferredVal(), "tokenid": tokenId},
		//	Query:      []string{},
		//}, ret)

		r1, err := me.getNep11PropertiesByContract(contractHash.TransferredVal(), tokenId.Val())

		if err != nil {
			return err
		}
		filter, err := me.Filter(r1, filter)
		if err != nil {
			return err
		}
		r4 = append(r4, filter)
	}
	r5, err := me.FilterArrayAndAppendCount(r4, int64(len(r4)), filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r5)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}

func (me *T) getNep11PropertiesByContract(asset string, tokenid string) (map[string]interface{}, error) {
	h := &joh.T{}
	c, err := h.OpenConfigFile()
	if err != nil {
		log2.Fatalf("Open config file error:%s", err)
	}
	nodes := c.Proxy.URI
	re := make(map[string]interface{})

	for _, item := range nodes {
		re, err = me.getPropertiesByRPC(item, asset, tokenid)
		if err != nil {
			continue
		}
		break
	}

	res := re["result"].(map[string]interface{})
	state := res["state"]
	exception := res["exception"]

	if state != "HALT" || exception != nil {
		return nil, stderr.ErrFind
	}

	result := res["stack"].([]interface{})[0].(map[string]interface{})["value"].([]interface{})
	properties := make(map[string]interface{})
	for _, item := range result {
		it := item.(map[string]interface{})
		key, err := base64.StdEncoding.DecodeString(it["key"].(map[string]interface{})["value"].(string))
		if err != nil {
			return nil, err
		}

		value_key := it["value"].(map[string]interface{})["type"].(string)

		//aa := it["value"].(map[string]interface{})["value"].(string)
		//fmt.Printf(aa)

		var value_result interface{}
		if value_key == "ByteString" {
			value_pre := it["value"].(map[string]interface{})["value"].(string)

			value_decode, err := crypto.Base64Decode(value_pre)
			if err != nil {
				return nil, err
			}
			value_reverse := "0x" + helper.BytesToHex(helper.ReverseBytes(value_decode))

			if h160.T(value_reverse).Valid1() {
				value_result = value_reverse
			} else {
				value, err := base64.StdEncoding.DecodeString(value_pre)
				value_result = string(value)
				if err != nil {
					return nil, err
				}
			}

		} else if value_key == "Map" {
			valueArray := it["value"].(map[string]interface{})["value"].([]interface{})

			mapresult := make([]map[string]interface{}, 0)
			for _, items := range valueArray {
				it_map := make(map[string]interface{})
				item2map := items.(map[string]interface{})
				item2map_key := item2map["value"].(map[string]interface{})["type"].(string)
				var value2_result interface{}
				it_key, err := base64.StdEncoding.DecodeString(item2map["key"].(map[string]interface{})["value"].(string))
				if err != nil {
					return nil, err
				}

				if item2map_key == "ByteString" {

					value2_pre := item2map["value"].(map[string]interface{})["value"].(string)

					value2_decode, _ := crypto.Base64Decode(value2_pre)
					value2_reverse := "0x" + helper.BytesToHex(helper.ReverseBytes(value2_decode))
					//address,_:=helper.UInt160FromString(value2_reverse)

					if h160.T(value2_reverse).Valid() {
						value2_result = value2_reverse
					} else {
						value2, err := base64.StdEncoding.DecodeString(value2_pre)
						value2_result = string(value2)
						if err != nil {
							return nil, err
						}
					}

					it_map[string(it_key)] = value2_result
				} else {
					it_map[string(it_key)] = item2map["value"].(map[string]interface{})["value"]
				}

				mapresult = append(mapresult, it_map)
			}

			value_result = mapresult
		} else {
			value_result = it["value"].(map[string]interface{})["value"]
		}

		properties[string(key)] = value_result
	}

	return properties, nil
}

func (me *T) getPropertiesByRPC(url string, asset string, tokenid string) (map[string]interface{}, error) {

	para := `{
		"jsonrpc": "2.0",
			"id": 1,
			"method": "invokefunction",
			"params": ["` + asset + `",
					"properties",[
						{"type":"ByteArray","value":"` + tokenid + `"}
					]]}`

	jsonData := []byte(para)
	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", url, response.StatusCode)
	}
	resbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, stderr.ErrPrice
	}

	re := make(map[string]interface{})
	err = json.Unmarshal(resbody, &re)
	if err != nil {
		return nil, err
	}

	return re, nil
}
