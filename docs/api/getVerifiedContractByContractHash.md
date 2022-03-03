# GetVerifiedContractByContractHash
Gets the verified contract by contract hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  contract script hash| required|
| UpdateCounter   | int|  the number of times the contract has been updated| required|


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetVerifiedContractByContractHash",
  "params": {"ContractHash":"0xfe924b7cfe89ddd271abaf7210a80a7e11178758","UpdateCounter":0},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "_id": "61700825eb743bed51ae9b1e",
    "hash": "0xfe924b7cfe89ddd271abaf7210a80a7e11178758",
    "id": -9,
    "updatecounter": 0
  },
  "error": null
}
```
