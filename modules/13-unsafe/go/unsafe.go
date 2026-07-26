// Package unsafeconv shows Go's `unsafe` escape hatch: zero-copy string<->[]byte
// conversion, struct layout introspection, and the uintptr GC hazard.
//
// `unsafe.Pointer` is Go's typed-hole: the compiler still type-checks the
// program, and the GC still runs, but a `uintptr` is JUST A NUMBER the GC does
// NOT track — so a pointer laundered through uintptr can be freed or moved out
// from under you. The safe conversions use unsafe.String / unsafe.Slice
// (Go 1.20+), which keep the GC informed.
package unsafeconv

import "unsafe"

// region:zerocopy:start

// BytesToString reinterprets a []byte as a string WITHOUT copying: the string
// header points at the slice's backing array. Safe only if b is not mutated
// afterwards (strings are supposed to be immutable). string(b) copies; this
// does not.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// StringToBytes reinterprets a string as a []byte WITHOUT copying. The result
// MUST NOT be written to (it aliases the string's immutable storage).
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// region:zerocopy:end

// Record is a deliberately padded struct so Offsetof/Sizeof show real padding.
//
//nolint:govet // intentionally field-misaligned to demonstrate layout
type Record struct {
	Flag bool
	ID   int64
	Kind uint8
}

// Layout returns (size, offset of ID, alignment) using unsafe intrinsics.
func Layout() (size, idOffset, align uintptr) {
	var r Record
	return unsafe.Sizeof(r), unsafe.Offsetof(r.ID), unsafe.Alignof(r)
}
