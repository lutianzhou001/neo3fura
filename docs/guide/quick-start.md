# Quick Start

### Send a request refer to neo rpc doc

Neofura supports all neo rpc requests from [https://docs.neo.org/docs/zh-cn/reference/rpc/latest-version/api.html](https://docs.neo.org/docs/zh-cn/reference/rpc/latest-version/api.html). You can easily test it with Postman or any programming language you like.

You can send a request like this
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
		"jsonrpc": "2.0",
		"method":  "getcontractstate",
		"params": ["0xfe924b7cfe89ddd271abaf7210a80a7e11178758"],
		"id": 1
	}'
```

And you can easily get the response

```
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "id": -9,
        "updatecounter": 0,
        "hash": "0xfe924b7cfe89ddd271abaf7210a80a7e11178758",
        "nef": {
            "magic": 860243278,
            "compiler": "neo-core-v3.0",
            "source": "",
            "tokens": [],
            "script": "EEEa93tnQBBBGvd7Z0AQQRr3e2dAEEEa93tnQBBBGvd7Z0A=",
            "checksum": 2663858513
        },
        "manifest": {
            "name": "OracleContract",
            "groups": [],
            "features": {},
            "supportedstandards": [],
            "abi": {
                "methods": [
                    {
                        "name": "finish",
                        "parameters": [],
                        "returntype": "Void",
                        "offset": 0,
                        "safe": false
                    },
                    {
                        "name": "getPrice",
                        "parameters": [],
                        "returntype": "Integer",
                        "offset": 7,
                        "safe": true
                    },
                    {
                        "name": "request",
                        "parameters": [
                            {
                                "name": "url",
                                "type": "String"
                            },
                            {
                                "name": "filter",
                                "type": "String"
                            },
                            {
                                "name": "callback",
                                "type": "String"
                            },
                            {
                                "name": "userData",
                                "type": "Any"
                            },
                            {
                                "name": "gasForResponse",
                                "type": "Integer"
                            }
                        ],
                        "returntype": "Void",
                        "offset": 14,
                        "safe": false
                    },
                    {
                        "name": "setPrice",
                        "parameters": [
                            {
                                "name": "price",
                                "type": "Integer"
                            }
                        ],
                        "returntype": "Void",
                        "offset": 21,
                        "safe": false
                    },
                    {
                        "name": "verify",
                        "parameters": [],
                        "returntype": "Boolean",
                        "offset": 28,
                        "safe": true
                    }
                ],
                "events": [
                    {
                        "name": "OracleRequest",
                        "parameters": [
                            {
                                "name": "Id",
                                "type": "Integer"
                            },
                            {
                                "name": "RequestContract",
                                "type": "Hash160"
                            },
                            {
                                "name": "Url",
                                "type": "String"
                            },
                            {
                                "name": "Filter",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "name": "OracleResponse",
                        "parameters": [
                            {
                                "name": "Id",
                                "type": "Integer"
                            },
                            {
                                "name": "OriginalTx",
                                "type": "Hash256"
                            }
                        ]
                    }
                ]
            },
            "permissions": [
                {
                    "contract": "*",
                    "methods": "*"
                }
            ],
            "trusts": [],
            "extra": null
        }
    }
}
```

### Send a request refer to neofura doc

Refer to JSON-RPC Requests in this doc to make a request by yourself!
