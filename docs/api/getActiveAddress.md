# GetActiveAddresses
Gets the active address count in several days
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description |  Required |
| ---------- | --- |    ------    | --------|
| Days       | int|  the number of recent days| required|


#### Example
``` powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"Days":5},
    "method": "GetActiveAddresses"
}'
```
### Response
```json
{
    "id": 1,
    "result": {
        "result": [
            {
                "_id": "614bfc5f30669383446d6d30",
                "blockhash": "0x2d3ac96785404ad370f7063db1a11f5b4018ebdd6b80754394360740bcc90c95",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "from": null,
                "frombalance": "0",
                "timestamp": 1627871579237,
                "to": "0xd9294b3f248b47cca2aa24fb47ece44bb5f9c1fe",
                "tobalance": "171928647550",
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "value": "649260"
            },
            {
                "_id": "614bfc5f30669383446d6d31",
                "blockhash": "0x2d3ac96785404ad370f7063db1a11f5b4018ebdd6b80754394360740bcc90c95",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "from": null,
                "frombalance": "0",
                "timestamp": 1627871579237,
                "to": "0x4889dca2933fb56e297035c3ec921af1d0394ba6",
                "tobalance": "131357741660",
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "value": "50000000"
            },
            {
                "_id": "614bfc5f30669383446d6d32",
                "blockhash": "0x2d3ac96785404ad370f7063db1a11f5b4018ebdd6b80754394360740bcc90c95",
                "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
                "from": "0x59057af11833590dff6a8f736fcd5fca46e12289",
                "frombalance": "62598983980",
                "timestamp": 1627871579237,
                "to": null,
                "tobalance": "0",
                "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "value": "23100975"
            }
        ],
        "totalCount": 3
    },
    "error": null
}
```
