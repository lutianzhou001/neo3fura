# GetExecutionByTrigger
Gets the execution by the given trigger
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Trigger    | enum|  Triggers in OnPersist, PostPersist, Application, Verification, System, All| Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetExecutionByTrigger",
  "params": {"Trigger":"Application","Limit":2},
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
        "_id": "61d7e8f30506da8998f6f42f",
        "blockhash": "0xcf35068b43281d700c6c7fc160ab844e74afeda08e793d061bbd1bc1a1203bd4",
        "exception": null,
        "gasconsumed": 9977780,
        "stacks": [],
        "timestamp": 1626850227986,
        "trigger": "Application",
        "txid": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287",
        "vmstate": "HALT"
      },
      {
        "_id": "61d7e8f47e71d96663aa73dc",
        "blockhash": "0x59890b5eeb781f334a076383ff92c453a0c2f6194abe19a97d6bb1c66c15bd79",
        "exception": null,
        "gasconsumed": 9977780,
        "stacks": [],
        "timestamp": 1626850648657,
        "trigger": "Application",
        "txid": "0x615c4c7ece85ce7d6cfe6d5f6d3495b5f46b43e298b79166488dbe431f067ca7",
        "vmstate": "HALT"
      }
    ],
    "totalCount": 86342
  },
  "error": null
}
```
