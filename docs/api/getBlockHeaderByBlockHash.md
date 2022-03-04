# GetBlockHeaderByBlockHash
Gets the block header by the blockhash.

<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHash      | string|  The blockHash of the block| Required |

### Example

Request body

```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"BlockHash":"0x7688cf2521bbb5274c22363350539f402e4614a015d9e62b63694c049dec89d6"},
    "method": "GetBlockHeaderByBlockHash"
}'
```

Response body

```json5

{
    "id": 1,
        "result": {
        "_id": "6167edfd50025b0161276efe",
            "hash": "0x7688cf2521bbb5274c22363350539f402e4614a015d9e62b63694c049dec89d6",
            "index": 479350,
            "merkleroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "nextConsensus": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
            "prevhash": "0xd53f16289fed00b644610fc99faba3cc3858c2731ef2461650fb9f5d09197386",
            "primaryindex": 4,
            "size": 696,
            "timestamp": 1634201085495,
            "version": 0,
            "witnesses": [
            {
                "invocation": "DEAr5dDUtswZLWVgDcPrVrz7g1sXWvttTxE404oTqp29z+v1dhd3aVLYYKUVfEJnAAugIVvLQTnNyRdUOJvZr6WGDEDA/19IPGoePVuY9hDIpGMgqo1O/63LBoedrt9oqS5t5bYMegiAJJt0EkRj6hYIMIXmpM5lXcK6J4QzWmZhSyXvDEApL0zL24jCCZYPN/8hkJ8yGB/aMTLrFuFcrmY01Cx/oLK+W2ts1zj3foBmkRGV9Gwz/nN5bX4sotKJqBKUbbp7DEAdgG4R2sz2btsqNELvxpzDIuBfHrn1by4vSpRujOsIzLKexNKWr1z8ZLz0DAykZLNLjUel/KT8/UXzUyC2Txb3DECyNk2SVLAaclwreQaIO6ngvzLZeejzgXKJIshrqAfZ2aLNo1pJjOvXDb69qre2nEARMP0842SM71dTk0qco8NG",
                "verification": "FQwhAwCbdUDhDyVi5f2PrJ6uwlFmpYsm5BI0j/WoaSe/rCKiDCEDAgXpzvrqWh38WAryDI1aokaLsBSPGl5GBfxiLIDmBLoMIQIUuvDO6jpm8X5+HoOeol/YvtbNgua7bmglAYkGX0T/AQwhAj6bMuqJuU0GbmSbEk/VDjlu6RNp6OKmrhsRwXDQIiVtDCEDQI3NQWOW9keDrFh+oeFZPFfZ/qiAyKahkg6SollHeAYMIQKng0vpsy4pgdFXy1u9OstCz9EepcOxAiTXpE6YxZEPGwwhAroscPWZbzV6QxmHBYWfriz+oT4RcpYoAHcrPViKnUq9F0Ge0Nw6"
            }
        ]
    },
    "error": null
}
```
