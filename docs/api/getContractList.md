# GetContractList
Gets the contract list
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters
|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetContractList",
  "params": {"Limit":2,"Skip":2},
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
                "Transaction": [],
                "createtime": 1468595301000,
                "hash": "0xcc5e4edd9f5f8dba8bb65734541df7a1c081c67b",
                "id": -7,
                "name": "Policy",
                "updatecounter": 0
            },
            {
                "Transaction": [],
                "createtime": 1468595301000,
                "hash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "id": -6,
                "name": "GasToken",
                "updatecounter": 0
            }
        ],
            "totalCount": 278
    },
    "error": null
}
```
