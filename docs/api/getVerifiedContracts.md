# GetVerifiedContracts
Gets the verified contract
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetVerifiedContracts",
  "params": {"Skip":2,"Limit":2},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": [
    {
      "_id": "61700825eb743bed51ae9b20",
      "hash": "0xcc5e4edd9f5f8dba8bb65734541df7a1c081c67b",
      "id": -7,
      "updatecounter": 0
    },
    {
      "_id": "61700825eb743bed51ae9b21",
      "hash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
      "id": -6,
      "updatecounter": 0
    }
  ],
  "error": null
}
```
