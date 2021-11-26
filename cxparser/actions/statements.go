package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

// used for selection_statement to layout its outputs
type SelectStatement struct {
	Condition []*ast.CXExpression
	Then      []*ast.CXExpression
	Else      []*ast.CXExpression
}

func SelectionStatement(prgrm *ast.CXProgram, predExprs []*ast.CXExpression, thenExprs []*ast.CXExpression, elseifExprs []SelectStatement, elseExprs []*ast.CXExpression, op int) []*ast.CXExpression {
	switch op {
	case SEL_ELSEIFELSE:
		var lastElse []*ast.CXExpression = elseExprs

		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(prgrm, elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(prgrm, elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(prgrm, predExprs, thenExprs, lastElse)
	case SEL_ELSEIF:
		var lastElse []*ast.CXExpression
		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(prgrm, elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(prgrm, elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(prgrm, predExprs, thenExprs, lastElse)
	}

	panic("")
}
