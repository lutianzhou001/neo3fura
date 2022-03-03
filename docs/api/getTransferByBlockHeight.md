# GetTransferByBlockHeight
Gets the transfer info by blockheight
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHeight    | int|  the block height| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetTransferByBlockHeight",
  "params": {"BlockHeight":53429,"Limit":2,"Skip":2},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614bf8fca141118435569d14",
        "blockhash": "0x0073a37aeb1d27b2d34f95e15d1eb63c96f8a1b89d05c4d045195f5400806ce9",
        "contract": "0xef4073a0f2b305a38ec4050e4d3d28bc40ea63f5",
        "from": "0xfa03cb7b40072c69ca41f0ad3606a548f1d59966",
        "frombalance": "5112",
        "timestamp": 1627616538985,
        "to": "0xdf378f38d705999148da9d8355871b908804f14c",
        "tobalance": "6776",
        "txid": "0x8b41d3989f90795be966fe02d412630f0227950396985ca703a4d3c5467683bc",
        "value": "1"
      },
      {
        "_id": "614bf8fca141118435569d1e",
        "blockhash": "0x0073a37aeb1d27b2d34f95e15d1eb63c96f8a1b89d05c4d045195f5400806ce9",
        "contract": "0x1415ab3b409a95555b77bc4ab6a7d9d7be0eddbd",
        "from": "0xfa03cb7b40072c69ca41f0ad3606a548f1d59966",
        "frombalance": "15905933693121826",
        "timestamp": 1627616538985,
        "to": "0xdf378f38d705999148da9d8355871b908804f14c",
        "tobalance": "83654314128",
        "txid": "0x7953cd79ee763d517f7308ee0708ab9875aa737027be353d55e112b09663409b",
        "value": "12345678"
      }
    ],
    "totalCount": 2424
  },
  "error": null
}
```
