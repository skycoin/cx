package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

// DefineNewScope marks the first and last expressions to define the boundaries of a scope.

/*
./cxgo/actions/scope.go:11:	SCOPE_UNUSED = iota //if this value appears; program should crash
./cxgo/actions/scope.go:12:	SCOPE_NEW           //= iota + 1 // 1
./cxgo/actions/scope.go:13:	SCOPE_REM           // 2
./cxgo/actions/scope.go:19:		exprs[0].ScopeOperation = SCOPE_NEW
./cxgo/actions/scope.go:21:		exprs[len(exprs)-1].ScopeOperation = SCOPE_REM
./cxgo/actions/functions.go:130:		if expr.ScopeOperation == SCOPE_NEW {
./cxgo/actions/functions.go:162:		if expr.ScopeOperation == SCOPE_REM {
*/

// create a new scope or return to the previous scope
const (
	SCOPE_UNUSED = iota //if this value appears; program should crash
	SCOPE_NEW           //= iota + 1 // 1
	SCOPE_REM           // 2
)

func DefineNewScope(exprs []*ast.CXExpression) {
	if len(exprs) > 1 {
		// initialize new scope
		exprs[0].ExpressionType = ast.CXEXPR_SCOPE_NEW
		// remove last scope
		exprs[len(exprs)-1].ExpressionType = ast.CXEXPR_SCOPE_DEL
	}
}
