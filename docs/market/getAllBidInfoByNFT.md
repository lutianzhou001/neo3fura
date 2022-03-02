# GetAllBidInfoByNFT
get nft historical bid info 
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string|  the asset hash | |
| MarketHash     | string|  | |
| TokenId     | string| nft token id | |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAllBidInfoByNFT",
  "params": {      
      "AssetHash": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
      "TokenId":"b7mzAd/hhpBYX95Gq8eJwkoZdS9JssMHHhJztAQNCKs=",
      "MarketHash":""
  },
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "result": [
      {
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "bidAmount": [
          27,
          20
        ],
        "bidder": [
          "0xc65e19cfa66b61800ce582d1b55f4e93fa214b17",
          "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e"
        ],
        "nonce": 11,
        "tokenid": "b7mzAd/hhpBYX95Gq8eJwkoZdS9JssMHHhJztAQNCKs="
      },
      {
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "bidAmount": [
          20,
          15,
          10
        ],
        "bidder": [
          "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e",
          "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e",
          "0xf0a33d62f32528c25e68951286f238ad24e30032"
        ],
        "nonce": 9,
        "tokenid": "b7mzAd/hhpBYX95Gq8eJwkoZdS9JssMHHhJztAQNCKs="
      }
    ],
    "totalCount": 2
  },
  "error": null
}
```

