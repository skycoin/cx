package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// CXTypeSignature_TYPE enum contains CXTypeSignature types.
type CXTypeSignature_TYPE int

type CXTypeSignatureIndex int

const (
	TYPE_UNUSED CXTypeSignature_TYPE = iota
	TYPE_ATOMIC
	TYPE_POINTER_ATOMIC
	TYPE_ARRAY_ATOMIC
	TYPE_ARRAY_POINTER_ATOMIC
	TYPE_SLICE_ATOMIC
	TYPE_SLICE_POINTER_ATOMIC

	TYPE_STRUCT
	TYPE_POINTER_STRUCT
	TYPE_ARRAY_STRUCT
	TYPE_ARRAY_POINTER_STRUCT
	TYPE_SLICE_STRUCT
	TYPE_SLICE_POINTER_STRUCT

	TYPE_COMPLEX
	TYPE_POINTER_COMPLEX
	TYPE_ARRAY_COMPLEX
	TYPE_ARRAY_POINTER_COMPLEX
	TYPE_SLICE_COMPLEX
	TYPE_SLICE_POINTER_COMPLEX

	// For CXArgument usage
	// To be deprecated
	TYPE_CXARGUMENT_DEPRECATE
)

type CXTypeSignature struct {
	Index CXTypeSignatureIndex
	// NameStringID int
	Name    string // temporary
	Offset  types.Pointer
	Type    CXTypeSignature_TYPE
	Package CXPackageIndex

	// if type is complex, meta is complex id
	// if type is struct, meta is struct id
	// if type is array, meta is CXTypeSignature_Array id
	// if type is atomic, meta is the atomic type
	// types.BOOL
	// types.I8
	// types.I16
	// types.I32
	// types.I64
	// types.UI8
	// types.UI16
	// types.UI32
	// types.UI64
	// types.F32
	// types.F64
	// types.STR
	Meta int
}

type CXTypeSignature_Array struct {
	Type   int
	Length int
}

// ----------------------------------------------------------------
//                             `CXTypeSignature` Handlers

func (typeSignature *CXTypeSignature) GetSize(prgrm *CXProgram) types.Pointer {
	switch typeSignature.Type {
	case TYPE_ATOMIC:
		return types.Code(typeSignature.Meta).Size()
	case TYPE_POINTER_ATOMIC:
		return types.POINTER.Size()
	case TYPE_ARRAY_ATOMIC:
		typeSignatureForArray := prgrm.GetTypeSignatureArrayFromArray(typeSignature.Meta)
		return types.Code(typeSignatureForArray.Type).Size()
	case TYPE_ARRAY_POINTER_ATOMIC:
		return types.POINTER.Size()
	case TYPE_SLICE_ATOMIC:
		return types.POINTER.Size()
	case TYPE_SLICE_POINTER_ATOMIC:
		return types.POINTER.Size()
	case TYPE_STRUCT:
		return prgrm.GetStructFromArray(CXStructIndex(typeSignature.Meta)).GetStructSize(prgrm)
	case TYPE_POINTER_STRUCT:
	case TYPE_ARRAY_STRUCT:
	case TYPE_ARRAY_POINTER_STRUCT:
	case TYPE_SLICE_STRUCT:
	case TYPE_SLICE_POINTER_STRUCT:

	case TYPE_COMPLEX:
	case TYPE_POINTER_COMPLEX:
	case TYPE_ARRAY_COMPLEX:
	case TYPE_ARRAY_POINTER_COMPLEX:
	case TYPE_SLICE_COMPLEX:
	case TYPE_SLICE_POINTER_COMPLEX:

	case TYPE_CXARGUMENT_DEPRECATE:
		argIdx := typeSignature.Meta
		return GetArgSize(prgrm, &prgrm.CXArgs[argIdx])
	}

	return 0
}

func (typeSignature *CXTypeSignature) GetArrayLength(prgrm *CXProgram) types.Pointer {
	switch typeSignature.Type {
	case TYPE_ARRAY_ATOMIC:
		typeSignatureForArray := prgrm.GetTypeSignatureArrayFromArray(typeSignature.Meta)
		return types.Pointer(typeSignatureForArray.Length)
	}

	return 0
}

func (typeSignature *CXTypeSignature) GetCXArgFormat(prgrm *CXProgram) *CXArgument {
	var arg *CXArgument = &CXArgument{}
	if typeSignature.Type == TYPE_ATOMIC {
		arg.Type = types.Code(typeSignature.Meta)
		arg.StructType = nil
		arg.Size = typeSignature.GetSize(prgrm)
		arg.TotalSize = typeSignature.GetSize(prgrm)

		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_BASIC}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_BASIC}
		}

	} else if typeSignature.Type == TYPE_ARRAY_ATOMIC {
		typeSignatureArray := prgrm.GetTypeSignatureArrayFromArray(typeSignature.Meta)
		arg.Type = types.Code(typeSignatureArray.Type)
		arg.StructType = nil
		arg.Size = typeSignature.GetSize(prgrm)
		arg.Lengths = []types.Pointer{typeSignature.GetArrayLength(prgrm)}
		arg.TotalSize = typeSignature.GetSize(prgrm) * arg.Lengths[0]

		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_ARRAY}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_BASIC}
		}

	} else if typeSignature.Type == TYPE_SLICE_ATOMIC {
		arg.Type = types.Code(typeSignature.Meta)
		arg.StructType = nil
		arg.Size = types.Code(typeSignature.Meta).Size()
		arg.Lengths = []types.Pointer{0}
		arg.TotalSize = typeSignature.GetSize(prgrm)
		arg.IsSlice = true

		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_SLICE}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_BASIC}
		}

		arg.DereferenceOperations = append([]int{constants.DEREF_POINTER}, arg.DereferenceOperations...)

	} else if typeSignature.Type == TYPE_STRUCT {
		arg.Type = types.STRUCT
		arg.StructType = prgrm.GetStructFromArray(CXStructIndex(typeSignature.Meta))
		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_STRUCT}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_STRUCT}
		}

	}
	arg.Offset = typeSignature.Offset

	return arg
}

// Added _ForStructs for now so it will be separate from others for smooth transition.
// TODO: Remove _ForStructs
func GetCXTypeSignatureRepresentationOfCXArg_ForStructs(prgrm *CXProgram, cxArgument *CXArgument) *CXTypeSignature {
	newCXTypeSignature := CXTypeSignature{
		Name:    cxArgument.Name,
		Package: cxArgument.Package,
	}

	fieldType := cxArgument.Type
	// If atomic type. i.e. i8, i16, i32, f32, etc.
	if !cxArgument.IsSlice && len(cxArgument.Lengths) == 0 && fieldType.IsPrimitive() {
		newCXTypeSignature.Type = TYPE_ATOMIC
		newCXTypeSignature.Meta = int(fieldType)

		// If pointer atomic, i.e. *i32, *f32, etc.
	} else if fieldType == types.POINTER && cxArgument.PointerTargetType.IsPrimitive() {
		newCXTypeSignature.Type = TYPE_POINTER_ATOMIC
		newCXTypeSignature.Meta = int(cxArgument.PointerTargetType)

		// If simple array atomic type, i.e. [5]i32, [2]f64, etc.
	} else if !cxArgument.IsSlice && len(cxArgument.Lengths) == 1 && len(cxArgument.Indexes) == 0 && fieldType.IsPrimitive() {
		newCXTypeSignature.Type = TYPE_ARRAY_ATOMIC

		typeSignatureForArray := &CXTypeSignature_Array{
			Type:   int(fieldType),
			Length: int(cxArgument.Lengths[0]),
		}
		typeSignatureForArrayIdx := prgrm.AddTypeSignatureArrayInArray(typeSignatureForArray)

		newCXTypeSignature.Meta = typeSignatureForArrayIdx

		// If slice atomic type, i.e. []i32, []f64, etc.
	} else if cxArgument.IsSlice && len(cxArgument.Lengths) == 1 && (fieldType.IsPrimitive() || fieldType == types.STR) {
		newCXTypeSignature.Type = TYPE_SLICE_ATOMIC
		newCXTypeSignature.Meta = int(fieldType)

		// If type is struct
	} else if !cxArgument.IsSlice && len(cxArgument.Lengths) == 0 && fieldType == types.STRUCT {
		newCXTypeSignature.Type = TYPE_STRUCT
		newCXTypeSignature.Meta = cxArgument.StructType.Index
	} else {
		fldIdx := prgrm.AddCXArgInArray(cxArgument)

		// All are TYPE_CXARGUMENT_DEPRECATE for now.
		// FieldIdx or the CXArg ID is in Meta field.
		newCXTypeSignature.Type = TYPE_CXARGUMENT_DEPRECATE
		newCXTypeSignature.Meta = int(fldIdx)
	}

	newCXTypeSignature.Offset = newCXTypeSignature.GetSize(prgrm)

	return &newCXTypeSignature
}

func GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm *CXProgram, cxArgument *CXArgument) *CXTypeSignature {
	newCXTypeSignature := CXTypeSignature{
		Name:    cxArgument.Name,
		Package: cxArgument.Package,
	}

	// fieldType := cxArgument.Type
	// // If atomic type. i.e. i8, i16, i32, f32, etc.
	// if !cxArgument.IsSlice && len(cxArgument.Lengths) == 0 && fieldType.IsPrimitive() {
	// 	newCXTypeSignature.Type = TYPE_ATOMIC
	// 	newCXTypeSignature.Meta = int(fieldType)

	// 	// If pointer atomic, i.e. *i32, *f32, etc.
	// } else if fieldType == types.POINTER && cxArgument.PointerTargetType.IsPrimitive() {
	// 	newCXTypeSignature.Type = TYPE_POINTER_ATOMIC
	// 	newCXTypeSignature.Meta = int(cxArgument.PointerTargetType)

	// 	// If simple array atomic type, i.e. [5]i32, [2]f64, etc.
	// } else if !cxArgument.IsSlice && len(cxArgument.Lengths) == 1 && len(cxArgument.Indexes) == 0 && fieldType.IsPrimitive() {
	// 	newCXTypeSignature.Type = TYPE_ARRAY_ATOMIC

	// 	typeSignatureForArray := &CXTypeSignature_Array{
	// 		Type:   int(fieldType),
	// 		Length: int(cxArgument.Lengths[0]),
	// 	}
	// 	typeSignatureForArrayIdx := prgrm.AddTypeSignatureArrayInArray(typeSignatureForArray)

	// 	newCXTypeSignature.Meta = typeSignatureForArrayIdx

	// 	// If slice atomic type, i.e. []i32, []f64, etc.
	// } else if cxArgument.IsSlice && len(cxArgument.Lengths) == 1 && (fieldType.IsPrimitive() || fieldType == types.STR) {
	// 	newCXTypeSignature.Type = TYPE_SLICE_ATOMIC
	// 	newCXTypeSignature.Meta = int(fieldType)

	// 	// If type is struct
	// } else if !cxArgument.IsSlice && len(cxArgument.Lengths) == 0 && fieldType == types.STRUCT {
	// 	newCXTypeSignature.Type = TYPE_STRUCT
	// 	newCXTypeSignature.Meta = cxArgument.StructType.Index
	// } else {
	fldIdx := cxArgument.Index

	// All are TYPE_CXARGUMENT_DEPRECATE for now.
	// FieldIdx or the CXArg ID is in Meta field.
	newCXTypeSignature.Type = TYPE_CXARGUMENT_DEPRECATE
	newCXTypeSignature.Meta = int(fldIdx)
	// }

	// newCXTypeSignature.Offset = newCXTypeSignature.GetSize(prgrm)

	return &newCXTypeSignature
}
