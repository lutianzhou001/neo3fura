package api

import (
	"encoding/json"
	"fmt"
	"neo3fura_http/lib/type/strval"
	"net/http"
)

func (me *T) GetNFSImgStatus(args struct {
	Url    strval.T
	Filter map[string]interface{}
}, ret *json.RawMessage) error {

	result := make(map[string]interface{})

	resp, err := http.Get(args.Url.Val())
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	if err != nil {
		fmt.Println("Search imageId error")

	}
	if resp.StatusCode == 200 {
		result["ImageStatus"] = true
	} else {
		result["ImageStatus"] = false
	}
	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
