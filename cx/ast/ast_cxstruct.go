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

// CXTypeSignature_META enum contains CXTypeSignature metas.
type CXTypeSignature_META int

const (
	META_UNUSED CXTypeSignature_META = iota
	META_ATOMIC_BOOL
	META_ATOMIC_I8
	META_ATOMIC_I16
	META_ATOMIC_I32
	META_ATOMIC_I64
	META_ATOMIC_UI8
	META_ATOMIC_UI16
	META_ATOMIC_UI32
	META_ATOMIC_UI64
	META_ATOMIC_F32
	META_ATOMIC_F64
	META_ATOMIC_STR
)

// type NewCXStruct struct {
// 	StructID     int
// 	NameStringID int
// 	Package      CXPackageIndex
// 	Fields       []CXTypeSignature
// }

type CXTypeSignature struct {
	// NameStringID int
	Name   string // temporary
	Offset types.Pointer
	Type   CXTypeSignature_TYPE

	// if type is complex, meta is complex id
	// if type is struct, meta is struct id
	// if type is array, meta is CXTypeSignature_Array id
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
	// All are TYPE_CXARGUMENT_DEPRECATE for now.
	// FieldIdx or the CXArg ID is in Meta field.
	for _, typeSignature := range strct.Fields {
		if typeSignature.Name == name {
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
func (strct *CXStruct) AddField(prgrm *CXProgram, fieldType CXTypeSignature_TYPE, cxArgument *CXArgument, cxStruct *CXStruct) *CXStruct {
	// All are TYPE_CXARGUMENT_DEPRECATEfor now.
	// FieldIdx or the CXArg ID is in Meta field.
	for _, typeSignature := range strct.Fields {
		fldIdx := typeSignature.Meta
		if prgrm.CXArgs[fldIdx].Name == cxArgument.Name {
			fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	numFlds := len(strct.Fields)
	fldIdx := prgrm.AddCXArgInArray(cxArgument)

	// All are TYPE_CXARGUMENT_DEPRECATE for now.
	// FieldIdx or the CXArg ID is in Meta field.
	newCXTypeSignature := CXTypeSignature{
		Name:   cxArgument.Name,
		Offset: cxArgument.Offset,
		Type:   TYPE_CXARGUMENT_DEPRECATE,
		Meta:   int(fldIdx),
	}

	strct.Fields = append(strct.Fields, newCXTypeSignature)
	if numFlds != 0 {
		// Pre-compiling the offset of the field.
		lastTypeSignature := strct.Fields[numFlds-1]
		lastFldIdx := lastTypeSignature.Meta
		prgrm.CXArgs[fldIdx].Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	}

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
		fldIdx := typeSignature.Meta
		fld := prgrm.CXArgs[fldIdx]
		structSize += GetSize(prgrm, &fld)
	}

	return structSize
}
