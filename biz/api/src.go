package api

import (
	"fmt"
	"neo3fura/biz/data"
)

// T ...
type T struct {
	Data *data.T
}

// Ping ...
func (me *T) Ping(args struct{}, ret *string) error {
	*ret = "pong"
	return nil
}

func (me *T) Filter(data map[string]interface{}, filter map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if filter == nil {
		return data, nil
	}
	if len(filter) == 0 {
		return data, nil
	}
	for k, _ := range filter {
		if data[k] != nil {
			switch data[k].(type) {
			case map[string]interface{}:
				r, err := me.Filter(data[k].(map[string]interface{}), filter[k].(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				res[k] = r
			default:
				res[k] = data[k]
			}
		}
	}
	return res, nil
}

func (me *T) Deduplicate(data []map[string]interface{}) ([]map[string]interface{}, error) {
	set := make(map[string]map[string]interface{})
	for _, v := range data {
		con := fmt.Sprintf("%v", v["contract"]) + fmt.Sprintf("%v", v["tokenId"])
		set[con] = v
	}
	res := make([]map[string]interface{}, 0)
	for _, v := range set {
		res = append(res, v)
	}
	return res, nil
}

func (me *T) FilterArrayAndAppendCount(data []map[string]interface{}, count int64, filter map[string]interface{}) (map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	for _, item := range data {
		r, err := me.Filter(item, filter)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	res2 := make(map[string]interface{})
	res2["totalCount"] = count
	res2["result"] = res
	return res2, nil
}

func (me *T) FilterAggragateAndAppendCount(data []map[string]interface{}, count interface{}, filter map[string]interface{}) (map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	for _, item := range data {
	   r, err := me.Filter(item, filter)
	   if err != nil {
		  return nil, err
	   }
	   res = append(res, r)
	}
	res2 := make(map[string]interface{})
	res2["totalCount"] = count
	res2["result"] = res
	return res2, nil
 }
