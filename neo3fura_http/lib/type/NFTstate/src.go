package NFTstate

import (
	"sort"
)

// T ...
type T string

const (
	Auction   T = "auction"
	Sale      T = "sale"
	NotListed T = "notlisted"
	Unclaimed T = "unclaimed"
	Expired   T = "expired"
	//Collection
	BuyNow     T = "Buy Now"       //优先级 1
	CurrentBid T = "Current Bid"   //优先级 1
	LastSold   T = "Last Sold"     //优先级 2
	Offer      T = "Highest Offer" //优先级 2

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
