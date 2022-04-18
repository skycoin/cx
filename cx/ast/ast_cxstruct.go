package ast

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

type CXStructIndex int

// CXStruct is used to represent a CX struct.
type CXStruct struct {
	// Metadata
	Index   int
	Name    string         // Name of the struct
	Package CXPackageIndex // The package this struct belongs to

	// Contents
	Fields []CXArgumentIndex // The fields of the struct
}

// ----------------------------------------------------------------
//                             `CXStruct` Getters

// GetField ...
func (strct *CXStruct) GetField(prgrm *CXProgram, name string) (*CXArgument, error) {
	for _, fldIdx := range strct.Fields {
		if prgrm.CXArgs[fldIdx].Name == name {
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
func (strct *CXStruct) AddField(prgrm *CXProgram, fld *CXArgument) *CXStruct {
	for _, fldIdx := range strct.Fields {
		if prgrm.CXArgs[fldIdx].Name == fld.Name {
			fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	numFlds := len(strct.Fields)
	fldIdx := prgrm.AddCXArgInArray(fld)
	strct.Fields = append(strct.Fields, fldIdx)
	if numFlds != 0 {
		// Pre-compiling the offset of the field.
		lastFldIdx := strct.Fields[numFlds-1]
		prgrm.CXArgs[fldIdx].Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	}

	return strct
}

// RemoveField ...
func (strct *CXStruct) RemoveField(prgrm *CXProgram, fldName string) {
	if len(strct.Fields) > 0 {
		lenFlds := len(strct.Fields)
		for i, fldIdx := range strct.Fields {
			if prgrm.CXArgs[fldIdx].Name == fldName {
				if i == lenFlds-1 {
					strct.Fields = strct.Fields[:len(strct.Fields)-1]
				} else {
					strct.Fields = append(strct.Fields[:i], strct.Fields[i+1:]...)
				}
				break
			}
		}
	}
}

func (strct *CXStruct) GetStructSize(prgrm *CXProgram) types.Pointer {
	var structSize types.Pointer
	for _, fldIdx := range strct.Fields {
		fld := prgrm.CXArgs[fldIdx]
		structSize += GetSize(prgrm, &fld)
	}

	return structSize
}

// ---------------- NEW CXStruct def ----------------

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

type NewCXStruct struct {
	StructID     int
	NameStringID int
	Package      CXPackageIndex
	Fields       []CXTypeSignature
}

type CXTypeSignature struct {
	NameStringID int
	Offset       int
	Type         CXTypeSignature_TYPE

	// if type is complex, meta is complex id
	// if type is struct, meta is struct id
	Meta int
}
