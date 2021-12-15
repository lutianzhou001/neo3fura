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

func GroupByString(maps []map[string]interface{}, key string) map[string][]map[string]interface{} {
	groups := make(map[string][]map[string]interface{})
	for _, m := range maps {

		v := m[key].(int64)

		k := strconv.FormatInt(v, 10)

		groups[k] = append(groups[k], m)
	}
	return groups
}
