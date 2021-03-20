package actions

import (
	"github.com/skycoin/cx/cx"
)

// used for selection_statement to layout its outputs
type SelectStatement struct {
	Condition []*cxcore.CXExpression
	Then      []*cxcore.CXExpression
	Else      []*cxcore.CXExpression
}

func SelectionStatement(predExprs []*cxcore.CXExpression, thenExprs []*cxcore.CXExpression, elseifExprs []SelectStatement, elseExprs []*cxcore.CXExpression, op int) []*cxcore.CXExpression {
	switch op {
	case SEL_ELSEIFELSE:
		var lastElse []*cxcore.CXExpression = elseExprs

		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(predExprs, thenExprs, lastElse)
	case SEL_ELSEIF:
		var lastElse []*cxcore.CXExpression
		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(predExprs, thenExprs, lastElse)
	}

	panic("")
}
