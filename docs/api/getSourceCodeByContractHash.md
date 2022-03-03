# GetSourceCodeByContractHash
Gets the source code of contract by the contract script hash
<hr>

### Request

> POST https://testneofura.ngd.network:444

#### Body Parameters

|    Name    | Type | Description | Required |
| ---------- | --- |    ------    | ----|
| ContractHash     | string|  the contract script hash| required|
| UpdateCounter     | string|  the number of times the contract hash been updated| required|
| Limit    | int|  the number of items to return| optional|
| Skip    | int|  the number of items to return| optional |


#### Example
```powershell
curl --location --request POST 'https://testneofura.ngd.network:444' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method":"GetSourceCodeByContractHash",
    "params": {
        "ContractHash":"0x04349971c9e5db2411e6c85dcf4d759510e72dcf",
        "UpdateCounter" :0
    }
}'
```
### Response
```json
{
  "id": null,
  "result": {
    "result": [
      {
        "_id": "61dcfc3e915b2f62ee1bfb5e",
        "code": "from boa3.builtin import NeoMetadata, metadata, public\nfrom boa3.builtin.interop.storage import put\n\n\n@public\ndef Main():\n    put('hello', 'world')\n\n\n@metadata\ndef manifest() -> NeoMetadata:\n    meta = NeoMetadata()\n    meta.author = \"COZ in partnership with Simpli\"\n    meta.email = \"contact@coz.io\"\n    meta.description = 'This is a contract example'\n    return meta\n",
        "filename": "helloworld.py",
        "hash": "0x04349971c9e5db2411e6c85dcf4d759510e72dcf",
        "updatecounter": 0
      }
    ],
    "totalCount": 1
  },
  "error": null
}
```
