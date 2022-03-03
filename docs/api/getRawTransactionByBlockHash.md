# GetRawTransactionByBlockHash
Gets the transaction by blockhash
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
  "jsonrpc": "2.0",
  "method": "GetRawTransactionByBlockHash",
  "params": {"BlockHash":"0xe19cdbf573086552cf4e9a1dd0cc3402bef246acbf2810822fa4a03d1ca05edc"},
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
        "_id": "61dd02b98f5f4a753c341119",
        "attributes": [],
        "blockIndex": 968949,
        "blockhash": "0xe19cdbf573086552cf4e9a1dd0cc3402bef246acbf2810822fa4a03d1ca05edc",
        "blocktime": 1641874102823,
        "hash": "0x3817acb8256387f557dd1cc8ea1e3e15d28e5c608ae8a8b57ccf09edff521b7b",
        "netfee": 122862,
        "nonce": 132167099,
        "script": "CxoMFPj5LiizJYNRp/4eFTCx2l6TliqSDBSUGubUQGFEB9XUjEJMKgNuectOfxTAHwwIdHJhbnNmZXIMFPVj6kC8KD1NDgXEjqMFs/Kgc0DvQWJ9W1I=",
        "sender": "NZR5RJpeRqjP3aHFGoiKMDkCLfgxRuMLzh",
        "signers": [
          {
            "account": "0x7f4ecb796e032a4c428cd4d507446140d4e61a94",
            "scopes": "CalledByEntry"
          }
        ],
        "size": 244,
        "sysfee": 11119242,
        "validUntilBlock": 968979,
        "version": 0,
        "witnesses": [
          {
            "invocation": "DEBCJc+zFpa1dyJTnxXTHPJUHyUrBIvQR3jICAc5C3qnEpysEuWp2sp7rlCVCOGaEGVyfnoTcMoyRJAS+OzKuBcE",
            "verification": "DCECt0BATPr7LE+oXPh9U2qUxCOImTPmF243WuXMO3ep9wVBVuezJw=="
          }
        ]
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
