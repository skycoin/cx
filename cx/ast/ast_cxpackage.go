package ast

import (
	"errors"
	"fmt"
)

type CXPackageIndex int

// CXPackage is used to represent a CX package.
type CXPackage struct {
	// Metadata
	Name  string // Name of the package
	Index int    // Index of package inside the CXPackage array

	// Contents
	Imports   map[string]CXPackageIndex // imported packages
	Functions map[string]*CXFunction    // declared functions in this package
	Structs   map[string]*CXStruct      // declared structs in this package
	Globals   []*CXArgument             // declared global variables in this package

	// Used by the REPL and cxgo
	CurrentFunction *CXFunction
	CurrentStruct   *CXStruct
}

// Only Used by Affordances in op_aff.go
func (pkg *CXPackage) GetFunction(prgrm *CXProgram, fnName string) (*CXFunction, error) {
	if fn := pkg.Functions[fnName]; fn != nil {
		return fn, nil
	}

	// now checking in imported packages
	for _, impIdx := range pkg.Imports {
		imp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}

		if fn := imp.Functions[fnName]; fn != nil {
			return fn, nil
		}
	}

	return nil, fmt.Errorf("function '%s' not found in package '%s' or its imports", fnName, pkg.Name)
}

// GetImport ...
func (pkg *CXPackage) GetImport(prgrm *CXProgram, impName string) (*CXPackage, error) {
	if _, ok := pkg.Imports[impName]; !ok {
		return nil, fmt.Errorf("package '%s' not imported", impName)
	}

	imp, err := prgrm.GetPackageFromArray(pkg.Imports[impName])
	if err != nil {
		panic(err)
	}
	return imp, nil
}

// GetMethod ...
func (pkg *CXPackage) GetMethod(fnName string, receiverType string) (*CXFunction, error) {

	if fn := pkg.Functions[fnName]; fn != nil && len(fn.Inputs) > 0 && fn.Inputs[0].StructType != nil && fn.Inputs[0].StructType.Name == receiverType {
		return fn, nil
	}

	// Trying to find it in `Natives`.
	// Most likely a method from a core package.
	if opCode, found := OpCodes[pkg.Name+"."+fnName]; found {
		return Natives[opCode], nil
	}

	return nil, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
}

// GetStruct ...
func (pkg *CXPackage) GetStruct(prgrm *CXProgram, strctName string) (*CXStruct, error) {
	if strct := pkg.Structs[strctName]; strct != nil {
		return strct, nil
	}

	// looking in imports
	for _, impIdx := range pkg.Imports {
		imp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		if strct := imp.Structs[strctName]; strct != nil {
			return strct, nil
		}
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
	if pkg.CurrentFunction == nil {
		return nil, errors.New("current function is nil")
	}

	return pkg.CurrentFunction, nil
}

// ----------------------------------------------------------------
//                             `CXPackage` Selectors

// SelectFunction ...
func (pkg *CXPackage) SelectFunction(name string) (*CXFunction, error) {
	fn := pkg.Functions[name]
	if fn == nil {
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}

	pkg.CurrentFunction = fn
	return fn, nil
}

// MakePackage creates a new empty CXPackage.
//
// It can be filled in later with imports, structs, globals and functions.
//
func MakePackage(name string) *CXPackage {
	return &CXPackage{
		Name:      name,
		Globals:   make([]*CXArgument, 0, 10),
		Imports:   make(map[string]CXPackageIndex, 0),
		Structs:   make(map[string]*CXStruct, 0),
		Functions: make(map[string]*CXFunction, 0),
	}
}

// ----------------------------------------------------------------
//                             `CXPackage` Getters

// GetCurrentStruct ...
func (pkg *CXPackage) GetCurrentStruct() (*CXStruct, error) {
	if pkg.CurrentStruct == nil {
		return nil, errors.New("current struct is nil")
	}

	return pkg.CurrentStruct, nil
}

// ----------------------------------------------------------------
//                     `CXPackage` Member handling

// AddImport ...
func (pkg *CXPackage) AddImport(prgrm *CXProgram, imp *CXPackage) *CXPackage {
	if _, ok := pkg.Imports[imp.Name]; !ok {
		// impIdx := prgrm.AddPackageInArray(imp)
		pkg.Imports[imp.Name] = CXPackageIndex(imp.Index)
	}

	return pkg
}

// RemoveImport ...
func (pkg *CXPackage) RemoveImport(impName string) {
	if _, ok := pkg.Imports[impName]; ok {
		delete(pkg.Imports, impName)
	}
}

// AddFunction ...
func (pkg *CXPackage) AddFunction(fn *CXFunction) *CXPackage {
	fn.Package = CXPackageIndex(pkg.Index)

	if pkg.Functions[fn.Name] != nil {
		println(CompilationError(fn.FileName, fn.FileLine), "function redeclaration")
	}

	pkg.Functions[fn.Name] = fn
	pkg.CurrentFunction = fn

	return pkg
}

// RemoveFunction ...
func (pkg *CXPackage) RemoveFunction(fnName string) {
	if pkg.Functions[fnName] == nil {
		return
	}

	delete(pkg.Functions, fnName)
}

// AddStruct ...
func (pkg *CXPackage) AddStruct(strct *CXStruct) *CXPackage {
	strct.Package = CXPackageIndex(pkg.Index)
	pkg.Structs[strct.Name] = strct
	pkg.CurrentStruct = strct

	return pkg
}

// RemoveStruct ...
func (pkg *CXPackage) RemoveStruct(strctName string) {
	if pkg.Structs[strctName] == nil {
		return
	}
	delete(pkg.Structs, strctName)
}

// AddGlobal ...
func (pkg *CXPackage) AddGlobal(def *CXArgument) *CXPackage {
	def.Package = CXPackageIndex(pkg.Index)
	found := false
	for i, df := range pkg.Globals {
		if df.Name == def.Name {
			pkg.Globals[i] = def
			found = true
			break
		}
	}
	if !found {
		pkg.Globals = append(pkg.Globals, def)
	}

	return pkg
}

// RemoveGlobal ...
func (pkg *CXPackage) RemoveGlobal(defName string) {
	lenGlobals := len(pkg.Globals)
	for i, def := range pkg.Globals {
		if def.Name == defName {
			if i == lenGlobals-1 {
				pkg.Globals = pkg.Globals[:len(pkg.Globals)-1]
			} else {
				pkg.Globals = append(pkg.Globals[:i], pkg.Globals[i+1:]...)
			}
			break
		}
	}
}
