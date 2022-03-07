# GetCandidateByAddress
Gets the candidate information by the candidate address.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address      | string|  The candidate address| Required |

### Example

Request body

```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "GetCandidateByAddress",
    "params": {"Address": "0xaa606e99a6d1cb45ba34872864a3578c8a668143"},
    "id": 1
}'
```

Response body

```json
{
  "id": 1,
  "result": {
    "_id": "614bef0ea14111843551a820",
    "candidate": "0xaa606e99a6d1cb45ba34872864a3578c8a668143",
    "isCommittee": true,
    "state": true,
    "votesOfCandidate": "2006803"
  },
  "error": null
}
```
