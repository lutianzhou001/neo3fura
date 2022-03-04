# GetApplicationLogByTransactionHash
Gets the applicationlog by the given transactionhash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash       | string|  TransactionHash| Required |


### Example

Request body

```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetApplicationLogByTransactionHash",
  "params": {"TransactionHash": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287" },
  "id": 1
}'
```

Response body

```json5

{
    "id": 1,
        "result": {
        "_id": "61d7e8f30506da8998f6f42f",
            "blockhash": "0xcf35068b43281d700c6c7fc160ab844e74afeda08e793d061bbd1bc1a1203bd4",
            "exception": null,
            "gasconsumed": 9977780,
            "notifications": [
            {
                "Vmstate": "HALT",
                "_id": "61d7e8f30506da8998f6f432",
                "blockhash": "0xcf35068b43281d700c6c7fc160ab844e74afeda08e793d061bbd1bc1a1203bd4",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "eventname": "Transfer",
                "index": 0,
                "state": {
                    "type": "Array",
                    "value": [
                        {
                            "type": "ByteString",
                            "value": "krOcd6pg8ptXwXPO2Rfxf9Mhpus="
                        },
                        {
                            "type": "ByteString",
                            "value": "wJjkrPCyCQ3Rbss9WN5CaocVhRs="
                        },
                        {
                            "type": "Integer",
                            "value": "100000000000000"
                        }
                    ]
                },
                "timestamp": 1626850227986,
                "txid": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287"
            }
        ],
            "stacks": [],
            "timestamp": 1626850227986,
            "trigger": "Application",
            "txid": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287",
            "vmstate": "HALT"
    },
    "error": null
}
```
