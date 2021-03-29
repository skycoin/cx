package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

// RegisterPackage registers a package on the CX standard library. This does not create a `CXPackage` structure,
// it only tells the CX runtime that `pkgName` will exist by the time a CX program is run.
//func RegisterPackage(pkgName string) {
//	constants.CorePackages = append(constants.CorePackages, pkgName)
//}

// RegisterOpCode ...
func RegisterOpCode(code int, name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
	if code >= len(ast.OpcodeHandlers) {
		ast.OpcodeHandlers = append(ast.OpcodeHandlers, make([]ast.OpcodeHandler, code+1)...)
	}
	if ast.OpcodeHandlers[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, ast.OpNames[code]))
	}
	ast.OpcodeHandlers[code] = handler

	ast.OpNames[code] = name
	ast.OpCodes[name] = code
	//ast.OpVersions[code] = 1
	if inputs == nil {
		inputs = []*ast.CXArgument{}
	}
	if outputs == nil {
		outputs = []*ast.CXArgument{}
	}
	ast.Natives[code] = ast.MakeNativeFunction(code, inputs, outputs)
}

/*
// Debug helper function used to find opcodes when they are not registered
func dumpOpCodes(opCode int) {
	var keys []int
	for k := range OpNames {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("%5d : %s\n", k, OpNames[k])
	}

	fmt.Printf("opCode : %d\n", opCode)
}*/

// Struct helper for creating a struct parameter. It creates a
// `CXArgument` named `argName`, that represents a structure instance of
// `strctName`, from package `pkgName`.
func Struct(pkgName, strctName, argName string) *ast.CXArgument {
	pkg, err := ast.PROGRAM.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		panic(err)
	}

	arg := ast.MakeArgument(argName, "", -1).AddType(constants.TypeNames[constants.TYPE_CUSTOM])
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
	arg.Size = strct.Size
	arg.TotalSize = strct.Size
	arg.CustomType = strct

	return arg
}

func opDebug(*ast.CXExpression, int) {
	ast.PROGRAM.PrintStack()
}

