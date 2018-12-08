package base

import ()

func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}
