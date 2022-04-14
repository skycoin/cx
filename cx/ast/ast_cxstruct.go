package ast

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// CXStruct is used to represent a CX struct.
type CXStruct struct {
	// Metadata
	Name    string         // Name of the struct
	Package CXPackageIndex // The package this struct belongs to
	Size    types.Pointer  // The size in memory that this struct takes.

	// Contents
	Fields []*CXArgument // The fields of the struct
}

// ----------------------------------------------------------------
//                             `CXStruct` Getters

// GetField ...
func (strct *CXStruct) GetField(name string) (*CXArgument, error) {
	for _, fld := range strct.Fields {
		if fld.Name == name {
			return fld, nil
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
	for _, fl := range strct.Fields {
		if fl.Name == fld.Name {
			fmt.Printf("%s : duplicate field", CompilationError(fl.ArgDetails.FileName, fl.ArgDetails.FileLine))
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	numFlds := len(strct.Fields)
	strct.Fields = append(strct.Fields, fld)
	if numFlds != 0 {
		// Pre-compiling the offset of the field.
		lastFld := strct.Fields[numFlds-1]
		fld.Offset = lastFld.Offset + lastFld.TotalSize
	}
	strct.Size += GetSize(prgrm, fld)

	return strct
}

// RemoveField ...
func (strct *CXStruct) RemoveField(fldName string) {
	if len(strct.Fields) > 0 {
		lenFlds := len(strct.Fields)
		for i, fld := range strct.Fields {
			if fld.Name == fldName {
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

// ---------------- NEW CXStruct def ----------------

// CXTypeSignature_TYPE enum contains CXTypeSignature types.
type CXTypeSignature_TYPE int

const (
	TYPE_UNUSED CXTypeSignature_TYPE = iota
	TYPE_ATOMIC_BOOL
	TYPE_ATOMIC_I8
	TYPE_ATOMIC_I8_POINTER
	TYPE_ATOMIC_I8_ARRAY
	TYPE_ATOMIC_I16
	TYPE_ATOMIC_I16_POINTER
	TYPE_ATOMIC_I16_ARRAY
	TYPE_ATOMIC_I32
	TYPE_ATOMIC_I32_POINTER
	TYPE_ATOMIC_I32_ARRAY
	TYPE_ATOMIC_I64
	TYPE_ATOMIC_I64_POINTER
	TYPE_ATOMIC_I64_ARRAY
	TYPE_ATOMIC_UI8
	TYPE_ATOMIC_UI8_POINTER
	TYPE_ATOMIC_UI8_ARRAY
	TYPE_ATOMIC_UI16
	TYPE_ATOMIC_UI16_POINTER
	TYPE_ATOMIC_UI16_ARRAY
	TYPE_ATOMIC_UI32
	TYPE_ATOMIC_UI32_POINTER
	TYPE_ATOMIC_UI32_ARRAY
	TYPE_ATOMIC_UI64
	TYPE_ATOMIC_UI64_POINTER
	TYPE_ATOMIC_UI64_ARRAY
	TYPE_ATOMIC_F32
	TYPE_ATOMIC_F32_POINTER
	TYPE_ATOMIC_F32_ARRAY
	TYPE_ATOMIC_F64
	TYPE_ATOMIC_F64_POINTER
	TYPE_ATOMIC_F64_ARRAY
	TYPE_ATOMIC_STR
	TYPE_ATOMIC_STR_POINTER
	TYPE_ATOMIC_STR_ARRAY
	TYPE_STRUCT
	TYPE_STRUCT_POINTER
	TYPE_STRUCT_ARRAY
	TYPE_COMPLEX
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
	META_STRUCT
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
	Meta         CXTypeSignature_META
}
