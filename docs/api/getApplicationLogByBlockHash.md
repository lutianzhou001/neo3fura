# GetApplicationLogByBlockHash
Gets the applicationlog by blockhash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHash     | string|  blockHash of a transaction| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetApplicationLogByBlockHash",
  "params": {"BlockHash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f" },
  "id": 1
}'
```
### Response
```json5{
    "id": 1,
    "result": {
        "result": [
            {
                "_id": "61d7e8f50506da8998f6f681",
                "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
                "exception": null,
                "gasconsumed": 0,
                "notifications": [
                    {
                        "Vmstate": "HALT",
                        "_id": "61d7e8f50506da8998f6f683",
                        "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
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
                                    "value": "gSvnpwgsmEL8Ks76QRi4vFpbkYs="
                                },
                                {
                                    "type": "Integer",
                                    "value": "50000000"
                                }
                            ]
                        },
                        "timestamp": 1626851177411,
                        "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
                    }
                ],
                "stacks": [],
                "timestamp": 1626851177411,
                "trigger": "PostPersist",
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "vmstate": "HALT"
            },
            {
                "_id": "61d7e8f50506da8998f6f682",
                "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
                "exception": null,
                "gasconsumed": 0,
                "notifications": [
                    {
                        "Vmstate": "HALT",
                        "_id": "61d7e8f50506da8998f6f683",
                        "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
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
                                    "value": "gSvnpwgsmEL8Ks76QRi4vFpbkYs="
                                },
                                {
                                    "type": "Integer",
                                    "value": "50000000"
                                }
                            ]
                        },
                        "timestamp": 1626851177411,
                        "txid": "0x0000000000000000000000000000000000000000000000000000000000000000"
                    }
                ],
                "stacks": [],
                "timestamp": 1626851177411,
                "trigger": "OnPersist",
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "vmstate": "HALT"
            }
        ],
        "totalCount": 2
    },
    "error": null
}
```
