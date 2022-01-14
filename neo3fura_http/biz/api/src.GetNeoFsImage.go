package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (me *T) GetNeoFsImage(args struct{ImageId string}, ret *json.RawMessage) error {
	result := make(map[string]interface{})
	fmt.Println(me.Client.NeoFs+args.ImageId)
	resp, err := http.Get("http://"+me.Client.NeoFs+args.ImageId)
	fmt.Println(me.Client.NeoFs+args.ImageId)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Status)
	if err != nil {
		fmt.Println("Search imageId error")

	}
	if resp.StatusCode != 200  {
		result["ImageUrl"] = "Image object not found"
	} else {
		result["ImageUrl"] = me.Client.NeoFs+args.ImageId
	}
	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
