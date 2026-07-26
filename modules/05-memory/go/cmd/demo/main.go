// Command demo prints the deterministic layout facts.
package main

import (
	"fmt"

	mem "github.com/dannywillems/go-for-systems-engineers/modules/05-memory/go"
)

func main() {
	padded, packed, saved := mem.Sizes()
	fmt.Printf("sizeof(Padded) = %d bytes\n", padded)
	fmt.Printf("sizeof(Packed) = %d bytes\n", packed)
	fmt.Printf("saved per element = %d bytes (%d MiB per 1M elements)\n",
		saved, saved*1_000_000>>20)

	before, after := mem.AliasParent()
	fmt.Printf("slice aliasing: orig was %v, after append(orig[:2], 99) it is %v\n",
		before, after)
}
