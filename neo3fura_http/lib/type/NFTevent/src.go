package NFTevent

import (
	"sort"
)

// T ...
type T string

const (
	Cancel T = "cancel" //卖家下架

	//直买直卖
	Sell_Listed  T = "sell_listed"  //卖家上架
	Sell_Expired T = "sell_expired" //卖家上架过期
	Sell_Sold    T = "sell_sold"    //卖家售出
	Sell_Buy     T = "sell_buy"     //买家购买

	//拍卖
	Auction_Listed  T = "auction_listed"  //卖家拍卖
	Auction_Expired T = "auction_expired" //卖家拍卖过期
	Aucion_Deal     T = "auction_deal"    //卖家成交

	Auction_Bid      T = "auction_bid"      //买家出价
	Auction_Return   T = "auction_return"   //买家出价退回
	Auction_Bid_Deal T = "auction_bid_deal" //买家竞价成功
	Auction_Withdraw T = "auction_withdraw" //买家领取

	//TRANSFER
	Send    T = "send"    //发送
	Receive T = "receive" //接收
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
