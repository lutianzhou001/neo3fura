package api

import (
	"encoding/json"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"strconv"
	"strings"
)

func (me *T) GetNNSNameByAdmin(args struct {
	Asset  h160.T
	Admin  h160.T
	Limit  int64
	Skip   int64
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Admin.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	adminstr := string(args.Admin)
	little_endian := helper.HexToBytes(adminstr[2:len(adminstr)])

	rea := helper.ReverseBytes(little_endian)
	encodeAdmin := crypto.Base64Encode(rea)
	var arr = []interface{}{}

	if strings.Index(encodeAdmin, "+") >= 0 {
		splitArr := strings.Split(encodeAdmin, "+")
		strlen := len(splitArr)

		before := splitArr[0] + "+"
		arr = append(arr, bson.M{"properties": bson.M{"$regex": "admin\": \"" + before, "$options": "$i"}})
		for i := 1; i < strlen; i++ {
			arr = append(arr, bson.M{"properties": bson.M{"$regex": splitArr[i], "$options": "$i"}})
		}

	} else {
		arr = append(arr, bson.M{"properties": bson.M{"$regex": "admin\": \"" + encodeAdmin, "$options": "$i"}})
		//arr = append(arr, bson.M{"properties": bson.M{"$regex": "admin\":\"" + encodeAdmin, "$options": "$i"}})
	}

	pipe := []bson.M{}
	pipe = append(pipe, bson.M{"$match": bson.M{"asset": args.Asset}})
	pipe = append(pipe, bson.M{"$match": bson.M{"$and": arr}})
	//pipe = append(pipe, bson.M{"$skip": args.Skip})
	//pipe = append(pipe, bson.M{"$limit": args.Limit})

	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Nep11Properties",
			Index:      "GetNFTByWords",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipe,
			Query:      []string{},
		}, ret)

	if err != nil {
		return err
	}

	for _, item := range r1 {
		//获取nft 属性
		asset := item["asset"].(string)
		tokenid := item["tokenid"].(string)
		if item["properties"] != nil {
			extendData := item["properties"].(string)
			if extendData != "" {
				var data map[string]interface{}
				if err2 := json.Unmarshal([]byte(extendData), &data); err2 == nil {

					tokenuri, ok := data["tokenURI"]
					if ok {
						ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)
						if err != nil {
							return err
						}
						for key, value := range ppjson {
							item[key] = value
						}
					}
					if item["name"] == "" || item["name"] == nil {
						name, ok := data["name"]
						if ok {
							item["name"] = name
						}
					}
					admin, ok2 := data["admin"]
					if ok2 {
						item["admin"] = admin
					}
					expiration, ok3 := data["expiration"]
					if ok3 {
						time, err := strconv.ParseInt(expiration.(string), 10, 64)
						if err != nil {
							return err
						}
						item["expiration"] = time
					}

				} else {
					return err2
				}

			}
		}

		delete(item, "properties")
		delete(item, "_id")
		delete(item, "admin")
		delete(item, "asset")
		delete(item, "tokenid")
	}

	r1 = mapsort.MapSort2(r1, "expiration")

	if args.Limit > 0 {
		pagedName := make([]map[string]interface{}, 0)
		for i, item := range r1 {
			if int64(i) < args.Skip {
				continue
			} else if int64(i) > args.Skip+args.Limit-1 {
				continue
			} else {
				pagedName = append(pagedName, item)
			}
		}
		r, err := json.Marshal(pagedName)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)

	} else {
		r, err := json.Marshal(r1)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
	}

	return nil
}
