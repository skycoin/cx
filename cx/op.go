package cxcore

import (
	//"fmt"
	//"math"

	//"github.com/skycoin/skycoin/src/cipher/encoder"
)

// GetSize ...
func GetSize(arg *CXArgument) int {
	if len(arg.Fields) > 0 {
		return GetSize(arg.Fields[len(arg.Fields)-1])
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		if deref == DEREF_SLICE || deref == DEREF_ARRAY {
			return arg.Size
		}
	}

	for decl := range arg.DeclarationSpecifiers {
		if decl == DECL_POINTER {
			return arg.TotalSize
		}
	}

	if arg.CustomType != nil {
		return arg.CustomType.Size
	}

	return arg.TotalSize
}

// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}
