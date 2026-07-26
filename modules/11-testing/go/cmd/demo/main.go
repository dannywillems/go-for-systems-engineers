// Command demo shows the pure subject under test, deterministically.
package main

import (
	"fmt"

	tk "testkit"
)

func main() {
	for _, s := range []string{"  Hello   World  ", "MiXeD\tCase", "   ", "a\n\nb"} {
		fmt.Printf("Normalize(%q) = %q\n", s, tk.Normalize(s))
	}
}
