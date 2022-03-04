# GetNep11BalanceByContractHashAddressTokenId
Gets the Nep11 balance by the contract script hash, user's address, and tokenId of the Nep11 standard.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required |
| Address | string | The user's address | Required |
| TokenId    | string| The tokenId of the Nep11 standard | Required |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep11BalanceByContractHashAddressTokenId",
  "params": {"ContractHash": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40","Address":"0x2e9a0e6a68a4acce23ca14408bb4d0b803425394","tokenId":"QmxpbmQgQm94IDIxNQ=="},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": {
        "_id": "614bf99ba14111843558cb34",
            "address": "0x2e9a0e6a68a4acce23ca14408bb4d0b803425394",
            "asset": "0xb3b65e5c0d2af3f98cac6e80083f6c2b90476f40",
            "balance": "1",
            "tokenid": "QmxpbmQgQm94IDIxNQ=="
    },
    "error": null
}
```
