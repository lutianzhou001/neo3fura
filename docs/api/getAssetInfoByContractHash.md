# GetAssetInfoByContractHash

Gets the asset info by the contact script hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the scrip hash of the asset want to query| required|



#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "GetAssetInfoByContractHash",
    "params": {
      "ContractHash": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f"
    }
}'
```
### Response
```json5
{
    "id": null,
        "result": {
        "_id": "61d82a010506da89981d7a69",
            "decimals": 0,
            "firsttransfertime": 1630901464602,
            "hash": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
            "holders": 3554,
            "ispopular": false,
            "symbol": "N3",
            "tokenname": "Neoverse",
            "totalsupply": "3554",
            "type": "NEP11"
    },
    "error": null
}
```
