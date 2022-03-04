# GetNep11TransferByTransactionHash
Gets the Nep11 transfer information by the transaction hash
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash     | string| The transaction hash | Required |
| Limit    | int|  The number of items to return| Optional|
| Skip    | int|  The number of items to return| Optional |



### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNep11TransferByTransactionHash",
  "params": {"TransactionHash": "0xa15ed65858d1e73a45c5f0f9d29462fe00e1d608a8f471a293eeda80ac28294b",
  "Limit":3 },
  "id": 1
}'
```

Response body

```json
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "614bf4b630669383446822b0",
        "blockhash": "0x645cc9744225cc8222ff18ec112e2e260cbeae0efad1094b9bc98930afb84304",
        "contract": "0x4f628a187e133fa98a5fd0795df3065f219e414e",
        "decimals": 8,
        "from": null,
        "frombalance": "0",
        "timestamp": 1627284349560,
        "to": "0x59057af11833590dff6a8f736fcd5fca46e12289",
        "tobalance": "1",
        "tokenId": "X3685uxIZRNoROOSfzBXJtUWSDMF8jEONGInRzb0KDg=",
        "tokenname": "FBACC.NFT",
        "txid": "0xa15ed65858d1e73a45c5f0f9d29462fe00e1d608a8f471a293eeda80ac28294b",
        "value": "1"
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
