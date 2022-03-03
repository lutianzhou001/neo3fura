# GetVmStateByTransactionHash
Gets the vm state by transaction hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| TransactionHash     | string|  the transactionHash| required|


#### Example
```
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
  "jsonrpc": "2.0",
  "method": "GetVmStateByTransactionHash",
  "params": {"TransactionHash":"0xa15ed65858d1e73a45c5f0f9d29462fe00e1d608a8f471a293eeda80ac28294b"},
  "id": 1
}'
```
### Response
```json5
{
  "id": 1,
  "result": {
    "vmstate": "HALT"
  },
  "error": null
}
```
