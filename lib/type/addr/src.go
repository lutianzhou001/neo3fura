package addr

import (
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
)

// T ...
type T string

// Valid ...
func (me T) Valid() bool {
	return len(base58.Decode(me.Val())) >= 22
}

// Val ...
func (me T) Val() string {
	return string(me)
}

// H160 ...
func (me T) H160() string {
	data := base58.Decode(me.Val())
	return hex.EncodeToString(data[1:21])
}
