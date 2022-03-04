# GetHourlyTransactions
Gets the number of transactions per hour
<hr>

### Parameters

|    Name    | Type | Description |
| ---------- | --- |    ------    |
| Hours      | int| Recent hours to query |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method":"GetHourlyTransactions",
    "params": {
        "Hours":5
    }
}'
```

Response body

```json
{
    "id": null,
        "result": [
        {
            "HourlyTransactions": 7,
            "_id": "61cab5f099f8679cd237b8de"
        },
        {
            "HourlyTransactions": 7,
            "_id": "61cac400184a3f74eb6374b9"
        },
        {
            "HourlyTransactions": 22,
            "_id": "61ca99d099f8679cd237b8d6"
        },
        {
            "HourlyTransactions": 9,
            "_id": "61ca8bc099f8679cd237b8d2"
        },
        {
            "HourlyTransactions": 9,
            "_id": "61caa7e099f8679cd237b8db"
        }
    ],
        "error": null
}
```
