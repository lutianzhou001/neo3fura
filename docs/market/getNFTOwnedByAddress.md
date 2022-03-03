# GetNFTOwnedByAddress
Gets the nft token list and nft token state by user's address
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contractHash| optional |
| AssetHash     | string|  the asset scriptHash| optional |
| Address     | string|  the user's address| required|
| MarketHash     | string| the marketplace hash | optional |
| Nftstate     | string| 3 types: "auction","sale" or "notlisted"| optional |
| Sort     | string| 4 types: "timestamp", "price", "deadline" or "unClaimed"| optional |
| Order     | int|  descending sort: -1, ascending sort: +1| optional |
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |



#### Example
```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTOwnedByAddress",
  "params": {
      "Address":"0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2", 
      "ContractHash":"0xc7b11b46f97bda7a8c82793841abba120e96695b",
      "AssetHash":"",
     " MarketHash":"0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180",
      "NFTstate":"notlisted",
      "Sort":"",
      "Order":1,       
      "Skip":0,
      "Limit":10
     
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
        "_id": "61d96d19b145511ecccef234",
        "amount": "1",
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "date": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "",
        "listedTimestamp": 0,
        "market": null,
        "name": "1",
        "number": -1,
        "owner": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "properties": {},
        "state": "notlisted",
        "timestamp": 1637824807645,
        "tokenid": "gV8ed6v25JKo46osCoMZn47h+6PJjOJW6/8nLWxMCvk="
      },
      {
        "_id": "61d96d1bb145511ecccef456",
        "amount": "1",
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "date": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "",
        "listedTimestamp": 0,
        "market": null,
        "name": "1",
        "number": -1,
        "owner": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "properties": {},
        "state": "notlisted",
        "timestamp": 1638153362135,
        "tokenid": "2alxWdHMQJ5HM52ia1XmGPeIZ/p1+J/r5bNm00hsd/M="
      },
      {
        "_id": "61d977f2b145511eccd937d1",
        "amount": "1",
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "date": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "",
        "listedTimestamp": 0,
        "market": null,
        "name": "test1",
        "number": -1,
        "owner": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "properties": {},
        "state": "notlisted",
        "timestamp": 1639018515616,
        "tokenid": "guBBkbEd4SYIN39QedDMAABzZTpF24HVGWOTY+bgc7s="
      },
      {
        "_id": "61d9780db145511eccd95624",
        "amount": "1",
        "asset": "0xc7b11b46f97bda7a8c82793841abba120e96695b",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "date": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "",
        "listedTimestamp": 0,
        "market": null,
        "name": "test3",
        "number": -1,
        "owner": "0xdd58b7a05fd9b58a6ec36d6401a89ff2cda224a2",
        "properties": {},
        "state": "notlisted",
        "timestamp": 1639031005657,
        "tokenid": "skmHnC2EQuTXH4E5q8RtSTF1FbEY2IDXVaFa5gJL88s="
      }
    ],
    "totalCount": 4
  },
  "error": null
}
```
###  
