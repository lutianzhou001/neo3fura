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
	"strings"
)

func (me *T) GetNFTByAssetClass(args struct {
	Asset     h160.T
	Class     string
	ClassName string
	Limit     int64
	Skip      int64
	Filter    map[string]interface{}
	Raw       *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 50

	}
	rt := os.ExpandEnv("${RUNTIME}")

	var primaryMarket, nns, polemen, genesis string
	if rt == "staging" {
		nns = Contract.Main_NNS.Val()
		//metapanacea = Contract.Main_MetaPanacea.Val()
		genesis = Contract.Main_ILEXGENESIS.Val()
		polemen = Contract.Main_ILEXPOLEMEN.Val()
		primaryMarket = Contract.Main_PrimaryMarket.Val()

	} else if rt == "test2" {
		nns = Contract.Test_NNS.Val()
		//metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Test_ILEXGENESIS.Val()
		polemen = Contract.Test_ILEXPOLEMEN.Val()
		primaryMarket = Contract.Test_PrimaryMarket.Val()
	} else {
		nns = Contract.Main_NNS.Val()
		//metapanacea = Contract.Test_MetaPanacea.Val()
		genesis = Contract.Main_ILEXGENESIS.Val()
		polemen = Contract.Main_ILEXPOLEMEN.Val()
		primaryMarket = Contract.Test_PrimaryMarket.Val()
	}

	r1, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "SelfControlNep11Properties",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"asset": args.Asset}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", nns}}, "then": "$asset",
					"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", genesis}}, "then": "$image",
						"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", polemen}}, "then": "$tokenid",
							"else": "$name"}}}}}}}},
				bson.M{"$match": bson.M{"class": args.ClassName}},
				bson.M{"$lookup": bson.M{
					"from": "Market",
					"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"amount": bson.M{"$gt": 0}, "market": bson.M{"$ne": primaryMarket}}},
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						//bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
					},
					"as": "market"},
				},
				bson.M{"$sort": bson.M{"tokenid": 1}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		tokenidArr := []string{tokenid}
		if item["market"] == nil || len(item["market"].(primitive.A)) == 0 {
			continue
		}
		//if item["image"] == nil {
		if item["properties"] != nil { //
			jsonData := make(map[string]interface{})
			properties := item["properties"].(string)
			if properties != "" {
				err := json.Unmarshal([]byte(properties), &jsonData)
				if err != nil {
					return err
				}

				tokenURI, ok := jsonData["tokenURI"]
				if ok {
					ppjson, err := GetImgFromTokenURL(tokenurl(tokenURI.(string)), asset, tokenid)
					if err != nil {
						return err
					}
					for key, value := range ppjson {
						item[key] = value
						if key == "image" {
							img := value.(string)
							thumbnail := ImagUrl(asset, img, "thumbnail")
							flag := strings.HasSuffix(thumbnail, ".mp4")
							if flag {
								thumbnail = strings.Replace(thumbnail, ".mp4", "mp4", -1)
							}
							item["thumbnail"] = thumbnail
							item["image"] = ImagUrl(asset, img, "images")
						}
						if key == "name" {
							item["name"] = value
						}

					}

				}

				image, ok := jsonData["image"]
				if ok {
					item["image"] = ImagUrl(item["asset"].(string), image.(string), "images")
				} else {
					item["image"] = ""
				}

				thumbnail, ok1 := jsonData["thumbnail"]
				if ok1 {
					tb, err2 := base64.URLEncoding.DecodeString(thumbnail.(string))
					if err2 != nil {
						return err2
					}
					ss := string(tb[:])
					if ss == "" {
						item["thumbnail"] = ImagUrl(item["asset"].(string), item["image"].(string), "thumbnail")
					} else {
						item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
					}

				} else {
					if item["thumbnail"] == nil {
						if item["image"] != nil && item["image"] != "" {
							if image == nil {
								item["thumbnail"] = item["image"]
							} else {
								item["thumbnail"] = ImagUrl(item["asset"].(string), image.(string), "thumbnail")
							}
						}
					}
				}
			}

		}
		if item["tokenURI"] != nil {
			tokenUrl := item["tokenURI"].(string)
			ppjson, err := GetImgFromTokenURL(tokenurl(tokenUrl), asset, tokenid)
			if err != nil {
				return err
			}
			for key, value := range ppjson {
				//item[key] = value
				if key == "image" {
					img := value.(string)
					thumbnail := ImagUrl(asset, img, "thumbnail")
					flag := strings.HasSuffix(thumbnail, ".mp4")
					if flag {
						thumbnail = strings.Replace(thumbnail, ".mp4", "mp4", -1)
					}
					item["thumbnail"] = thumbnail
					item["image"] = ImagUrl(asset, img, "images")
				}
				if key == "name" {
					item["name"] = value
				}

			}
		}

		if item["name"] != nil && item["name"].(string) == "Nuanced Floral Symphony" {
			item["video"] = item["image"]
			delete(item, "image")
		}
		//}

		re := map[string]interface{}{}
		err := me.GetInfoByNFT(struct {
			Asset   h160.T
			Tokenid []string
			Filter  map[string]interface{}
			Raw     *map[string]interface{}
		}{Asset: h160.T(asset), Tokenid: tokenidArr, Raw: &re}, ret)

		if err != nil {
			return stderr.ErrGetNFTInfo
		}

		marketInfo := re["result"]
		if marketInfo != nil {
			marketItem := marketInfo.([]map[string]interface{})
			info := marketItem[0]
			for key, value := range info {
				item[key] = value
			}
			delete(item, "properties")
		}

		result = append(result, item)
	}
	// totalcount
	r2, err := me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "SelfControlNep11Properties",
			Index:      "GetAssetInfo",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline: []bson.M{
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", nns}}, "then": "$asset",
					"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", genesis}}, "then": "$image",
						"else": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", polemen}}, "then": "$tokenid",
							"else": "$name"}}}}}}}},
				bson.M{"$match": bson.M{"class": args.ClassName}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	count := len(r2)
	r3, err := me.FilterAggragateAndAppendCount(result, count, args.Filter)

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
