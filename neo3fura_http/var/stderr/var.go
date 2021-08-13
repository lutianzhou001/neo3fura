package stderr

import "errors"

// ErrUnknown ...
var ErrUnknown = errors.New("unknown error")

// ErrInvalidArgs ...
var ErrInvalidArgs = errors.New("invalid args")

// ErrUnsupportedMethod ...
var ErrUnsupportedMethod = errors.New("unsupported method")

// ErrNotFound ...
var ErrNotFound = errors.New("not found")

// ErrZero
var ErrZero = errors.New("txid cannot be zero")
