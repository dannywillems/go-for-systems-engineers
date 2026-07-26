package exercises

import (
	"testing"
	"unsafe"
)

func TestToString(t *testing.T) {
	b := []byte("zero copy please")
	if s := ToString(b); s != "zero copy please" {
		t.Fatalf("ToString = %q", s)
	}
	if ToString(nil) != "" {
		t.Fatal("empty case")
	}
	if a := testing.AllocsPerRun(100, func() { _ = ToString(b) }); a != 0 {
		t.Fatalf("ToString allocated %v, want 0 (did you copy?)", a)
	}
}

func TestUint64(t *testing.T) {
	var x uint64 = 0x0102030405060708
	b := (*[8]byte)(unsafe.Pointer(&x))[:] // native-endian bytes of x
	if got := Uint64(b); got != x {
		t.Fatalf("Uint64 = %#x, want %#x", got, x)
	}
}
