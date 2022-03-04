# GetRawTransactionByTransactionHash
Gets the raw transaction by transactionhash.
<hr>

### Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash     | string| The transaction hash | Required |


### Example

Request body

```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"TransactionHash":"0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287"},
    "method": "GetRawTransactionByTransactionHash"
}'
```

Response body

```json
{
  "id": 1,
  "result": {
    "_id": "614befd9a1411184355217dd",
    "attributes": [],
    "blockIndex": 3762,
    "blockhash": "0xcf35068b43281d700c6c7fc160ab844e74afeda08e793d061bbd1bc1a1203bd4",
    "blocktime": 1626850227986,
    "hash": "0x85b55479fc43668077821234f547824d3111343aec21988f8c0aa1ff9b2ee287",
    "netfee": 8727740,
    "nonce": 1564834642,
    "script": "CwMAQHoQ81oAAAwUwJjkrPCyCQ3Rbss9WN5CaocVhRsMFJKznHeqYPKbV8FzztkX8X/TIabrFMAfDAh0cmFuc2ZlcgwUz3bii9AGLEpHjuNVYQETGfPPpNJBYn1bUjk=",
    "sender": "NL4PXTc8dxjBca8FEkJYCEgDWL98ZnzcnV",
    "signers": [
      {
        "account": "0x0bf916d727c75f2e51e1ab2c476304513da59701",
        "scopes": "None"
      },
      {
        "account": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
        "scopes": "CalledByEntry"
      }
    ],
    "size": 860,
    "sysfee": 9977780,
    "validUntilBlock": 9506,
    "version": 0,
    "vmstate": "HALT",
    "witnesses": [
      {
        "invocation": "DEDEsEdacsBv9knpBjGyAzEADODx3N0cxYRN6nri3f08zl2AmqxRFY/BwVrbKHDhUxU2rouo7jXfOlmHYsHBn+hv",
        "verification": "DCECPpsy6om5TQZuZJsST9UOOW7pE2no4qauGxHBcNAiJW1BVuezJw=="
      },
      {
        "invocation": "DEBPi8c9SXTTzM4ZgHEnrPevPZfXTLxVTSNj34IzTUg7Nha5cyI9AQxb5uaIe6ACF2BzK+rp8hzx5HKQurzJQU4lDEB4ZU/Nf/RrWgAFOhRYMNbB+7uwIfZEUUbAYUbjFSbZb5DweECXOMXOFoy1hagUgbORcd/YXJJnlvQQa6F6KbOKDEDHimMY65M2+6ZEh7/0DDUMNYoybjuHfGS+TF0LuGEruNOyU2pOz0zBpxzi9T0sa/9j2y5Rau50uCCf6c6eVZrGDECxsfC/xRtriRDfn2aV6TtK1BLT7hpX9qm2VFLK3K1X7s09XcA/+cTgTvxFzrHsTnB3Lcnh5/T8hsQLCdd8QOmSDED9G2Ac3l7jK8C3GJvfM35/ePiZXwuUKDa/niS6i5i6RrjH/srg8/6YsrtODqa4vB6uySueMwDqw8iMZlyUtIH4",
        "verification": "FQwhAwCbdUDhDyVi5f2PrJ6uwlFmpYsm5BI0j/WoaSe/rCKiDCEDAgXpzvrqWh38WAryDI1aokaLsBSPGl5GBfxiLIDmBLoMIQIUuvDO6jpm8X5+HoOeol/YvtbNgua7bmglAYkGX0T/AQwhAj6bMuqJuU0GbmSbEk/VDjlu6RNp6OKmrhsRwXDQIiVtDCEDQI3NQWOW9keDrFh+oeFZPFfZ/qiAyKahkg6SollHeAYMIQKng0vpsy4pgdFXy1u9OstCz9EepcOxAiTXpE6YxZEPGwwhAroscPWZbzV6QxmHBYWfriz+oT4RcpYoAHcrPViKnUq9F0Ge0Nw6"
      }
    ]
  },
  "error": null
}
```
