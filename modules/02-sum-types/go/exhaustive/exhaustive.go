// Package exhaustive is a go/analysis pass that restores the exhaustiveness
// check Go's type switch lacks. For any interface marked with a `//sumtype:decl`
// directive (a "sealed" sum type: its variants are closed by an unexported
// marker method), it reports every `switch x.(type)` on that interface that
// omits a variant and has no `default` clause.
//
// This is the check a Rust or OCaml engineer expects from `match`, moved from
// the compiler into an external analyzer because Go's `switch` is not a total
// eliminator.
package exhaustive

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const marker = "sumtype:decl"

// Analyzer is the exported pass, usable with singlechecker or `go vet -vettool`.
var Analyzer = &analysis.Analyzer{
	Name: "exhaustive",
	Doc: "check that every type switch over a //sumtype:decl sealed interface " +
		"covers all variants or has a default clause",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	sealed := sealedInterfaces(pass)
	if len(sealed) == 0 {
		return nil, nil
	}

	// variants[iface] is the set of concrete types (by TypeName object) that
	// satisfy the sealed interface within this package.
	variants := map[*types.Named]map[*types.TypeName]bool{}
	for _, iface := range sealed {
		variants[iface] = implementers(pass, iface)
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.TypeSwitchStmt)(nil)}, func(n ast.Node) {
		ts := n.(*ast.TypeSwitchStmt)
		subj := typeSwitchSubjectType(pass, ts)
		named := asNamed(subj)
		if named == nil {
			return
		}
		want, ok := variants[named]
		if !ok {
			return
		}

		covered := map[*types.TypeName]bool{}
		hasDefault := false
		for _, cl := range ts.Body.List {
			cc := cl.(*ast.CaseClause)
			if len(cc.List) == 0 {
				hasDefault = true
				continue
			}
			for _, e := range cc.List {
				if tn := caseTypeName(pass, e); tn != nil {
					covered[tn] = true
				}
			}
		}
		if hasDefault {
			return
		}

		var missing []string
		for tn := range want {
			if !covered[tn] {
				missing = append(missing, tn.Name())
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			pass.Reportf(ts.Pos(),
				"non-exhaustive type switch on %s: missing cases %s (add them or a default clause)",
				named.Obj().Name(), strings.Join(missing, ", "))
		}
	})
	return nil, nil
}

// sealedInterfaces returns the named interface types in this package whose
// declaration carries a //sumtype:decl directive.
func sealedInterfaces(pass *analysis.Pass) []*types.Named {
	var out []*types.Named
	for _, f := range pass.Files {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !hasMarker(gd.Doc) && !hasMarker(ts.Doc) {
					continue
				}
				obj := pass.TypesInfo.Defs[ts.Name]
				if obj == nil {
					continue
				}
				named, ok := obj.Type().(*types.Named)
				if !ok {
					continue
				}
				if _, ok := named.Underlying().(*types.Interface); ok {
					out = append(out, named)
				}
			}
		}
	}
	return out
}

func hasMarker(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}
	for _, c := range cg.List {
		if strings.TrimSpace(strings.TrimPrefix(c.Text, "//")) == marker {
			return true
		}
	}
	return false
}

// implementers returns the concrete named types in this package that satisfy
// the sealed interface (via value or pointer receiver).
func implementers(pass *analysis.Pass, iface *types.Named) map[*types.TypeName]bool {
	ifaceT, _ := iface.Underlying().(*types.Interface)
	out := map[*types.TypeName]bool{}
	scope := pass.Pkg.Scope()
	for _, name := range scope.Names() {
		tn, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}
		named, ok := tn.Type().(*types.Named)
		if !ok || named == iface {
			continue
		}
		if _, isIface := named.Underlying().(*types.Interface); isIface {
			continue
		}
		if types.Implements(named, ifaceT) ||
			types.Implements(types.NewPointer(named), ifaceT) {
			out[tn] = true
		}
	}
	return out
}

// typeSwitchSubjectType returns the static type of the value being switched on.
func typeSwitchSubjectType(pass *analysis.Pass, ts *ast.TypeSwitchStmt) types.Type {
	var assert *ast.TypeAssertExpr
	switch s := ts.Assign.(type) {
	case *ast.AssignStmt: // switch v := x.(type)
		if len(s.Rhs) == 1 {
			assert, _ = s.Rhs[0].(*ast.TypeAssertExpr)
		}
	case *ast.ExprStmt: // switch x.(type)
		assert, _ = s.X.(*ast.TypeAssertExpr)
	}
	if assert == nil {
		return nil
	}
	if tv, ok := pass.TypesInfo.Types[assert.X]; ok {
		return tv.Type
	}
	return nil
}

// caseTypeName maps a case-clause type expression to the TypeName of the
// concrete variant it names (dereferencing a pointer), or nil for `nil`/other.
func caseTypeName(pass *analysis.Pass, e ast.Expr) *types.TypeName {
	tv, ok := pass.TypesInfo.Types[e]
	if !ok {
		return nil
	}
	if named := asNamed(tv.Type); named != nil {
		return named.Obj()
	}
	return nil
}

func asNamed(t types.Type) *types.Named {
	switch t := t.(type) {
	case *types.Named:
		return t
	case *types.Pointer:
		return asNamed(t.Elem())
	default:
		return nil
	}
}
