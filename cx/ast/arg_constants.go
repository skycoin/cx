package ast

import "github.com/skycoin/cx/cx/constants"

//Integer Types, CxArg Constants
var ConstCXArg_I8 = NewCXArgument(constants.TYPE_I8)
var ConstCxArg_I16 = NewCXArgument(constants.TYPE_I16)
var ConstCxArg_I32 = NewCXArgument(constants.TYPE_I32)
var ConstCxArg_I64 = NewCXArgument(constants.TYPE_I64)

//Unsigned Integer Types, CxArg Constants
var ConstCxArg_UI8 = NewCXArgument(constants.TYPE_UI8)
var ConstCxArg_UI16 = NewCXArgument(constants.TYPE_UI16)
var ConstCxArg_UI32 = NewCXArgument(constants.TYPE_UI32)
var ConstCxArg_UI64 = NewCXArgument(constants.TYPE_UI64)

//Floating point, f32, f64, CxArg Constants
var ConstCxArg_F32 = NewCXArgument(constants.TYPE_F32)
var ConstCxArg_F64 = NewCXArgument(constants.TYPE_F64)

// ConstCxArg_BOOL Default bool parameter
var ConstCxArg_BOOL = NewCXArgument(constants.TYPE_BOOL)

// ConstCxArg_STR Default str parameter
var ConstCxArg_STR = NewCXArgument(constants.TYPE_STR)

// ConstCxArg_UND_TYPE Default und parameter
var ConstCxArg_UND_TYPE = NewCXArgument(constants.TYPE_UNDEFINED)

// ConstCxArg_Affordance Default aff parameter
var ConstCxArg_Affordance = NewCXArgument(constants.TYPE_AFF)
