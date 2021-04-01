package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

// AddNativeExpressionToFunction adds a native
// expression to a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// expressionType - the type of the expression to add.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//
// We use AddNativeExpressionToFunction(cxprogram, "TestFunction", cxconstants.OP_ADD).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add()
//
// Note the new expression add() added in TestFunction.
func AddNativeExpressionToFunction(cxprogram *cxast.CXProgram, functionName string, expressionType int) error {
	exp := cxast.MakeExpression(cxast.Natives[expressionType], "", -1)
	exp.Operator.Name = cxast.OpNames[expressionType]

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.AddExpression(exp)
	return nil
}

// RemoveExpressionFromFunction removes an expression from function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// expressionLineNumber - the line number of expression to be removed
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add()
//
// We use RemoveExpressionFromFunction(cxprogram, "TestFunction", 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//
// Note the add() expression was removed.
func RemoveExpressionFromFunction(cxprogram *cxast.CXProgram, functionName string, expressionLineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveExpression(expressionLineNumber)
	return nil
}

// AddNativeExpressionToFunctionByLineNumber adds a native expression
// on a specific N line number of a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// expressionType - the type of the expression to add.
// lineNumber - the line number where the expression is to be added.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add()
//             1.- Expression: div()
//
// We use AddNativeExpressionToFunctionByLineNumber(cxprogram, "TestFunction", cxconstants.OP_SUB, 1).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add()
//             1.- Expression: sub()
//             2.- Expression: div()
//
// Note the new expression sub() is added in line number 1.
func AddNativeExpressionToFunctionByLineNumber(cxprogram *cxast.CXProgram, functionName string, expressionType, lineNumber int) error {
	exp := cxast.MakeExpression(cxast.Natives[expressionType], "", -1)
	exp.Operator.Name = cxast.OpNames[expressionType]
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.AddExpressionByLineNumber(exp, lineNumber)
	return nil
}
