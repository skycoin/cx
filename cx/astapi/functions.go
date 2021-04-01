package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// AddEmptyFunctionToPackage adds an empty function to a package in cx program.
func AddEmptyFunctionToPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	fn := cxast.MakeFunction(functionName, "", -1)

	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	pkg.AddFunction(fn)

	return nil
}

// RemoveFunctionFromPackage removes a function from a package in cx program.
func RemoveFunctionFromPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}
	pkg.RemoveFunction(functionName)

	return nil
}

// AddNativeInputToFunction adds a native input to a function in cx program.
func AddNativeInputToFunction(cxprogram *cxast.CXProgram, packageName, functionName, inputName string, inputType int) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	arg := cxast.MakeField(inputName, inputType, "", -1).AddType(cxconstants.TypeNames[inputType])
	arg.Package = pkg
	fn.AddInput(arg)

	return nil
}

// RemoveFunctionInput removes an input from a function in cx program.
func RemoveFunctionInput(cxprogram *cxast.CXProgram, functionName, inputName string) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveInput(inputName)
	return nil
}

// AddNativeOutputToFunction adds a native output to a function in cx program.
func AddNativeOutputToFunction(cxprogram *cxast.CXProgram, packageName, functionName, outputName string, outputType int) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}
	arg := cxast.MakeField(outputName, outputType, "", -1).AddType(cxconstants.TypeNames[outputType])
	arg.Package = pkg
	fn.AddOutput(arg)

	return nil
}

// RemoveFunctionOutput removes an output from a function in cx program.
func RemoveFunctionOutput(cxprogram *cxast.CXProgram, functionName, outputName string) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveOutput(outputName)
	return nil
}
