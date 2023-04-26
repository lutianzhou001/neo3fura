# GetNNSNameByAdmin
Gets the list of NNS names by admin 
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Asset     | string|  The asset hash | Required |
| Admin     | string| The Admin hash | Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |

### Example

Request body

```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNNSNameByAdmin",
  "params": {
      "Asset":"0x152fa9ceeb2c83f40e3d3d6da6c1f8898dd4891a",
      "Admin":"0xf0a33d62f32528c25e68951286f238ad24e30032",
      "Skip":0,
      "Limit":20   
       },
  "id": 1
}'
```
Response body

```json
{
  "id": 1,
  "result": [
    {
      "name": [
        "hhh.neo"
      ]
    }
  ],
  "error": null
}
```

