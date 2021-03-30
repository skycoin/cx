package cxcore

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

//TODO: FIX

//TODO: Deprecate, ParamData is only use by http package
type ParamData struct {
	TypCode   int               // The type code of the parameter.
	ParamType int               // Type of the parameter (struct, slice, etc.).
	strctName string            // Name of the struct in case we're handling a struct instance.
	Pkg       *ast.CXPackage    // To what package does this param belongs to.
	inputs    []*ast.CXArgument // Input parameters to a TYPE_FUNC parameter.
	outputs   []*ast.CXArgument // Output parameters to a TYPE_FUNC parameter.
}

// ParamEx Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func ParamEx(paramData ParamData) *ast.CXArgument {
	var arg *ast.CXArgument
	switch paramData.ParamType {
	case constants.PARAM_DEFAULT:
		arg = ast.NewCXArgument(paramData.TypCode)
	case constants.PARAM_SLICE:
		arg = ast.Slice(paramData.TypCode)
	case constants.PARAM_STRUCT:
		arg = Struct(paramData.Pkg.Name, paramData.strctName, "")
	}
	arg.Inputs = paramData.inputs
	arg.Outputs = paramData.outputs
	arg.Package = paramData.Pkg
	return arg
}

// In Returns a slice of arguments from an argument list
func In(params ...*ast.CXArgument) []*ast.CXArgument {
	return params
}

// Out Returns a slice of arguments from an argument list
func Out(params ...*ast.CXArgument) []*ast.CXArgument {
	return params
}

