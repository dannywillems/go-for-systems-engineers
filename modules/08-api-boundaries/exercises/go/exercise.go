// Package exercises: Module 08 reader task. RED until you implement the stubs.
// Build an OPAQUE Ratio whose invariant (always in lowest terms, denominator
// > 0) cannot be broken from outside: no exported fields, one validating
// constructor, and operations that return reduced results.
package exercises

// Ratio is opaque: num and den are unexported, so the only way to obtain a valid
// Ratio is NewRatio, and the invariant it establishes holds forever.
type Ratio struct {
	num, den int
}

// TODO(reader): return an error if den == 0. Otherwise store the ratio in
// LOWEST TERMS (divide num and den by gcd(|num|,|den|)) with the sign carried on
// the numerator (den must end up > 0). NewRatio(2,4) => 1/2; NewRatio(3,-6) =>
// -1/2.
func NewRatio(num, den int) (Ratio, error) {
	return Ratio{}, nil // replace me
}

// TODO(reader): return the numerator and denominator of the reduced form.
func (r Ratio) Num() int { return 0 }
func (r Ratio) Den() int { return 0 }

// TODO(reader): return the product, also in lowest terms.
func (r Ratio) Mul(o Ratio) Ratio { return Ratio{} }
