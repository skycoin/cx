package ast

import (
	"errors"
	"fmt"

	"github.com/skycoin/cx/cx/types"
)

// CXFunction is used to represent a CX function.
type CXFunction struct {
	// Metadata
	Name         string     // Name of the function
	Package      *CXPackage // The package it's a member of
	AtomicOPCode int

	// Contents
	Inputs      []*CXArgument   // Input parameters to the function
	Outputs     []*CXArgument   // Output parameters from the function
	Expressions []*CXExpression // Expressions, including control flow statements, in the function

	//TODO: Better Comment for this
	LineCount int // number of expressions, pre-computed for performance

	//TODO: Better Comment for this
	Size types.Pointer // automatic memory size

	// Debugging
	FileName string
	FileLine int

	// Used by the GC
	ListOfPointers []*CXArgument // Root pointers for the GC algorithm

	// Used by the REPL and parser
	CurrentExpression *CXExpression
}

// IsBuiltIn determines if opcode is not 0
// True if the function is native to CX, e.g. int32.add()
func (cxf CXFunction) IsBuiltIn() bool {
	return cxf.AtomicOPCode != 0
}

func (cxf CXFunction) IsAtomic() bool {
	return cxf.AtomicOPCode != 0
}

// MakeFunction creates an empty function.
// Later, parameters and contents can be added.
//
func MakeFunction(name string, fileName string, fileLine int) *CXFunction {
	return &CXFunction{
		Name:     name,
		FileName: fileName,
		FileLine: fileLine,
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
	fn.LineCount++
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
	fn.LineCount++
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

// SelectExpression ...
func (fn *CXFunction) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			if fn, err := mod.GetCurrentFunction(); err == nil {
	// 				fn.SelectExpression(line)
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)
	if len(fn.Expressions) == 0 {
		return nil, errors.New("There are no expressions in this function")
	}

	if line >= len(fn.Expressions) {
		line = len(fn.Expressions) - 1
	}

	if line < 0 {
		line = 0
	}

	expr := fn.Expressions[line]
	fn.CurrentExpression = expr

	return expr, nil
}
