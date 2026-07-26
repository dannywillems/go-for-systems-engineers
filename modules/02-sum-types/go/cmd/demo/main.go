// Command demo shows a total evaluator and the silent-wrong-answer bug a
// non-exhaustive type switch produces. stdout is injected into the README.
package main

import (
	"fmt"

	expr "github.com/dannywillems/go-for-systems-engineers/modules/02-sum-types/go"
	"github.com/dannywillems/go-for-systems-engineers/modules/02-sum-types/go/examples/incomplete"
)

func main() {
	// (2 * 3) + (-4) = 2
	e := expr.Add{
		L: expr.Mul{L: expr.Lit{V: 2}, R: expr.Lit{V: 3}},
		R: expr.Neg{X: expr.Lit{V: 4}},
	}
	fmt.Printf("Eval((2*3) + -4) = %d\n", expr.Eval(e))

	// The exhaustiveness hole: incomplete.Name omits the Blue case and still
	// compiles. At runtime it silently returns the wrong answer for Blue, which
	// is worse than a panic because nothing signals the mistake.
	fmt.Printf("incomplete.Name(Red)  = %q\n", incomplete.Name(incomplete.Red{}))
	fmt.Printf("incomplete.Name(Blue) = %q  (silently wrong: Blue unhandled)\n",
		incomplete.Name(incomplete.Blue{}))
}
