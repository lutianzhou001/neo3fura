# GetExecutionByTransactionHash
Gets the execution by the transaction hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash     | string| The transaction hash | Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"TransactionHash":"0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287"},
    "method": "GetExecutionByTransactionHash"
}'
```

Response body

```json
{
    "id": 1,
        "result": {
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
    "error": null
}
```
