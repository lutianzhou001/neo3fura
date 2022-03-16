# GetTransferByAddress
Gets the transfer by the user's address
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
  "method": "GetTransferByAddress",
  "params": {"Limit":2,"Skip":2,"Address":"0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab"},
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
        "_id": "614c3932306693834484328c",
        "blockhash": "0xc4bfb7ec80cb47fc60f2bf123fdd24150c3499e5c99b553b7bb4131857f3564a",
        "contract": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "from": null,
        "frombalance": "0",
        "timestamp": 1632278180438,
        "to": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "tobalance": "1",
        "tokenId": "QmxpbmQgQm94ICMxMzY=",
        "txid": "0x4e2cadb3f2db7071aaf2aaa483ed27d65797f5e226855ff5e155a30a110b44f3",
        "value": "1"
      },
      {
        "_id": "614c3932306693834484328d",
        "blockhash": "0xc4bfb7ec80cb47fc60f2bf123fdd24150c3499e5c99b553b7bb4131857f3564a",
        "contract": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "from": null,
        "frombalance": "0",
        "timestamp": 1632278180438,
        "to": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "tobalance": "1",
        "tokenId": "QmxpbmQgQm94ICMxNDE=",
        "txid": "0x4e2cadb3f2db7071aaf2aaa483ed27d65797f5e226855ff5e155a30a110b44f3",
        "value": "1"
      }
    ],
    "totalCount": 537
  },
  "error": null
}
```
