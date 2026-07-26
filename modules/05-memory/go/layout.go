// Package mem covers struct layout/padding, escape analysis, slice aliasing,
// and GC behavior — the places where Go's managed memory model diverges from a
// Rust engineer's stack-first intuition and an OCaml engineer's minor-heap one.
package mem

import "unsafe"

// region:padding:start

// Padded orders fields worst-case: a bool, then an 8-byte int64, then a bool.
// Alignment forces padding, so the struct is 24 bytes even though its data is
// 10. fieldalignment would rewrite it to the Packed order below.
//
//nolint:govet // intentionally field-misaligned to measure padding
type Padded struct {
	A bool
	B int64
	C bool
}

// Packed puts the 8-byte field first, then the two bools pack together: 16
// bytes. Same fields, 33% less memory — multiplied across a large slice.
type Packed struct {
	B int64
	A bool
	C bool
}

// region:padding:end

// Sizes returns the two struct sizes and the bytes saved per element.
func Sizes() (padded, packed, savedPerElem uintptr) {
	p := unsafe.Sizeof(Padded{})
	q := unsafe.Sizeof(Packed{})
	return p, q, p - q
}
