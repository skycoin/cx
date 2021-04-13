package ast

import (
	"errors"
	"fmt"
)

// MakeFunction creates an empty function.
// Later, parameters and contents can be added.
//
func MakeFunction(name string, fileName string, fileLine int) *CXFunction {
	return &CXFunction{
		Name:      name,
		FileName:  fileName,
		FileLine:  fileLine,
		IsBuiltin: false,
	}
}

// ----------------------------------------------------------------
//                             `CXFunction` Getters

// GetExpressions is not used
func (fn *CXFunction) GetExpressions() ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	}
	return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)

}

// GetExpressionByLabel
func (fn *CXFunction) GetExpressionByLabel(lbl string) (*CXExpression, error) {
	if fn.Expressions == nil {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}
	for _, expr := range fn.Expressions {
		if expr.Label == lbl {
			return expr, nil
		}
	}
	return nil, fmt.Errorf("expression '%s' not found in function '%s'", lbl, fn.Name)
}

// GetExpressionByLine ...
func (fn *CXFunction) GetExpressionByLine(line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		}
		return nil, fmt.Errorf("expression line number '%d' exceeds number of expressions in function '%s'", line, fn.Name)

	}
	return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)

}

// GetCurrentExpression ...
func (fn *CXFunction) GetCurrentExpression() (*CXExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("current expression is nil")
	}
}

// ----------------------------------------------------------------
//                     `CXFunction` Member handling

// AddInput ...
func (fn *CXFunction) AddInput(param *CXArgument) *CXFunction {
	found := false
	for _, inp := range fn.Inputs {
		if inp.ArgDetails.Name == param.ArgDetails.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Inputs = append(fn.Inputs, param)
	}

	return fn
}

// RemoveInput ...
func (fn *CXFunction) RemoveInput(inpName string) {
	if len(fn.Inputs) > 0 {
		lenInps := len(fn.Inputs)
		for i, inp := range fn.Inputs {
			if inp.ArgDetails.Name == inpName {
				if i == lenInps {
					fn.Inputs = fn.Inputs[:len(fn.Inputs)-1]
				} else {
					fn.Inputs = append(fn.Inputs[:i], fn.Inputs[i+1:]...)
				}
				break
			}
		}
	}
}

// AddOutput ...
func (fn *CXFunction) AddOutput(param *CXArgument) *CXFunction {
	found := false
	for _, out := range fn.Outputs {
		if out.ArgDetails.Name == param.ArgDetails.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Outputs = append(fn.Outputs, param)
	}

	param.ArgDetails.Package = fn.Package

	return fn
}

// RemoveOutput ...
func (fn *CXFunction) RemoveOutput(outName string) {
	if len(fn.Outputs) > 0 {
		lenOuts := len(fn.Outputs)
		for i, out := range fn.Outputs {
			if out.ArgDetails.Name == outName {
				if i == lenOuts {
					fn.Outputs = fn.Outputs[:len(fn.Outputs)-1]
				} else {
					fn.Outputs = append(fn.Outputs[:i], fn.Outputs[i+1:]...)
				}
				break
			}
		}
	}
}

// AddExpression ...
func (fn *CXFunction) AddExpression(expr *CXExpression) *CXFunction {
	// expr.Program = fn.Program
	expr.Package = fn.Package
	expr.Function = fn
	fn.Expressions = append(fn.Expressions, expr)
	fn.CurrentExpression = expr
	fn.Length++
	return fn
}

func (fn *CXFunction) AddExpressionByLineNumber(expr *CXExpression, line int) *CXFunction {
	expr.Package = fn.Package
	expr.Function = fn

	lenExprs := len(fn.Expressions)
	if lenExprs == line {
		fn.Expressions = append(fn.Expressions, expr)
	} else {
		fn.Expressions = append(fn.Expressions[:line+1], fn.Expressions[line:]...)
		fn.Expressions[line] = expr
	}

	fn.CurrentExpression = expr
	fn.Length++
	return fn
}

// RemoveExpression ...
func (fn *CXFunction) RemoveExpression(line int) {
	if len(fn.Expressions) > 0 {
		lenExprs := len(fn.Expressions)
		if line >= lenExprs-1 || line < 0 {
			fn.Expressions = fn.Expressions[:len(fn.Expressions)-1]
		} else {
			fn.Expressions = append(fn.Expressions[:line], fn.Expressions[line+1:]...)
		}
		// for i, expr := range fn.Expressions {
		// 	expr.Index = i
		// }
	}
}

// ----------------------------------------------------------------
//                             `CXFunction` Selectors

// MakeExpression ...
func MakeExpression(op *CXFunction, fileName string, fileLine int) *CXExpression {
	return &CXExpression{
		Operator: op,
		FileLine: fileLine,
		FileName: fileName}
}
