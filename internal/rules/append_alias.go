package rules

import (
	"go/ast"

	"github.com/enetx/glint/checker"
)

func AppendAliasRule(ctx *checker.Context, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}

		lhsIdent, ok := assign.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}

		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		switch fun := call.Fun.(type) {
		case *ast.Ident:
			if fun.Name == "append" && len(call.Args) > 0 {
				arg, ok := call.Args[0].(*ast.Ident)
				if !ok || lhsIdent.Name == arg.Name {
					return true
				}

				if callExpr, ok := call.Args[0].(*ast.CallExpr); ok {
					if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Clone" {
						return true
					}
				}

				ctx.Reportf(assign.Pos(),
					"slice '%s' may be aliased by append, clone it first using slices.Clone(%s)",
					arg.Name, arg.Name)
			}

		case *ast.SelectorExpr:
			if fun.Sel.Name == "Append" || fun.Sel.Name == "Push" {
				recv := fun.X

				if recvCall, ok := recv.(*ast.CallExpr); ok {
					if sel, ok := recvCall.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Clone" {
						return true
					}
				}

				if recvIdent, ok := recv.(*ast.Ident); ok && lhsIdent.Name != recvIdent.Name {
					ctx.Reportf(assign.Pos(),
						"slice '%s' may be aliased by %s, clone it first using %s.Clone()",
						recvIdent.Name, fun.Sel.Name, recvIdent.Name)
				}
			}
		}

		return true
	})
}
