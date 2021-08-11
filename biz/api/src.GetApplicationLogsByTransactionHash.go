package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"neo3fura/lib/type/h256"
	"neo3fura/var/stderr"
)

func (me *T) GetApplicationLogsByTransactionHash(args struct {
	TransactionHash h256.T
	Filter          map[string]interface{}
}, ret *json.RawMessage) error {
	if args.TransactionHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.TransactionHash.IsZero() == true {
		return stderr.ErrZero
	}
	r1, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Execution",
		Index:      "GetApplicationLogByTransactionHash",
		Sort:       bson.M{},
		Filter:     bson.M{"txid": args.TransactionHash.Val()},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	 r2, err :=me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "Notification",
			Index:      "someIndex",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:  []bson.M{
				bson.M{"$match": bson.M{"txid": r1["txid"].(string),"blockhash": r1["blockhash"].(string)}},
				bson.M{"$lookup": bson.M{
					"from": "Contract",
					"localField": "contract",
					"foreignField": "hash",
					"as": "Contract"}},

				bson.M{"$project": bson.M{
					"_id":0,
					"Contract.nef":1,
					"contract":1,
					"eventname":1,
					"state":1,
					"timestamp":1,
					"Vmstate":1,
				},
				},},
			Query: []string{},
		}, ret)
	if err != nil {
		return err
	}

	//目前输出格式，此处只用r2中的一条数据作为例子


	//{
	//	"Vmstate": "HALT",
	//	"contract": "0x618d44dc3af16c6120dbf65402024f40a04f772a",
	//	"eventname": "event",
	//	"manifest": [
	//{
	//"manifest": "{\"name\":\"CCMC\",\"groups\":[],\"features\":{},\"supportedstandards\":[],\"abi\":{\"methods\":[{\"name\":\"verify\",\"parameters\":[],\"returntype\":\"Boolean\",\"offset\":0,\"safe\":true},{\"name\":\"update\",\"parameters\":[{\"name\":\"nefFile\",\"type\":\"ByteArray\"},{\"name\":\"manifest\",\"type\":\"String\"}],\"returntype\":\"Void\",\"offset\":66,\"safe\":false},{\"name\":\"isOwner\",\"parameters\":[{\"name\":\"scriptHash\",\"type\":\"Hash160\"}],\"returntype\":\"Boolean\",\"offset\":104,\"safe\":true},{\"name\":\"getOwner\",\"parameters\":[],\"returntype\":\"Hash160\",\"offset\":14,\"safe\":true},{\"name\":\"setOwner\",\"parameters\":[{\"name\":\"ownerScriptHash\",\"type\":\"Hash160\"}],\"returntype\":\"Boolean\",\"offset\":124,\"safe\":false},{\"name\":\"tryDeserializeHeader\",\"parameters\":[{\"name\":\"Source\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":175,\"safe\":true},{\"name\":\"getBookKeepers\",\"parameters\":[],\"returntype\":\"Any\",\"offset\":834,\"safe\":true},{\"name\":\"changeBookKeeper\",\"parameters\":[{\"name\":\"rawHeader\",\"type\":\"ByteArray\"},{\"name\":\"pubKeyList\",\"type\":\"ByteArray\"},{\"name\":\"signList\",\"type\":\"ByteArray\"}],\"returntype\":\"Boolean\",\"offset\":860,\"safe\":false},{\"name\":\"isGenesised\",\"parameters\":[],\"returntype\":\"Boolean\",\"offset\":1324,\"safe\":true},{\"name\":\"tryVerifyPubkey\",\"parameters\":[{\"name\":\"pubKeyList\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":2274,\"safe\":true},{\"name\":\"compressMCPubKey\",\"parameters\":[{\"name\":\"key\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":1719,\"safe\":true},{\"name\":\"getCompressPubKey\",\"parameters\":[{\"name\":\"key\",\"type\":\"ByteArray\"}],\"returntype\":\"PublicKey\",\"offset\":2001,\"safe\":true},{\"name\":\"crossChain\",\"parameters\":[{\"name\":\"toChainID\",\"type\":\"Integer\"},{\"name\":\"toChainAddress\",\"type\":\"ByteArray\"},{\"name\":\"functionName\",\"type\":\"ByteArray\"},{\"name\":\"args\",\"type\":\"ByteArray\"},{\"name\":\"caller\",\"type\":\"ByteArray\"}],\"returntype\":\"Boolean\",\"offset\":2290,\"safe\":false},{\"name\":\"verifyAndExecuteTx\",\"parameters\":[{\"name\":\"proof\",\"type\":\"ByteArray\"},{\"name\":\"RawHeader\",\"type\":\"ByteArray\"},{\"name\":\"headerProof\",\"type\":\"ByteArray\"},{\"name\":\"currentRawHeader\",\"type\":\"ByteArray\"},{\"name\":\"signList\",\"type\":\"ByteArray\"}],\"returntype\":\"Boolean\",\"offset\":2645,\"safe\":false},{\"name\":\"verifySigWithOrder\",\"parameters\":[{\"name\":\"rawHeader\",\"type\":\"ByteArray\"},{\"name\":\"signList\",\"type\":\"ByteArray\"},{\"name\":\"keepers\",\"type\":\"Array\"}],\"returntype\":\"Boolean\",\"offset\":2104,\"safe\":false},{\"name\":\"verifySigWithOrderForHashTest\",\"parameters\":[{\"name\":\"rawHeader\",\"type\":\"ByteArray\"},{\"name\":\"signList\",\"type\":\"ByteArray\"}],\"returntype\":\"Boolean\",\"offset\":4194,\"safe\":false},{\"name\":\"tryDeserializeMerkleValue\",\"parameters\":[{\"name\":\"Source\",\"type\":\"ByteArray\"}],\"returntype\":\"Boolean\",\"offset\":4454,\"safe\":false},{\"name\":\"merkleProve\",\"parameters\":[{\"name\":\"path\",\"type\":\"ByteArray\"},{\"name\":\"root\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":3491,\"safe\":true},{\"name\":\"hashChildren\",\"parameters\":[{\"name\":\"v\",\"type\":\"ByteArray\"},{\"name\":\"hash\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":3662,\"safe\":true},{\"name\":\"hashLeaf\",\"parameters\":[{\"name\":\"value\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":3640,\"safe\":true},{\"name\":\"hash256\",\"parameters\":[{\"name\":\"message\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":2253,\"safe\":true},{\"name\":\"hash160\",\"parameters\":[{\"name\":\"message\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":2069,\"safe\":true},{\"name\":\"readHash\",\"parameters\":[{\"name\":\"Source\",\"type\":\"ByteArray\"},{\"name\":\"offset\",\"type\":\"Integer\"}],\"returntype\":\"ByteArray\",\"offset\":558,\"safe\":true},{\"name\":\"writeUint16\",\"parameters\":[{\"name\":\"Source\",\"type\":\"ByteArray\"},{\"name\":\"value\",\"type\":\"Integer\"}],\"returntype\":\"ByteArray\",\"offset\":1562,\"safe\":true},{\"name\":\"writeVarBytes\",\"parameters\":[{\"name\":\"Source\",\"type\":\"ByteArray\"},{\"name\":\"Content\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":1779,\"safe\":true},{\"name\":\"readVarBytes\",\"parameters\":[{\"name\":\"buffer\",\"type\":\"ByteArray\"},{\"name\":\"offset\",\"type\":\"Integer\"}],\"returntype\":\"Array\",\"offset\":590,\"safe\":true},{\"name\":\"readVarInt\",\"parameters\":[{\"name\":\"buffer\",\"type\":\"ByteArray\"},{\"name\":\"offset\",\"type\":\"Integer\"}],\"returntype\":\"Array\",\"offset\":633,\"safe\":true},{\"name\":\"writeVarInt\",\"parameters\":[{\"name\":\"value\",\"type\":\"Integer\"},{\"name\":\"source\",\"type\":\"ByteArray\"}],\"returntype\":\"ByteArray\",\"offset\":1792,\"safe\":true},{\"name\":\"readBytes\",\"parameters\":[{\"name\":\"buffer\",\"type\":\"ByteArray\"},{\"name\":\"offset\",\"type\":\"Integer\"},{\"name\":\"count\",\"type\":\"Integer\"}],\"returntype\":\"Array\",\"offset\":754,\"safe\":true},{\"name\":\"padRight\",\"parameters\":[{\"name\":\"value\",\"type\":\"ByteArray\"},{\"name\":\"length\",\"type\":\"Integer\"}],\"returntype\":\"ByteArray\",\"offset\":1597,\"safe\":true},{\"name\":\"convertUintToByteArray\",\"parameters\":[{\"name\":\"unsignNumber\",\"type\":\"Integer\"}],\"returntype\":\"ByteArray\",\"offset\":1678,\"safe\":false},{\"name\":\"_initialize\",\"parameters\":[],\"returntype\":\"Void\",\"offset\":4483,\"safe\":false}],\"events\":[{\"name\":\"CrossChainLockEvent\",\"parameters\":[{\"name\":\"arg1\",\"type\":\"ByteArray\"},{\"name\":\"arg2\",\"type\":\"ByteArray\"},{\"name\":\"arg3\",\"type\":\"Integer\"},{\"name\":\"arg4\",\"type\":\"ByteArray\"},{\"name\":\"arg5\",\"type\":\"ByteArray\"}]},{\"name\":\"CrossChainUnlockEvent\",\"parameters\":[{\"name\":\"arg1\",\"type\":\"ByteArray\"},{\"name\":\"arg2\",\"type\":\"ByteArray\"},{\"name\":\"arg3\",\"type\":\"ByteArray\"}]},{\"name\":\"ChangeBookKeeperEvent\",\"parameters\":[{\"name\":\"arg1\",\"type\":\"Integer\"},{\"name\":\"arg2\",\"type\":\"ByteArray\"}]},{\"name\":\"event\",\"parameters\":[{\"name\":\"obj\",\"type\":\"String\"}]}]},\"permissions\":[{\"contract\":\"*\",\"methods\":\"*\"}],\"trusts\":[],\"extra\":{}}"
	//}
	//],
	//"state": {
	//"type": "Array",
	//"value": [
	//{
	//"type": "ByteString",
	//"value": "TmV3IFR4IGV4ZWN1dGluZw=="
	//}
	//]
	//},
	//"timestamp": 1628574824926
	//},

	//执行操作：
	//  在manifest的abi中的event中有合约的所有方法，可以与eventname对应上，将该event 的parameters的name，type，进行返回。


	//返回结果
	//{
	//	"Vmstate": "HALT",
	//	"contract": "0x618d44dc3af16c6120dbf65402024f40a04f772a",
	//	"eventname": "event",
	//manifest": {{\"name\":\"event\",\"parameters\":[{\"name\":\"obj\",\"type\":\"String\"}]}}
	//"state": {
	//"type": "Array",
	//"value": [
	//{
	//"type": "ByteString",
	//"value": "TmV3IFR4IGV4ZWN1dGluZw=="
	//}
	//]
	//},
	//"timestamp": 1628574824926
	//},

	r1["notification"] = r2
	r1, err = me.Filter(r1, args.Filter)
	if err != nil {
		return nil
	}
	r, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
