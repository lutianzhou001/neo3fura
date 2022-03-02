# GetMarketTokenidList
get the list of on sale tokenId  in the primary marketplace
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Account     | string|  the user's address| required |
| AssetHash     | string| the token scriptHash | required |
| MarketHash     | string| the marketplace hash | required |
| SubClass     | Array| the nft's tokenid class | optional |




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetMarketTokenidList",
  "params": {
      "Account":"0x2296bd323004b439f46d1557b7b58a8f6cfe36af",
      "AssetHash":"0x0e35991d4eaeea0ff35b4a849342d59c8091de18",
      "MarketHash":"0x0b92cf1c2f308d8084085dc446f7c033b753e959",     
      "SubClass": [
			[
				"TWV0YVBhbmFjZWEgIzEtMDE=",
				"TWV0YVBhbmFjZWEgIzEtMDU="
			],
			[
				"TWV0YVBhbmFjZWEgIzItMDE=",
				"TWV0YVBhbmFjZWEgIzItMDU="
			],
			[
				"TWV0YVBhbmFjZWEgIzMtMDE=",
				"TWV0YVBhbmFjZWEgIzMtMDU="
			],
			[
				"TWV0YVBhbmFjZWEgIzQtMDE=",
				"TWV0YVBhbmFjZWEgIzQtMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzUtMDE=",
				"TWV0YVBhbmFjZWEgIzUtMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzYtMDE=",
				"TWV0YVBhbmFjZWEgIzYtMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzctMDE=",
				"TWV0YVBhbmFjZWEgIzctMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzgtMDE=",
				"TWV0YVBhbmFjZWEgIzgtMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzktMDE=",
				"TWV0YVBhbmFjZWEgIzktMTU="
			],
			[
				"TWV0YVBhbmFjZWEgIzEwLTAx",
				"TWV0YVBhbmFjZWEgIzEwLTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzExLTAx",
				"TWV0YVBhbmFjZWEgIzExLTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzEyLTAx",
				"TWV0YVBhbmFjZWEgIzEyLTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzEzLTAx",
				"TWV0YVBhbmFjZWEgIzEzLTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzE0LTAx",
				"TWV0YVBhbmFjZWEgIzE0LTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzE1LTAx",
				"TWV0YVBhbmFjZWEgIzE1LTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzE2LTAx",
				"TWV0YVBhbmFjZWEgIzE2LTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzE3LTAx",
				"TWV0YVBhbmFjZWEgIzE3LTMz"
			],
			[
				"TWV0YVBhbmFjZWEgIzE4LTAx",
				"TWV0YVBhbmFjZWEgIzE4LTMz"
			]
		]
     
      
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
        "id": 0,
        "tokenid": [
          "TWV0YVBhbmFjZWEgIzEtMDU=",
          "TWV0YVBhbmFjZWEgIzEtMDQ=",
          "TWV0YVBhbmFjZWEgIzEtMDI=",
          "TWV0YVBhbmFjZWEgIzEtMDM=",
          "TWV0YVBhbmFjZWEgIzEtMDE="
        ]
      },
      {
        "id": 1,
        "tokenid": [
          "TWV0YVBhbmFjZWEgIzItMDU=",
          "TWV0YVBhbmFjZWEgIzItMDM=",
          "TWV0YVBhbmFjZWEgIzItMDE=",
          "TWV0YVBhbmFjZWEgIzItMDI=",
          "TWV0YVBhbmFjZWEgIzItMDQ="
        ]
      },
      {
        "id": 2,
        "tokenid": [
          "TWV0YVBhbmFjZWEgIzMtMDI=",
          "TWV0YVBhbmFjZWEgIzMtMDU=",
          "TWV0YVBhbmFjZWEgIzMtMDM=",
          "TWV0YVBhbmFjZWEgIzMtMDQ=",
          "TWV0YVBhbmFjZWEgIzMtMDE="
        ]
      },
      {
        "id": 3,
        "tokenid": [
          "TWV0YVBhbmFjZWEgIzQtMTM=",
          "TWV0YVBhbmFjZWEgIzQtMDk=",
          "TWV0YVBhbmFjZWEgIzQtMTA=",
          "TWV0YVBhbmFjZWEgIzQtMTE=",
          "TWV0YVBhbmFjZWEgIzQtMTI=",
          "TWV0YVBhbmFjZWEgIzQtMDQ=",
          "TWV0YVBhbmFjZWEgIzQtMDg=",
          "TWV0YVBhbmFjZWEgIzQtMDU=",
          "TWV0YVBhbmFjZWEgIzQtMTQ=",
          "TWV0YVBhbmFjZWEgIzQtMDE=",
          "TWV0YVBhbmFjZWEgIzQtMDc=",
          "TWV0YVBhbmFjZWEgIzQtMDM=",
          "TWV0YVBhbmFjZWEgIzQtMDI=",
          "TWV0YVBhbmFjZWEgIzQtMTU=",
          "TWV0YVBhbmFjZWEgIzQtMDY="
        ]
      },
      ......
    ],
    "totalCount": 18
  },
  "error": null
}
```

