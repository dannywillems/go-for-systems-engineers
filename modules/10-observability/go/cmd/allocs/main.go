// Command allocs prints the DETERMINISTIC allocation count of the two string
// builders. testing.AllocsPerRun sets GOMAXPROCS=1, forces a GC, and counts
// mallocs, so the number is machine-independent: it is a fact about the code,
// not a timing measurement. This is the falsifiable claim the README injects.
package main

import (
	"fmt"
	"testing"

	obs "observability"
)

func main() {
	parts := make([]string, 64)
	for i := range parts {
		parts[i] = fmt.Sprintf("chunk-%02d;", i)
	}

	plus := testing.AllocsPerRun(200, func() { _ = obs.ConcatPlus(parts) })
	grow := testing.AllocsPerRun(200, func() { _ = obs.BuilderGrow(parts) })

	fmt.Printf("ConcatPlus  (64 parts): %.0f allocs/op\n", plus)
	fmt.Printf("BuilderGrow (64 parts): %.0f allocs/op\n", grow)
	fmt.Printf("reduction: %.0fx fewer allocations\n", plus/grow)
}
