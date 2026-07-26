// Command demo prints deterministic results; the GCShape mechanism and the
// benchmark are shown separately in the README.
package main

import (
	"fmt"

	gen "github.com/dannywillems/go-for-systems-engineers/modules/03-generics/go"
)

func main() {
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.5, 2.5, 3.0}
	adders := []gen.Adder{gen.Int(1), gen.Int(2), gen.Int(3), gen.Int(4), gen.Int(5)}

	fmt.Printf("Sum[int](%v)       = %d\n", ints, gen.Sum(ints))
	fmt.Printf("Sum[float64](%v) = %g\n", floats, gen.Sum(floats))
	fmt.Printf("SumIntConcrete(%v) = %d\n", ints, gen.SumIntConcrete(ints))
	fmt.Printf("SumInterface(1..5)    = %d\n", gen.SumInterface(adders))
}
