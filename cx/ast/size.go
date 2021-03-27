package ast

import "github.com/skycoin/cx/cx/constants"

// GetSize ...
func GetSize(arg *CXArgument) int {
	if len(arg.Fields) > 0 {
		return GetSize(arg.Fields[len(arg.Fields)-1])
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		if deref == constants.DEREF_SLICE || deref == constants.DEREF_ARRAY {
			return arg.Size
		}
	}

	for decl := range arg.DeclarationSpecifiers {
		if decl == constants.DECL_POINTER {
			return arg.TotalSize
		}
	}

	if arg.CustomType != nil {
		return arg.CustomType.Size
	}

	return arg.TotalSize
}
