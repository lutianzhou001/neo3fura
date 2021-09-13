package api

import (
	"encoding/json"
)

func (me *T) Getbestblockhash(args []interface{}, ret *json.RawMessage) error {
	var raw1 map[string]interface{}
	err := me.GetBestBlockHash(struct {
		Filter map[string]interface{}
		Raw    *map[string]interface{}
	}{Filter: nil, Raw: &raw1}, ret)
	if err != nil {
		return err
	}
	r, err := json.Marshal(raw1["hash"])
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
