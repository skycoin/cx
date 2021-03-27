package cxcore

import (
	//"fmt"
	//"math"

	//"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/cx/cx/ast"
)

// GetDerefSize ...
func GetDerefSize(arg *ast.CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}
