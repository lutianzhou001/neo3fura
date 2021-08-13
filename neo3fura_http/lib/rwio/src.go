package rwio

import "io"

// T ...
type T struct {
	R io.Reader
	W io.Writer
}

// Read ...
func (me *T) Read(p []byte) (int, error) {
	return me.R.Read(p)
}

// Write ...
func (me *T) Write(p []byte) (int, error) {
	return me.W.Write(p)
}

// Close ...
func (me *T) Close() error {
	return nil
}
