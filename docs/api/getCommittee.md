# GetCommittee
Gets the committee of the blockchain
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
  "method": "GetCommittee",
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
        "_id": "614bef0ea14111843551a811",
        "candidate": "0x8b915b5abcb81841face2afc42982c08a7e72b81",
        "isCommittee": true,
        "state": true,
        "votesOfCandidate": "2000000"
      },
      {
        "_id": "614bef0ea14111843551a812",
        "candidate": "0xa4887b48371fe7727d9f96f4922f464c9c457d89",
        "isCommittee": true,
        "state": true,
        "votesOfCandidate": "2000000"
      }
    ],
    "totalCount": 21
  },
  "error": null
}
```
