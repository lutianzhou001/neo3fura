# GetNotificationByContractHash
Gets the notification by the contract script hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string| The contract script hash | Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetNotificationByContractHash",
  "params": {"ContractHash":"0xd2a4cff31913016155e38e474a2c06d08be276cf","Limit":2},
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
                "_id": "61712d40770c0b0a6a3ceb24",
                "blockhash": "0xcfab421df7e15976a8878303bfb5cb0d5703b77a8a9c15e8a33cb82c49a9a5db",
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
                            "value": "ZlPI6+wR2Aei/ElqjviDj9PpLeU="
                        },
                        {
                            "type": "Integer",
                            "value": "50000000"
                        }
                    ]
                },
                "timestamp": 1634807103949,
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
            },
            {
                "Vmstate": "HALT",
                "_id": "61712d3050025b01612b7efb",
                "blockhash": "0xbab929e8de4cd215eb6244e5ef17ce4b128234362a9e00d2cbbf83a11383156f",
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
                            "value": "lPAUr0QjnhYz8fCJaBQQm47mzIU="
                        },
                        {
                            "type": "Integer",
                            "value": "50000000"
                        }
                    ]
                },
                "timestamp": 1634807088845,
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
            }
        ],
            "totalCount": 593943
    },
    "error": null
}
```
