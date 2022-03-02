# GetMarketWhiteList
gets the market white list
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| MarketHash     | string| the marketplace hash | required |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetMarketWhiteList",
  "params": {
      "MarketHash":"0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180"     
       },
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "market": "0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180",
    "whiteList": [
      "0x15130d478ec0baaee86f98a75310e431490c3441",
      "0xd2a4cff31913016155e38e474a2c06d08be276cf",
      "0x5af33f9a57d96255e2e705aa0e1942f51b658e2b",
      "0x48c40d4666f93408be1bef038b6722404d9a4c2a",
      "0x1415ab3b409a95555b77bc4ab6a7d9d7be0eddbd",
      "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f"
    ]
  },
  "error": null
}
```
