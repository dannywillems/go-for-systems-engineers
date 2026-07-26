// Package exercises: Module 04 reader tasks. RED until you implement the stubs.
package exercises

import (
	"errors"
	"io"
)

// ErrDivByZero is the sentinel callers match with errors.Is.
var ErrDivByZero = errors.New("division by zero")

// TODO(reader): return a/b, or (0, ErrDivByZero) when b == 0. Wrap the sentinel
// with %w so a caller can errors.Is it even after adding context.
func Divide(a, b int) (int, error) {
	return 0, nil // replace me
}

// TODO(reader): read r to completion WITHOUT losing the bytes that arrive on the
// same Read as io.EOF (the product trap). Return the full contents. Treat io.EOF
// as success (nil error); return any other error.
func Drain(r io.Reader) ([]byte, error) {
	return nil, nil // replace me
}
