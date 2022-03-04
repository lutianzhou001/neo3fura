# GetAssetsHeldByAddress
Gets the information of assets owned by the user's address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address    | string|  The user's address| Required|
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetsHeldByAddress",
  "params": {"Address":"0xeba621d37ff117d9ce73c1579bf260aa779cb392"},
  "id": 1
}'
```

Response body

```json5

{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614bef0ea14111843551a804",
        "address": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "asset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "balance": "5099945401484000",
        "tokenid": ""
      },
      {
        "_id": "614bef0ea14111843551a802",
        "address": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "asset": "0xef4073a0f2b305a38ec4050e4d3d28bc40ea63f5",
        "balance": "79000000",
        "tokenid": ""
      }
    ],
    "totalCount": 2
  },
  "error": null
}
```
