package h256

import (
	"encoding/hex"
	"regexp"
)

// T ...
type T string

// Valid ...
func (me T) Valid() bool {
	if len(me.Val()) != 66 {
		return false
	} else {
		content := me.Val()[2:len(me.Val())]
		return re.MatchString(content)
	}
}

func (me T) IsZero() bool {
	if me.Val() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return true
	} else {
		return false
	}
}

// Val ...
func (me T) Val() string {
	return string(me)
}

// RevVal ...
func (me T) RevVal() string {
	bytes, _ := hex.DecodeString(me.Val())
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return hex.EncodeToString(bytes)
}

var re = regexp.MustCompile(`^[0-9a-f]{64}$`)
