# GetNep11TransferCountByAddress
Gets the Nep11 transfer count by the user's address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  The user's address| Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep11TransferCountByAddress",
  "params": {"Address":"0x2e9a0e6a68a4acce23ca14408bb4d0b803425394"},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": 55,
        "error": null
}
```
