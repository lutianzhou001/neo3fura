# GetAssetHoldersListByContractHash
Gets all NEP11 assets and its holders with the contract script hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  contract script hash| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |
| Balance    | int|  equal to 1 | required|


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAssetHoldersListByContractHash",
  "params": {"ContractHash":"0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f","Limit":10,"balance":1},
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
        "_id": "61d98d6da24c739532d796e2",
        "address": "0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgSCAjMjc2"
      },
      {
        "_id": "61d98152b145511ecce2382e",
        "address": "0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgRyAjMTAzOA=="
      },
      {
        "_id": "61d980bbb145511ecce1a96c",
        "address": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgQyAjMTY1NA=="
      },
      {
        "_id": "61d97767b145511eccd8ad85",
        "address": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgQiAjMTA3Nw=="
      },
      {
        "_id": "61d970edb145511eccd28b19",
        "address": "0x6835f6961eadbad3e75f2ea2f7a52d04deb82005",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgRSAjMTIy"
      },
      {
        "_id": "61d96adcb145511eccccc7ef",
        "address": "0xc73b2693a3c6d125d2cdfa24a2ef74f11de2a128",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "TyAjMzI5IFNjYWxhYmlsaXR5"
      },
      {
        "_id": "61d96adcb145511eccccc763",
        "address": "0xc73b2693a3c6d125d2cdfa24a2ef74f11de2a128",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgRiAjMjMzOA=="
      },
      {
        "_id": "61d96adcb145511eccccc762",
        "address": "0xc73b2693a3c6d125d2cdfa24a2ef74f11de2a128",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgRyAjMjI5MA=="
      },
      {
        "_id": "61d96adcb145511eccccc761",
        "address": "0xc73b2693a3c6d125d2cdfa24a2ef74f11de2a128",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgSSAjNDYw"
      },
      {
        "_id": "61d96adcb145511eccccc75e",
        "address": "0xc73b2693a3c6d125d2cdfa24a2ef74f11de2a128",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "balance": "1",
        "percentage": "0.0002813731007315700619",
        "tokenid": "RnJhZ21lbnQgQSAjMTgyNQ=="
      }
    ],
    "totalCount": 3554
  },
  "error": null
}
```
