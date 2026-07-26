// Package solutions is the corrigé for Module 08. Run via `make solutions M=08`.
package solutions

import "errors"

// ErrZeroDenominator is returned by NewRatio when den == 0.
var ErrZeroDenominator = errors.New("ratio: zero denominator")

// Ratio is opaque: num and den are unexported, so the reduced-form invariant
// established by NewRatio cannot be broken from another package.
type Ratio struct {
	num, den int
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func NewRatio(num, den int) (Ratio, error) {
	if den == 0 {
		return Ratio{}, ErrZeroDenominator
	}
	if den < 0 { // carry the sign on the numerator
		num, den = -num, -den
	}
	if num == 0 {
		return Ratio{num: 0, den: 1}, nil
	}
	g := gcd(num, den)
	return Ratio{num: num / g, den: den / g}, nil
}

func (r Ratio) Num() int { return r.num }
func (r Ratio) Den() int { return r.den }

func (r Ratio) Mul(o Ratio) Ratio {
	// den values are already > 0, so the product's denominator stays > 0.
	res, _ := NewRatio(r.num*o.num, r.den*o.den)
	return res
}
