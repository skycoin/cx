package base

import (
	// "errors"
	"fmt"
	. "github.com/satori/go.uuid"
)

/* The CXStruct struct contains information about a CX struct.
 */

type CXStruct struct {
	Fields    []*CXArgument
	Name      string
	Size      int
	Package   *CXPackage
	ElementID UUID
}

func MakeStruct(name string) *CXStruct {
	return &CXStruct{
		ElementID: MakeElementID(),
		Name:      name,
	}
}

// ----------------------------------------------------------------
//                             Getters

func (strct *CXStruct) GetFields() ([]*CXArgument, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	} else {
		return nil, fmt.Errorf("structure '%s' has no fields", strct.Name)
	}
}

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

func (strct *CXStruct) AddField(fld *CXArgument) *CXStruct {
	found := false
	for _, fl := range strct.Fields {
		if fl.Name == fld.Name {
			found = true
			break
		}
	}
	if !found {
		strct.Fields = append(strct.Fields, fld)
		strct.Size += fld.TotalSize
	}
	return strct
}

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
