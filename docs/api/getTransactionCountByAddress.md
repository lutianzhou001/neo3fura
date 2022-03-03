# GetTransactionCountByAddress
Gets the transaction count by given user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  the user's address| required|


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetTransactionCountByAddress",
  "params": {"Address":"0x0bf916d727c75f2e51e1ab2c476304513da59701"},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "total counts": 10
  },
  "error": null
}
```
