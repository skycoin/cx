package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

// AddEmptyFunctionToPackage adds an empty function to a package in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the
// 				 function to be added will be located.
// functionName - the name of the function to be added.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//
// We use AddEmptyFunctionToPackage(cxprogram, "main", "TestFunction").
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction()()
//
// Note the new function added to the package.
func AddEmptyFunctionToPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	fn := cxast.MakeFunction(functionName, "", -1)

	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	pkg.AddFunction(cxprogram, fn)

	return nil
}

// RemoveFunctionFromPackage removes a function from a package in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the
// 				  function to be removed is located.
// functionName - the name of the function to be removed.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction()()
//
// We use RemoveFunctionFromPackage(cxprogram, "main", "TestFunction").
// The Result will be:
// 0.- Package: main
//
// Note the TestFunction was removed from the package.
func RemoveFunctionFromPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}
	pkg.RemoveFunction(functionName)

	return nil
}

// AddNativeInputToFunction adds a native input to a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the
// 				 function is located.
// functionName - the name of the function where the
// 				  input is to be added.
// inputName - the name of the input to be added.
// inputType - the type of the input to be added.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction()()
//
// We use AddNativeInputToFunction(cxprogram, "main", "TestFunction", "inputOne", cxtypes.I32).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction(inputOne i32) ()
//
// Note the new inputOne i32 input to the TestFunction.
func AddNativeInputToFunction(cxprogram *cxast.CXProgram, packageName, functionName, inputName string, inputType types.Code) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	arg := cxast.MakeField(inputName, inputType, "", -1).SetType(types.Code(inputType))
	arg.Package = cxast.CXPackageIndex(pkg.Index)
	fn.AddInput(cxprogram, arg)

	return nil
}

// RemoveFunctionInput removes an input from a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the
// 				  input is to be removed.
// inputName - the name of the input to be removed.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction(inputOne i32) ()
//
// We use RemoveFunctionInput(cxprogram, "TestFunction", "inputOne").
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction() ()
//
// Note the inputOne i32 input was removed from the TestFunction.
// func RemoveFunctionInput(cxprogram *cxast.CXProgram, functionName, inputName string) error {
// 	fn, err := FindFunction(cxprogram, functionName)
// 	if err != nil {
// 		return err
// 	}

// 	fn.RemoveInput(cxprogram, inputName)
// 	return nil
// }

// AddNativeOutputToFunction adds a native output to a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// packageName - the name of the package where the
// 				 function is located.
// functionName - the name of the function where the
// 				  output is to be added.
// outputName - the name of the output to be added.
// outputType - the type of the output to be added.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction()()
//
// We use AddNativeOutputToFunction(cxprogram, "main", "TestFunction", "outputOne", cxtypes.I32).
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction() (outputOne i32)
//
// Note the new outputOne i32 output of the TestFunction.
func AddNativeOutputToFunction(cxprogram *cxast.CXProgram, packageName, functionName, outputName string, outputType types.Code) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}
	arg := cxast.MakeField(outputName, outputType, "", -1).SetType(types.Code(outputType))
	arg.Package = cxast.CXPackageIndex(pkg.Index)
	fn.AddOutput(cxprogram, arg)

	return nil
}

// RemoveFunctionOutput removes an output from a function in cx program.
//
// The inputs are
// cxprogram - in the form of cxast.CXProgram.
// functionName - the name of the function where the
// 				  output is to be removed.
// outputName - the name of the output to be removed.
//
// Example:
// We have this CX Program:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction() (outputOne i32)
//
// We use RemoveFunctionOutput(cxprogram, "TestFunction", "outputOne").
// The Result will be:
// 0.- Package: main
//     Functions
//         0.- Function: TestFunction() ()
//
// Note the outputOne i32 output was removed from the TestFunction.
// func RemoveFunctionOutput(cxprogram *cxast.CXProgram, functionName, outputName string) error {
// 	fn, err := FindFunction(cxprogram, functionName)
// 	if err != nil {
// 		return err
// 	}

// 	fn.RemoveOutput(cxprogram, outputName)
// 	return nil
// }
