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
	TYPE_POINTER_ARRAY_ATOMIC
	TYPE_SLICE_ATOMIC
	TYPE_SLICE_POINTER_ATOMIC

	TYPE_STRUCT
	TYPE_POINTER_STRUCT
	TYPE_ARRAY_STRUCT
	TYPE_POINTER_ARRAY_STRUCT
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

	PassBy  int  // pass by value or reference
	IsDeref bool // true if it is dereferencing a pointer, array or slice

	// if type is complex, meta is complex id
	// if type is struct, meta is struct id
	// if type is array, meta is CXTypeSignature_Array id
	// if type is atomic, meta is the atomic type
	Meta int
}

type CXTypeSignature_Array struct {
	Type int // atomic type

	// length of the array
	// 1-dimensional array [x]
	// 2-dimensional array [x y]
	// 3-dimensional array [x y z]
	Lengths []types.Pointer

	// Index of the element which
	// we want to access.
	Indexes []CXTypeSignatureIndex
}

// ----------------------------------------------------------------
//                             `CXTypeSignature` Handlers

func (typeSignature *CXTypeSignature) GetSize(prgrm *CXProgram, IsForUpdateSymbolsTable bool) types.Pointer {
	switch typeSignature.Type {
	case TYPE_ATOMIC:
		return types.Code(typeSignature.Meta).Size()
	case TYPE_POINTER_ATOMIC:
		if typeSignature.IsDeref {
			return types.Code(typeSignature.Meta).Size()
		}

		return types.POINTER.Size()
	case TYPE_ARRAY_ATOMIC:
		if IsForUpdateSymbolsTable {
			return typeSignature.GetArraySize(prgrm)
		}

		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
		if len(arrDetails.Indexes) > 0 {
			return types.Code(arrDetails.Type).Size()
		}

		return typeSignature.GetArraySize(prgrm)
	case TYPE_POINTER_ARRAY_ATOMIC:
		arrDetails := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
		if typeSignature.IsDeref && len(arrDetails.Indexes) > 0 {
			return types.Code(arrDetails.Type).Size()
		} else if typeSignature.IsDeref {
			return typeSignature.GetArraySize(prgrm)
		}

		return types.POINTER.Size()
	case TYPE_SLICE_ATOMIC:
		return types.POINTER.Size()
	case TYPE_SLICE_POINTER_ATOMIC:
		return types.POINTER.Size()
	case TYPE_STRUCT:
		return prgrm.GetStructFromArray(CXStructIndex(typeSignature.Meta)).GetStructSize(prgrm)
	case TYPE_POINTER_STRUCT:
	case TYPE_ARRAY_STRUCT:
	case TYPE_POINTER_ARRAY_STRUCT:
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
		typeSignatureForArray := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
		return TotalLength(typeSignatureForArray.Lengths)
	}

	return 0
}

func (typeSignature *CXTypeSignature) GetArraySize(prgrm *CXProgram) types.Pointer {
	switch typeSignature.Type {
	case TYPE_ARRAY_ATOMIC, TYPE_POINTER_ARRAY_ATOMIC:
		typeSignatureForArray := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
		arrType := types.Code(typeSignatureForArray.Type)
		return TotalLength(typeSignatureForArray.Lengths) * arrType.Size()
	}

	return 0
}

func (typeSignature *CXTypeSignature) GetCXArgFormat(prgrm *CXProgram) *CXArgument {
	var arg *CXArgument = &CXArgument{}
	if typeSignature.Type == TYPE_ATOMIC {
		arg.Type = types.Code(typeSignature.Meta)
		arg.StructType = nil
		arg.Size = typeSignature.GetSize(prgrm, false)

		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_BASIC}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_BASIC}
		}

	} else if typeSignature.Type == TYPE_POINTER_ATOMIC {
		arg.Type = types.POINTER
		arg.PointerTargetType = types.Code(typeSignature.Meta)
		arg.StructType = nil
		arg.Size = types.POINTER.Size()

		// TODO: this should not be needed.
		if len(arg.DeclarationSpecifiers) > 0 {
			arg.DeclarationSpecifiers = append([]int{constants.DECL_BASIC, constants.DECL_POINTER}, arg.DeclarationSpecifiers[1:]...)
		} else {
			arg.DeclarationSpecifiers = []int{constants.DECL_POINTER}
		}
	} else if typeSignature.Type == TYPE_ARRAY_ATOMIC {
		typeSignatureArray := prgrm.GetCXTypeSignatureArrayFromArray(typeSignature.Meta)
		arg.Type = types.Code(typeSignatureArray.Type)
		arg.StructType = nil
		arg.Size = typeSignature.GetSize(prgrm, false)
		arg.Lengths = typeSignatureArray.Lengths
		arg.Indexes = typeSignatureArray.Indexes

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
		PassBy:  cxArgument.PassBy,
	}

	fieldType := cxArgument.Type
	if IsTypeAtomic(cxArgument) {
		// If atomic type. i.e. i8, i16, i32, f32, etc.

		newCXTypeSignature.Type = TYPE_ATOMIC
		newCXTypeSignature.Meta = int(fieldType)

	} else if IsTypePointerAtomic(cxArgument) {
		// If pointer atomic, i.e. *i32, *f32, etc.

		newCXTypeSignature.Type = TYPE_POINTER_ATOMIC
		newCXTypeSignature.Offset = cxArgument.Offset

		newCXTypeSignature.Meta = int(cxArgument.Type)
		if cxArgument.Type == types.STR || cxArgument.StructType != nil {
			newCXTypeSignature.Meta = int(cxArgument.PointerTargetType)
		}

	} else if IsTypeArrayAtomic(cxArgument) {
		// If simple array atomic type, i.e. [5]i32, [2]f64, etc.

		newCXTypeSignature.Type = TYPE_ARRAY_ATOMIC

		typeSignatureForArray := &CXTypeSignature_Array{
			Type:    int(fieldType),
			Lengths: cxArgument.Lengths,
			Indexes: cxArgument.Indexes,
		}
		typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)

		newCXTypeSignature.Meta = typeSignatureForArrayIdx

	} else if IsTypePointerArrayAtomic(cxArgument) {
		// If pointer array atomic type, i.e. *[5]i32, *[2]f64, etc.

		newCXTypeSignature.Type = TYPE_POINTER_ARRAY_ATOMIC

		typeSignatureForArray := &CXTypeSignature_Array{
			Type:    int(fieldType),
			Lengths: cxArgument.Lengths,
			Indexes: cxArgument.Indexes,
		}
		typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)

		newCXTypeSignature.Meta = typeSignatureForArrayIdx

	} else if IsTypeSliceAtomic(cxArgument) {
		// If slice atomic type, i.e. []i32, []f64, etc.
		newCXTypeSignature.Type = TYPE_SLICE_ATOMIC
		newCXTypeSignature.Meta = int(fieldType)

	} else if IsTypeStruct(cxArgument) {
		// If type is struct

		newCXTypeSignature.Type = TYPE_STRUCT
		newCXTypeSignature.Meta = cxArgument.StructType.Index
	} else {
		fldIdx := prgrm.AddCXArgInArray(cxArgument)

		// All are TYPE_CXARGUMENT_DEPRECATE for now.
		// FieldIdx or the CXArg ID is in Meta field.
		newCXTypeSignature.Type = TYPE_CXARGUMENT_DEPRECATE
		newCXTypeSignature.Meta = int(fldIdx)
	}

	newCXTypeSignature.Offset = newCXTypeSignature.GetSize(prgrm, false)

	return &newCXTypeSignature
}

func GetCXTypeSignatureRepresentationOfCXArg(prgrm *CXProgram, cxArgument *CXArgument) *CXTypeSignature {
	newCXTypeSignature := CXTypeSignature{
		Name:    cxArgument.Name,
		Package: cxArgument.Package,
		PassBy:  cxArgument.PassBy,
	}

	if IsTypeAtomic(cxArgument) {
		// If atomic type. i.e. i8, i16, i32, f32, etc.
		newCXTypeSignature.Type = TYPE_ATOMIC
		newCXTypeSignature.Meta = int(cxArgument.Type)
		newCXTypeSignature.Offset = cxArgument.Offset
	} else if IsTypePointerAtomic(cxArgument) {
		// If pointer atomic, i.e. *i32, *f32, etc.
		newCXTypeSignature.Type = TYPE_POINTER_ATOMIC
		newCXTypeSignature.Offset = cxArgument.Offset

		newCXTypeSignature.Meta = int(cxArgument.Type)
		if cxArgument.Type == types.STR || cxArgument.StructType != nil {
			newCXTypeSignature.Meta = int(cxArgument.PointerTargetType)
		}
	} else if IsTypeArrayAtomic(cxArgument) {
		// If simple array atomic type, i.e. [5]i32, [2]f64, etc.
		newCXTypeSignature.Type = TYPE_ARRAY_ATOMIC

		typeSignatureForArray := &CXTypeSignature_Array{
			Type:    int(cxArgument.Type),
			Lengths: cxArgument.Lengths,
			Indexes: cxArgument.Indexes,
		}
		typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)

		newCXTypeSignature.Meta = typeSignatureForArrayIdx
	} else if IsTypePointerArrayAtomic(cxArgument) {
		// If pointer array atomic type, i.e. *[5]i32, *[2]f64, etc.

		newCXTypeSignature.Type = TYPE_POINTER_ARRAY_ATOMIC

		typeSignatureForArray := &CXTypeSignature_Array{
			Type:    int(cxArgument.Type),
			Lengths: cxArgument.Lengths,
			Indexes: cxArgument.Indexes,
		}
		typeSignatureForArrayIdx := prgrm.AddCXTypeSignatureArrayInArray(typeSignatureForArray)

		newCXTypeSignature.Meta = typeSignatureForArrayIdx
	} else {
		// TYPE_CXARGUMENT_DEPRECATE
		// FieldIdx or the CXArg ID is in Meta field.
		cxArgumentIndex := cxArgument.Index
		if cxArgument.Index == -1 || cxArgument.Index == 0 {
			cxArgumentIndex = int(prgrm.AddCXArgInArray(cxArgument))
		}
		newCXTypeSignature.Type = TYPE_CXARGUMENT_DEPRECATE
		newCXTypeSignature.Meta = cxArgumentIndex
		newCXTypeSignature.Offset = cxArgument.Offset
	}

	return &newCXTypeSignature
}
