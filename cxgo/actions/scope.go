package actions

import cxcore "github.com/skycoin/cx/cx"

// DefineNewScope marks the first and last expressions to define the boundaries of a scope.

//func DefineNewScope(exprs []*cxcore.CXExpression) {}

// create a new scope or return to the previous scope
const (
	SCOPE_UNUSED = iota //if this value appears; program should crash
	SCOPE_NEW           //= iota + 1 // 1
	SCOPE_REM           // 2
)

func DefineNewScope(exprs []*cxcore.CXExpression) {
	if len(exprs) > 1 {
		// initialize new scope
		exprs[0].ScopeOperation = SCOPE_NEW
		// remove last scope
		exprs[len(exprs)-1].ScopeOperation = SCOPE_REM
	}
}


