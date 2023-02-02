# GetAssetInfoByContracts

Gets the asset information by the contact script hash.
<hr>

### Parameters

|    Name    | Type     | Description | Required |
| ---------- |----------|    ------    | ----|
| ContractHash     | []string | The scrip hash of the asset to query | Required|

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "GetAssetInfoByContracts",
    "params": {
      "ContractHash": ["0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
                       "0xd2a4cff31913016155e38e474a2c06d08be276cf"
                    ]
    }
}'
```

Response body

```json5

{
  "id": 1,
  "result": {
    "result": [     
      {
        "_id": "63c50098f53308831ae9d346",
        "decimals": 8,
        "firsttransfertime": 1468595301000,
        "hash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "holders": 42835,
        "ispopular": true,
        "symbol": "GAS",
        "tokenname": "GasToken",
        "totalsupply": "5981680925124462",
        "type": "NEP17"
      },
      {
        "_id": "63c7ae118a54c045d4f79e3c",
        "decimals": 0,
        "firsttransfertime": 1663705756342,
        "hash": "0xf456649d0b8f331596035a07f977cb8d8dbf0122",
        "holders": 119,
        "ispopular": false,
        "symbol": "CUTIE",
        "tokenname": "CutieToken",
        "totalsupply": "22788",
        "type": "NEP11"
      }
    ],
    "totalCount": 4
  },
  "error": null
}
```
