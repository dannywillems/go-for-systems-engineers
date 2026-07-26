// Package solutions is the corrigé for Module 13. Run via `make solutions`.
package solutions

import "unsafe"

func ToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

func Uint64(b []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&b[0]))
}
