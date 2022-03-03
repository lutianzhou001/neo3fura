# GetBlockByBlockHeight
Gets the block the blockheight
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| BlockHeight      | int|  the blockHeight| required |


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "params": {"BlockHeight":3823},
    "method": "GetBlockByBlockHeight"
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "_id": "614befe4a141118435521a42",
    "hash": "0xf6ba8db5c013834890903a30a4ce0d65ec5da2addaf4799f15efbedaff42c56f",
    "index": 3823,
    "merkleroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "networkFee": 0,
    "nextConsensus": "0xeba621d37ff117d9ce73c1579bf260aa779cb392",
    "nonce": "2433887357192660069",
    "prevhash": "0x51c62c117ce9e1aa6669a0842d899d637d9cad848caf2b3550dab28445879abb",
    "primary": 1,
    "size": 697,
    "systemFee": 0,
    "timestamp": 1626851177411,
    "version": 0,
    "witnesses": [
      {
        "invocation": "DEDBfA7yCX8/k3DnUVYosDJPpfGi6jo4T2sNPYcTRsNYzLMviMcbeBnWJx0UbZGeIFj6NM2C0PswVo7ELkjmgWqGDEAQAN6Ur9VLqLSLigM3QEb2MBptfTOTFlEq7DYQ/yukMqkgYamz7o0ECiksTUxSK3B7A9/GtmI7dmc2WlQ8AVNgDECSi3z+UotVFrOyM8Q57uIJ5s+jbKl0l3qn5aYNPbKkpcxCVmZe1gKiAIkvq0M+HKYJnmNyVjvMP45MZ1isrwFUDECECDYBDaU1WcvMMnbq7YrpSeSyBj7xRtgaD4ISvIKqrA6LYhea96YlsQDuSjuHBlZ/tH3I1AQzZpMBf24yyU4dDECvTQ7dUAr/B+sVnW5CYb3mzuGowtFQ5XJH6R3KUBiLhw6aVD3SrtoE+Z39vlWCrPjTwX0DNlP4iN5INwKJcU3p",
        "verification": "FQwhAwCbdUDhDyVi5f2PrJ6uwlFmpYsm5BI0j/WoaSe/rCKiDCEDAgXpzvrqWh38WAryDI1aokaLsBSPGl5GBfxiLIDmBLoMIQIUuvDO6jpm8X5+HoOeol/YvtbNgua7bmglAYkGX0T/AQwhAj6bMuqJuU0GbmSbEk/VDjlu6RNp6OKmrhsRwXDQIiVtDCEDQI3NQWOW9keDrFh+oeFZPFfZ/qiAyKahkg6SollHeAYMIQKng0vpsy4pgdFXy1u9OstCz9EepcOxAiTXpE6YxZEPGwwhAroscPWZbzV6QxmHBYWfriz+oT4RcpYoAHcrPViKnUq9F0Ge0Nw6"
      }
    ]
  },
  "error": null
}
```
