# GetCumulativeFeeBurn
Gets the cumulative systemFee burnt in total and the last 10 blocks systemFee burnt.
<hr>

### Parameters

None

### Example

``` powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetCumulativeFeeBurn",
  "params": {},
  "id": 1
}'
```

Response body

```json
{
    "id": 1,
        "result": [
        {
            "_id": "",
            "feeburn": 4438971290893,
            "result": [
                {
                    "index": 484849,
                    "systemFee": 0
                },
                {
                    "index": 484848,
                    "systemFee": 0
                },
                {
                    "index": 484847,
                    "systemFee": 0
                },
                {
                    "index": 484846,
                    "systemFee": 0
                },
                {
                    "index": 484845,
                    "systemFee": 0
                },
                {
                    "index": 484844,
                    "systemFee": 0
                },
                {
                    "index": 484843,
                    "systemFee": 0
                },
                {
                    "index": 484842,
                    "systemFee": 997775
                },
                {
                    "index": 484841,
                    "systemFee": 0
                },
                {
                    "index": 484840,
                    "systemFee": 0
                }
            ]
        }
    ],
        "error": null
}
```
