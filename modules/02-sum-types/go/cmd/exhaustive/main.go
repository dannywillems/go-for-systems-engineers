// Command exhaustive runs the sealed-interface exhaustiveness analyzer as a
// standalone checker (also usable via `go vet -vettool=$(which exhaustive)`).
package main

import (
	"github.com/dannywillems/go-for-systems-engineers/modules/02-sum-types/go/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
