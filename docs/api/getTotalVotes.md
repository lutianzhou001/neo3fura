# GetTotalVotes
Gets the total votes of all candidates
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters
none

#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetTotalVotes",
  "params": {},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "totalvotes": 20141013
  },
  "error": null
}
```
