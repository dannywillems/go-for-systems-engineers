package observability

import (
	"fmt"
	"testing"
)

func testParts() []string {
	parts := make([]string, 64)
	for i := range parts {
		parts[i] = fmt.Sprintf("chunk-%02d;", i)
	}
	return parts
}

// TestBuildEquivalent: the two builders must produce byte-identical output, so
// the allocation comparison is between two implementations of the SAME function.
func TestBuildEquivalent(t *testing.T) {
	parts := testParts()
	if ConcatPlus(parts) != BuilderGrow(parts) {
		t.Fatal("ConcatPlus and BuilderGrow disagree")
	}
}

// TestBuilderAllocatesLess turns the README's claim into a RED/GREEN gate: the
// pre-sized Builder must allocate strictly fewer times than +=. AllocsPerRun is
// deterministic, so this is not flaky.
func TestBuilderAllocatesLess(t *testing.T) {
	parts := testParts()
	plus := testing.AllocsPerRun(200, func() { _ = ConcatPlus(parts) })
	grow := testing.AllocsPerRun(200, func() { _ = BuilderGrow(parts) })
	if !(grow < plus) {
		t.Fatalf("expected BuilderGrow (%.0f allocs) < ConcatPlus (%.0f allocs)", grow, plus)
	}
	if grow != 1 {
		t.Fatalf("expected BuilderGrow to be a single allocation, got %.0f", grow)
	}
}

var benchSink string

func BenchmarkConcatPlus(b *testing.B) {
	parts := testParts()
	b.ReportAllocs()
	for b.Loop() {
		benchSink = ConcatPlus(parts)
	}
}

func BenchmarkBuilderGrow(b *testing.B) {
	parts := testParts()
	b.ReportAllocs()
	for b.Loop() {
		benchSink = BuilderGrow(parts)
	}
}
