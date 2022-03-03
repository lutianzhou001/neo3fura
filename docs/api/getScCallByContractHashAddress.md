# GetScCallByContractHashAddress
Gets the ScCall by contract script hash and user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contract script hash| required|
| Address     | string|  the user's address| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetScCallByContractHashAddress",
  "params": {"Address":"0xeba621d37ff117d9ce73c1579bf260aa779cb392","ContractHash":"0xd2a4cff31913016155e38e474a2c06d08be276cf","Limit":2},
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
        "_id": "61d7e8f90506da8998f6fb67",
        "callFlags": "All",
        "contractHash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "hexStringParams": [
          "92b39c77aa60f29b57c173ced917f17fd321a6eb",
          "812be7a7082c9842fc2acefa4118b8bc5a5b918b",
          "00e8764817000000",
          ""
        ],
        "method": "transfer",
        "originSender": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "stack": "HALT",
        "txid": "0x88142a83918d35f30930dc88e370e69db5c9573acf2010a8e0aa5b2094094020"
      },
      {
        "_id": "61d7e8f90506da8998f6fb85",
        "callFlags": "All",
        "contractHash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "hexStringParams": [
          "92b39c77aa60f29b57c173ced917f17fd321a6eb",
          "b14858f18c76837415a61521c9cf69776e751f55",
          "00e8764817000000",
          ""
        ],
        "method": "transfer",
        "originSender": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "stack": "HALT",
        "txid": "0x4b4df96e27b2d763ebbf3f89b422c2f9f3eccc863f422dcda7e2f36c936d0bbd"
      }
    ],
    "totalCount": 6
  },
  "error": null
}
```
