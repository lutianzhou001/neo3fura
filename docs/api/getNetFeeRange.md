# GetNetFeeRange
Gets the range of net fee
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
  "method": "GetNetFeeRange",
  "params": {},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "fast": 0,
    "fastest": 0,
    "slow": 0
  },
  "error": null
}
```
