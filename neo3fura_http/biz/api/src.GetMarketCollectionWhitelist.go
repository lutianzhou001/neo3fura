package api

import (
	"encoding/json"
	"neo3fura_http/var/stderr"
)

func (me *T) GetMarketCollectionWhitelist(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	//var hashArr []interface{}

	result, err := me.Client.QueryLastJob(struct {
		Collection string
	}{Collection: "MarketCollectionWhitelist"})
	if err != nil {
		return err
	}
	list := result["CollectionWhitelist"]

	r, err := json.Marshal(list)
	if err != nil {
		return stderr.ErrInsertDocument
	}
	*ret = json.RawMessage(r)
	return nil
}
