// Command demo prints the deterministic, synchronized results. The racy result
// is nondeterministic (that is the point) and is demonstrated by the -race
// report in the README, not by a captured number.
package main

import (
	"context"
	"fmt"

	conc "github.com/dannywillems/go-for-systems-engineers/modules/06-memory-model/go"
)

func main() {
	fmt.Printf("AtomicInc() = %d (correct)\n", conc.AtomicInc())
	fmt.Printf("MutexInc()  = %d (correct)\n", conc.MutexInc())

	sq, err := conc.ParallelSquares(context.Background(), []int{1, 2, 3, 4, 5}, 2)
	fmt.Printf("ParallelSquares([1..5], limit=2) = %v err=%v\n", sq, err)
}
