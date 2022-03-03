# GetNFTByWords
Fuzzy search by name of NFT 
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| PrimaryMarket | string | the Primary marketplace hash | optional |
| SecondaryMarket    | string| the Secondary marketplace hash | optional |
| Words | string |    search item  | required |
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |



#### Example
```powershell
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFTByWords",
  "params": {
      "SecondaryMarket":"0x1f594c26a50d25d22d8afc3f1843b4ddb17cf180",
	  "PrimaryMarket":"0x22231899d6946802f66a0fb06ce0960ae88e9eb6",
      "Words":"Blind Box",
      "Skip":0,
      "Limit":2   
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
        "_id": "61d83d880506da899828e21a",
        "amount": "1",
        "asset": "0xd9e2093de3dc2ef7cf5704ceec46ab7fadd48e7f",
        "auctionAmount": "0",
        "auctionAsset": null,
        "auctionType": 0,
        "auctor": null,
        "bidAmount": "0",
        "bidder": null,
        "deadline": 0,
        "image": "https://neo.org/BlindBox.png",
        "market": null,
        "name": "Blind Box #97",
        "number": 97,
        "owner": "0xed369077652ddd55bd7696df93fe49c0bb40d3bc",
        "properties": {
          "image": "https://neo.org/BlindBox.png",
          "number": 97,
          "video": "aHR0cHM6Ly9uZW8ub3JnL0JsaW5kQm94Lm1wNA=="
        },
        "state": "notlisted",
        "timestamp": 1632237992486,
        "tokenid": "QmxpbmQgQm94ICM5Nw=="
      },
      ......
    ],
    "totalCount": 1485
  },
  "error": null
}

```
