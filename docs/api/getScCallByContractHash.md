# GetScCallByContractHash

Gets the ScCall by the contract script hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  contract script hash| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetScCallByContractHash",
  "params": {"ContractHash":"0xd2a4cff31913016155e38e474a2c06d08be276cf","Limit":2},
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
        "_id": "61d7e8f30506da8998f6f42e",
        "callFlags": "All",
        "contractHash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "hexStringParams": [
          "92b39c77aa60f29b57c173ced917f17fd321a6eb",
          "c098e4acf0b2090dd16ecb3d58de426a8715851b",
          "00407a10f35a0000",
          ""
        ],
        "method": "transfer",
        "originSender": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
        "stack": "HALT",
        "txid": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287"
      },
      {
        "_id": "61d7e8f47e71d96663aa73db",
        "callFlags": "All",
        "contractHash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "hexStringParams": [
          "c098e4acf0b2090dd16ecb3d58de426a8715851b",
          "1aefc0ce56187f47d0043290b63a45817933f681",
          "0010a5d4e8000000",
          ""
        ],
        "method": "transfer",
        "originSender": "0x1b8515876a42de583dcb6ed10d09b2f0ace498c0",
        "stack": "HALT",
        "txid": "0x615c4c7ece85ce7d6cfe6d5f6d3495b5f46b43e298b79166488dbe431f067ca7"
      }
    ],
    "totalCount": 30451
  },
  "error": null
}
```
