package utils

import (
	"strconv"
)

// 自定义map分组
func GroupBy(maps []map[string]interface{}, key string) map[string][]map[string]interface{} {
	groups := make(map[string][]map[string]interface{})
	for _, m := range maps {
		k := m[key].(string) // XXX: will panic if m[key] is not a string.
		groups[k] = append(groups[k], m)
	}
	return groups
}

// nft  根据Asset 分类
func GroupByAsset(maps []map[string]interface{}) map[string][]string {
	groups := make(map[string][]string)

	for i := 0; i < len(maps); i++ {
		k := maps[i]["asset"].(string)
		v := maps[i]["tokenid"].(string)
		tokenids := groups[k]
		tokenids = append(tokenids, v)
		groups[k] = tokenids
	}
	//for _, m := range maps {
	//	k := m["asset"].(string)
	//	v := m["tokenid"].(string)
	//
	//	tokenids := groups[k]
	//	tokenids = append(tokenids, v)
	//
	//	groups[k] = append(groups[k], tokenids)
	//}
	return groups
}

func GroupByString(maps []map[string]interface{}, key string) map[string][]map[string]interface{} {
	groups := make(map[string][]map[string]interface{})
	for _, m := range maps {

		v := m[key].(int64)

		k := strconv.FormatInt(v, 10)

		groups[k] = append(groups[k], m)
	}
	return groups
}
