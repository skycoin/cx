package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// GetAssignmentElement ...
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return GetAssignmentElement(arg.Fields[len(arg.Fields)-1])
	}
	return arg

}

// GetType ...
func GetType(arg *CXArgument) types.Code {
	fieldCount := len(arg.Fields)
	if fieldCount > 0 {
		return GetType(arg.Fields[fieldCount-1])
	}
	return arg.Type
}

// Pointer takes an already defined `CXArgument` and turns it into a pointer.
//Only used once, deprecate
//TODO: only used by HTTP, create a better module system
func Pointer(arg *CXArgument) *CXArgument {
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_POINTER)
	arg.IsPointer = true
	arg.Size = types.POINTER_SIZE
	arg.TotalSize = types.POINTER_SIZE

	return arg
}

// Struct helper for creating a struct parameter. It creates a
// `CXArgument` named `argName`, that represents a structure instane of
// `strctName`, from package `pkgName`.
func Struct(pkgName, strctName, argName string) *CXArgument {
	pkg, err := PROGRAM.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		panic(err)
	}

	arg := MakeArgument(argName, "", -1).AddType(types.STRUCT)
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
	arg.Size = strct.Size
	arg.TotalSize = strct.Size
	arg.StructType = strct

	return arg
}

// Slice Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func Slice(typeCode types.Code) *CXArgument {
	arg := Param(typeCode)
	arg.IsSlice = true
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_SLICE)
	return arg
}

// Func Helper function for creating function parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func Func(pkg *CXPackage, inputs []*CXArgument, outputs []*CXArgument) *CXArgument {
	arg := Param(types.FUNC)
	arg.ArgDetails.Package = pkg
	arg.Inputs = inputs
	arg.Outputs = outputs
	return arg
}

// Param ...
func Param(typeCode types.Code) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(typeCode)
	arg.IsLocalDeclaration = true
	return arg
}
