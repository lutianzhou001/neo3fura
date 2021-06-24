package strval

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// T ...
type T string

// Valid ...
func (me T) Valid() bool {
	return true
}

// Val ...
func (me T) Val() string {
	return string(me)
}

// Bytes ...
func (me T) Bytes() []byte {
	return []byte(me.Val())
}

// Hex ...
func (me T) Hex() string {
	return hex.EncodeToString(me.Bytes())
}

// H256 ...
func (me T) H256() string {
	data := me.Bytes()
	l1 := sha256.Sum256(data)
	l2 := sha256.Sum256(l1[:])
	return hex.EncodeToString(l2[:])
}

func (me T) In(str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, me.Val())
	if index < len(str_array) && str_array[index] == me.Val() { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}
