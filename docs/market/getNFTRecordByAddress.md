# GetNFTRecordByAddress
get nft record by user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  the user's address| required|
| MarketHash     | string|  | |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTRecordByAddress",
  "params": {
      "Address":"0xf0a33d62f32528c25e68951286f238ad24e30032",
      "MarketHash": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2"    
  },
  "id": 1
} '
```
### Response
```json5
{
  "id": 1,
  "result": {
    "result": [
      {
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": 10,
        "auctionAsset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "event": "Claim",
        "from": "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e",
        "image": "",
        "name": "sell-1",
        "nonce": 5,
        "state": "sell_sold",
        "timestamp": 1639392588921,
        "to": "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e",
        "tokenid": "b7mzAd/hhpBYX95Gq8eJwkoZdS9JssMHHhJztAQNCKs=",
        "user": "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e"
      },
      .....
    ],
    "totalCount": 7
  },
  "error": null
}
```
### Response Analyse
```
Status Condition: 
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
```
