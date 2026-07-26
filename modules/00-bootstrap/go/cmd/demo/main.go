// Command demo prints the two deterministic facts Module 00 checks across all
// three languages. Its stdout is injected verbatim into the module README by the
// capture tool; the docs-fresh CI gate fails if the committed output drifts.
package main

import (
	"fmt"

	bootstrap "github.com/dannywillems/go-for-systems-engineers/modules/00-bootstrap/go"
)

const n = 1_000_000

func main() {
	fmt.Printf("sum(1..%d) = %d\n", n, bootstrap.Sum(n))
	fmt.Printf("word size (bytes) = %d\n", bootstrap.WordSizeBytes())
}
