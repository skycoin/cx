package base

import (
	"errors"
	// "fmt"
	. "github.com/satori/go.uuid" //nolint golint
)

/* The CXExpression struct contains information about a CX expression.
 */

// CXExpression ...
type CXExpression struct {
	Inputs   []*CXArgument
	Outputs  []*CXArgument
	Label    string
	FileName string
	Operator *CXFunction
	// debugging
	FileLine int
	// used for jmp statements
	ThenLines       int
	ElseLines       int
	Function        *CXFunction
	Package         *CXPackage
	ElementID       UUID
	IsMethodCall    bool
	IsStructLiteral bool
	IsArrayLiteral  bool
	IsUndType       bool
	IsBreak         bool
	IsContinue      bool
}

// MakeExpression ...
func MakeExpression(op *CXFunction, fileName string, fileLine int) *CXExpression {
	return &CXExpression{
		ElementID: MakeElementID(),
		Operator:  op,
		FileLine:  fileLine,
		FileName:  fileName}
}

// ----------------------------------------------------------------
//                             Getters

// GetInputs ...
func (expr *CXExpression) GetInputs() ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	}
	return nil, errors.New("expression has no arguments")

}

// ----------------------------------------------------------------
//                     Member handling

// AddInput ...
func (expr *CXExpression) AddInput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Inputs = append(expr.Inputs, param)
	if param.Package == nil {
		param.Package = expr.Package
	}
	return expr
}

// RemoveInput ...
func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

// AddOutput ...
func (expr *CXExpression) AddOutput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Outputs = append(expr.Outputs, param)
	param.Package = expr.Package
	return expr
}

// RemoveOutput ...
func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}

// AddLabel ...
func (expr *CXExpression) AddLabel(lbl string) *CXExpression {
	expr.Label = lbl
	return expr
}
