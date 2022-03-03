# GetNFTRecordByAddress
get nft record by user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  the user's address| required|
| MarketHash     | string| the marketplace hash | optional |
| Skip | int | the number of items to return | optional |
| Limit | int | the number of items to return | optional |






#### Example
```powershell
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
```json
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
###  
