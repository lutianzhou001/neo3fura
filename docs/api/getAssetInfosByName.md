# GetAssetInfosByName
Gets the asset infos with token name
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| Name     | string|  token name (fuzzy search supported)| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |

#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetInfosByName",
  "params": {"Name":"GasTo"},
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
                "_id": "614bef0ea14111843551a7fd",
                "decimals": 8,
                "firsttransfertime": 1468595301000,
                "hash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "holders": 619,
                "ispopular": true,
                "symbol": "GAS",
                "tokenname": "GasToken",
                "totalsupply": "5221960958138950",
                "type": "NEP17"
            }
        ],
            "totalCount": 1
    },
    "error": null
}
```
