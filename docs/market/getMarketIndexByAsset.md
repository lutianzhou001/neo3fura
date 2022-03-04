# GetMarketIndexByAsset
Gets the primary market classification.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string| The asset script hash | Required |
| MarketHash     | string| The marketplace hash | Required|

### Example

Request body

```powershell
{
  "jsonrpc": "2.0",
  "method": "GetMarketIndexByAsset",
  "params": {     
      "MarketHash":"0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
      "AssetHash":"0xc7b11b46f97bda7a8c82793841abba120e96695b"
      },      
  "id": 1
}
```
Response body

```json

{
    "id": 1,
    "result": {
        "auctionAmount":  "5",   
        "auctionAsset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "conAmount": 3.2919624466782094,    //The floor price (usd)  
        "totalowner": 6,   
        "totalsupply": 28, 
        "totaltxamount":277.03192125227014   
    },
    "error": null
}
```

