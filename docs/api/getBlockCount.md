# GetBlockCount
Gets the total blocks of the chain
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
    "id": 1,
    "params": {},
    "method": "GetBlockCount"
}'
```
### Response
```json5
{
    "id": 1,
    "result": {
        "index": 479389
    },
    "error": null
}
```
