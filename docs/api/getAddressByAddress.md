# GetAddressByAddress
Gets the details of the given address.
<hr>

### Parameters

|    Name    | Type | Description |  Required |
| ---------- | --- |    ------    | -----|
| Address       | string|  The user's address| Required|

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAddressByAddress",
  "params": {"Address": "0x0bf916d727c75f2e51e1ab2c476304513da59701"},
  "id": 1
}'
```

Response body

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
