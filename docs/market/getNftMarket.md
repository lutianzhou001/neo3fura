# GetNFTMarket
Gets the nft token list by contractHash, asset and nftState
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contractHash|optional|
| SkipAssetHash | intstring | the number of items to returnthe asset scriptHash |optional|
| SecondaryMarket     | string| the Secondary marketplace hash | optional |
| SkipPrimaryMarket | intstring | the number of items to returnthe PrimaryMarket marketplace hash | optional |
| Nftstate     | string| 3 types: "auction","sale" or "notlisted"| optional |
| SkipSort | intstring | the number of items to return3 types: "timestamp", "price" or "deadline" | optional |
| Order     | int|  descending sort: -1, ascending sort: +1| optional |
| SkipLimit | intint | the number of items to returnthe number of items to return | optional |
| Skip    | int|  the number of items to return| optional |



#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTMarket",
  "params": {     
      "ContractHash":"",
      "AssetHash":"",
      "SecondaryMarket":"0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180",
	  "PrimaryMarket":"0x1ba667322022693c8629d87b804f5d7730d10779",
      "NFTstate":"",
      "Sort":"",
      "Order":-1,        
      "Skip":0,
      "Limit":2
     
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
        "_id": "62022566491328a3d65b427b",
        "amount": "1",
        "asset": "0x15130d478ec0baaee86f98a75310e431490c3441",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "https://http.fs.neo.org/GD5YUdHWFQfSVmSzgZ55y9akuqHQ8oXVhXnArtv1fLKr/BL6fUdVjxfDysdutBxcB9VkURcrYwUQJt9ttbbjKjg31",
        "listedTimestamp": 0,
        "market": null,
        "name": "MetaPanacea #1",
        "number": 1,
        "owner": "0x17c6a10042e1a92d96ee9544092280a2a6f123e9",
        "properties": {
          "image": "https://http.fs.neo.org/GD5YUdHWFQfSVmSzgZ55y9akuqHQ8oXVhXnArtv1fLKr/JAmFarFV5Pwt83k9rKd8LFwctdkGR6PRNFk3hewLCFJh",
          "number": 1,
          "series": "RWxpeGly",
          "supply": "NQ=="
        },
        "state": "notlisted",
        "timestamp": 1644307812007,
        "tokenid": "TWV0YVBhbmFjZWEgIzEtMDM="
      },
      {
        "_id": "6204ae16491328a3d65bf55a",
        "amount": "1",
        "asset": "0x15130d478ec0baaee86f98a75310e431490c3441",
        "auctionAmount": "0",
        "auctionAmountCond": null,
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "deadline": 0,
        "deadlineCond": null,
        "image": "https://http.fs.neo.org/GD5YUdHWFQfSVmSzgZ55y9akuqHQ8oXVhXnArtv1fLKr/7swpUWx2B9KoeU9ut6eh3nUhPfF2JfckCMec9sX6KWuR",
        "listedTimestamp": 0,
        "market": null,
        "name": "MetaPanacea #8",
        "number": 8,
        "owner": "0x17c6a10042e1a92d96ee9544092280a2a6f123e9",
        "properties": {
          "image": "https://http.fs.neo.org/GD5YUdHWFQfSVmSzgZ55y9akuqHQ8oXVhXnArtv1fLKr/3dR5Xep5iA9Jq9T1EJJutmzWvxYfsw57tqpvBkUSQwND",
          "number": 8,
          "series": "SW5nZW51aXR5",
          "supply": "MTU="
        },
        "state": "notlisted",
        "timestamp": 1644473866421,
        "tokenid": "TWV0YVBhbmFjZWEgIzgtMDE="
      }
    ],
    "totalCount": 3984
  },
  "error": null
}
```
