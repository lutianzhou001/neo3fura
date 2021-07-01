package bins

import (
	"encoding/binary"
	"math/big"
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
