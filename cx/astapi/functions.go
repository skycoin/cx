package astapi

import (
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

func AddEmptyFunctionToPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	fn := cxast.MakeFunction(functionName, "", -1)

	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}

	pkg.AddFunction(fn)

	return nil
}

func RemoveFunctionFromPackage(cxprogram *cxast.CXProgram, packageName, functionName string) error {
	pkg, err := FindPackage(cxprogram, packageName)
	if err != nil {
		return err
	}
	pkg.RemoveFunction(functionName)

	return nil
}

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

func RemoveFunctionInput(cxprogram *cxast.CXProgram, functionName, inputName string) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveInput(inputName)
	return nil
}

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

func RemoveFunctionOutput(cxprogram *cxast.CXProgram, functionName, outputName string) error {
	fn, err := FindFunction(cxprogram, functionName)
	if err != nil {
		return err
	}

	fn.RemoveOutput(outputName)
	return nil
}
