package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (me *T) GetRawMemPool(args struct {
	Filter map[string]interface{}
}, ret *json.RawMessage) error {
	rpcCli := me.Client.RpcCli
	ports := me.Client.RpcPorts
	results := make([][]string, len(ports))
	for i, port := range ports {
		u, err := url.Parse(port)
		if err != nil {
			return nil
		}
		rpcCli.Endpoint = u
		res := rpcCli.GetRawMemPool()
		if res.HasError() {
			return fmt.Errorf(res.GetErrorInfo())
		}
		results[i] = res.Result // tx hash strings
	}
	// intersect
	r, err := json.Marshal(intersect(results))
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}

// intersect gets the intersection of all arrays
func intersect(values [][]string) []string {
	m := make(map[string]int)
	for _, value := range values {
		for _, s := range value {
			m[s]++
		}
	}
	r := make([]string, 0, maxArraySize(values))
	for k, v := range m {
		if v >= len(values) { // means tx is in each mem pool
			r = append(r, k)
		}
	}
	return r
}

func maxArraySize(values [][]string) int {
	max := 0
	for _, v := range values {
		if len(v) > max {
			max = len(v)
		}
	}
	return max
}
