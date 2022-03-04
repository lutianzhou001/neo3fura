# GetVotesByCandidateAddress
Gets the votes by the candidate address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| CandidateAddress     | string|  The candidate address| Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetVotesByCandidateAddress",
  "params": {"CandidateAddress":"0x0bf916d727c75f2e51e1ab2c476304513da59701","Limit":2},
  "id": 1
}'
```

Response body

```json
{
  "id": 1,
  "result": {
    "_id": "614bef0ea14111843551a818",
    "candidate": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
    "isCommittee": true,
    "state": true,
    "votesOfCandidate": "3000530"
  },
  "error": null
}
```
