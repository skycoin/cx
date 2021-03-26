package cxcore

import (
	"errors"
	"fmt"
)

// SelectPackage ...
func (cxprogram *CXProgram) SelectPackage(name string) (*CXPackage, error) {

	var found *CXPackage
	for _, mod := range cxprogram.Packages {
		if mod.Name == name {
			cxprogram.CurrentPackage = mod
			found = mod
		}
	}

	if found == nil {
		return nil, fmt.Errorf("Package '%s' does not exist", name)
	}

	return found, nil
}

// GetFunction ...
func (pkg *CXPackage) GetFunction(fnName string) (*CXFunction, error) {
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

// Getfunction2 ...
func (cxprogram *CXProgram) Getfunction2(name string) (*CXFunction, error) {
	mod, err := cxprogram.GetCurrentPackage()
	if err == nil {
		return mod.SelectFunction(name)
	}
	return nil, err
}

// SelectStruct ...
func (cxprogram *CXProgram) SelectStruct(name string) (*CXStruct, error) {
	mod, err := cxprogram.GetCurrentPackage()
	if err == nil {
		return mod.SelectStruct(name)
	}
	return nil, err
}

// SelectExpression ...
func (cxprogram *CXProgram) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxprogram *CXProgram) {
	// 		cxprogram.SelectExpression(line)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxprogram)

	mod, err := cxprogram.GetCurrentPackage()
	if err == nil {
		return mod.SelectExpression(line)
	}
	return nil, err

}

// GetImport ...
func (pkg *CXPackage) GetImport(impName string) (*CXPackage, error) {
	for _, imp := range pkg.Imports {
		if imp.Name == impName {
			return imp, nil
		}
	}
	return nil, fmt.Errorf("package '%s' not imported", impName)
}

// GetFunctions ...
func (pkg *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if pkg.Functions != nil {
		return pkg.Functions, nil
	}
	return nil, fmt.Errorf("package '%s' has no functions", pkg.Name)

}

// GetMethod ...
func (pkg *CXPackage) GetMethod(fnName string, receiverType string) (*CXFunction, error) {
	for _, fn := range pkg.Functions {
		if fn.Name == fnName && len(fn.Inputs) > 0 && fn.Inputs[0].CustomType != nil && fn.Inputs[0].CustomType.Name == receiverType {
			return fn, nil
		}
	}

	// Trying to find it in `Natives`.
	// Most likely a method from a core package.
	if opCode, found := OpCodes[pkg.Name+"."+fnName]; found {
		return Natives[opCode], nil
	}

	return nil, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
}

// GetStruct ...
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
	}
	return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, pkg.Name)

}

// GetGlobal ...
func (pkg *CXPackage) GetGlobal(defName string) (*CXArgument, error) {
	var foundDef *CXArgument
	for _, def := range pkg.Globals {
		if def.Name == defName {
			foundDef = def
			break
		}
	}

	if foundDef != nil {
		return foundDef, nil
	}
	return nil, fmt.Errorf("global '%s' not found in package '%s'", defName, pkg.Name)

}

// GetCurrentFunction ...
func (pkg *CXPackage) GetCurrentFunction() (*CXFunction, error) {
	if pkg.CurrentFunction != nil {
		return pkg.CurrentFunction, nil
	}

	return nil, errors.New("current function is nil")
}

// ----------------------------------------------------------------
//                             `CXPackage` Selectors

// SelectFunction ...
func (pkg *CXPackage) SelectFunction(name string) (*CXFunction, error) {
	var found *CXFunction
	for _, fn := range pkg.Functions {
		if fn.Name == name {
			pkg.CurrentFunction = fn
			found = fn
		}
	}
	if found == nil {
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}
	return found, nil
}

// SelectStruct ...
func (pkg *CXPackage) SelectStruct(name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if pkg, err := cxt.GetCurrentPackage(); err == nil {
	// 			pkg.SelectStruct(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, pkg.Context)

	var found *CXStruct
	for _, strct := range pkg.Structs {
		if strct.Name == name {
			pkg.CurrentStruct = strct
			found = strct
		}
	}

	if found == nil {
		return nil, errors.New("Desired structure does not exist")
	}

	return found, nil
}

// SelectExpression ...
func (pkg *CXPackage) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if pkg, err := cxt.GetCurrentPackage(); err == nil {
	// 			pkg.SelectExpression(line)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, pkg.Context)
	fn, err := pkg.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	}
	return nil, err
}
