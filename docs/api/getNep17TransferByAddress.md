# GetNep17TransferByAddress
Gets the Nep17 transfer information by the user's address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  The user's address| Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep17TransferByAddress",
  "params": {"Address": "NbbBtdAbiCdvCaAhdT5dCgrZsAn1ZaUdot","ExcludeBonusAndBurn": true,"Limit":5},
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
        "_id": "614c43bfa1411184356c59af",
        "blockhash": "0x8a7fcdfd1227c2d425707605f010b0211dfb63670f6c998a4ca421197db6f0a0",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "frombalance": "821276368267",
        "netfee": 123262,
        "sysfee": 47662170,
        "timestamp": 1632388029414,
        "to": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "tobalance": "242600000005",
        "txid": "0x7c8663d61b247dbe822ac5b7c924204b3ac5a070c3e54ff9d3241ec24808e77c",
        "value": "2000000000",
        "vmstate": "HALT"
      },
      {
        "_id": "614c39ca3066938344847c14",
        "blockhash": "0xaf42456f93380c032bcb1b2efdde962d9c0eeb6e4fe2dda24d03f442b2ef1bc9",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "frombalance": "823588636617",
        "netfee": 1251030,
        "sysfee": 137341466,
        "timestamp": 1632302528857,
        "to": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "tobalance": "216000000000",
        "txid": "0x82f524bdcd9e3caf48873e28f181ee4cc042c0630a1cc784d7750b34a4a39a67",
        "value": "6000000000",
        "vmstate": "HALT"
      },
      {
        "_id": "614c3988306693834484570c",
        "blockhash": "0x8bd8bbe12914ec3cd8b39b24df0f7664ecc4a990e354f4f1e127f8a0391af062",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "frombalance": "829966492226",
        "netfee": 1250630,
        "sysfee": 46302170,
        "timestamp": 1632292412564,
        "to": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "tobalance": "117000000000",
        "txid": "0x9b07f9c43a0b0ada92dd6ed1434778c9d6eb127478dd494e9159bdb34c4e86c5",
        "value": "2000000000",
        "vmstate": "HALT"
      },
      {
        "_id": "614c3942306693834484379c",
        "blockhash": "0x428825af84441f0894770c0fc4887ba9e30bee1ef43277252e2e06c598cdacc4",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "frombalance": "832368130696",
        "netfee": 1306630,
        "sysfee": 231510850,
        "timestamp": 1632279489828,
        "to": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "tobalance": "32000000000",
        "txid": "0xa5f6cb26d6d1e8558f80e26457fbf19283e0d64bc06223d1b5503f3b623a744d",
        "value": "2000000000",
        "vmstate": "HALT"
      },
      {
        "_id": "614c3940a1411184356bcef3",
        "blockhash": "0x914e0c6a04ee936d2f2baa023b7d074e7d502f493404863224c7f2d2d7d9b7bb",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": "0xb31b1ef4b504f5413dbed7e6e58fd11dedb6f4ab",
        "frombalance": "834955133846",
        "netfee": 1306630,
        "sysfee": 231510850,
        "timestamp": 1632279324155,
        "to": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "tobalance": "30000000000",
        "txid": "0x83821bf768efe417840349f716ba241a9b5ab1c4ef6e106e403888a747994647",
        "value": "2000000000",
        "vmstate": "HALT"
      }
    ],
    "totalCount": 66
  },
  "error": null
}
```
