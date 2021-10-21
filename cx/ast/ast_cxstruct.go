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
	Name    string        // Name of the struct
	Package *CXPackage    // The package this struct belongs to
	Size    types.Pointer // The size in memory that this struct takes.

	// Contents
	Fields []*CXArgument // The fields of the struct
}

// ----------------------------------------------------------------
//                             `CXStruct` Getters

// GetField ...
func (strct *CXStruct) GetField(name string) (*CXArgument, error) {
	for _, fld := range strct.Fields {
		if fld.ArgDetails.Name == name {
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
func (strct *CXStruct) AddField(fld *CXArgument) *CXStruct {
	for _, fl := range strct.Fields {
		if fl.ArgDetails.Name == fld.ArgDetails.Name {
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
	strct.Size += GetSize(fld)

	return strct
}

// RemoveField ...
func (strct *CXStruct) RemoveField(fldName string) {
	if len(strct.Fields) > 0 {
		lenFlds := len(strct.Fields)
		for i, fld := range strct.Fields {
			if fld.ArgDetails.Name == fldName {
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
