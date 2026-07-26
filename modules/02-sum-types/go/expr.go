// Package expr shows Go's lack of coproducts (sum types) and the "sealed
// interface" idiom used to fake one, plus the exhaustiveness hole that idiom
// leaves open.
//
// Type-theoretically: Go has products (structs, A x B) and a bounded-existential
// form of subtyping (interfaces), but no coproduct A + B with a compiler-checked
// eliminator. The nearest encoding is a "sealed" interface: an interface with an
// UNEXPORTED marker method, so only types in THIS package can satisfy it. That
// closes the set of variants, but Go's `switch x.(type)` is not an eliminator
// the compiler checks for totality: a missing case is legal and silently does
// nothing (or hits `default`). The analyzer in ./exhaustive restores the check.
package expr

// Expr is a sealed sum type: the unexported exprNode() marker means no type
// outside this package can add a variant, so the set {Lit, Add, Mul, Neg} is
// closed. See ./exhaustive for the check Go itself does not perform.
//
//sumtype:decl
type Expr interface {
	exprNode()
}

// region:variants:start

type Lit struct{ V int }

type Add struct{ L, R Expr }

type Mul struct{ L, R Expr }

type Neg struct{ X Expr }

func (Lit) exprNode() {}
func (Add) exprNode() {}
func (Mul) exprNode() {}
func (Neg) exprNode() {}

// region:variants:end

// Eval is TOTAL: it handles every variant. The analyzer confirms this; if a new
// variant were added to Expr without a case here, the analyzer would flag it,
// which the Go compiler alone would not.
func Eval(e Expr) int {
	switch e := e.(type) {
	case Lit:
		return e.V
	case Add:
		return Eval(e.L) + Eval(e.R)
	case Mul:
		return Eval(e.L) * Eval(e.R)
	case Neg:
		return -Eval(e.X)
	}
	panic("unreachable: Expr is sealed and the switch is exhaustive")
}
