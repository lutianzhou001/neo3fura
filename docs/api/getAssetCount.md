# GetAssetCount
Gets the number of all assets.
<hr>

### Parameters
None

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetCount",
  "params": {},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": {
        "total counts": 79
    },
    "error": null
}
```
