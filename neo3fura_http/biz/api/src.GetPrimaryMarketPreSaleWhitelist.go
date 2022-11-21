package api

import (
	"encoding/json"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) GetPrimaryMarketPreSaleWhitelist(args struct {
	Filter     map[string]interface{}
	MarketHash h160.T
}, ret *json.RawMessage) error {
	if args.MarketHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	result, err := me.Client.QueryLastJob(struct{ Collection string }{Collection: "PrimaryMarketPreSaleWhitelist"})
	if err != nil {
		return err
	}
	//result := make(map[string]interface{})
	market := result["market"].(string)
	if market != args.MarketHash.Val() {
		return stderr.ErrNotFound
	}
	r, err := json.Marshal(result)
	if err != nil {
		return stderr.ErrInsertDocument
	}
	*ret = json.RawMessage(r)
	return nil
}
