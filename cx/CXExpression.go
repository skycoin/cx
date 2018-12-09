package base

import (
	"errors"
	"fmt"
)

/* The CXExpression struct contains information about a CX expression.
 */

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

func MakeExpression(op *CXFunction, fileName string, fileLine int) *CXExpression {
	return &CXExpression{
		ElementID: MakeElementID(),
		Operator:  op,
		FileLine:  fileLine,
		FileName:  fileName}
}

// ----------------------------------------------------------------
//                             Getters

func (expr *CXExpression) GetInputs() ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	} else {
		return nil, errors.New("expression has no arguments")
	}
}

// ----------------------------------------------------------------
//                     Member handling

func (expr *CXExpression) AddInput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Inputs = append(expr.Inputs, param)
	if param.Package == nil {
		param.Package = expr.Package
	}
	return expr
}

func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

func (expr *CXExpression) AddOutput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Outputs = append(expr.Outputs, param)
	param.Package = expr.Package
	return expr
}

func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}

func (expr *CXExpression) AddLabel(lbl string) *CXExpression {
	expr.Label = lbl
	return expr
}

// func (expr *CXExpression) AddOutputName(outName string) *CXExpression {
// 	if len(expr.Operator.Outputs) > 0 {
// 		nextOutIdx := len(expr.Outputs)

// 		var typ string
// 		if expr.Operator.Name == ID_FN || expr.Operator.Name == INIT_FN {
// 			var tmp string
// 			encoder.DeserializeRaw(*expr.Inputs[0].Value, &tmp)

// 			if expr.Operator.Name == INIT_FN {
// 				// then tmp is the type (e.g. initDef("i32") to initialize an i32)
// 				typ = tmp
// 			} else {
// 				var err error
// 				// then tmp is an identifier
// 				if typ, err = GetIdentType(tmp, expr.FileLine, expr.FileName, expr.Program); err == nil {
// 				} else {
// 					panic(err)
// 				}
// 			}
// 		} else {
// 			typ = expr.Operator.Outputs[nextOutIdx].Typ
// 		}

// 		outDef := MakeArgument(outName, "", -1).AddValue(MakeDefaultValue(expr.Operator.Outputs[nextOutIdx].Typ)).AddType(typ)

// 		outDef.Package = expr.Package
// 		outDef.Program = expr.Program

// 		expr.Outputs = append(expr.Outputs, outDef)
// 	}

// 	return expr
// }
