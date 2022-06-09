package api

import (
	"encoding/json"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetVoteRecordByProjectId(args struct {
	ContractHash h160.T
	ProjectId    string
	Limit        int64
	Skip         int64
	Filter       map[string]interface{}
}, ret *json.RawMessage) error {
	if args.Limit == 0 {
		args.Limit = 512
	}
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	r1, count, err := me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Notification",
		Index:      "GetVoteRecordByProjectId",
		Sort:       bson.M{},
		Filter:     bson.M{"contract": args.ContractHash.Val(), "eventname": "Vote", "state.value.value": args.ProjectId},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {
		it := make(map[string]interface{})
		state := item["state"].(map[string]interface{})["value"].(primitive.A)
		projectId := state[0].(map[string]interface{})["value"]
		votor := state[1].(map[string]interface{})["value"]
		votes := state[2].(map[string]interface{})["value"]

		votor_decode, err := crypto.Base64Decode(votor.(string))
		if err != nil {
			return err
		}
		votor_reverse := "0x" + helper.BytesToHex(helper.ReverseBytes(votor_decode))

		it["project_id"] = projectId
		it["voter_address"] = votor_reverse
		it["tx_hash"] = item["txid"]
		it["votes"] = votes
		it["timestamp"] = item["timestamp"]

		r2, _, err := me.Client.QueryAll(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
			Limit      int64
			Skip       int64
		}{
			Collection: "Notification",
			Index:      "GetVoteRecordByProjectId",
			Sort:       bson.M{},
			Filter:     bson.M{"contract": args.ContractHash.Val(), "eventname": "OnTransfer", "txid": item["txid"]},
			Query:      []string{},
			Limit:      args.Limit,
			Skip:       args.Skip,
		}, ret)
		if err != nil {
			return err
		}

		state2 := r2[0]["state"].(map[string]interface{})["value"].(primitive.A)
		//token := state2[0].(map[string]interface{})["value"]
		from := state2[1].(map[string]interface{})["value"]
		//to := state2[2].(map[string]interface{})["value"]
		amount := state2[3].(map[string]interface{})["value"]
		it["value"] = ""
		if from == votor {
			it["value"] = amount
		}

		result = append(result, it)

	}

	r2, err := me.FilterArrayAndAppendCount(result, count, args.Filter)
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
