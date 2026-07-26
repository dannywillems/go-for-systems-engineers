// Command demo shows reflection at work: encoding/json marshals a Person by
// walking its type at run time, the unexported field is silently dropped, and
// Describe inspects the same value generically. Deterministic output.
package main

import (
	"encoding/json"
	"fmt"

	rg "reflectgen"
)

func main() {
	p := rg.NewPerson("Ada", 36, "top-secret")

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json.Marshal (reflection): %s\n", b)
	fmt.Printf("secret field still set:    %q (dropped from JSON, no error)\n", p.Secret())

	fmt.Println("Describe (reflection walk):")
	for _, line := range rg.Describe(p) {
		fmt.Printf("  %s\n", line)
	}

	// Reflection also round-trips back into a typed value.
	var q rg.Person
	if err := json.Unmarshal(b, &q); err != nil {
		panic(err)
	}
	fmt.Printf("round-trip: %s is %d\n", q.Name, q.Age)
}
