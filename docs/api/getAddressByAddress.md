# GetAddressByAddress
Gets the address info by the address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description |  Required |
| ---------- | --- |    ------    | -----|
| Address       | string|  the user's address| required|


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAddressByAddress",
  "params": {"Address": "0x0bf916d727c75f2e51e1ab2c476304513da59701"},
  "id": 1
}'
```
### Response
```json5
{
    "id": 1,
        "result": {
        "_id": "614bef0ea14111843551a800",
            "address": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
            "firstusetime": 1468595301000
    },
    "error": null
}
```
