package base

import (
	"errors"
	"fmt"
)

func (fn *CXFunction) GetCurrentExpression() (*CXExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("current expression is nil")
	}
}

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

func (fn *CXFunction) GetExpressions () ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	} else {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}
}

func (fn *CXFunction) GetExpressionByLabel (lbl string) (*CXExpression, error) {
	if fn.Expressions != nil {
		for _, expr := range fn.Expressions {
			if expr.Label == lbl {
				return expr, nil
			}
		}

		return nil, fmt.Errorf("expression '%s' not found in function '%s'", lbl, fn.Name)
	} else {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}
}

func (fn *CXFunction) GetExpressionByLine (line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		} else {
			return nil, fmt.Errorf("expression line number '%d' exceeds number of expressions in function '%s'", line, fn.Name)
		}

	} else {
		return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)
	}
}

func (expr *CXExpression) GetInputs () ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	} else {
		return nil, errors.New("expression has no arguments")
	}
}
