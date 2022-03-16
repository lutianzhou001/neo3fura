# GetRawTransactionByAddress
Gets the raw transaction by address
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Address     | string|  The user's address| Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{  
  "jsonrpc": "2.0",
  "method": "GetRawTransactionByAddress",
  "params": {"Address":"NfnYGUoC6TKhevumjTUFyYnfQtZRjMaXch","Limit":2},
  "id": 1
}
'
```

Response body

```json
{
    "id": 1,
        "result": {
        "result": [
            {
                "_id": "614c32aaa1411184356a37a3",
                "attributes": [],
                "blockIndex": 328660,
                "blockhash": "0x91fd15076cfb5c4b7cf4390b21543afbeed3623667415329f7a11da4a4add1b9",
                "blocktime": 1631869867718,
                "faultdetail": {
                    "_id": "61445babbf012c83f64e7243",
                    "callFlags": "All",
                    "contractHash": "0xe1b2d188b37e8ed26df4e88c2f5165429bd530a0",
                    "hexStringParams": [
                        ""
                    ],
                    "method": "stake",
                    "originSender": "0x530e099ac607559db8424896ca159c7090cdfad9",
                    "txid": "0x79ed5544ffa6e28e5ac171e3bc5e027da7982695c1b0565bd3ab06e9d1aac9fd",
                    "vmstate": "FAULT"
                },
                "hash": "0x79ed5544ffa6e28e5ac171e3bc5e027da7982695c1b0565bd3ab06e9d1aac9fd",
                "netfee": 1240235,
                "nonce": 1129874227,
                "script": "ERHAHwwFc3Rha2UMFKAw1ZtCZVEvjOj0bdKOfrOI0bLhQWJ9W1I6",
                "sender": "NfnYGUoC6TKhevumjTUFyYnfQtZRjMaXch",
                "signers": [
                    {
                        "account": "0x530e099ac607559db8424896ca159c7090cdfad9",
                        "scopes": "Global"
                    }
                ],
                "size": 197,
                "sysfee": 15988875,
                "validUntilBlock": 334245,
                "version": 0,
                "vmstate": "FAULT",
                "witnesses": [
                    {
                        "invocation": "DEChvJFRAd0132LO9JWx4fSxarEI1tIaoTx0KjT/G4YuY5g8yDfSSWwYQDYXTkxg/mZESVVlHhPy6MIz5GPTpcPE",
                        "verification": "DCED9ocyvgPQfEj4x6JFJKy4jWXlbdZBmiB6H013ZKiB7TRBVuezJw=="
                    }
                ]
            },
            {
                "_id": "614c32aaa1411184356a377d",
                "attributes": [],
                "blockIndex": 328657,
                "blockhash": "0xa0ea92a6e831c758aceb245626cedf10ec8560c52cdec8f9f2fa47f88cd8c560",
                "blocktime": 1631869822520,
                "hash": "0x5396b96e614f44c41a1a0240188d4b8b30c5bd9c9c41557bea05b616eae6f306",
                "netfee": 1240235,
                "nonce": 1873101081,
                "script": "ERHAHwwFc3Rha2UMFKAw1ZtCZVEvjOj0bdKOfrOI0bLhQWJ9W1I=",
                "sender": "NfnYGUoC6TKhevumjTUFyYnfQtZRjMaXch",
                "signers": [
                    {
                        "account": "0x530e099ac607559db8424896ca159c7090cdfad9",
                        "scopes": "Global"
                    }
                ],
                "size": 196,
                "sysfee": 15988875,
                "validUntilBlock": 334245,
                "version": 0,
                "vmstate": "HALT",
                "witnesses": [
                    {
                        "invocation": "DEAumLMvdLgWXTL+l6RhbeOHndU1GW9nRt4RYvgFnXZfmgfIspDd9DmMNAfguHiiRytbvn2vmb8C3ELcYyyneMxk",
                        "verification": "DCED9ocyvgPQfEj4x6JFJKy4jWXlbdZBmiB6H013ZKiB7TRBVuezJw=="
                    }
                ]
            }
        ],
            "totalCount": 20
    },
    "error": null
}
```
