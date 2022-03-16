# GetNotificationByEvent
Gets the notification by event
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Event    | string|  The name of event eg:"transfer"| Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNotificationByEvent",
  "params": {"Event": "Transfer","Limit":2},
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
                "Vmstate": "HALT",
                "_id": "614bef0ea14111843551a80c",
                "blockhash": "0x9d3276785e7306daf59a3f3b9e31912c095598bbfb8a4476b821b0e59be4c57a",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "eventname": "Transfer",
                "index": 0,
                "state": {
                    "type": "Array",
                    "value": [
                        {
                            "type": "Any",
                            "value": null
                        },
                        {
                            "type": "ByteString",
                            "value": "AZelPVEEY0csq+FRLl/HJ9cW+Qs="
                        },
                        {
                            "type": "Integer",
                            "value": "50000000"
                        }
                    ]
                },
                "timestamp": 1468595301000,
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
            },
            {
                "Vmstate": "HALT",
                "_id": "614bef0ea14111843551a80a",
                "blockhash": "0x9d3276785e7306daf59a3f3b9e31912c095598bbfb8a4476b821b0e59be4c57a",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "eventname": "Transfer",
                "index": 1,
                "state": {
                    "type": "Array",
                    "value": [
                        {
                            "type": "Any",
                            "value": null
                        },
                        {
                            "type": "ByteString",
                            "value": "krOcd6pg8ptXwXPO2Rfxf9Mhpus="
                        },
                        {
                            "type": "Integer",
                            "value": "5200000000000000"
                        }
                    ]
                },
                "timestamp": 1468595301000,
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
            }
        ],
            "totalCount": 692396
    },
    "error": null
}
```
