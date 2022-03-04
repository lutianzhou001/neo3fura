# GetNFTRecordByAddress
Gets the NFT record by the user's address.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  The user's address| Required|
| MarketHash     | string| The marketplace hash | Optional |
| Skip | int | The number of items to return | Optional |
| Limit | int | The number of items to return | Optional |



### Example

Request body

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
Response body

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
