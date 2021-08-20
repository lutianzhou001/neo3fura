package stderr

import (
	"neo3fura_http/lib/jsonrpc2"
)

// ErrUnknown ...
var ErrUnknown = jsonrpc2.NewError(-32001, "unknown error")

// ErrInvalidArgs ...
var ErrInvalidArgs = jsonrpc2.NewError(-32001, "invalid args")

// ErrUnsupportedMethod ...
var ErrUnsupportedMethod = jsonrpc2.NewError(-32001, "unsupported method")

// ErrNotFound ...
var ErrNotFound = jsonrpc2.NewError(-32001, "not found")

// ErrZero
var ErrZero = jsonrpc2.NewError(-32001, "txid cannot be zero")

// FindDocumentErr
var ErrFind = jsonrpc2.NewError(-32001, "find document(s) error")

// InsertJobErr
var ErrInsert = jsonrpc2.NewError(-32001, "insert job error")
