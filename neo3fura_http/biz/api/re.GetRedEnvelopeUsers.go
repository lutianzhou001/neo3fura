package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/sc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/Contract"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetRedEnvelopeUsers(args struct {
	Asset     h160.T
	StartTime uint64
	EndTime   uint64
	Filter    map[string]interface{}
	Raw       *[]map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Asset.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	if args.EndTime < args.StartTime {
		return stderr.ErrInvalidArgs
	}
	rt := os.ExpandEnv("${RUNTIME}")
	NetEndPoint := "http://seed2.neo.org:10332"
	nnsContract := Contract.Test_NNS
	switch rt {
	case "test":
		NetEndPoint = "http://seed2t5.neo.org:20332"
		nnsContract = Contract.Test_NNS
	case "test2":
		NetEndPoint = "http://seed2t5.neo.org:20332"
		nnsContract = Contract.Test_NNS
	case "staging":
		NetEndPoint = "http://seed2.neo.org:10332"
		nnsContract = Contract.Main_NNS
	default:
		log2.Fatalf("runtime environment mismatch")
	}

	flag, err := isExpiresNNS(NetEndPoint, nnsContract, "Y3J5cHRvem9tYmllLm5lbw==")
	fmt.Println(flag, err)
	r1, err := me.Client.QueryAggregate(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Pipeline   []bson.M
		Query      []string
	}{Collection: "Nep11TransferNotification",
		Index:  "GetRedEnvelopeUsers",
		Sort:   bson.M{},
		Filter: bson.M{},
		Pipeline: []bson.M{
			bson.M{"$match": bson.M{"contract": args.Asset,
				"$and": []interface{}{
					bson.M{"timestamp": bson.M{"$gte": args.StartTime}},
					bson.M{"timestamp": bson.M{"$lte": args.EndTime}},
				}}},

			bson.M{"$lookup": bson.M{
				"from": "Address-Asset",
				"let":  bson.M{"to": "$to"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"asset": nnsContract.Val(), "$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$address", "$$to"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "address": 1}}, //  bson.M{"$eq": []interface{}{"$address", "$$to"}},
				},
				"as": "nns"},
			},

			bson.M{"$group": bson.M{"_id": "$tokenId", "transferList": bson.M{"$push": "$$ROOT"}}},
			bson.M{"$sort": bson.M{"from": 1}},
		},
		Query: []string{},
	}, ret)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	for _, item := range r1 {
		var minter string
		var minterInfo []primitive.A
		transferList := item["transferList"].(primitive.A)
		if len(transferList) > 1 {
			for _, transfer := range transferList {
				transferItem := transfer.(map[string]interface{})
				if transferItem["from"] == nil { //mint
					minter = transferItem["to"].(string)
					if transferItem["nns"] != nil {
						nnsList := transferItem["nns"].(primitive.A)
						for _, nns := range nnsList {
							nnsItem := nns.(map[string]interface{})
							tokenid, err := base64.URLEncoding.DecodeString(nnsItem["tokenid"].(string))
							if err != nil {
								return fmt.Errorf("tokenid base64.URLEncoding.DecodeString error %s", err)
							}
							isExpired, _ := isExpiresNNS(NetEndPoint, nnsContract, string(tokenid))
							if isExpired {
								minterInfo = append(minterInfo, transferList)
								goto endfor
							}
						}
					}
				}

			}
		endfor:
		}
		result[minter] = minterInfo
	}

	//r2, err := me.FilterArrayAndAppendCount(result, 0, args.Filter)
	r2, err := me.Filter(result, args.Filter)
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

func isExpiresNNS(endPoint string, contract Contract.T, nns string) (bool, error) {
	client := rpc.NewClient(endPoint)
	sb := sc.NewScriptBuilder()
	sh, err := helper.UInt160FromString(contract.Val())
	if err != nil {
		return false, err
	}
	var arg = []interface{}{nns}
	sb.EmitDynamicCall(sh, "ownerOf", arg)
	script, err := sb.ToArray()
	if err != nil {
		return false, err
	}

	response := client.InvokeScript(crypto.Base64Encode(script), nil)
	if response.Result.State == "HALT" {
		return true, nil
	} else {
		return false, nil
	}
	return false, nil
}
