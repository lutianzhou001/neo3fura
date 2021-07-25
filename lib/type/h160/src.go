package h160

import (
	"encoding/hex"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"regexp"
)

// T ...
type T string

// Valid ...
func (me T) Valid() bool {
	if (len(me.Val()) != 42) && (len(me.Val()) != 34) {
		return false
	} else if len(me.Val()) != 34 {
		content := me.Val()[2:len(me.Val())]
		return re.MatchString(content)
	} else {
		return rx.MatchString(me.Val())
	}
}

// Val ...
func (me T) Val() string {
	return string(me)
}

// TransferredVal
func (me T) TransferredVal() string {
	if len(me.Val()) == 42 {
		return me.Val()
	} else {
		transferredVal, _ := me.AddressToScriptHash()
		return transferredVal
	}
}

func (me T) ToByte() []byte {
	return []byte(me.Val()[2:len(me.Val())])
}

// ScriptHashToAddress ...
func (me T) ScriptHashToAddress() (string, error) {
	// be
	u, err := helper.UInt160FromString(me.Val())
	if err != nil {
		return "", err
	}
	return crypto.ScriptHashToAddress(u, 0x35), nil
}

// AddressToScriptHash
func (me T) AddressToScriptHash() (string, error) {
	// be
	u, err := crypto.AddressToScriptHash(me.Val(), 0x35)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// RevVal ...
func (me T) RevVal() string {
	bytes, _ := hex.DecodeString(me.Val())
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return hex.EncodeToString(bytes)
}

var re = regexp.MustCompile(`^[0-9a-f]{40}$`)
var rx = regexp.MustCompile(`^[0-9A-Za-z]{34}$`)