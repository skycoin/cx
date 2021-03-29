package ast

import "github.com/skycoin/cx/cx/constants"

// GetAssignmentElement ...
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}

// GetType ...
func GetType(arg *CXArgument) int {
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
	arg.Size = constants.TYPE_POINTER_SIZE
	arg.TotalSize = constants.TYPE_POINTER_SIZE

	return arg
}

// Slice Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func Slice(typCode int) *CXArgument {
	arg := NewCXArgument(typCode)
	arg.IsSlice = true
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_SLICE)
	return arg
}

// NewCXArgument ...
// Was "Param"
func NewCXArgument(typCode int) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(constants.TypeNames[typCode])
	arg.IsLocalDeclaration = true
	return arg
}

