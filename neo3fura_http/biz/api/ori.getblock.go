package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h256"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"
)

func (me *T) Getblock(args []interface{}, ret *json.RawMessage) error {
	if len(args) <= 1 {
		return stderr.ErrInvalidArgs
	}
	if args[1] != true {
		return stderr.ErrInvalidArgs
	}
	switch args[0].(type) {
	case string:
		blockHash := h256.T(args[0].(string))
		if blockHash.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		var raw1 map[string]interface{}
		var raw2 []map[string]interface{}
		err := me.GetBlockByBlockHash(struct {
			BlockHash h256.T
			Filter    map[string]interface{}
			Raw       *map[string]interface{}
		}{
			BlockHash: blockHash,
			Filter:    nil,
			Raw:       &raw1,
		}, ret)
		if err != nil {
			return err
		}
		err = me.GetRawTransactionByBlockHash(struct {
			BlockHash h256.T
			Limit     int64
			Skip      int64
			Filter    map[string]interface{}
			Raw       *[]map[string]interface{}
		}{BlockHash: blockHash, Raw: &raw2}, ret)
		if err != nil {
			return err
		}
		raw1["tx"] = raw2
		r, err := json.Marshal(raw1)
		if err != nil {
			return err
		}
		*ret = json.RawMessage(r)
		return nil
	case float64:
		blockHeight := uintval.T(uint64(args[0].(float64)))
		if blockHeight.Valid() == false {
			return stderr.ErrInvalidArgs
		}
		var raw1 map[string]interface{}
		var raw2 []map[string]interface{}
		err := me.GetBlockByBlockHeight(struct {
			BlockHeight uintval.T
			Filter      map[string]interface{}
			Raw         *map[string]interface{}
		}{
			BlockHeight: blockHeight,
			Filter:      nil,
			Raw:         &raw1,
		}, ret)
		if err != nil {
			return err
		}
		err = me.GetRawTransactionByBlockHeight(struct {
			BlockHeight uintval.T
			Limit       int64
			Skip        int64
			Filter      map[string]interface{}
			Raw         *[]map[string]interface{}
		}{BlockHeight: blockHeight, Raw: &raw2}, ret)
		if err != nil {
			return err
		}
		raw1["tx"] = raw2
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
