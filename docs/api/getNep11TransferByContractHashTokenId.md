# GetNep11TransferByContractHashTokenId
Gets the nep11 transfer information by the contract script hash and tokenid
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep11TransferByContractHashTokenId",
  "params": {"ContractHash":"0xb137c83610d3f0331a48d8d6283864120b4f23a1","tokenId":"1wA="},
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
        "_id": "61d978e1a24c739532d61bed",
        "blockhash": "0x81ea161bc4f8228b1dc92e16031242d83e03420f1051477ad802d443e321b43b",
        "contract": "0xb137c83610d3f0331a48d8d6283864120b4f23a1",
        "from": null,
        "frombalance": "0",
        "timestamp": 1639124530142,
        "to": "0x1372bd39e447a94c8cd0491cfbd2703b5e39bd16",
        "tobalance": "1",
        "tokenId": "1wA=",
        "txid": "0x432f9252c4b2146e94ef3bcd44b04072366f8b9a8efa7de8af11fe39866b6d02",
        "value": "1"
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
