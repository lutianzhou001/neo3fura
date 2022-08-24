package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (me *T) GetUserInfoTwitter(args struct {
	AccessToken string

	Filter map[string]interface{}
	Raw    *map[string]interface{}
}, ret *json.RawMessage) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		//log.Errorf("make request error:%v", err)
		return err
	}
	var bearer = "Bearer " + args.AccessToken
	req.Header.Add("Authorization", bearer)

	fmt.Println("GetUserInfoTwitter para:", req)
	resp, err := client.Do(req)
	if err != nil {
		//log.Errorf("send request error:%v", err)
		return err
	}
	defer resp.Body.Close()
	reader := resp.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("GetUserInfoTwitter err:", err)
		return err
	}

	fmt.Println("body: ", string(body))
	var data map[string]interface{}
	if err1 := json.Unmarshal(body, &data); err1 != nil {
		return err
	}
	username := ""
	if data["data"] != nil {
		uname := data["data"].(map[string]interface{})["username"]

		if uname != nil {
			username = data["data"].(map[string]interface{})["username"].(string)
		}
	}

	fmt.Println(username)

	r2, err := me.Filter(data, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil

}
