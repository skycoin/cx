package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

const (
	SEL_RESERVED = iota
	SEL_ELSEIF
	SEL_ELSEIFELSE
)

// used for selection_statement to layout its outputs
type SelectStatement struct {
	Condition []ast.CXExpression
	Then      []ast.CXExpression
	Else      []ast.CXExpression
}

// CreateSelectionStatement creates series of expressions that will create
// an if elseif else condition.
//
// Input arguments description:
// 	prgrm - a CXProgram that contains all the data and arrays of the program.
//  conditionExprs - contains the condition of the first if statement.
// 	thenExprs - contains the statements if condition is true.
// 	elseifExprs - contains the statements for another if condition for an
// 				  elseif statement.
//  elseExprs - contains the statements if condition is false and there are
// 			 	no elseif conditions.
//  statementType - determines if the statement is an if+elseif+else or if+elseif only.
func CreateSelectionStatement(prgrm *ast.CXProgram, conditionExprs []ast.CXExpression, thenExprs []ast.CXExpression, elseifExprs []SelectStatement, elseExprs []ast.CXExpression, statementType int) []ast.CXExpression {
	var lastElse []ast.CXExpression
	switch statementType {
	case SEL_ELSEIFELSE:
		lastElse = elseExprs
	case SEL_ELSEIF:
		lastElse = nil
	}

	for c := len(elseifExprs) - 1; c >= 0; c-- {
		lastElse = SelectionExpressions(prgrm, elseifExprs[c].Condition, elseifExprs[c].Then, lastElse)
	}

	return SelectionExpressions(prgrm, conditionExprs, thenExprs, lastElse)
}
