package cxcore

import (
	"fmt"
)

// CXStruct is used to represent a CX struct.
//
type CXStruct struct {
	// Metadata
	Name    string     // Name of the struct
	Package *CXPackage // The package this struct belongs to
	Size    int        // The size in memory that this struct takes.

	// Contents
	Fields []*CXArgument // The fields of the struct
}

// MakeStruct ...
func MakeStruct(name string) *CXStruct {
	return &CXStruct{
		Name: name,
	}
}

// ----------------------------------------------------------------
//                             Getters

// GetFields ...
func (strct *CXStruct) GetFields() ([]*CXArgument, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	}
	return nil, fmt.Errorf("structure '%s' has no fields", strct.Name)

}

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
//                     Member handling

// AddField ...
func (strct *CXStruct) AddField(fld *CXArgument) *CXStruct {
	found := false
	for _, fl := range strct.Fields {
		if fl.Name == fld.Name {
			found = true
			break
		}
	}

	// FIXME: Shouldn't it be a compilation error if we define a new field
	// 	  with the same name as another field?
	if !found {
		numFlds := len(strct.Fields)
		strct.Fields = append(strct.Fields, fld)
		if numFlds != 0 {
			// Pre-compiling the offset of the field.
			// TODO: I think this would cause conflicts with affordances. For instance, affordances could remove a field and we'd need to re-calculate the offsets. We need to check this.
			lastFld := strct.Fields[numFlds-1]
			fld.Offset = lastFld.Offset + lastFld.TotalSize
		}
		strct.Size += GetSize(fld)
	} else {
		panic("duplicate field")
	}
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
