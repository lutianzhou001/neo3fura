# GetMarketIndexByAsset
Gets the primary market classification.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string| The asset script hash | Required |
| PrimaryMarket | string| The primary marketplace hash | Required|
| SecondaryMarket | string | The secondary marketplace hash | Required |



### Example

Request body

```powershell
{
  "jsonrpc": "2.0",
  "method": "GetMarketIndexByAsset",
  "params": {     
      "PrimaryMarket":"0xa41600dec34741b143c66f2d3448d15c7d79a0b7",
      "SecondaryMarket":"0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29",
      "AssetHash":"0x19ed09dadac28e6b6a2f76588516ef681aff29b1"
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

