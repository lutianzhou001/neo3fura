package api

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetNFTByAssetSeries(args struct {
	Asset  h160.T
	Series string
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
	Raw    *map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit == 0 {
		args.Limit = 50

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
				bson.M{"$set": bson.M{"class": bson.M{"$ifNull": []interface{}{"$series", "$tokenid"}}}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x50ac1c37690cc2cfc594472833cf57505d5f46de"}}, "then": "$asset", "else": "$class"}}}},
				bson.M{"$match": bson.M{"class": args.Series}},
				bson.M{"$skip": args.Skip},
				bson.M{"$limit": args.Limit},
				//bson.M{"$lookup": bson.M{
				//	"from": "MarketNotification",
				//	"let":  bson.M{"asset": "$asset","tokenid":"$tokenid"},
				//"pipeline": []bson.M{	//
				//	bson.M{"$match":bson.M{"eventname":bson.M{"$in":[]interface{}{"Auction","OfferCollection","CompleteOfferCollection","Offer","CompleteOffer","Claim"}}} },
				//	bson.M{"$match":bson.M{"$expr":bson.M{"$and":[]interface{}{
				//		bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
				//		bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
				//	}}}},
				//	bson.M{"$sort":bson.M{"nonce":-1}},
				//	bson.M{"$group":bson.M{"_id":"$eventname","eventname":bson.M{"$last":"$eventname"},"market":bson.M{"$last":"$market"},"timestamp":bson.M{"$last":"$timestamp"},"extendData":bson.M{"$last":"$extendData"}}},
				//	bson.M{"$project": bson.M{"eventname":1,"market":1,"extendData":1,"timestamp":1}},
				//	//bson.M{"$limit":args.Limit},
				//},
				//"as": "notification"}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	for _, item := range r1 {
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		tokenidArr := []string{tokenid}

		//获取属性
		if item["image"] != nil {
			item["image"] = ImagUrl(item["asset"].(string), item["image"].(string), "images")
		}
		if item["thumbnail"] != nil {
			tb, err2 := base64.URLEncoding.DecodeString(item["thumbnail"].(string))
			if err2 != nil {
				return err2
			}
			item["thumbnail"] = ImagUrl(item["asset"].(string), string(tb[:]), "thumbnail")

		}
		if item["image"] == nil {

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
						item["thumbnail"] = ImagUrl(asset, img, "thumbnail")
						item["image"] = ImagUrl(asset, img, "images")
					}
				}
			} else if item["properties"] != nil {
				//
				jsonData := make(map[string]interface{})
				properties := item["properties"].(string)
				if properties != "" {
					err := json.Unmarshal([]byte(properties), &jsonData)
					if err != nil {
						return err
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
						item["thumbnail"] = ImagUrl(item["asset"].(string), string(tb[:]), "thumbnail")
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

			} else {
				item["image"] = ""
				item["thumbnail"] = ""
			}
		}

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
				bson.M{"$set": bson.M{"class": bson.M{"$ifNull": []interface{}{"$series", "$tokenid"}}}},
				bson.M{"$set": bson.M{"class": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$asset", "0x50ac1c37690cc2cfc594472833cf57505d5f46de"}}, "then": "$asset", "else": "$class"}}}},
				bson.M{"$match": bson.M{"class": args.Series}},
			},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}
	count := len(r2)
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
