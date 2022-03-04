# GetWhiteListByMarketHash
Gets the white list and details by the market hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| MarketHash     | string| The marketplace hash | Required |

### Example

Request body

```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetWhiteListByMarketHash",
  "params": {
      "MarketHash":"0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180"
     
       },
  "id": 1
}'
```
Response body

```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "asset": "0x48c40d4666f93408be1bef038b6722404d9a4c2a",
        "decimal": 8,
        "feeRate": "0",
        "rewardRate": "0",
        "rewardReceiveAddress": "0x481b5f71a738d3d43c5a9b621f93aa00f2a5acfd",
        "symbol": "bNEO",
        "tokenname": "BurgerNEO",
        "totalsupply": "2118188748700",
        "type": "NEP17"
      },
      ....
    ],
    "totalCount": 4
  },
  "error": null
}

```
