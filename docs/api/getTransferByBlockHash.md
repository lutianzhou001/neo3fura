# GetTransferByBlockHash
Gets the transfer by blockhash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHash      | string|  the blockHash of the block| required |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc":"2.0",
    "params":{"BlockHash":"0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f"},
    "method": "GetTransferByBlockHash",
    "id":1
}'
```
### Response
```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614befe4a141118435521a3c",
        "blockhash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
        "contract": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "from": null,
        "frombalance": "0",
        "timestamp": 1626851177411,
        "to": "0x8b915b5abcb81841face2afc42982c08a7e72b81",
        "tobalance": "9151230520",
        "txid": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "value": "50000000"
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
