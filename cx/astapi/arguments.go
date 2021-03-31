package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

func AddNativeInputToExpression(cxprogram *cxast.CXProgram, packageName, functionName, inputName string, inputType, lineNumber int) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	arg := cxast.MakeField(inputName, inputType, "", -1).AddType(cxconstants.TypeNames[inputType])
	arg.Package = pkg
	expr.AddInput(arg)

	return nil
}

func RemoveInputFromExpression(cxprogram *cxast.CXProgram, functionName string, lineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	expr.RemoveInput()
	return nil
}

func AddNativeOutputToExpression(cxprogram *cxast.CXProgram, packageName, functionName, outputName string, outputType, lineNumber int) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	arg := cxast.MakeField(outputName, outputType, "", -1).AddType(cxconstants.TypeNames[outputType])
	arg.Package = pkg
	expr.AddOutput(arg)

	return nil
}

func RemoveOutputFromExpression(cxprogram *cxast.CXProgram, functionName string, lineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	expr.RemoveOutput()
	return nil
}

func MakeInputExpressionAPointer(cxprogram *cxast.CXProgram, functionName string, lineNumber, expressionNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	cxast.Pointer(expr.Inputs[expressionNumber])
	return nil
}

func MakeOutputExpressionAPointer(cxprogram *cxast.CXProgram, functionName string, lineNumber, expressionNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	cxast.Pointer(expr.Outputs[expressionNumber])
	return nil
}

func GetAccessibleArgsForFunctionByType(cxprogram *cxast.CXProgram, packageLocationName, functionName string, argType int) ([]*cxast.CXArgument, error) {
	var argsList []*cxast.CXArgument

	// Get all globals
	pkg, err := FindPackage(cxprogram, packageLocationName)
	if err != nil {
		return nil, err
	}

	for _, global := range pkg.Globals {
		if global.IsStruct {
			for _, field := range global.CustomType.Fields {
				if field.Type == argType {
					argsList = append(argsList, field)
				}
			}
		} else if global.Type == argType {
			argsList = append(argsList, global)
		}
	}

	for _, imp := range pkg.Imports {
		for _, global := range imp.Globals {
			if global.IsStruct {
				for _, field := range global.CustomType.Fields {
					if field.Type == argType {
						argsList = append(argsList, field)
					}
				}
			} else if global.Type == argType {
				argsList = append(argsList, global)
			}
		}
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return nil, err
	}

	// Get all args from expression inputs
	for _, expr := range fn.Expressions {
		for _, arg := range expr.Inputs {
			if arg.IsStruct {
				for _, field := range arg.CustomType.Fields {
					if field.Type == argType {
						argsList = append(argsList, field)
					}
				}
			} else if arg.Type == argType {
				argsList = append(argsList, arg)
			}
		}
	}

	return argsList, nil
}
