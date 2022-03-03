# GetAssetInfos
Gets asset infos by the contract script hash array
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| Addresses     | string[]|  the script hash array of the asset want to query| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetInfos",
  "params": {"Addresses":["0xd2a4cff31913016155e38e474a2c06d08be276cf"]},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614bef0ea14111843551a7fd",
        "decimals": 8,
        "firsttransfertime": 1468595301000,
        "hash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "holders": 619,
        "ispopular": true,
        "symbol": "GAS",
        "tokenname": "GasToken",
        "totalsupply": "5221959738128399",
        "type": "NEP17"
      }
    ],
    "totalCount": 1
  },
  "error": null
}

```
