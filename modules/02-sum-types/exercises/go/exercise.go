// Package exercises: Module 02 reader tasks. RED until you implement the stubs.
// The lesson: because Go will not force you to handle every variant, YOU must,
// and a table test is your only safety net. Compare with ../../solutions/go.
package exercises

//sumtype:decl
type Expr interface{ exprNode() }

type Lit struct{ V int }

type Add struct{ L, R Expr }

type Mul struct{ L, R Expr }

func (Lit) exprNode() {}
func (Add) exprNode() {}
func (Mul) exprNode() {}

// TODO(reader): return the height of the expression tree (a Lit is height 1;
// Add/Mul are 1 + max(height L, height R)). Handle EVERY variant; the Go
// compiler will not remind you if you miss one.
func Height(e Expr) int {
	return 0 // replace me
}

// TODO(reader): report whether the literal value v appears anywhere in the tree.
func Contains(e Expr, v int) bool {
	return false // replace me
}
