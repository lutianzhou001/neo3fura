package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/Contract"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"os"
)

func (me *T) GetCollectionsByAsset(args struct {
	MarketHash h160.T
	Assets     []h160.T
	Filter     map[string]interface{}
	Raw        *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var list []interface{}
	for _, item := range args.Assets {
		if item.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		list = append(list, item)
	}
	rt := os.ExpandEnv("${RUNTIME}")
	var nns, genesis, polemen string
	if rt == "staging" {
		nns = Contract.Main_NNS.Val()
		//  metapanacea = Contract.Main_MetaPanacea.Val()
		genesis = Contract.Main_ILEXGENESIS.Val()
		polemen = Contract.Main_ILEXPOLEMEN.Val()

	} else if rt == "test2" {
		nns = Contract.Test_NNS.Val()
		//	metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Test_ILEXGENESIS.Val()
		polemen = Contract.Test_ILEXPOLEMEN.Val()
	} else {
		nns = Contract.Test_NNS.Val()
		//	metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Test_ILEXGENESIS.Val()
		polemen = Contract.Test_ILEXPOLEMEN.Val()
	}
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
						bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", nns}}, "then": "$tokenid",
							"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", genesis}}, "then": "$image",
								"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", polemen}}, "then": "$tokenid",
									"else": "$name"}}}}}}}},
						bson.M{"$group": bson.M{"_id": bson.M{"asset": "$asset", "class": "$class"}, "asset": bson.M{"$last": "$asset"}, "tokenid": bson.M{"$last": "$tokenid"}, "properties": bson.M{"$push": "$$ROOT"}}},
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
					count := 3
					if index < count {
						p := it.(map[string]interface{})["properties"]
						if p == nil {
							continue
						}
						pro := p.(primitive.A)[0]
						pitem := pro.(map[string]interface{})
						asset := pitem["asset"].(string)
						if pitem["tokenid"] == nil {
							count++
							continue
						}
						tokenid := pitem["tokenid"].(string)
						proMap["name"] = pitem["name"]
						proMap["asset"] = asset
						proMap["tokenid"] = tokenid
						tokenidArr = append(tokenidArr, tokenid)
						//获取属性
						if pitem["image"] != nil {
							proMap["image"] = ImagUrl(pitem["asset"].(string), pitem["image"].(string), "images")
						}
						if pitem["thumbnail"] != nil && pitem["thumbnail"] != "" {
							tb, err2 := base64.URLEncoding.DecodeString(pitem["thumbnail"].(string))
							if err2 != nil {
								return err2
							}
							ss := string(tb[:])
							if ss == "" {
								proMap["thumbnail"] = ImagUrl(proMap["asset"].(string), pitem["image"].(string), "thumbnail")
							} else {
								proMap["thumbnail"] = ImagUrl(proMap["asset"].(string), string(tb[:]), "thumbnail")
							}

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
									if key == "name" {
										proMap["name"] = value
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

						if proMap["name"] != nil && proMap["name"].(string) == "Nuanced Floral Symphony" {
							proMap["video"] = proMap["image"]
							delete(proMap, "image")
						}

						if proMap["image"] == "" && proMap["video"] == "" {
							count++
							continue
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
					tokeniditem["name"] = it["name"]
					tokeniditem["owner"] = it["owner"]
					tokeniditem["nns"] = it["nns"]
					if tokeniditem["name"] != nil && tokeniditem["name"].(string) == "Nuanced Floral Symphony" {
						tokeniditem["video"] = tokeniditem["image"]
						delete(tokeniditem, "image")
					}
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
