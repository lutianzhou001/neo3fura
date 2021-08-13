package req

import (
	"encoding/json"
	"errors"
)

// T ...
type T struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	ID     *json.RawMessage `json:"id"`
}

// Reset ...
func (me *T) Reset() {
	me.Method = ""
	me.Params = nil
	me.ID = nil
}

var errMissingParams = errors.New("jsonrpc: request body missing params")
var null = json.RawMessage([]byte("null"))

type serverResponse struct {
	ID     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}
