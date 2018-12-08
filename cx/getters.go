package base

import (
	"errors"
	"fmt"
)

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

func (expr *CXExpression) GetInputs () ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	} else {
		return nil, errors.New("expression has no arguments")
	}
}
