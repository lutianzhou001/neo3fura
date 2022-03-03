# GetBlockInfoList
Gets the blockinfos of the recent blocks
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    name    | type | description | Required |
| ---------- | --- |    ------    | ----|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetBlockInfoList",
  "params": {"Limit":2},
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
        "_id": "6168fe8150025b016127dd62",
        "hash": "0xc6d735b33dad4298f3397ef5d77454e68dcba064ce8897f3ac58a21e442db339",
        "index": 483825,
        "size": 697,
        "timestamp": 1634270848950,
        "transactioncount": 0
      },
      {
        "_id": "6168fe7150025b016127dd59",
        "hash": "0x5c877b9fd7b87955ad98ec63cb843164e04bd92fae5ca88e5a35eb339d47bbcc",
        "index": 483824,
        "size": 697,
        "timestamp": 1634270833849,
        "transactioncount": 0
      }
    ],
    "totalCount": 483826
  },
  "error": null
}
```
