# GetBlockHashByBlockHeight
Gets the blockhash by the block height
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHeight     | int| The block height | Required |

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"BlockHeight":3823},
    "method": "GetBlockHashByBlockHeight"
}'
```

Response body

```json5
{
  "id": 1,
  "result": {
    "hash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f"
  },
  "error": null
}
```
