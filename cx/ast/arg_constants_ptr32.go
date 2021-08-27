// +build ptr32

package ast

import "github.com/skycoin/cx/cx/types"

// Pointer type, CXArg Constants
var ConstCxArg_PTR = Param(types.I32) // TODO:PTR use right type when we have ptr alias in cx
