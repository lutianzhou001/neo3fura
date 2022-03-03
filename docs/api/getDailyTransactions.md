# GetDailyTransactions
Gets the daily transaction counts in last several days
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ---- |
| Days       | int|  the number of recent days| required |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetDailyTransactions",
  "params": {"Days":2},
  "id": 1
}'
```
### Response
```json
{
    "id": 1,
        "result": [
        {
            "DailyTransactions": 98,
            "_id": "6135961f8fb1d7b4f1f4bce9"
        },
        {
            "DailyTransactions": 128,
            "_id": "6136ab808fb1d7b4f1f4bd29"
        }
    ],
        "error": null
}
```
