package uintval

import "fmt"

// T ...
type T uint64

// Valid ...
func (me T) Valid() bool {
	return true
}

// Val ...
func (me T) Val() uint64 {
	return uint64(me)
}

// Hex ...
func (me T) Hex() string {
	return fmt.Sprintf("%016x", me.Val())
}
