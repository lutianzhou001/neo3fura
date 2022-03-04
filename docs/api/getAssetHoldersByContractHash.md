# GetAssetHoldersByContractHash
Gets the asset holders with the contract script hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required|
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetHoldersByContractHash",
  "params": {"ContractHash":"0xd2a4cff31913016155e38e474a2c06d08be276cf","Limit":2},
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
        "percentage": "0.9766352648440508376",
        "tokenid": ""
      },
      {
        "_id": "614befe8a141118435521dd8",
        "address": "0x4487494dac2f7eb68bdae009cacd6de88243e542",
        "asset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "balance": "39408603911904",
        "percentage": "0.00754671458001047158",
        "tokenid": ""
      }
    ],
    "totalCount": 596
  },
  "error": null
}
```
