# GetCandidateByAddress
Gets the candidate by the candidate address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address      | string|  the candidate address| required |


#### Example
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
### Response
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
