# GetAddressList
Gets the list of addresses
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    |------|
| Limit      | int|  the number of items to return| optional|
| Skip      | int|  the number of items to return| optional|


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetAddressList",
  "params": {"Limit":2,"Skip":2},
  "id": 1
}'
```
### Response
```json5
{
    "id": 1,
        "result": {
        "result": [
            {
                "_id": "61656ccb0f08664e4d486554",
                "address": "0xbb0d6102deb178ec62b56c163796bd3d33ff6884",
                "firstusetime": 1634036938160
            },
            {
                "_id": "616526240f08664e4d48451a",
                "address": "0x550f5098ea3647744d699c851733c397647c39b8",
                "firstusetime": 1634018852638
            }
        ],
            "totalCount": 721
    },
    "error": null
}
```
