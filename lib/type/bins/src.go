package bins

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"neo3fura/var/stderr"

	"github.com/neophora/neo2go/pkg/core/block"
	"github.com/neophora/neo2go/pkg/core/state"
	"github.com/neophora/neo2go/pkg/core/transaction"
	"github.com/neophora/neo2go/pkg/crypto/keys"
	"github.com/neophora/neo2go/pkg/encoding/address"
	"github.com/neophora/neo2go/pkg/io"
	"github.com/neophora/neo2go/pkg/util"
)

// T ...
type T []byte

// Valid ...
func (me T) Valid() bool {
	return me != nil
}

// Val ...
func (me T) Val() []byte {
	return []byte(me)
}

// Uint64 ...
func (me T) Uint64() uint64 {
	if len(me) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(me)
}

// BigString ...
func (me T) BigString() string {
	return big.NewInt(0).SetBytes(me.Val()).String()
}

// JSONViaBlock ...
func (me T) JSONViaBlock() (json.RawMessage, error) {
	var blk block.Block
	reader := io.NewBinReaderFromBuf(me.Val())
	blk.DecodeBinary(reader)
	ret, err := json.Marshal(blk)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaUTXO ...
func (me T) JSONViaUTXO() (json.RawMessage, error) {
	var op transaction.Output
	reader := io.NewBinReaderFromBuf(me.Val())
	op.DecodeBinary(reader)
	ret, err := json.Marshal(map[string]interface{}{
		"asset":   op.AssetID,
		"value":   op.Amount,
		"address": address.Uint160ToString(op.ScriptHash),
	})
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaContract ...
func (me T) JSONViaContract() (json.RawMessage, error) {
	bytes := me.Val()
	if len(bytes) < 1 {
		return nil, stderr.ErrInvalidArgs
	}
	var cs state.Contract
	obj := make(map[string]interface{})
	reader := io.NewBinReaderFromBuf(bytes[1:])
	cs.DecodeBinary(reader)
	obj["author"] = cs.Author
	obj["properties"] = cs.Properties
	obj["email"] = cs.Email
	obj["parameters"] = cs.ParamList
	obj["hash"] = cs.ScriptHash().StringBE()
	obj["script"] = hex.EncodeToString(cs.Script)
	obj["returntype"] = cs.ReturnType
	obj["name"] = cs.Name
	obj["code_version"] = cs.CodeVersion
	obj["description"] = cs.Description
	ret, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaAccount ...
func (me T) JSONViaAccount() (json.RawMessage, error) {
	obj := make(map[string]interface{})
	var version byte
	var sh util.Uint160
	var frozen bool
	var votes []*keys.PublicKey
	var balances []map[string]interface{}
	reader := io.NewBinReaderFromBuf(me.Val())
	version = reader.ReadB()
	reader.ReadBytes(sh[:])
	frozen = reader.ReadBool()
	reader.ReadArray(&votes)
	n := int(reader.ReadVarUint())
	balances = make([]map[string]interface{}, 0, n)
	for i := 0; i < n; i++ {
		var asset util.Uint256
		var value util.Fixed8
		balance := make(map[string]interface{})
		reader.ReadBytes(asset[:])
		value.DecodeBinary(reader)
		balance["asset"] = asset.StringBE()
		balance["value"] = value.String()
		balances = append(balances, balance)
	}
	obj["version"] = version
	obj["script_hash"] = "0x" + sh.StringLE()
	obj["frozen"] = frozen
	obj["votes"] = votes
	obj["balances"] = balances

	ret, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaAsset ...
func (me T) JSONViaAsset() (json.RawMessage, error) {
	if len(me) < 55 {
		return nil, stderr.ErrInvalidArgs
	}
	bytes := append(me[1:len(me)-46-8], me[len(me)-46:]...)
	var as state.Asset
	obj := make(map[string]interface{})
	reader := io.NewBinReaderFromBuf(bytes)
	as.DecodeBinary(reader)
	obj["id"] = as.ID
	switch as.AssetType {
	case 0x00:
		obj["type"] = "GoverningToken"
	case 0x01:
		obj["type"] = "UtilityToken"
	case 0x08:
		obj["type"] = "Currency"
	case 0x40:
		obj["type"] = "CreditFlag"
	case 0x80:
		obj["type"] = "DutyFlag"
	case 0x80 | 0x10:
		obj["type"] = "Share"
	case 0x80 | 0x18:
		obj["type"] = "Invoice"
	case 0x80 | 0x20:
		obj["type"] = "Token"
	}
	obj["name"] = as.Name
	obj["amount"] = as.Amount
	obj["available"] = as.Available
	obj["precision"] = as.Precision
	obj["owner"] = as.Owner
	obj["admin"] = as.Admin
	obj["issuer"] = as.Issuer
	obj["admin"] = as.Admin
	obj["expiration"] = as.Expiration
	obj["frozen"] = as.IsFrozen
	ret, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaCoinState ...
func (me T) JSONViaCoinState() (json.RawMessage, error) {
	if len(me) == 8 {
		return json.RawMessage(`"Untracted State: ` + hex.EncodeToString(me) + `"`), nil
	}
	return json.RawMessage(`"Invalid State: ` + hex.EncodeToString(me) + `"`), nil
}

// JSONViaSpentNEO ...
func (me T) JSONViaSpentNEO() (json.RawMessage, error) {
	if len(me) != 16 {
		return nil, stderr.ErrInvalidArgs
	}
	mint := binary.BigEndian.Uint64(me[0:8])
	burn := binary.BigEndian.Uint64(me[8:16])
	js, err := json.Marshal(map[string]interface{}{
		"mint": mint,
		"burn": burn,
	})
	if err != nil {
		return nil, stderr.ErrInvalidArgs
	}
	return js, nil
}

// JSONViaTransaction ...
func (me T) JSONViaTransaction() (json.RawMessage, error) {
	var tx transaction.Transaction
	reader := io.NewBinReaderFromBuf(me)
	tx.DecodeBinary(reader)
	ret, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

// JSONViaHeader ...
func (me T) JSONViaHeader() (json.RawMessage, error) {
	var hd block.Header
	reader := io.NewBinReaderFromBuf(me)
	hd.DecodeBinary(reader)
	ret, err := json.Marshal(hd)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(ret), nil
}

func (me T) JSONViaApplicationLog() (json.RawMessage, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(me.Val(), &raw); err != nil {
		return nil, stderr.ErrInvalidArgs
	}
	ret, err := json.Marshal(raw)
	if err != nil {
		return nil, stderr.ErrInvalidArgs
	}
	return json.RawMessage(ret), nil
}
