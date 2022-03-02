# GetNep17TransferByTransactionHash
Gets the Nep17 transfer by transaction hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash    | string|  the transaction hash| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep17TransferByTransactionHash",
  "params": {"TransactionHash": "0x237aae2efdb459ade601d07db39b1b0134c19b40912933414217bc1494cd009b","Limit":2},
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
        "_id": "61d7eb870506da8998f87aa4",
        "blockhash": "0x8d0fffcc1938987c01b460ddcbd6c9fc01a9b83300d07c3cd680e4036b680dfc",
        "contract": "0xef4073a0f2b305a38ec4050e4d3d28bc40ea63f5",
        "decimals": 0,
        "from": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "frombalance": "79000000",
        "symbol": "NEO",
        "timestamp": 1627027447114,
        "to": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
        "tobalance": "3000000",
        "tokenname": "NeoToken",
        "txid": "0x237aae2efdb459ade601d07db39b1b0134c19b40912933414217bc1494cd009b",
        "value": "1000000",
        "vmstate": "HALT"
      },
      {
        "_id": "61d7eb870506da8998f87aa3",
        "blockhash": "0x8d0fffcc1938987c01b460ddcbd6c9fc01a9b83300d07c3cd680e4036b680dfc",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "decimals": 8,
        "from": null,
        "frombalance": "0",
        "symbol": "GAS",
        "timestamp": 1627027447114,
        "to": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
        "tobalance": "368013189000",
        "tokenname": "GasToken",
        "txid": "0x237aae2efdb459ade601d07db39b1b0134c19b40912933414217bc1494cd009b",
        "value": "331624000000",
        "vmstate": "HALT"
      }
    ],
    "totalCount": 3
  },
  "error": null
}
```
