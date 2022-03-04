# GetCandidateByVoterAddress

Gets the candidate voted by the voter's address.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| VoterAddress    | string |  the voter address| Required |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetCandidateByVoterAddress",
  "params": {"VoterAddress":"0xaa606e99a6d1cb45ba34872864a3578c8a668143"},
  "id": 1
}'
```

Response body

```json
{
  "id": 1,
  "result": {
    "_id": "614befeda1411184355220b5",
    "balanceOfVoter": "2000000",
    "blockNumber": 3967,
    "candidate": "0xaa606e99a6d1cb45ba34872864a3578c8a668143",
    "candidatePubKey": "02a7834be9b32e2981d157cb5bbd3acb42cfd11ea5c3b10224d7a44e98c5910f1b",
    "lastTransferTxid": null,
    "lastVoteTxid": "0x29c2fd2fd9c5c796e0bf10069269251e9433924e0b90ecd86614fe583ff23c75",
    "voter": "0xaa606e99a6d1cb45ba34872864a3578c8a668143"
  },
  "error": null
}
```
