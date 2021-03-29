package ast

import (
	"errors"
	"fmt"
	"strings"
)

//cxprogram.CurrentPackag
//current package is only used by affordances
//also used by serialize
//Should be moved to AstWalker

// Only two useres, both in cx/execute.go
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

// GetCurrentPackage ...
func (cxprogram *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	if cxprogram.CurrentPackage != nil {
		return cxprogram.CurrentPackage, nil
	}
	return nil, errors.New("current package is nil")

}

// GetCurrentStruct ...
func (cxprogram *CXProgram) GetCurrentStruct() (*CXStruct, error) {
	if cxprogram.CurrentPackage != nil {
		if cxprogram.CurrentPackage.CurrentStruct != nil {
			return cxprogram.CurrentPackage.CurrentStruct, nil
		}
		return nil, errors.New("current struct is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentFunction ...
func (cxprogram *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	if cxprogram.CurrentPackage != nil {
		if cxprogram.CurrentPackage.CurrentFunction != nil {
			return cxprogram.CurrentPackage.CurrentFunction, nil
		}
		return nil, errors.New("current function is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentExpression ...
func (cxprogram *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	if cxprogram.CurrentPackage != nil &&
		cxprogram.CurrentPackage.CurrentFunction != nil &&
		cxprogram.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return cxprogram.CurrentPackage.CurrentFunction.CurrentExpression, nil
	}
	return nil, errors.New("current package, function or expression is nil")
}

// GetGlobal ...
/*
func (cxprogram *CXProgram) GetGlobal(name string) (*CXArgument, error) {
	mod, err := cxprogram.GetCurrentPackage()
	if err != nil {
		return nil, err
	}

	var foundArgument *CXArgument
	for _, def := range mod.Globals {
		if def.Name == name {
			foundArgument = def
			break
		}
	}

	for _, imp := range mod.Imports {
		for _, def := range imp.Globals {
			if def.Name == name {
				foundArgument = def
				break
			}
		}
	}

	if foundArgument == nil {
		return nil, fmt.Errorf("global '%s' not found", name)
	}
	return foundArgument, nil
}
*/

// Refactor to return nil on error
func (cxprogram *CXProgram) GetPackage(packageNameToFind string) (*CXPackage, error) {
	//iterate packages looking for package; same as GetPackage?
	for _, cxpackage := range cxprogram.Packages {
		if cxpackage.Name == packageNameToFind {
			return cxpackage, nil //can return once found
		}
	}
	//not found
	return nil, fmt.Errorf("package '%s' not found", packageNameToFind)
}

// GetStruct ...
func (cxprogram *CXProgram) GetStruct(strctName string, modName string) (*CXStruct, error) {
	var foundPkg *CXPackage
	for _, mod := range cxprogram.Packages {
		if modName == mod.Name {
			foundPkg = mod
			break
		}
	}

	var foundStrct *CXStruct

	if foundPkg != nil {
		for _, strct := range foundPkg.Structs {
			if strct.Name == strctName {
				foundStrct = strct
				break
			}
		}
	}

	if foundStrct == nil {
		//looking in imports
		typParts := strings.Split(strctName, ".")

		if mod, err := cxprogram.GetPackage(modName); err == nil {
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
	}
	return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, modName)

}

// GetFunction ...
func (cxprogram *CXProgram) GetFunction(functionNameToFind string, pkgName string) (*CXFunction, error) {
	// I need to first look for the function in the current package


	//TODO: WHEN WOULD CurrentPackage not be in cxprogram.Packages?
	//TODO: Add assert to crash if CurrentPackage is not in cxprogram.Packages
	if pkg, err := cxprogram.GetCurrentPackage(); err == nil {
		for _, fn := range pkg.Functions {
			if fn.Name == functionNameToFind {
				return fn, nil
			}
		}
	}

	//iterate packages until the package is found
	//Same as GetPackage? Use GetPackage
	var foundPkg *CXPackage
	for _, pkg := range cxprogram.Packages {
		if pkgName == pkg.Name {
			foundPkg = pkg
			break
		}
	}
	if foundPkg == nil {
		return nil, fmt.Errorf("package '%s' not found", pkgName)
	}

	//iterates package to find function
	//same as GetFunction?
	for _, fn := range foundPkg.Functions {
		if fn.Name == functionNameToFind {
			return fn, nil //can return when found
		}
	}
	return nil, fmt.Errorf("function '%s' not found in package '%s'", functionNameToFind, pkgName)
}



// GetCurrentCall returns the current CXCall
//TODO: What does this do?
//TODO: Only used in OP_JMP
func (cxprogram *CXProgram) GetCurrentCall() *CXCall {
	return &cxprogram.CallStack[cxprogram.CallCounter]
}

/*
// GetCurrentOpCode returns the current OpCode
func (cxprogram *CXProgram) GetCurrentOpCode() int {
	return cxprogram.GetCurrentExpression2().Operator.OpCode
}
*/

/*
//not used
func (cxprogram *CXProgram) GetFramePointer() int {
	return cxprogram.GetCurrentCall().FramePointer
}
*/
