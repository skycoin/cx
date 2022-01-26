package ast

import (
	"errors"
	"fmt"

	"github.com/skycoin/cx/cx/types"
)

type CXFunctionIndex int

// CXFunction is used to represent a CX function.
type CXFunction struct {
	// Metadata
	Index        int
	Name         string         // Name of the function
	Package      CXPackageIndex // The package it's a member of
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
		Index:    -1,
		FileName: fileName,
		FileLine: fileLine,
	}
}

// ----------------------------------------------------------------
//                             `CXFunction` Getters

// GetExpressions is not used
func (fn *CXFunction) GetExpressions() ([]*CXExpression, error) {
	if fn.Expressions == nil {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}

	return fn.Expressions, nil

}

// GetExpressionByLabel
func (fn *CXFunction) GetExpressionByLabel(prgrm *CXProgram, lbl string) (*CXExpression, error) {
	if fn.Expressions == nil {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}

	for _, expr := range fn.Expressions {
		if expr.GetLabel(prgrm) == lbl {
			return expr, nil
		}
	}
	return nil, fmt.Errorf("expression '%s' not found in function '%s'", lbl, fn.Name)
}

// GetExpressionByLine ...
func (fn *CXFunction) GetExpressionByLine(line int) (*CXExpression, error) {
	if fn.Expressions == nil {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}

	if line > len(fn.Expressions) {
		return nil, fmt.Errorf("expression line number '%d' exceeds number of expressions in function '%s'", line, fn.Name)
	}

	return fn.Expressions[line], nil
}

// GetCurrentExpression ...
func (fn *CXFunction) GetCurrentExpression() (*CXExpression, error) {
	if fn.CurrentExpression == nil {
		return nil, errors.New("current expression is nil")
	}

	return fn.CurrentExpression, nil
}

// ----------------------------------------------------------------
//                     `CXFunction` Member handling

// AddInput ...
func (fn *CXFunction) AddInput(param *CXArgument) *CXFunction {
	found := false
	for _, inp := range fn.Inputs {
		if inp.Name == param.Name {
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
			if inp.Name == inpName {
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
		if out.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Outputs = append(fn.Outputs, param)
	}

	param.Package = fn.Package

	return fn
}

// RemoveOutput ...
func (fn *CXFunction) RemoveOutput(outName string) {
	if len(fn.Outputs) > 0 {
		lenOuts := len(fn.Outputs)
		for i, out := range fn.Outputs {
			if out.Name == outName {
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
func (fn *CXFunction) AddExpression(prgrm *CXProgram, expr *CXExpression) *CXFunction {
	if expr.Type == CX_ATOMIC_OPERATOR {
		cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Package = fn.Package
		cxAtomicOp.Function = fn
	}

	fn.Expressions = append(fn.Expressions, expr)

	fn.CurrentExpression = expr
	fn.LineCount++
	return fn
}

func (fn *CXFunction) AddExpressionByLineNumber(prgrm *CXProgram, expr *CXExpression, line int) *CXFunction {
	cxAtomicOp, _, _, err := prgrm.GetOperation(expr)
	if err != nil {
		panic(err)
	}
	cxAtomicOp.Package = fn.Package
	cxAtomicOp.Function = fn

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
	}
}

// ----------------------------------------------------------------
//                             `CXFunction` Selectors

// MakeAtomicOperatorExpression ...
func MakeAtomicOperatorExpression(prgrm *CXProgram, op *CXFunction) *CXExpression {
	index := prgrm.AddCXAtomicOp(&CXAtomicOperator{
		Operator: op,
	})

	return &CXExpression{
		Index: index,
		Type:  CX_ATOMIC_OPERATOR,
	}
}

// MakeCXLineExpression ...
func MakeCXLineExpression(prgrm *CXProgram, fileName string, lineNo int, lineStr string) *CXExpression {
	index := prgrm.AddCXLine(&CXLine{
		FileName:   fileName,
		LineNumber: lineNo,
		LineStr:    lineStr,
	})

	return &CXExpression{
		Index: index,
		Type:  CX_LINE,
	}
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
