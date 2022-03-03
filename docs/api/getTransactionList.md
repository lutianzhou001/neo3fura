# GetTransactionList
Gets the transaction list
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetTransactionList",
  "params": {"Limit":2,"Skip":2},
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
        "_id": "6176712950025b01612dcfe2",
        "attributes": [],
        "blockIndex": 540820,
        "blockhash": "0x3a5fdece2a3a2d05667a0a35fba09572517bb46dd57b1b99d4121463b16ca6a0",
        "blocktime": 1635152168885,
        "hash": "0xf72f07cf3354792a763a867d319ec8e9620ef8c7aa34a39fa16218ec71a21649",
        "netfee": 123262,
        "nonce": 2538454220,
        "script": "CwKAlpgADBQWEbxUqVrpYVXCKu9c8UopEu0wYAwUq/S27R3Rj+Xm1749QfUEtfQeG7MUwB8MCHRyYW5zZmVyDBTPduKL0AYsSkeO41VhARMZ88+k0kFifVtS",
        "sender": "NbbBtdAbiCdvCaAhdT5dCgrZsAn1ZaUdot",
        "signers": [
          {
            "account": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
            "scopes": "CalledByEntry"
          }
        ],
        "size": 248,
        "sysfee": 997775,
        "validUntilBlock": 540850,
        "version": 0,
        "witnesses": [
          {
            "invocation": "DEBqd9tHd114hh8g51uviZ2BUOypG41AgBpNku4ZK3MRbQfD7AVmlW5uZjwZJwtjtEUO56w9yocDwmoaP84Tn+Cj",
            "verification": "DCECbr4UsjAVwo4lsFiy/8PkbdMprchw04+ZaJl297KKxVNBVuezJw=="
          }
        ]
      },
      {
        "_id": "61766bf750025b01612dcd59",
        "attributes": [],
        "blockIndex": 540734,
        "blockhash": "0x46593cff5f1532aebad88700fe5e3f3c7572b64e1cf9b071018eb9d9a3e04b01",
        "blocktime": 1635150839749,
        "hash": "0x872fec7575dbbedb99d72055efcb96517d3165c098129211aeeeea2e6cacff8c",
        "netfee": 123662,
        "nonce": 3343936879,
        "script": "CwMA5AtUAgAAAAwUFhG8VKla6WFVwirvXPFKKRLtMGAMFKv0tu0d0Y/l5te+PUH1BLX0HhuzFMAfDAh0cmFuc2ZlcgwUKkyaTUAiZ4sD7xu+CDT5ZkYNxEhBYn1bUg==",
        "sender": "NbbBtdAbiCdvCaAhdT5dCgrZsAn1ZaUdot",
        "signers": [
          {
            "account": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
            "scopes": "CalledByEntry"
          }
        ],
        "size": 252,
        "sysfee": 3262030,
        "validUntilBlock": 540764,
        "version": 0,
        "witnesses": [
          {
            "invocation": "DEAixl8mUfiIlSZlhyNgwJ3KPX1UA6cWXCIX3WsA3Bqb2+s2yRbEotqMcBXGee2wAQHAB0qe7VsX5WJgCOTk4Z0k",
            "verification": "DCECbr4UsjAVwo4lsFiy/8PkbdMprchw04+ZaJl297KKxVNBVuezJw=="
          }
        ]
      }
    ],
    "totalCount": 35885
  },
  "error": null
}
```
