# Account API 接口文档

```GO
POST:/GetAccountList:
    functions: get the account list (10 accounts data)
    params:{
        "jsonrpc": "2.0",
        "method": "GetAccountsList",
        "params": {},
        "id": 1
    }
    return sample
    {
        "id": 1,
        "result": {
            "result": [
                {
                    "_id": "60f7d191e95d8c6665d31145",
                    "address": "0x55a314d822f6427ef972a6ab873e8d73f559e907",
                    "firstusetime": 1626853776472
                },
                {
                    "_id": "60f7cba3e95d8c6665d30ee0",
                    "address": "0x259096ffbbb264d0aa0c8638d915409f37ca0c32",
                    "firstusetime": 1626852259206
                },
                {
                    "_id": "60f7c1b9e95d8c6665d30aea",
                    "address": "0xaa7c8fdf05434ed49cdf0878318e16d0b48165fb",
                    "firstusetime": 1626849720384
                },
                {
                    "_id": "60f7bd99e95d8c6665d3094d",
                    "address": "0x996f18bb9845d91b435e31c4e8e342901cef0adb",
                    "firstusetime": 1626848664546
                },
                {
                    "_id": "60f7b854e95d8c6665d30740",
                    "address": "0x6a9d04940ef4ea53e00899ea756a4d815c71ff4a",
                    "firstusetime": 1626847315496
                },
                {
                    "_id": "60f7b17de95d8c6665d304a6",
                    "address": "0x351d77987503e1587eda1596ebe3d4343f5a34b6",
                    "firstusetime": 1626845565054
                },
                {
                    "_id": "60f6ff3ce95d8c6665d2c168",
                    "address": "0x85474174b8a69d802f902885ceb5a6b4fbd62adb",
                    "firstusetime": 1626799931147
                },
                {
                    "_id": "60f6cd5fe95d8c6665d2ae89",
                    "address": "0xd74fc31fd8ccdfcb39e62aee277eafecd6c0405c",
                    "firstusetime": 1626787165203
                },
                {
                    "_id": "60f6b72ce95d8c6665d2a5ce",
                    "address": "0xe0bc595e7d3dd11e0a3c06655c389fd8c29364eb",
                    "firstusetime": 1626781483195
                },
                {
                    "_id": "60f6aeb9e95d8c6665d2a24e",
                    "address": "0xc374dd30b24acff8a62eca276b941e3d32e130f3",
                    "firstusetime": 1626779320394
                },
                {
                    "_id": "60f6a617e95d8c6665d29e7d",
                    "address": "0x87aff3f339ada6833d84cda12d847d49de2a39a7",
                    "firstusetime": 1626777110906
                },
                {
                    "_id": "60f6901ae95d8c6665d295a2",
                    "address": "0xf9b9def532bb4c27ec7cdd239a50b38e0c884592",
                    "firstusetime": 1626771481337
                },
                {
                    "_id": "60f68fd9e95d8c6665d29580",
                    "address": "0x329f0f21416c154239ed37e4a41ec189465b5404",
                    "firstusetime": 1626771416485
                },
                {
                    "_id": "60f68fc9e95d8c6665d29569",
                    "address": "0x9245880c8eb3509a23dd7cec274cbb32f5deb9f9",
                    "firstusetime": 1626771400603
                },
                {
                    "_id": "60f666d4e95d8c6665d230fc",
                    "address": "0x5ca28f4afb47d220aa46056d44102644c114d8d4",
                    "firstusetime": 1626703527812
                },
                {
                    "_id": "60f666abe95d8c6665d223a5",
                    "address": "0x7b8802d127b3db23d5efbf215b1976993d2e53df",
                    "firstusetime": 1626694579921
                },
                {
                    "_id": "60f666aae95d8c6665d22389",
                    "address": "0x118ba6f59931a56ec469770f7fc790ece96df00d",
                    "firstusetime": 1626694532413
                },
                {
                    "_id": "60f66653e95d8c6665d2072d",
                    "address": "0x03357a220f33770a928284a1ac68639e88a734dc",
                    "firstusetime": 1626675490557
                },
                {
                    "_id": "60f66645e95d8c6665d202b5",
                    "address": "0x122296f233885ed95e392425cca60e162d493e24",
                    "firstusetime": 1626672414918
                },
                {
                    "_id": "60f6661ae95d8c6665d1f417",
                    "address": "0x851bb742cb2069d86d44ce7b5323fc74fe223b46",
                    "firstusetime": 1626662419016
                }
            ],
            "totalCount": 568
        },
        "error": null
    }
    
    GetAccountInfoByAddress:
    functions: get a account info by its address
    params:{
        "jsonrpc": "2.0",
        "method": "GetAccountInfoByAddress",
        "params": {
            "address": "0x55a314d822f6427ef972a6ab873e8d73f559e907"
         },
        "id": 1
    }
    return sample
    {
        "id": 1,
        "result": {
            "_id": "60f7d191e95d8c6665d31145",
            "firstusetime": 1626853776472
        },
        "error": null
    }


```	
	
