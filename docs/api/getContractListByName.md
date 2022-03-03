# GetContractListByName
Gets the contract list by name (fuzzy search supported)
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| Name     | string|  the contract name (fuzzy search supported)| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
``` powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetContractListByName",
  "params": {"Name":"PriceFeedService"},
  "id": 1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "Transaction": [
          {
            "sender": "NTAv9Q5p9Vsckku56sbWSHQBkg3c5ntBk1"
          }
        ],
        "createtime": 1631174537099,
        "hash": "0x89d9839aa840a0bc55b64501faeac3ab037f471d",
        "id": 142,
        "name": "PriceFeedService",
        "updatecounter": 0
      },
      {
        "Transaction": [
          {
            "sender": "NN38jUtTP68pBjUx1pXEAFbZqU9anjqGT1"
          }
        ],
        "createtime": 1631165194829,
        "hash": "0xd30ed1c087d8b8077275f2c7be90f80b5a9c5d8d",
        "id": 137,
        "name": "PriceFeedService",
        "updatecounter": 0
      }
    ],
    "totalCount": 2
  },
  "error": null
}
```
