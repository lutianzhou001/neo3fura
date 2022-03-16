# GetAllBidInfoByNFT
get nft historical bid info 
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string|  the asset hash | required |
| MarketHash     | string| the marketplace hash | optional |
| TokenId     | string| nft token id | optional |




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
      "MarketHash":["0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29","0xa41600dec34741b143c66f2d3448d15c7d79a0b7"]
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

