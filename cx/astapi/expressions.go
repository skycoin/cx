package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

// AddNativeExpressionToFunction adds a native expression to a function in cx program.
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
