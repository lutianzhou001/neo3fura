package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

// 广告位NFT

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