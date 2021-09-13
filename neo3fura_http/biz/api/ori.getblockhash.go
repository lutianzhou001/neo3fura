package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"
)

func (me *T) Getblockhash(args []interface{}, ret *json.RawMessage) error {
	var raw1 map[string]interface{}
	if len(args) == 0 {
		return stderr.ErrInvalidArgs
	}
	switch args[0].(type) {
	case float64:
		err := me.GetBlockHashByBlockHeight(struct {
			BlockHeight uintval.T
			Filter      map[string]interface{}
			Raw         *map[string]interface{}
		}{BlockHeight: uintval.T(uint64(args[0].(float64))), Raw: &raw1}, ret)
		if err != nil {
			return err
		}
		r, err := json.Marshal(raw1["hash"])
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}
