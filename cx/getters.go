package base

import (
	"errors"
	"fmt"
)

func (pkg *CXPackage) GetImport(impName string) (*CXPackage, error) {
	for _, imp := range pkg.Imports {
		if imp.Name == impName {
			return imp, nil
		}
	}
	return nil, fmt.Errorf("package '%s' not imported", impName)
}

func (mod *CXPackage) GetCurrentStruct() (*CXStruct, error) {
	if mod.CurrentStruct != nil {
		return mod.CurrentStruct, nil
	} else {
		return nil, errors.New("current struct is nil")
	}

}

func (mod *CXPackage) GetCurrentFunction() (*CXFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	} else {
		return nil, errors.New("current function is nil")
	}
}

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

func (mod *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		return mod.Functions, nil
	} else {
		return nil, fmt.Errorf("package '%s' has no functions", mod.Name)
	}
}

func (pkg *CXPackage) GetStruct(strctName string) (*CXStruct, error) {
	var foundStrct *CXStruct
	for _, strct := range pkg.Structs {
		if strct.Name == strctName {
			foundStrct = strct
			break
		}
	}

	if foundStrct == nil {
		//looking in imports
		for _, imp := range pkg.Imports {
			for _, strct := range imp.Structs {
				if strct.Name == strctName {
					foundStrct = strct
					break
				}
			}
		}
	}
	
	if foundStrct != nil {
		return foundStrct, nil
	} else {
		return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, pkg.Name)
	}
}

func (pkg *CXPackage) GetGlobal(defName string) (*CXArgument, error) {
	var foundDef *CXArgument
	for _, def := range pkg.Globals {
		if def.Name == defName {
			foundDef = def
			break
		}
	}

	// for _, imp := range pkg.Imports {
	// 	for _, def := range imp.Globals {
	// 		if def.Name == defName {
	// 			foundDef = def
	// 			break
	// 		}
	// 	}
	// }

	if foundDef != nil {
		return foundDef, nil
	} else {
		return nil, fmt.Errorf("global '%s' not found in package '%s'", defName, pkg.Name)
	}
}

func (pkg *CXPackage) GetFunction (fnName string) (*CXFunction, error) {
	var found bool
	for _, fn := range pkg.Functions {
		if fn.Name == fnName {
			return fn, nil
		}
	}

	// now checking in imported packages
	if !found {
		for _, imp := range pkg.Imports {
			for _, fn := range imp.Functions {
				if fn.Name == fnName {
					return fn, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("function '%s' not found in package '%s' or its imports", fnName, pkg.Name)
}

func (pkg *CXPackage) GetMethod (fnName string, receiverType string) (*CXFunction, error) {
	for _, fn := range pkg.Functions {
		if fn.Name == fnName && len(fn.Inputs) > 0 && fn.Inputs[0].CustomType != nil && fn.Inputs[0].CustomType.Name == receiverType {
			return fn, nil
		}
	}
	
	return nil, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
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
