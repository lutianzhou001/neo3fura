# GetNFSImgStatus
get image state by url 
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| Url   | string|the image url  | required|




#### Example
```
curl --location --request GET 'https://testneofura.ngd.network:444' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetNFSImgStatus",
  "params": {"Url":"https://http.testnet.fs.neo.org/C1UKxuvGNNjEHgtGi3YAFSthsfTC9zxJtBh8eXhCmMoi/9cWgnZe75d8X1jhbkYVSVa7ZkmDT5KkeQiiCNwjfJtxC2"},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "ImageStatus": false
  },
  "error": null
}

```
