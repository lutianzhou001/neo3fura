# GetTotalSentAndReceivedByContractHash
Gets the total sent and received amount by the contract hash.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetTotalSentAndReceivedByContractHashAddress",
  "params": {"ContractHash":"0xd2a4cff31913016155e38e474a2c06d08be276cf","Address":"0xfa03cb7b40072c69ca41f0ad3606a548f1d59966"},
  "id": 1
  }'
```

Response body

```json
{
    "id": 1,
        "result": {
        "Address": "0xfa03cb7b40072c69ca41f0ad3606a548f1d59966",
            "ContractHash": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
            "received": 2439410926498,
            "sent": 2028508871288
    },
    "error": null
}
```
