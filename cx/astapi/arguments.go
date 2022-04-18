package astapi

import (
	"errors"

	"github.com/skycoin/cx/cx/ast"
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
	cxparseractions "github.com/skycoin/cx/cxparser/actions"
)

// AddNativeInputToExpression adds native input to an expression
// in a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the function is located.
// functionName - the name of the function where the expression is located.
// inputName - the name of input to be added.
// inputType - the type of the input to be added.
// lineNumber - the line number of the expression in the function.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16)
//
// We use AddNativeInputToExpression(cxprogram, "main", "TestFunction", "y", cxtypes.I16, 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// Note the y i16 added as an input to the expression in line 0.
func AddNativeInputToExpression(cxprogram *cxast.CXProgram, packageName, functionName, inputName string, inputType types.Code, lineNumber int) error {
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

	arg := cxast.MakeField(inputName, inputType, "", -1).SetType(types.Code(inputType))
	arg.Package = cxast.CXPackageIndex(pkg.Index)
	argIdx := cxprogram.AddCXArgInArray(arg)
	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.AddInput(cxprogram, argIdx)

	return nil
}

// RemoveInputFromExpression removes an input from an
// expression in a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// lineNumber - the line number of the expression in the function.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// We use RemoveInputFromExpression(cxprogram, "TestFunction", 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16)
//
// Note the y i16 removed from the expression in line 0.
func RemoveInputFromExpression(cxprogram *cxast.CXProgram, functionName string, lineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.RemoveInput()
	return nil
}

// AddNativeOutputToExpression adds native output to
// an expression in a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the function is located.
// functionName - the name of the function where the expression is located.
// outputName - the name of output to be added.
// outputType - the type of the output to be added.
// lineNumber - the line number of the expression in the function.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add(x i16, y i16)
//
// We use AddNativeOutputToExpression(cxprogram, "main", "TestFunction", "z", cxtypes.I16, 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// Note the z i16 added as an output of the expression in line 0.
func AddNativeOutputToExpression(cxprogram *cxast.CXProgram, packageName, functionName, outputName string, outputType types.Code, lineNumber int) error {
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

	arg := cxast.MakeField(outputName, outputType, "", -1).SetType(types.Code(outputType))
	arg.Package = cxast.CXPackageIndex(pkg.Index)
	argIdx := cxprogram.AddCXArgInArray(arg)
	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.AddOutput(cxprogram, argIdx)

	return nil
}

// RemoveOutputFromExpression removes an output
// from an expression in a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// lineNumber - the line number of the expression in the function.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// We use RemoveOutputFromExpression(cxprogram, "TestFunction", 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: add(x i16, y i16)
//
// Note the z i16 removed as an output from the expression in line 0.
func RemoveOutputFromExpression(cxprogram *cxast.CXProgram, functionName string, lineNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxAtomicOp.RemoveOutput()
	return nil
}

// MakeInputExpressionAPointer makes an input of an
// expression a pointer.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// lineNumber - the line number of the expression in the function.
// expressionNumber
// inputNumber - the input number of the expression to be made as a pointer.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// We use MakeInputExpressionAPointer(cxprogram, "TestFunction", 0, 1).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y *i16)
//
// Note that y *i16 became a pointer.
func MakeInputExpressionAPointer(cxprogram *cxast.CXProgram, functionName string, lineNumber, inputNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	if expr.Type == ast.CX_LINE {
		return errors.New("Expression is a CXLine")
	}
	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxast.MakePointer(cxprogram.GetCXArgFromArray(cxAtomicOp.Inputs[inputNumber]))
	return nil
}

// MakeOutputExpressionAPointer makes an output
// of an expression a pointer.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the expression is located.
// lineNumber - the line number of the expression in the function.
// expressionNumber
// outputNumber - the output number of the expression to be made as a pointer.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z i16 = add(x i16, y i16)
//
// We use MakeOutputExpressionAPointer(cxprogram, "TestFunction", 0, 0).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction (inputOne i32) (outputOne i16)
//             0.- Expression: z *i16 = add(x i16, y i16)
//
// Note that z *i16 became a pointer.
func MakeOutputExpressionAPointer(cxprogram *cxast.CXProgram, functionName string, lineNumber, outputNumber int) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	if expr.Type == ast.CX_LINE {
		return errors.New("Expression is a CXLine")
	}
	cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
	if err != nil {
		panic(err)
	}

	cxast.MakePointer(cxprogram.GetCXArgFromArray(cxAtomicOp.Outputs[outputNumber]))
	return nil
}

// GetAccessibleArgsForFunctionByType gets all accessible
// arguments in cx program for a function by specified
// argument type.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageLocationName - the name of the package where the function is located.
// functionName - the name of the function where the argument will be used.
// argType - the type of argument we are looking for.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//         Globals
//                 0.- Global: testGlobal
//         Functions
//                 0.- Function: TestFunction (inputOne i32) (outputOne i16)
//                         0.- Expression: add(x i16)
//
// We use GetAccessibleArgsForFunctionByType(cxprogram, "main", "TestFunction", cxtypes.I16).
// We then get a result of
// 1. testGlobal i16 arg in the form of cxast.CXArgument.
// 2. x i16 arg in th eform of cxast.CXArgument.
func GetAccessibleArgsForFunctionByType(cxprogram *cxast.CXProgram, packageLocationName, functionName string, argType types.Code) ([]*cxast.CXArgument, error) {
	var argsList []*cxast.CXArgument

	// Get all globals
	pkg, err := FindPackage(cxprogram, packageLocationName)
	if err != nil {
		return nil, err
	}

	for _, globalIdx := range pkg.Globals {
		global := cxprogram.GetCXArg(globalIdx)
		if global.IsStruct {
			for _, fieldIdx := range global.StructType.Fields {
				field := cxprogram.CXArgs[fieldIdx]
				if field.Type == argType {
					argsList = append(argsList, &field)
				}
			}
		} else if global.Type == argType {
			argsList = append(argsList, global)
		}
	}

	for _, impIdx := range pkg.Imports {
		imp, err := cxprogram.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		for _, globalIdx := range imp.Globals {
			global := cxprogram.GetCXArg(globalIdx)
			if global.IsStruct {
				for _, fieldIdx := range global.StructType.Fields {
					field := cxprogram.CXArgs[fieldIdx]
					if field.Type == argType {
						argsList = append(argsList, &field)
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
		cxAtomicOp, _, _, err := cxprogram.GetOperation(&expr)
		if err != nil {
			panic(err)
		}
		for _, argIdx := range cxAtomicOp.Inputs {
			arg := cxprogram.GetCXArgFromArray(argIdx)
			if arg.IsStruct {
				for _, fieldIdx := range arg.StructType.Fields {
					field := cxprogram.CXArgs[fieldIdx]
					if field.Type == argType {
						argsList = append(argsList, &field)
					}
				}
			} else if arg.Type == argType {
				argsList = append(argsList, arg)
			}
		}
	}

	return argsList, nil
}

func AddLiteralInputToExpression(cxprogram *cxast.CXProgram, packageName, functionName string, bytes []byte, argType types.Code, lineNumber int) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}
	cxprogram.CurrentPackage = cxprogram.Packages[packageName]
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	expr, err := fn.GetExpressionByLine(lineNumber)
	if err != nil {
		return err
	}

	cxparseractions.AST = cxprogram
	litArg := cxparseractions.WritePrimary(cxprogram, argType, bytes, false)

	cxAtomicOp1, err := cxprogram.GetCXAtomicOpFromExpressions(litArg, 0)
	if err != nil {
		panic(err)
	}

	arg := cxprogram.GetCXArgFromArray(cxAtomicOp1.Outputs[0])
	argIdx := cxAtomicOp1.Outputs[0]
	cxAtomicOp2, _, _, err := cxprogram.GetOperation(expr)
	if err != nil {
		panic(err)
	}

	arg.Package = cxast.CXPackageIndex(pkg.Index)
	cxAtomicOp2.AddInput(cxprogram, argIdx)

	return nil
}
