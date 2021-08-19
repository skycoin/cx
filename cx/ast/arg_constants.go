package ast

import "github.com/skycoin/cx/cx/types"

//Integer Types, CxArg Constants
var ConstCXArg_I8 = Param(types.I8)
var ConstCxArg_I16 = Param(types.I16)
var ConstCxArg_I32 = Param(types.I32)
var ConstCxArg_I64 = Param(types.I64)

//Unsigned Integer Types, CxArg Constants
var ConstCxArg_UI8 = Param(types.UI8)
var ConstCxArg_UI16 = Param(types.UI16)
var ConstCxArg_UI32 = Param(types.UI32)
var ConstCxArg_UI64 = Param(types.UI64)

//Floating point, f32, f64, CxArg Constants
var ConstCxArg_F32 = Param(types.F32)
var ConstCxArg_F64 = Param(types.F64)

// ConstCxArg_BOOL Default bool parameter
var ConstCxArg_BOOL = Param(types.BOOL)

// ConstCxArg_STR Default str parameter
var ConstCxArg_STR = Param(types.STR)

// ConstCxArg_UND_TYPE Default und parameter
var ConstCxArg_UND_TYPE = Param(types.UNDEFINED)

// ConstCxArg_Affordance Default aff parameter
var ConstCxArg_Affordance = Param(types.AFF)
