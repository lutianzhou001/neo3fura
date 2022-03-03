#GetCandidateCount
Gets the count of candidates
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters
none


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetCandidateCount",
  "params": {},
  "id": 1
}'
```
### Response
```json
{
    "id": 1,
        "result": {
        "total counts": 49
    },
    "error": null
}
```
