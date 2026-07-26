package solutions

import (
	"testing"
	"unsafe"
)

func TestToString(t *testing.T) {
	b := []byte("zero copy")
	if ToString(b) != "zero copy" || ToString(nil) != "" {
		t.Fatal("ToString")
	}
	if a := testing.AllocsPerRun(100, func() { _ = ToString(b) }); a != 0 {
		t.Fatalf("allocated %v", a)
	}
}

func TestUint64(t *testing.T) {
	var x uint64 = 0x0102030405060708
	b := (*[8]byte)(unsafe.Pointer(&x))[:]
	if Uint64(b) != x {
		t.Fatal("Uint64")
	}
}
