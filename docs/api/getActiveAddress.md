# GetActiveAddresses
Gets the active address count in the specified days.
<hr>

### Parameters

|    Name    | Type | Description |  Required |
| ---------- | --- |    ------    | --------|
| Days       | int|  The number of recent days| Required|

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"Days":5},
    "method": "GetActiveAddresses"
}'
```

Response body

```json5
{
    "id": 1,
    "result": [
        {
            "ActiveAddresses": 72,
            "_id": "622005004b0d2d0e96331f63"
        },
        {
            "ActiveAddresses": 106,
            "_id": "621eb380795b34526419d05e"
        },
        {
            "ActiveAddresses": 217,
            "_id": "621d620146ca943cc00e82a8"
        },
        {
            "ActiveAddresses": 316,
            "_id": "621c1080e4cc4fe82e8a99cd"
        },
        {
            "ActiveAddresses": 160,
            "_id": "621abf00a5895c83a67ce121"
        }
    ],
    "error": null
}
```
