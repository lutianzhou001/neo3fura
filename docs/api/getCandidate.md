# GetCandidate
Gets the candidate(s)
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetCandidate",
  "params": {"Limit":2},
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
                "_id": "614bef0ea14111843551a818",
                "candidate": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
                "isCommittee": true,
                "state": true,
                "votesOfCandidate": "3000367"
            },
            {
                "_id": "614bef0ea14111843551a820",
                "candidate": "0xaa606e99a6d1cb45ba34872864a3578c8a668143",
                "isCommittee": true,
                "state": true,
                "votesOfCandidate": "2006803"
            }
        ],
            "totalCount": 49
    },
    "error": null
}
```