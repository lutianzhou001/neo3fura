package api

import (
	"encoding/json"
	"neo3fura/lib/type/uintval"
	"neo3fura/var/stderr"
)

func (me *T) Getblockhash(args []interface{}, ret *json.RawMessage) error {
	var raw1 map[string]interface{}
	switch args[0].(type) {
	case float64:
		err := me.GetBlockHashByBlockHeight(struct {
			BlockHeight uintval.T
			Filter      map[string]interface{}
		}{BlockHeight: uintval.T(uint64(args[0].(float64)))}, ret)
		if err != nil {
			return err
		}
		r, err := json.Marshal(raw1)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}
