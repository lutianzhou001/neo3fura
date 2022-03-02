# GetNep11TransferByBlockHeight
Gets the Nep11 transfer by block height
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHeight    | int|  the blockHeight| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```
{  
    "jsonrpc": "2.0",
    "method": "GetNep11TransferByBlockHeight",
    "params": {"BlockHeight":69981},
    "id": 1
}
```
### Response
```json5
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614bfc5f30669383446d6d2f",
        "blockhash": "0x2d3ac96785404ad370f7063db1a11f5b4018ebdd6b80754394360740bcc90c95",
        "contract": "0x4f628a187e133fa98a5fd0795df3065f219e414e",
        "from": null,
        "frombalance": "0",
        "timestamp": 1627871579237,
        "to": "0x0000000000000000000000000000000000000000",
        "tobalance": "1",
        "tokenId": "U9qaPp0ehFf6afitJ+msIqcM/3T+wEWfvFz1/HOYUzw=",
        "txid": "0x0a8809c34ff72f9ce8b670c584ed6416e1b9ad80ab678b29e400c9dc37bde6be",
        "value": "1"
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
