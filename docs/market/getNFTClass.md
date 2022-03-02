# GetNFTClass
get PrimaryMarket classification
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| AssetHash     | string|  the asset scriptHash| |
| MarketHash     | string|  | required|
| SubClass     | Array|  | |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTClass",
  "params": {  
      "AssetHash":"0xc7b11b46f97bda7a8c82793841abba120e96695b",
      "SubClass":[["VbdQL2cl8ngkJjITK8aNzeY07PLKiEyiXCORcgw+lfI=","sNU/EpLlV1GuiH4P0zet1rz+SlCb1/2YNucEanpVWIA="],["79WdS6cDK2ZC74UPFlILgiZlus49WkhYo5z8XpR+ckg=","GSDIwJTkjsqbWMQG4eAkPkzCXrTv/390QciVb/B3cow="]],
      "MarketHash":"0xf63cccfe6cfac7ee776dada552b976c74fe5b51a" 
      },
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
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "claimed": 6,  //已经卖掉的数量
        "image": "",
        "name": "sell-1",
        "price": "5",   //售卖价格
        "sellAsset": "0xd2a4cff31913016155e38e474a2c06d08be276cf"  //售卖资产
      },
      ....
    ],
    "totalCount": 2
  },
  "error": null
}
```

