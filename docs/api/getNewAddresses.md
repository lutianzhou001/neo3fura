# GetNewAddresses
Gets the new addresses
<hr>

### Parameters

|    Name    | Type | Description |
| ---------- | --- |    ------    |
| Days       | int| The days in which new addresses were generated |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNewAddresses",
  "params": {"Days":2},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": [
        {
            "NewAddresses": 10,
            "_id": "6179e8006ea23aa8c2ac9d81"
        },
        {
            "NewAddresses": 8,
            "_id": "617b39809a0938a32fe195c3"
        }
    ],
        "error": null
}
```
