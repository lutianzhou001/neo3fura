# GetNep11OwnedByContractHashAddress
Gets the Nep11 owned by contracthash and user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash    | string|  contract script hash| required|
| Address   | string|  the user's address| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetNep11OwnedByContractHashAddress",
  "params": {"Address":"0x2e9a0e6a68a4acce23ca14408bb4d0b803425394","ContractHash":"0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40","limit":3},
  "id": 1
}'
```
### Response
```json5
{
    "id": 1,
        "result": {
        "result": [
            {
                "_id": "614bf7dba141118435551610",
                "blockhash": "0x3bebc4e090a3f1e7d2dc6f466c08377a407dd685e0eea84a64233af0411d9aa1",
                "contract": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
                "from": null,
                "frombalance": "0",
                "timestamp": 1627540007545,
                "to": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
                "tobalance": "1",
                "tokenId": "QmxpbmQgQm94IDg4",
                "txid": "0x5581a8020fad2a422e75b7993ee3202be0a46350831a41e060a10cfe18bad877",
                "value": "1"
            },
            {
                "_id": "614bf7dba1411184355515f5",
                "blockhash": "0x3bebc4e090a3f1e7d2dc6f466c08377a407dd685e0eea84a64233af0411d9aa1",
                "contract": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
                "from": null,
                "frombalance": "0",
                "timestamp": 1627540007545,
                "to": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
                "tobalance": "1",
                "tokenId": "QmxpbmQgQm94IDUx",
                "txid": "0x5581a8020fad2a422e75b7993ee3202be0a46350831a41e060a10cfe18bad877",
                "value": "1"
            },
            {
                "_id": "614bf7dba1411184355515fc",
                "blockhash": "0x3bebc4e090a3f1e7d2dc6f466c08377a407dd685e0eea84a64233af0411d9aa1",
                "contract": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
                "from": null,
                "frombalance": "0",
                "timestamp": 1627540007545,
                "to": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
                "tobalance": "1",
                "tokenId": "QmxpbmQgQm94IDcx",
                "txid": "0x5581a8020fad2a422e75b7993ee3202be0a46350831a41e060a10cfe18bad877",
                "value": "1"
            }
        ],
            "totalCount": 304
    },
    "error": null
}
```
