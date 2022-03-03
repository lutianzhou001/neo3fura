# GetNFTRecordByContractHashTokenId
get nft token transfer record 
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contractHash| required |
| MarketHash     | string| the marketplace hash | required |
| TokenId     | string| nft token id | optional |




#### Example
```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTRecordByContractHashTokenId",
  "params": {
      "ContractHash":"0xc7b11b46f97bda7a8c82793841abba120e96695b",   
      "TokenId":"BoN2dx2fSFeRuT7kp87u3e1Jewc3ZIqQ5U0dQSdxofA=",
      "MarketHash":""
  },
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
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": 100,
        "auctionAsset": "0x0daba9cbfa59cf4d43ff1b76d3691725da278450",
        "from": "0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
        "image": "",
        "name": "1",
        "timestamp": 1639380705777,
        "to": "0x78fed05e0ed095b47826bd7461da11c8281195f6",
        "tokenid": "BoN2dx2fSFeRuT7kp87u3e1Jewc3ZIqQ5U0dQSdxofA="
      },
      ......
    ],
    "totalCount": 7
  },
  "error": null
}
```

