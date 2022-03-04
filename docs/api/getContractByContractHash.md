# GetContractByContractHash
Gets the contract information by the contract script hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  contract script hash| Required |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetContractByContractHash",
  "params": {"ContractHash": "0xcc5e4edd9f5f8dba8bb65734541df7a1c081c67b" },
   "id":1
}'
```

Response body

```json
{
    "id": 1,
        "result": {
        "_id": "614bef0ea14111843551a7f6",
            "createTxid": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "createtime": 1468595301000,
            "hash": "0xcc5e4edd9f5f8dba8bb65734541df7a1c081c67b",
            "id": -7,
            "manifest": "{\"name\":\"PolicyContract\",\"groups\":[],\"features\":{},\"supportedstandards\":[],\"abi\":{\"methods\":[{\"name\":\"blockAccount\",\"parameters\":[{\"name\":\"account\",\"type\":\"Hash160\"}],\"returntype\":\"Boolean\",\"offset\":0,\"safe\":false},{\"name\":\"getExecFeeFactor\",\"parameters\":[],\"returntype\":\"Integer\",\"offset\":7,\"safe\":true},{\"name\":\"getFeePerByte\",\"parameters\":[],\"returntype\":\"Integer\",\"offset\":14,\"safe\":true},{\"name\":\"getStoragePrice\",\"parameters\":[],\"returntype\":\"Integer\",\"offset\":21,\"safe\":true},{\"name\":\"isBlocked\",\"parameters\":[{\"name\":\"account\",\"type\":\"Hash160\"}],\"returntype\":\"Boolean\",\"offset\":28,\"safe\":true},{\"name\":\"setExecFeeFactor\",\"parameters\":[{\"name\":\"value\",\"type\":\"Integer\"}],\"returntype\":\"Void\",\"offset\":35,\"safe\":false},{\"name\":\"setFeePerByte\",\"parameters\":[{\"name\":\"value\",\"type\":\"Integer\"}],\"returntype\":\"Void\",\"offset\":42,\"safe\":false},{\"name\":\"setStoragePrice\",\"parameters\":[{\"name\":\"value\",\"type\":\"Integer\"}],\"returntype\":\"Void\",\"offset\":49,\"safe\":false},{\"name\":\"unblockAccount\",\"parameters\":[{\"name\":\"account\",\"type\":\"Hash160\"}],\"returntype\":\"Boolean\",\"offset\":56,\"safe\":false}],\"events\":[]},\"permissions\":[{\"contract\":\"*\",\"methods\":\"*\"}],\"trusts\":[],\"extra\":null}",
            "name": "Policy",
            "nef": "{\"magic\":860243278,\"compiler\":\"neo-core-v3.0\",\"tokens\":[],\"script\":\"EEEa93tnQBBBGvd7Z0AQQRr3e2dAEEEa93tnQBBBGvd7Z0AQQRr3e2dAEEEa93tnQBBBGvd7Z0AQQRr3e2dA\",\"checksum\":3443651689}",
            "sender": null,
            "totalsccall": 6,
            "updatecounter": 0
    },
    "error": null
}
```
