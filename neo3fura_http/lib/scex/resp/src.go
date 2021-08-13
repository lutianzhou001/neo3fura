package resp

import (
	"encoding/json"
)

// T ...
type T struct {
	ID     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}
