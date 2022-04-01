package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

// create a new scope or return to the previous scope
const (
	SCOPE_UNUSED = iota //if this value appears; program should crash
	SCOPE_NEW
	SCOPE_REM
)

// DefineNewScope marks the first and last expressions to define the boundaries of a scope.
func DefineNewScope(exprs []ast.CXExpression) {
	if len(exprs) > 2 {
		for i := 0; i < len(exprs); i++ {
			if exprs[i].Type != ast.CX_LINE {
				// initialize new scope
				exprs[i].ExpressionType = ast.CXEXPR_SCOPE_NEW
				break
			}
		}

		// remove last scope
		exprs[len(exprs)-1].ExpressionType = ast.CXEXPR_SCOPE_DEL
	}
}
