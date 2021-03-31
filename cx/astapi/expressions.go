package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

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

func RemoveExpressionFromFunction(cxprogram *cxast.CXProgram, functionName string, expressionLineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveExpression(expressionLineNumber)
	return nil
}

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
