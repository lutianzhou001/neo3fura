# GetMarketIndexByAsset
get PrimaryMarket classification
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string|  the asset scriptHash| required |
| MarketHash     | string| the marketplace hash | required|




#### Example
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
### Response
```json

{
    "id": 1,
    "result": {
        "auctionAmount":  "5",   //地板价价格   nep17
        "auctionAsset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "conAmount": 3.2919624466782094,    //地板价价格     usd  
        "totalowner": 6,   //owner总量
        "totalsupply": 28, // NFT系列总量
        "totaltxamount":277.03192125227014   // 交易总额  usd
    },
    "error": null
}
```

