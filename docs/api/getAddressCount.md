# GetAddressCount
Gets the count of all addresses
<hr>

### Parameters

None


### Example

Request body

```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAddressCount",
  "params": {},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": {
        "total counts": 721
    },
    "error": null
}
```
