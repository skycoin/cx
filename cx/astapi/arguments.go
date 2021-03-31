package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// AddNativeInputToExpression adds native input to an expression
// in a function in cx program.
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

// RemoveInputFromExpression removes an input from an
// expression in a function in cx program.
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

// AddNativeOutputToExpression adds native output to
// an expression in a function in cx program.
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

// RemoveOutputFromExpression removes an output
// from an expression in a function in cx program.
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

// MakeInputExpressionAPointer makes an input of an
// expression a pointer.
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

// MakeOutputExpressionAPointer makes an output
// of an expression a pointer.
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

// GetAccessibleArgsForFunctionByType gets all accessible
// arguments in cx program for a function by specified
// argument type.
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
