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

// AI8 Default i8 parameter
var AI8 = ast.NewCXArgument(constants.TYPE_I8)

// AI16 Default i16 parameter
var AI16 = ast.NewCXArgument(constants.TYPE_I16)

// AI32 Default i32 parameter
var AI32 = ast.NewCXArgument(constants.TYPE_I32)

// AI64 Default i64 parameter
var AI64 = ast.NewCXArgument(constants.TYPE_I64)

// AUI8 Default ui8 parameter
var AUI8 = ast.NewCXArgument(constants.TYPE_UI8)

// AUI16 Default ui16 parameter
var AUI16 = ast.NewCXArgument(constants.TYPE_UI16)

// AUI32 Default ui32 parameter
var AUI32 = ast.NewCXArgument(constants.TYPE_UI32)

// AUI64 Default ui64 parameter
var AUI64 = ast.NewCXArgument(constants.TYPE_UI64)

// AF32 Default f32 parameter
var AF32 = ast.NewCXArgument(constants.TYPE_F32)

// AF64 Default f64 parameter
var AF64 = ast.NewCXArgument(constants.TYPE_F64)

// ASTR Default str parameter
var ASTR = ast.NewCXArgument(constants.TYPE_STR)

// ABOOL Default bool parameter
var ABOOL = ast.NewCXArgument(constants.TYPE_BOOL)

// AUND Default und parameter
var AUND = ast.NewCXArgument(constants.TYPE_UNDEFINED)

// AAFF Default aff parameter
var AAFF = ast.NewCXArgument(constants.TYPE_AFF)

// In Returns a slice of arguments from an argument list
func In(params ...*ast.CXArgument) []*ast.CXArgument {
	return params
}

// Out Returns a slice of arguments from an argument list
func Out(params ...*ast.CXArgument) []*ast.CXArgument {
	return params
}

