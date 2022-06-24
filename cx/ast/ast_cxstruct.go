package ast

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// CXTypeSignature_TYPE enum contains CXTypeSignature types.
type CXTypeSignature_TYPE int

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
	// NameStringID int
	Name   string // temporary
	Offset types.Pointer
	Type   CXTypeSignature_TYPE

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

type CXStructIndex int

// CXStruct is used to represent a CX struct.
type CXStruct struct {
	// Metadata
	Index   int
	Name    string         // Name of the struct
	Package CXPackageIndex // The package this struct belongs to

	// Contents
	Fields []CXTypeSignature // The fields of the struct
}

// ----------------------------------------------------------------
//                             `CXStruct` Getters

// GetField ...
func (strct *CXStruct) GetField(prgrm *CXProgram, name string) (*CXArgument, error) {
	for _, typeSignature := range strct.Fields {
		if typeSignature.Name == name {
			// If type is struct
			if typeSignature.Type == TYPE_STRUCT {
				return &CXArgument{
					Name:                  typeSignature.Name,
					Type:                  types.STRUCT,
					DeclarationSpecifiers: []int{constants.DECL_STRUCT},
					StructType:            prgrm.GetStructFromArray(CXStructIndex(typeSignature.Meta)),
				}, nil
				// If type is not cxargument deprecate
			} else if typeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
				return &CXArgument{
					Name:                  typeSignature.Name,
					Type:                  types.Code(typeSignature.Meta),
					DeclarationSpecifiers: []int{constants.DECL_BASIC},
					StructType:            nil,
				}, nil
			}

			// If type is cxargument deprecate
			fldIdx := typeSignature.Meta
			return &prgrm.CXArgs[fldIdx], nil
		}
	}
	return nil, fmt.Errorf("field '%s' not found in struct '%s'", name, strct.Name)
}

// ----------------------------------------------------------------
//                     `CXStruct` Member handling

// MakeStruct ...
func MakeStruct(name string) *CXStruct {
	return &CXStruct{
		Name: name,
	}
}

// AddField ...
func (strct *CXStruct) AddField(prgrm *CXProgram, cxArgument *CXArgument) *CXStruct {
	// Check if field already exist
	for _, typeSignature := range strct.Fields {
		if typeSignature.Name == cxArgument.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	newCXTypeSignature := GetCXTypeSignatureRepresentationOfCXArg_ForStructs(prgrm, cxArgument)
	strct.Fields = append(strct.Fields, *newCXTypeSignature)

	// TODO: Found out the effect and completely remove this.
	// Dont remove this yet as removing this gives another error on cxfx
	// even though all cx tests passed already.
	// numFlds := len(strct.Fields)
	// if numFlds > 2 {
	// 	// Pre-compiling the offset of the field.
	// 	lastTypeSignature := strct.Fields[numFlds-2]
	// 	currentTypeSignature := strct.Fields[numFlds-1]
	// 	if lastTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
	// 		lastFldIdx := lastTypeSignature.Meta
	// 		currFldIdx := currentTypeSignature.Meta
	// 		prgrm.CXArgs[currFldIdx].Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	// 	} else if lastTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
	// 		lastFldIdx := lastTypeSignature.Meta
	// 		currentTypeSignature.Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	// 	} else if lastTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
	// 		currFldIdx := currentTypeSignature.Meta
	// 		prgrm.CXArgs[currFldIdx].Offset = lastTypeSignature.Offset + lastTypeSignature.GetSize(prgrm)
	// 	} else if lastTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
	// 		currentTypeSignature.Offset = lastTypeSignature.Offset + lastTypeSignature.GetSize(prgrm)
	// 	}
	// }

	return strct
}

// AddField ...
func (strct *CXStruct) AddField_Function(prgrm *CXProgram, field *CXTypeSignature) *CXStruct {
	// Check if field already exist
	for _, typeSignature := range strct.Fields {
		if typeSignature.Name == field.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	strct.Fields = append(strct.Fields, *field)

	return strct
}

func (strct *CXStruct) AddField_CXAtomicOps(prgrm *CXProgram, field *CXTypeSignature) *CXStruct {
	strct.Fields = append(strct.Fields, *field)

	return strct
}

func (strct *CXStruct) AddField_Globals_CXAtomicOps(prgrm *CXProgram, cxArgIdx CXArgumentIndex) *CXStruct {
	cxArgument := prgrm.GetCXArgFromArray(cxArgIdx)

	// Check if field already exist
	for _, typeSignature := range strct.Fields {
		if typeSignature.Name == cxArgument.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	newCXTypeSignature := GetCXTypeSignatureRepresentationOfCXArg_ForGlobals_CXAtomicOps(prgrm, cxArgument)
	strct.Fields = append(strct.Fields, *newCXTypeSignature)

	return strct
}

// RemoveField ...
// func (strct *CXStruct) RemoveField(prgrm *CXProgram, fldName string) {
// 	if len(strct.Fields) > 0 {
// 		lenFlds := len(strct.Fields)
// 		for i, fldIdx := range strct.Fields {
// 			if prgrm.CXArgs[fldIdx].Name == fldName {
// 				if i == lenFlds-1 {
// 					strct.Fields = strct.Fields[:len(strct.Fields)-1]
// 				} else {
// 					strct.Fields = append(strct.Fields[:i], strct.Fields[i+1:]...)
// 				}
// 				break
// 			}
// 		}
// 	}
// }

func (strct *CXStruct) GetStructSize(prgrm *CXProgram) types.Pointer {
	var structSize types.Pointer
	for _, typeSignature := range strct.Fields {
		structSize += typeSignature.GetSize(prgrm)
	}

	return structSize
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
		Name: cxArgument.Name,
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
		Name: cxArgument.Name,
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

	newCXTypeSignature.Offset = newCXTypeSignature.GetSize(prgrm)

	return &newCXTypeSignature
}
