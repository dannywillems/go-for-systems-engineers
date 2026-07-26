// Package exercises: Module 13 reader tasks. RED until you implement the stubs.
package exercises

// TODO(reader): reinterpret b as a string WITHOUT copying, using unsafe.String.
// Return "" for an empty slice. The test checks the value AND that it does not
// allocate.
func ToString(b []byte) string {
	return "not implemented" // replace me
}

// TODO(reader): interpret the first 8 bytes of b as a uint64 in the machine's
// native byte order, via an unsafe pointer cast. Assume len(b) >= 8.
func Uint64(b []byte) uint64 {
	return 0 // replace me
}
