// Package solutions is the corrigé for Module 02. Run via `make solutions`.
package solutions

//sumtype:decl
type Expr interface{ exprNode() }

type Lit struct{ V int }

type Add struct{ L, R Expr }

type Mul struct{ L, R Expr }

func (Lit) exprNode() {}
func (Add) exprNode() {}
func (Mul) exprNode() {}

func Height(e Expr) int {
	switch e := e.(type) {
	case Lit:
		return 1
	case Add:
		return 1 + max(Height(e.L), Height(e.R))
	case Mul:
		return 1 + max(Height(e.L), Height(e.R))
	}
	return 0
}

func Contains(e Expr, v int) bool {
	switch e := e.(type) {
	case Lit:
		return e.V == v
	case Add:
		return Contains(e.L, v) || Contains(e.R, v)
	case Mul:
		return Contains(e.L, v) || Contains(e.R, v)
	}
	return false
}
