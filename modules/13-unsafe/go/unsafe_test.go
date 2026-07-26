package unsafeconv

import "testing"

func TestRoundTrip(t *testing.T) {
	b := []byte("the quick brown fox")
	if s := BytesToString(b); s != "the quick brown fox" {
		t.Fatalf("BytesToString = %q", s)
	}
	if got := string(StringToBytes("hello")); got != "hello" {
		t.Fatalf("StringToBytes round-trip = %q", got)
	}
	if BytesToString(nil) != "" || StringToBytes("") != nil {
		t.Fatal("empty cases")
	}
}

func TestZeroCopyAllocatesNothing(t *testing.T) {
	b := []byte("a moderately long string to make a copy visibly cost an allocation")
	// string(b) copies (1 alloc); BytesToString does not (0 allocs).
	if a := testing.AllocsPerRun(100, func() { _ = BytesToString(b) }); a != 0 {
		t.Fatalf("BytesToString allocated %v times, want 0", a)
	}
}

var (
	sinkS string
	sinkB []byte
)

func BenchmarkStdStringCopy(b *testing.B) {
	buf := []byte("a moderately long string to make a copy visibly cost an allocation")
	b.ReportAllocs()
	for b.Loop() {
		sinkS = string(buf) // copies
	}
}

func BenchmarkZeroCopyString(b *testing.B) {
	buf := []byte("a moderately long string to make a copy visibly cost an allocation")
	b.ReportAllocs()
	for b.Loop() {
		sinkS = BytesToString(buf) // no copy
	}
	_ = sinkB
}
