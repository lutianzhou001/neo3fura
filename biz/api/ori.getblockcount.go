package api

import (
	"encoding/json"
)

func (me *T) Getblockcount(args []interface{}, ret *json.RawMessage) error {
	var raw1 map[string]interface{}
	err := me.GetBlockCount(
		struct {
			Filter map[string]interface{}
			Raw    *map[string]interface{}
		}{Raw: &raw1}, ret)
	if err != nil {
		return err
	}
	r, err := json.Marshal(raw1["index"])
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
