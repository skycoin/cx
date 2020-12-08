package actions

import (
	. "github.com/skycoin/cx/cx"
)

// used for selection_statement to layout its outputs
type SelectStatement struct {
	Condition []*CXExpression
	Then      []*CXExpression
	Else      []*CXExpression
}

func SelectionStatement(predExprs []*CXExpression, thenExprs []*CXExpression, elseifExprs []SelectStatement, elseExprs []*CXExpression, op int) []*CXExpression {
	switch op {
	case SEL_ELSEIFELSE:
		var lastElse []*CXExpression = elseExprs

		for c := len(elseifExprs) - 1; c >= 0; c-- {
			if lastElse != nil {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
			} else {
				lastElse = SelectionExpressions(elseifExprs[c].Condition, elseifExprs[c].Then, nil)
			}
		}

		return SelectionExpressions(predExprs, thenExprs, lastElse)
	case SEL_ELSEIF:
		var lastElse []*CXExpression
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
