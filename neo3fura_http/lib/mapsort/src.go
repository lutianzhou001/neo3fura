package mapsort

import (
	"sort"
)

type MapsSort struct {
	Key     string
	MapList []map[string]interface {
	}
}

// Len 为集合内元素的总数
func (m *MapsSort) Len() int {
	return len(m.MapList)
}

//如果index为i的元素小于index为j的元素，则返回true，否则返回false
func (m *MapsSort) Less(i, j int) bool {
	return m.MapList[i][m.Key].(int64) > m.MapList[j][m.Key].(int64)
}

//Swap 交换索引为 i 和 j 的元素
func (m *MapsSort) Swap(i, j int) {
	m.MapList[i], m.MapList[j] = m.MapList[j], m.MapList[i]
}

func MapSort(ms []map[string]interface {
}, key string) []map[string]interface {
} {
	mapsSort := MapsSort{}
	mapsSort.Key = key
	mapsSort.MapList = ms
	sort.Sort(&mapsSort)
	return mapsSort.MapList
}
