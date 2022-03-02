# GetNFTRecordByContractHashTokenId
get nft token info by contractHash and tokenId
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contractHash| required |
| MarketHash     | string| the  marketplace hash | optional |
| TokenIds    | Array| Array of nft token id| optional |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTByContractHashTokenId",
  "params": {
      "ContractHash":"0xc7b11b46f97bda7a8c82793841abba120e96695b",     
      "TokenIds":["LzKk2aeLybZTv83Hzw8djcvJJyVldIyi8oly1qqmqUo="],
      "MarketHash":""
      
      },
  "id": 1
}
'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "result": [
      {
        "_id": "61b724ab0b1347a931f4cae2",
        "amount": "1",
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": "10",
        "auctionAsset": "0xd2a4cff31913016155e38e474a2c06d08be276cf",
        "auctionType": 2,
        "auctor": "0xc65e19cfa66b61800ce582d1b55f4e93fa214b17",
        "bidAmount": "20",
        "bidder": "0x6fd49ab2f14a6bd9a060bb91fdbf29799a885a9e",
        "deadline": 1639565840451,
        "market": "0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
        "owner": "0xf63cccfe6cfac7ee776dada552b976c74fe5b51a",
        "state": "auction",
        "timestamp": 1639393572133,
        "tokenid": "az2dNYa7xEzk2XAQoHnH22k6AbO5/RkyqMDK64VuuXE="
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```

