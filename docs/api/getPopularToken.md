# GetPopularToken
Gets the votes by the candidate address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Standard | string| The type of asset:NEP11  or NEP17 | Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetPopularToken",
  "params": {"Standard": "NEP11"},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
    "result": {
        "result": [
            {
                "decimals": 0,
                "hash": "0x9f344fe24c963d70f5dcf0cfdeb536dc9c0acb3a",
                "symbol": "ILEX POLEMEN",
                "tokenname": "ILEX POLEMEN",
                "totalsupply": "2000",
                "type": "NEP11"
            }
        ],
        "totalCount": 1
    },
    "error": null
}
```
