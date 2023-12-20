package Contract

import (
	"sort"
)

// T ...
type T string

const (
	Main_MetaPanacea     T = "0x19ed09dadac28e6b6a2f76588516ef681aff29b1"
	Test_MetaPanacea     T = "0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56"
	Main_ILEXPOLEMEN     T = "0x9f344fe24c963d70f5dcf0cfdeb536dc9c0acb3a"
	Test_ILEXPOLEMEN     T = "0xb13b57056775529e9461418a0a66b6dd97640ef8"
	Main_ILEXGENESIS     T = "0xc91b4becc7f4052a22e33990ed7696b4b175ec62"
	Test_ILEXGENESIS     T = "0x6a2893f97401e2b58b757f59d71238d91339856a"
	Main_NNS             T = "0x50ac1c37690cc2cfc594472833cf57505d5f46de"
	Test_NNS             T = "0x50ac1c37690cc2cfc594472833cf57505d5f46de"
	Main_TREE            T = "0x50ac1c37690cc2cfc594472833cf57505d5f46de"
	Test_TREE            T = "0xf6b4d6b3af093c15ff64cfc68a03faf31ad5ae92"
	Main_SecondaryMarket T = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	Main_PrimaryMarket   T = "0xa41600dec34741b143c66f2d3448d15c7d79a0b7"
	Test_SecondaryMarket T = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	Test_PrimaryMarket   T = "0x6f1ef5147a00ebbb7de1cf82420485674c5c55bc"
)

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

func (me T) In(str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, me.Val())
	if index < len(str_array) && str_array[index] == me.Val() { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}
