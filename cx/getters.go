package base

import (
	"errors"
	"fmt"
	"strings"
)

func (prgrm *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	if prgrm.CurrentPackage != nil {
		return prgrm.CurrentPackage, nil
	} else {
		return nil, errors.New("current package is nil")
	}
}

func (pkg *CXPackage) GetImport(impName string) (*CXPackage, error) {
	for _, imp := range pkg.Imports {
		if imp.Name == impName {
			return imp, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("package '%s' not imported", impName))
}

func (prgrm *CXProgram) GetCurrentStruct() (*CXStruct, error) {
	if prgrm.CurrentPackage != nil {
		if prgrm.CurrentPackage.CurrentStruct != nil {
			return prgrm.CurrentPackage.CurrentStruct, nil
		} else {
			return nil, errors.New("current struct is nil")
		}
	} else {
		return nil, errors.New("current package is nil")
	}
}

func (mod *CXPackage) GetCurrentStruct() (*CXStruct, error) {
	if mod.CurrentStruct != nil {
		return mod.CurrentStruct, nil
	} else {
		return nil, errors.New("current struct is nil")
	}

}

func (prgrm *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	if prgrm.CurrentPackage != nil {
		 if prgrm.CurrentPackage.CurrentFunction != nil {
			 return prgrm.CurrentPackage.CurrentFunction, nil
		 } else {
			 return nil, errors.New("current function is nil")
		 }
	} else {
		return nil, errors.New("current package is nil")
	}
}

func (mod *CXPackage) GetCurrentFunction() (*CXFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	} else {
		return nil, errors.New("current function is nil")
	}
}

func (prgrm *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	if prgrm.CurrentPackage != nil &&
		prgrm.CurrentPackage.CurrentFunction != nil &&
		prgrm.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return prgrm.CurrentPackage.CurrentFunction.CurrentExpression, nil
	} else {
		return nil, errors.New("current package, function or expression is nil")
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

// func (prgrm *CXProgram) GetCurrentDefinitions () ([]*CXArgument, error) {
// 	mod, err := prgrm.GetCurrentPackage()

// 	if err == nil {
// 		return mod.GetCurrentDefinitions()
// 	} else {
// 		return nil, err
// 	}
// }

// func (mod *CXPackage) GetCurrentDefinitions () ([]*CXArgument, error) {
// 	return mod.GetDefinitions()
// }

// func (mod *CXPackage) GetDefinitions () ([]*CXArgument, error) {
// 	if mod.Globals != nil {
// 		return mod.Globals, nil
// 	} else {
// 		return nil, errors.New("definitions array is nil")
// 	}
// }

func (prgrm *CXProgram) GetGlobal(name string) (*CXArgument, error) {
	if mod, err := prgrm.GetCurrentPackage(); err == nil {
		var found *CXArgument
		for _, def := range mod.Globals {
			if def.Name == name {
				found = def
				break
			}
		}

		for _, imp := range mod.Imports {
			for _, def := range imp.Globals {
				if def.Name == name {
					found = def
					break
				}
			}
		}

		if found == nil {
			return nil, errors.New(fmt.Sprintf("global '%s' not found", name))
		} else {
			return found, nil
		}
	} else {
		return nil, err
	}
}

func (strct *CXStruct) GetFields() ([]*CXArgument, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	} else {
		return nil, errors.New(fmt.Sprintf("structure '%s' has no fields", strct.Name))
	}
}

func (strct *CXStruct) GetField(name string) (*CXArgument, error) {
	for _, fld := range strct.Fields {
		if fld.Name == name {
			return fld, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("field '%s' not found in struct '%s'", name, strct.Name))
}

func (mod *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		return mod.Functions, nil
	} else {
		return nil, errors.New(fmt.Sprintf("package '%s' has no functions", mod.Name))
	}
}

func (prgrm *CXProgram) GetPackage(modName string) (*CXPackage, error) {
	if prgrm.Packages != nil {
		var found *CXPackage
		for _, mod := range prgrm.Packages {
			if modName == mod.Name {
				found = mod
				break
			}
		}
		if found != nil {
			return found, nil
		} else {
			return nil, errors.New(fmt.Sprintf("package '%s' not found", modName))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("package '%s' not found", modName))
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
		return nil, errors.New(fmt.Sprintf("struct '%s' not found in package '%s'", strctName, pkg.Name))
	}
}

func (prgrm *CXProgram) GetStruct(strctName string, modName string) (*CXStruct, error) {
	var foundPkg *CXPackage
	for _, mod := range prgrm.Packages {
		if modName == mod.Name {

			foundPkg = mod
			break
		}
	}
	
	var foundStrct *CXStruct
	
	for _, strct := range foundPkg.Structs {
		if strct.Name == strctName {
			foundStrct = strct
			break
		}
	}

	if foundStrct == nil {
		//looking in imports
		typParts := strings.Split(strctName, ".")

		if mod, err := prgrm.GetPackage(modName); err == nil {
			for _, imp := range mod.Imports {
				for _, strct := range imp.Structs {
					if strct.Name == typParts[0] {
						foundStrct = strct
						break
					}
				}
			}
		}
	}

	if foundPkg != nil && foundStrct != nil {
		return foundStrct, nil
	} else {
		return nil, errors.New(fmt.Sprintf("struct '%s' not found in package '%s'", strctName, modName))
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

	for _, imp := range pkg.Imports {
		for _, def := range imp.Globals {
			if def.Name == defName {
				foundDef = def
				break
			}
		}
	}

	if foundDef != nil {
		return foundDef, nil
	} else {
		return nil, errors.New(fmt.Sprintf("global '%s' not found in package '%s'", defName, pkg.Name))
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

	return nil, errors.New(fmt.Sprintf("function '%s' not found in package '%s' or its imports", fnName, pkg.Name))
}

func (pkg *CXPackage) GetMethod (fnName string, receiverType string) (*CXFunction, error) {
	for _, fn := range pkg.Functions {
		if fn.Name == fnName && len(fn.Inputs) > 0 && fn.Inputs[0].CustomType != nil && fn.Inputs[0].CustomType.Name == receiverType {
			return fn, nil
		}
	}
	
	return nil, errors.New(fmt.Sprintf("method '%s' not found in package '%s'", fnName, pkg.Name))
}

func (prgrm *CXProgram) GetFunction (fnName string, pkgName string) (*CXFunction, error) {
	// I need to first look for the function in the current package
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		for _, fn := range pkg.Functions {
			if fn.Name == fnName {
				return fn, nil
			}
		}
	}

	var foundPkg *CXPackage
	for _, pkg := range prgrm.Packages {
		if pkgName == pkg.Name {
			foundPkg = pkg
			break
		}
	}

	var foundFn *CXFunction
	if foundPkg != nil {
		for _, fn := range foundPkg.Functions {
			if fn.Name == fnName {
				foundFn = fn
				break
			}
		}
	} else {
		return nil, errors.New(fmt.Sprintf("package '%s' not found", pkgName))
	}

	if foundPkg != nil && foundFn != nil {
		return foundFn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("function '%s' not found in package '%s'", fnName, pkgName))
	}
}

func (fn *CXFunction) GetExpressions () ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	} else {
		return nil, errors.New(fmt.Sprintf("function '%s' has no expressions", fn.Name))
	}
}

func (fn *CXFunction) GetExpressionByLabel (lbl string) (*CXExpression, error) {
	if fn.Expressions != nil {
		for _, expr := range fn.Expressions {
			if expr.Label == lbl {
				return expr, nil
			}
		}

		return nil, errors.New(fmt.Sprintf("expression '%s' not found in function '%s'", lbl, fn.Name))
	} else {
		return nil, errors.New(fmt.Sprintf("function '%s' has no expressions", fn.Name))
	}
}

func (fn *CXFunction) GetExpressionByLine (line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		} else {
			return nil, errors.New(fmt.Sprintf("expression line number '%d' exceeds number of expressions in function '%s'", line, fn.Name))
		}

	} else {
		return nil, errors.New(fmt.Sprintf("function '%s' has no expressions", fn.Name))
	}
}

func (expr *CXExpression) GetInputs () ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	} else {
		return nil, errors.New("expression has no arguments")
	}
}
