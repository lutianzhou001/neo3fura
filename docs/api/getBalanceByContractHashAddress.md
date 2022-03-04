# GetBalanceByContractHashAddress
Gets the asset balance by the given asset contract script hash and user's address

<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required|
| Address   | string| The user's address | Required|

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetBalanceByContractHashAddress",
  "params": {"Address":"NUzy2Ns2D35BTdFVqDhUCRoZb1cmix2cXS","ContractHash":"0xef4073a0f2b305a38ec4050e4d3d28bc40ea63f5"},
  "id": 1
}'
```

Response body

```json5

{
  "id": 1,
  "result": {
    "_id": "614bf1c8306693834466cd91",
    "address": "0x96d5942028891de8e5d866f504b36ff5ae13ab63",
    "asset": "0xef4073a0f2b305a38ec4050e4d3d28bc40ea63f5",
    "balance": "43",
    "tokenid": ""
  },
  "error": null
}
```
