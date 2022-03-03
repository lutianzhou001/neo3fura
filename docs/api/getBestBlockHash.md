# GetBestBlockHash
Gets the best (latest) blockhash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

none


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetBestBlockHash",
  "params": {},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "hash": "0x95ac24ea866de870cf4d664e03c35cfd1a21d377284a1f22dfc8d04501b93a5b"
  },
  "error": null
}
```
