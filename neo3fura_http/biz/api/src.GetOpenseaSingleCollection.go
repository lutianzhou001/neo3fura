package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (me *T) GetOpenseaSingleCollection(args struct {
	CollectionSlug string
	ApiKey         string
	Filter         map[string]interface{}
}, ret *json.RawMessage) error {

	var requestGetURLNoParams = "https://api.opensea.io/api/v1/collection/" + args.CollectionSlug

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestGetURLNoParams, nil)
	if err != nil {
		panic(err)
		return err
	}
	req.Header.Set("X-API-KEY", args.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	re := make(map[string]interface{})
	err = json.Unmarshal(resbody, &re)
	if err != nil {
		return err
	}

	r, err := json.Marshal(re)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
