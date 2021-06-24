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
	if len(me.Val()) != 42 {
		return false
	} else {
		content := me.Val()[2:len(me.Val())]
		return re.MatchString(content)
	}
}

// Val ...
func (me T) Val() string {
	return string(me)
}

func (me T) ToByte() []byte {
	return []byte(me.Val()[2:len(me.Val())])
}

// ScriptHashToAddress ...
func (me T) ScriptHashToAddress() string {
	u := helper.UInt160FromBytes(me.ToByte())
	return crypto.ScriptHashToAddress(u, 0x35)
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
