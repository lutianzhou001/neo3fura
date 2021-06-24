package trans

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"neo3fura/var/stderr"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/neophora/neo2go/pkg/core/block"
	"github.com/neophora/neo2go/pkg/core/state"
	"github.com/neophora/neo2go/pkg/core/transaction"
	"github.com/neophora/neo2go/pkg/crypto/keys"
	"github.com/neophora/neo2go/pkg/io"
	"github.com/neophora/neo2go/pkg/util"
)

// T ...
type T struct {
	V interface{}
}

// AddressToHash ...
func (me *T) AddressToHash() error {
	switch address := me.V.(type) {
	case string:
		data := base58.Decode(address)
		if len(data) < 22 {
			return stderr.ErrInvalidArgs
		}
		me.V = data[1:21]
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToHex ...
func (me *T) BytesToHex() error {
	switch bytes := me.V.(type) {
	case []byte:
		me.V = hex.EncodeToString(bytes)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// HexToBytes ...
func (me *T) HexToBytes() error {
	switch enc := me.V.(type) {
	case string:
		var err error
		me.V, err = hex.DecodeString(enc)
		return err
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToHash ...
func (me *T) BytesToHash() error {
	switch bytes := me.V.(type) {
	case []byte:
		l1 := sha256.Sum256(bytes)
		l2 := sha256.Sum256(l1[:])
		me.V = l2[:]
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesReverse ...
func (me *T) BytesReverse() error {
	switch bytes := me.V.(type) {
	case []byte:
		for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
			bytes[i], bytes[j] = bytes[j], bytes[i]
		}
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// StringToLowerCase ...
func (me *T) StringToLowerCase() error {
	switch str := me.V.(type) {
	case string:
		me.V = strings.ToLower(str)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// Remove0xPrefix ...
func (me *T) Remove0xPrefix() error {
	switch str := me.V.(type) {
	case string:
		matches := libTransReg0x.FindStringSubmatch(str)
		if len(matches) != 3 {
			return stderr.ErrInvalidArgs
		}
		me.V = matches[2]
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// HexReverse ...
func (me *T) HexReverse() error {
	if err := me.HexToBytes(); err != nil {
		return err
	}
	if err := me.BytesReverse(); err != nil {
		return err
	}
	if err := me.BytesToHex(); err != nil {
		return err
	}
	return nil
}

// BytesToJSONViaTX ...
func (me *T) BytesToJSONViaTX() error {
	switch bytes := me.V.(type) {
	case []byte:
		var tx transaction.Transaction
		reader := io.NewBinReaderFromBuf(bytes)
		tx.DecodeBinary(reader)
		ret, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToJSONViaBlock ...
func (me *T) BytesToJSONViaBlock() error {
	switch bytes := me.V.(type) {
	case []byte:
		var blk block.Block
		reader := io.NewBinReaderFromBuf(bytes)
		blk.DecodeBinary(reader)
		ret, err := json.Marshal(blk)
		if err != nil {
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToJSONViaHeader ...
func (me *T) BytesToJSONViaHeader() error {
	switch bytes := me.V.(type) {
	case []byte:
		var hd block.Header
		reader := io.NewBinReaderFromBuf(bytes)
		hd.DecodeBinary(reader)
		ret, err := json.Marshal(hd)
		if err != nil {
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToJSONViaAsset ...
func (me *T) BytesToJSONViaAsset() error {
	switch bytes := me.V.(type) {
	case []byte:
		if len(bytes) < 55 {
			return stderr.ErrInvalidArgs
		}
		bytes = append(bytes[1:len(bytes)-46-8], bytes[len(bytes)-46:]...)
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
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToJSONViaAccount ...
func (me *T) BytesToJSONViaAccount() error {
	switch bytes := me.V.(type) {
	case []byte:
		obj := make(map[string]interface{})
		var version byte
		var sh util.Uint160
		var frozen bool
		var votes []*keys.PublicKey
		var balances []map[string]interface{}

		reader := io.NewBinReaderFromBuf(bytes)
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
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

// BytesToJSONViaContract ...
func (me *T) BytesToJSONViaContract() error {
	switch bytes := me.V.(type) {
	case []byte:
		if len(bytes) < 1 {
			return stderr.ErrInvalidArgs
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
			return err
		}
		me.V = json.RawMessage(ret)
		return nil
	default:
		return stderr.ErrInvalidArgs
	}
}

var libTransReg0x = regexp.MustCompile(`^(0x)?([0-9a-f]+)$`)
