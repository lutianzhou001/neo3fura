package bq

import (
	"neo3fura/var/stderr"
	"sync"
)

// T ...
type T struct {
	txs  [][]byte
	lock sync.Mutex
}

// Push ...
func (me *T) Push(tx []byte) error {
	me.lock.Lock()
	defer me.lock.Unlock()
	if len(tx) > 0x10000 {
		return stderr.ErrInvalidArgs
	}
	if len(me.txs) >= 0x10000 {
		return stderr.ErrUnknown
	}

	me.txs = append(me.txs, tx)
	return nil
}

// Pop ...
func (me *T) Pop() []byte {
	me.lock.Lock()
	defer me.lock.Unlock()
	if len(me.txs) == 0 {
		return nil
	}

	ret := me.txs[0]
	me.txs = me.txs[1:]
	return ret
}
