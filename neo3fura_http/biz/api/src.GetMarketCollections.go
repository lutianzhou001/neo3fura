package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math/big"
	"neo3fura_http/lib/joh"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"net/http"
	"strconv"
)

// 广告位NFT

func (me *T) GetMarketCollections(args struct {
	MarketHash h160.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	result, err := me.Client.QueryLastJob(struct {
		Collection string
	}{Collection: "MarketCollectionWhitelist"})
	if err != nil {
		return err
	}
	list := result["CollectionWhitelist"]
	//获取Collection基本信息
	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Asset",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"type": "NEP11", "hash": bson.M{"$in": list}}},
				bson.M{"$lookup": bson.M{
					"from": "SelfControlNep11Properties",
					"let":  bson.M{"asset": "$hash"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$asset", "$$asset"}}}},
						//bson.M{"$group": bson.M{"_id": "$asset","asset":bson.M{"$last":"$asset"}, "properities": bson.M{"$push": "$$ROOT"}}},
					},
					"as": "properties"},
				},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	//	collectioResult := make([]map[string]interface{},0)
	for _, item := range r1 {
		tokenidProperties := make([]map[string]interface{}, 0)
		hash := item["hash"].(string)
		tokenidArr := []string{}
		if item["properties"] != nil {
			properities := item["properties"].(primitive.A)
			if len(properities) > 0 {

				for index, it := range properities {
					proMap := make(map[string]interface{})
					if index < 3 {
						pitem := it.(map[string]interface{})
						asset := pitem["asset"].(string)
						tokenid := pitem["tokenid"].(string)

						proMap["asset"] = asset
						proMap["tokenid"] = tokenid
						tokenidArr = append(tokenidArr, tokenid)
						//获取属性
						if pitem["image"] != nil {
							proMap["image"] = ImagUrl(pitem["asset"].(string), pitem["image"].(string), "images")
						}
						if pitem["thumbnail"] != nil {
							tb, err2 := base64.URLEncoding.DecodeString(pitem["thumbnail"].(string))
							if err2 != nil {
								return err2
							}
							proMap["thumbnail"] = ImagUrl(pitem["asset"].(string), string(tb[:]), "thumbnail")

						}
						if proMap["image"] == nil {

							if pitem["tokenURI"] != nil {
								tokenUrl := pitem["tokenURI"].(string)
								ppjson, err := GetImgFromTokenURL(tokenurl(tokenUrl), asset, tokenid)
								if err != nil {
									return err
								}
								for key, value := range ppjson {
									//item[key] = value
									if key == "image" {
										img := value.(string)
										proMap["thumbnail"] = ImagUrl(asset, img, "thumbnail")
										proMap["image"] = ImagUrl(asset, img, "images")
									}
								}
							} else if pitem["properties"] != nil {
								//
								jsonData := make(map[string]interface{})
								properties := pitem["properties"].(string)
								if properties != "" {
									err := json.Unmarshal([]byte(properties), &jsonData)
									if err != nil {
										return err
									}

									image, ok := jsonData["image"]
									if ok {
										proMap["image"] = ImagUrl(pitem["asset"].(string), image.(string), "images")
									} else {
										proMap["image"] = ""
									}

									thumbnail, ok1 := jsonData["thumbnail"]
									if ok1 {
										tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
										if err2 != nil {
											return err2
										}
										proMap["thumbnail"] = ImagUrl(pitem["asset"].(string), string(tb[:]), "thumbnail")
									} else {
										if proMap["thumbnail"] == nil {
											if proMap["image"] != nil && proMap["image"] != "" {
												if image == nil {
													proMap["thumbnail"] = item["image"]
												} else {
													proMap["thumbnail"] = ImagUrl(pitem["asset"].(string), image.(string), "thumbnail")
												}
											}
										}
									}
								}

							} else {
								proMap["image"] = ""
								proMap["thumbnail"] = ""
							}
						}
						tokenidProperties = append(tokenidProperties, proMap)
					}

				}
			}
		}

		//获取nft 信息
		re := map[string]interface{}{}
		err := me.GetInfoByNFT(struct {
			Asset   h160.T
			Tokenid []string
			Filter  map[string]interface{}
			Raw     *map[string]interface{}
		}{Asset: h160.T(hash), Tokenid: tokenidArr, Raw: &re}, ret)

		if err != nil {
			return stderr.ErrGetNFTInfo
		}

		tokenidList := re["result"].([]map[string]interface{})
		for _, tokeniditem := range tokenidList {
			for _, it := range tokenidProperties {
				if tokeniditem["asset"] == it["asset"] && tokeniditem["tokenid"] == it["tokenid"] {
					tokeniditem["image"] = it["image"]
					tokeniditem["thumbnail"] = it["thumbnail"]
				}
			}

		}

		item["NFTList"] = tokenidList
		delete(item, "properties")

	}

	count := len(r1)
	r3, err := me.FilterAggragateAndAppendCount(r1, count, args.Filter)

	if err != nil {
		return err
	}
	r, err := json.Marshal(r3)
	if err != nil {
		return err
	}
	if args.Raw != nil {
		*args.Raw = r3
	}
	*ret = json.RawMessage(r)
	return nil
}

func (me *T) GetUserSavingsAmount(contract h160.T, user h160.T, asset h160.T) (*big.Int, error) {
	h := &joh.T{}
	c, err := h.OpenConfigFile()
	if err != nil {
		log2.Fatalf("Open config file error:%s", err)
	}
	nodes := c.Proxy.URI
	re := make(map[string]interface{})

	for _, item := range nodes {
		re, err = me.GetUserSavingsAmountByRPC(item, contract.Val(), user.Val(), asset.Val())
		if err != nil {
			continue
		}
		break
	}

	res := re["result"].(map[string]interface{})
	state := res["state"]
	exception := res["exception"]

	if state != "HALT" || exception != nil {
		return big.NewInt(0), stderr.ErrFind
	}

	result := res["stack"].([]interface{})[0].(map[string]interface{})["value"]
	savingAmount, err := strconv.ParseInt(result.(string), 10, 64)

	return big.NewInt(savingAmount), nil
}

func (me *T) GetUserSavingsAmountByRPC(node string, contract string, user string, asset string) (map[string]interface{}, error) {

	para := `{
		"jsonrpc": "2.0",
			"id": 1,
			"method": "invokefunction",
			"params": ["` + contract + `",
					"getUserSavingsAmount",[
						{"type":"Hash160","value":"` + user + `"},
						{"type":"Hash160","value":"` + asset + `"}
					]]}`

	jsonData := []byte(para)
	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(node, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", node, response.StatusCode)
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
