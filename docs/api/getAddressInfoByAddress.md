# GetAddressInfoByAddress
Gets the address info by the address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description |  Required |
| ---------- | --- |    ------    | -------|
| Address      | string|  the user's address| required|


#### Example
``` powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAddressInfoByAddress",
  "params": {"Address":"0x0bf916d727c75f2e51e1ab2c476304513da59701"},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "_id": "614bef0ea14111843551a800",
    "address": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
    "firstusetime": 1468595301000,
    "lastusetime": 1627572748907,
    "transactionssent": 10
  },
  "error": null
}
```
