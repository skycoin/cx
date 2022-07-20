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
	Inputs      *CXStruct      // Input parameters to the function
	Outputs     *CXStruct      // Output parameters from the function
	Expressions []CXExpression // Expressions, including control flow statements, in the function

	LocalVariables []string // contains the name of its local variables

	//TODO: Better Comment for this
	LineCount int // number of expressions, pre-computed for performance

	//TODO: Better Comment for this
	Size types.Pointer // automatic memory size

	// Debugging
	FileName string
	FileLine int

	// Used by the GC
	ListOfPointers []*CXArgument // Root pointers for the GC algorithm
}

// CXNativeFunction is used to represent built-in operations.
type CXNativeFunction struct {
	// Metadata
	AtomicOPCode int

	// Contents
	Inputs  []*CXArgument // Input parameters to the function
	Outputs []*CXArgument // Output parameters from the function
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
		Inputs: &CXStruct{
			Name: name + "_Input",
		},
		LocalVariables: []string{},
	}
}

// ----------------------------------------------------------------
//                             `CXFunction` Getters

// GetExpressions is not used
func (fn *CXFunction) GetExpressions() ([]CXExpression, error) {
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
			return &expr, nil
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

	return &fn.Expressions[line], nil
}

// ----------------------------------------------------------------
//                     `CXFunction` Member handling

// AddInput ...
func (fn *CXFunction) AddInput(prgrm *CXProgram, param *CXArgument) *CXFunction {
	fnInputs := fn.GetInputs(prgrm)
	for _, inputIdx := range fnInputs {
		input := prgrm.GetCXTypeSignatureFromArray(inputIdx)
		if input.Type == TYPE_CXARGUMENT_DEPRECATE {
			inp := prgrm.GetCXArgFromArray(CXArgumentIndex(input.Meta))
			if inp.Name == param.Name {
				return fn
			}
		}

	}

	param.Package = fn.Package

	var newField *CXTypeSignature
	// if atomic type
	if IsTypeAtomic(param) {
		newField = &CXTypeSignature{
			Name:    param.Name,
			Type:    TYPE_ATOMIC,
			Meta:    int(param.Type),
			Package: param.Package,
		}
	} else {
		paramIdx := prgrm.AddCXArgInArray(param)

		newField = &CXTypeSignature{
			Name:    param.Name,
			Type:    TYPE_CXARGUMENT_DEPRECATE,
			Meta:    int(paramIdx),
			Package: param.Package,
		}
	}

	if fn.Inputs == nil {
		fn.Inputs = &CXStruct{}
	}

	newFieldIdx := prgrm.AddCXTypeSignatureInArray(newField)
	fn.Inputs.AddField_TypeSignature(prgrm, newFieldIdx)

	return fn
}

func (fn *CXFunction) GetInputs(prgrm *CXProgram) []CXTypeSignatureIndex {
	if fn == nil || fn.Inputs == nil {
		return []CXTypeSignatureIndex{}
	}

	return fn.Inputs.Fields
}

// RemoveInput ...
// func (fn *CXFunction) RemoveInput(prgrm *CXProgram, inpName string) {
// 	if len(fn.Inputs) > 0 {
// 		lenInps := len(fn.Inputs)
// 		for i, inpIdx := range fn.Inputs {
// 			inp := prgrm.GetCXArgFromArray(inpIdx)
// 			if inp.Name == inpName {
// 				if i == lenInps {
// 					fn.Inputs = fn.Inputs[:len(fn.Inputs)-1]
// 				} else {
// 					fn.Inputs = append(fn.Inputs[:i], fn.Inputs[i+1:]...)
// 				}
// 				break
// 			}
// 		}
// 	}
// }

// AddOutput ...
func (fn *CXFunction) AddOutput(prgrm *CXProgram, param *CXArgument) *CXFunction {
	fnOutputs := fn.GetOutputs(prgrm)
	for _, outputIdx := range fnOutputs {
		output := prgrm.GetCXTypeSignatureFromArray(outputIdx)
		if output.Name == param.Name {
			return fn
		}
	}

	param.Package = fn.Package

	var newField *CXTypeSignature
	// if atomic type
	if IsTypeAtomic(param) {
		newField = &CXTypeSignature{
			Name:    param.Name,
			Type:    TYPE_ATOMIC,
			Meta:    int(param.Type),
			Package: param.Package,
		}
	} else {
		paramIdx := prgrm.AddCXArgInArray(param)

		newField = &CXTypeSignature{
			Name:    param.Name,
			Type:    TYPE_CXARGUMENT_DEPRECATE,
			Meta:    int(paramIdx),
			Package: param.Package,
		}
	}

	if fn.Outputs == nil {
		fn.Outputs = &CXStruct{}
	}

	newFieldIdx := prgrm.AddCXTypeSignatureInArray(newField)
	fn.Outputs.AddField_TypeSignature(prgrm, newFieldIdx)

	return fn
}

func (fn *CXFunction) GetOutputs(prgrm *CXProgram) []CXTypeSignatureIndex {
	if fn == nil || fn.Outputs == nil {
		return []CXTypeSignatureIndex{}
	}

	return fn.Outputs.Fields
}

// RemoveOutput ...
// func (fn *CXFunction) RemoveOutput(prgrm *CXProgram, outName string) {
// 	if len(fn.Outputs) > 0 {
// 		lenOuts := len(fn.Outputs)
// 		for i, outIdx := range fn.Outputs {
// 			out := prgrm.GetCXArgFromArray(outIdx)
// 			if out.Name == outName {
// 				if i == lenOuts {
// 					fn.Outputs = fn.Outputs[:len(fn.Outputs)-1]
// 				} else {
// 					fn.Outputs = append(fn.Outputs[:i], fn.Outputs[i+1:]...)
// 				}
// 				break
// 			}
// 		}
// 	}
// }

// AddExpression ...
func (fn *CXFunction) AddExpression(prgrm *CXProgram, expr *CXExpression) *CXFunction {
	if expr.Type == CX_ATOMIC_OPERATOR {
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOp.Package = fn.Package
		cxAtomicOp.Function = CXFunctionIndex(fn.Index)
	}

	fn.Expressions = append(fn.Expressions, *expr)
	fn.LineCount++
	return fn
}

func (fn *CXFunction) AddExpressionByLineNumber(prgrm *CXProgram, expr *CXExpression, line int) *CXFunction {
	cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}
	cxAtomicOp.Package = fn.Package
	cxAtomicOp.Function = CXFunctionIndex(fn.Index)

	lenExprs := len(fn.Expressions)
	if lenExprs == line {
		fn.Expressions = append(fn.Expressions, *expr)
	} else {
		fn.Expressions = append(fn.Expressions[:line+1], fn.Expressions[line:]...)
		fn.Expressions[line] = *expr
	}

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

func (fn *CXFunction) AddLocalVariableName(name string) error {
	fn.LocalVariables = append(fn.LocalVariables, name)

	return nil
}

func (fn *CXFunction) IsLocalVariable(name string) bool {
	for _, localVar := range fn.LocalVariables {
		if name == localVar {
			return true
		}
	}

	return false
}

// ----------------------------------------------------------------
//                             `CXFunction` Selectors

// MakeAtomicOperatorExpression ...
func MakeAtomicOperatorExpression(prgrm *CXProgram, op *CXNativeFunction) *CXExpression {
	var opIdx CXFunctionIndex = -1
	if op != nil {
		opIdx = prgrm.AddNativeFunctionInArray(op)
	}

	index := prgrm.AddCXAtomicOp(&CXAtomicOperator{
		Operator: opIdx,
		// Function: -1,
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

	return &expr, nil
}
