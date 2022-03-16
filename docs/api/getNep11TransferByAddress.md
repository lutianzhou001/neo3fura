# GetNep11TransferByAddress
Gets the Nep11 transfer information by the user's address (0x0 transaction not included)
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address    | string|  The user's address| Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep11TransferByAddress",
  "params": {"Address":"0x2e9a0e6a68a4acce23ca14408bb4d0b803425394","Limit":2},
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
        "_id": "614bf7dba1411184355515f5",
        "blockhash": "0x3bebc4e090a3f1e7d2dc6f466c08377a407dd685e0eea84a64233af0411d9aa1",
        "contract": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
        "from": null,
        "frombalance": "0",
        "netfee": 747760,
        "sysfee": 3195474690,
        "timestamp": 1627540007545,
        "to": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
        "tobalance": "1",
        "tokenId": "QmxpbmQgQm94IDUx",
        "txid": "0x5581a8020fad2a422e75b7993ee3202be0a46350831a41e060a10cfe18bad877",
        "value": "1",
        "vmstate": "HALT"
      },
      {
        "_id": "614bf7dba1411184355515fc",
        "blockhash": "0x3bebc4e090a3f1e7d2dc6f466c08377a407dd685e0eea84a64233af0411d9aa1",
        "contract": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
        "from": null,
        "frombalance": "0",
        "netfee": 747760,
        "sysfee": 3195474690,
        "timestamp": 1627540007545,
        "to": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
        "tobalance": "1",
        "tokenId": "QmxpbmQgQm94IDcx",
        "txid": "0x5581a8020fad2a422e75b7993ee3202be0a46350831a41e060a10cfe18bad877",
        "value": "1",
        "vmstate": "HALT"
      }
    ],
    "totalCount": 356
  },
  "error": null
}
```
