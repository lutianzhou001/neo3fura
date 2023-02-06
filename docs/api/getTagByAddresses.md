# GetTagByAddresses
Gets the tag information by the address in batches
<hr>

### Parameters

| Name    | Type     | Description | Required |
|---------|----------|    ------    | ----|
| Address | []string | The user's address| Required |

### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "GetTagByAddresses",
    "params": {
        "Address": [
                    "0x4d5a85b0c83777df72cfb665a933970e4e20c0ec",
                    "0x6da47e57c7f3d43a8023bfdf4bc932d87d9e734f",
                    "0xb8df435a39b9141b2e10b32af28bbde97c992550"
                   ]
    },
    "id": 1
}'
```

Response body

```json5
{
  "id": 1,
  "result": {
    "result": [
      {
        "address": "0x9f8f056a53e39585c7bb52886418c7bed83d126b",
        "bneoSum": 0,
        "ft_tag": "Semi-Experienced Trader",
        "ft_total": "71585",
        "neoSum": 71585,
        "nft_tag": "",
        "nft_total": 0
      },
      {
        "address": "0x6da47e57c7f3d43a8023bfdf4bc932d87d9e734f",
        "bneoSum": 0,
        "ft_tag": "",
        "ft_total": "28",
        "neoSum": 28,
        "nft_tag": "",
        "nft_total": 0
      },
      {
        "address": "0xb8df435a39b9141b2e10b32af28bbde97c992550",
        "bneoSum": 0,
        "ft_tag": "",
        "ft_total": "0",
        "neoSum": 0,
        "nft_tag": "NFT Player",
        "nft_total": 1
      }
    ],
    "totalCount": 3
  },
  "error": null
}
```
