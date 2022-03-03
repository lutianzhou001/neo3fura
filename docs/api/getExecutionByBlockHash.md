# GetExecutionByBlockHash

Gets the execution by blockhash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters
|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHash      | string|  the blockHash of the block| required


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"BlockHash":"0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f"},
    "method": "GetExecutionByBlockHash"
}'
```
### Response
```json5
{
    "id": 1,
        "result": {
        "result": [
            {
                "_id": "61d7e8f50506da8998f6f681",
                "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
                "exception": null,
                "gasconsumed": 0,
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
