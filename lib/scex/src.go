package scex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"neo3fura/lib/scex/req"
	"neo3fura/lib/scex/resp"
	"net/rpc"
	"strings"
	"sync"
)

// T ...
type T struct {
	dec     *json.Decoder
	enc     *json.Encoder
	c       io.Closer
	req     req.T
	mutex   sync.Mutex
	seq     uint64
	pending map[uint64]*json.RawMessage
}

// Init ...
func (me *T) Init(conn io.ReadWriteCloser) {
	me.dec = json.NewDecoder(conn)
	me.enc = json.NewEncoder(conn)
	me.c = conn
	me.pending = make(map[uint64]*json.RawMessage)
}

// ReadRequestHeader ...
func (me *T) ReadRequestHeader(r *rpc.Request) error {
	me.req.Reset()
	if err := me.dec.Decode(&me.req); err != nil {
		return err
	}
	r.ServiceMethod = fmt.Sprintf("T.%s", strings.Title(me.req.Method))

	me.mutex.Lock()
	me.seq++
	me.pending[me.seq] = me.req.ID
	me.req.ID = nil
	r.Seq = me.seq
	me.mutex.Unlock()

	return nil
}

// ReadRequestBody ...
func (me *T) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if me.req.Params == nil {
		return errMissingParams
	}
	return json.Unmarshal(*me.req.Params, x)
}

// WriteResponse ...
func (me *T) WriteResponse(r *rpc.Response, x interface{}) error {
	me.mutex.Lock()
	b, ok := me.pending[r.Seq]
	if !ok {
		me.mutex.Unlock()
		return errors.New("invalid sequence number in response")
	}
	delete(me.pending, r.Seq)
	me.mutex.Unlock()

	if b == nil {
		// Invalid request so no id. Use JSON null.
		b = &null
	}
	resp := resp.T{ID: b}
	if r.Error == "" {
		resp.Result = x
	} else {
		resp.Error = r.Error
	}
	return me.enc.Encode(resp)
}

// Close ...
func (me *T) Close() error {
	return me.c.Close()
}

var errMissingParams = errors.New("jsonrpc: request body missing params")
var null = json.RawMessage([]byte("null"))
