package api

import (
	"encoding/json"
)

func (me *T) Getcommittee(args []interface{}, ret *json.RawMessage) error {
	var raw []map[string]interface{}
	err := me.GetCommittee(struct {
		Filter map[string]interface{}
		Limit  int64
		Skip   int64
		Raw    *[]map[string]interface{}
	}{Raw: &raw}, ret)
	if err != nil {
		return err
	}

	var res []string
	for _, item := range raw {
		res = append(res, item["candidate"].(string))
	}
	r, err := json.Marshal(res)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
