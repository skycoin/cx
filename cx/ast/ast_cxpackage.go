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
	Imports   map[string]CXPackageIndex  // imported packages
	Functions map[string]CXFunctionIndex // declared functions in this package
	Structs   map[string]CXStructIndex   // declared structs in this package
	Globals   *CXStruct                  // declared global variables in this package

	// Used by the REPL and cxgo
	CurrentFunction CXFunctionIndex
}

// Only Used by Affordances in op_aff.go
func (pkg *CXPackage) GetFunction(prgrm *CXProgram, fnName string) (*CXFunction, error) {
	if fnIdx, ok := pkg.Functions[fnName]; ok {
		fn := prgrm.GetFunctionFromArray(fnIdx)

		return fn, nil
	}

	// now checking in imported packages
	for _, impIdx := range pkg.Imports {
		imp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}

		if fnIdx, ok := imp.Functions[fnName]; ok {
			fn := prgrm.GetFunctionFromArray(fnIdx)

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
func (pkg *CXPackage) GetMethod(prgrm *CXProgram, fnName string, receiverType string) (CXFunctionIndex, error) {
	if fnIdx, ok := pkg.Functions[fnName]; ok {
		fn := prgrm.GetFunctionFromArray(fnIdx)
		fnInputs := fn.GetInputs(prgrm)
		fnInputTypeSig := fnInputs[0]
		var fnInput *CXArgument
		if fnInputTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
			fnInput = prgrm.GetCXArgFromArray(CXArgumentIndex(fnInputTypeSig.Meta))
		}

		if len(fnInputs) > 0 && fnInput.StructType != nil && fnInput.StructType.Name == receiverType {
			return fnIdx, nil
		}
	}

	// Trying to find it in `Natives`.
	// Most likely a method from a core package.
	if opCode, found := OpCodes[pkg.Name+"."+fnName]; found {
		opFnIdx := prgrm.AddNativeFunctionInArray(Natives[opCode])
		return opFnIdx, nil
	}

	return -1, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
}

// GetStruct ...
func (pkg *CXPackage) GetStruct(prgrm *CXProgram, strctName string) (*CXStruct, error) {
	if strctIdx, ok := pkg.Structs[strctName]; ok {
		return &prgrm.CXStructs[strctIdx], nil
	}

	// looking in imports
	for _, impIdx := range pkg.Imports {
		imp, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		if strctIdx, ok := imp.Structs[strctName]; ok {
			return &prgrm.CXStructs[strctIdx], nil
		}
	}

	return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, pkg.Name)

}

// GetGlobal ...
func (pkg *CXPackage) GetGlobal(prgrm *CXProgram, defName string) (*CXTypeSignature, error) {
	for _, field := range pkg.Globals.Fields {
		if field.Name == defName {
			return field, nil
		}
	}

	return nil, fmt.Errorf("global '%s' not found in package '%s'", defName, pkg.Name)
}

// GetCurrentFunction ...
func (pkg *CXPackage) GetCurrentFunction(prgrm *CXProgram) (*CXFunction, error) {
	if pkg.CurrentFunction == -1 {
		return nil, errors.New("current function is nil")
	}

	return prgrm.GetFunctionFromArray(pkg.CurrentFunction), nil
}

// ----------------------------------------------------------------
//                             `CXPackage` Selectors

// SelectFunction ...
func (pkg *CXPackage) SelectFunction(prgrm *CXProgram, name string) (*CXFunction, error) {
	if _, ok := pkg.Functions[name]; !ok {
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}

	idx := pkg.Functions[name]
	pkg.CurrentFunction = idx

	return prgrm.GetFunctionFromArray(idx), nil
}

// MakePackage creates a new empty CXPackage.
//
// It can be filled in later with imports, structs, globals and functions.
//
func MakePackage(name string) *CXPackage {
	return &CXPackage{
		Name:            name,
		Globals:         &CXStruct{},
		Imports:         make(map[string]CXPackageIndex, 0),
		Structs:         make(map[string]CXStructIndex, 0),
		Functions:       make(map[string]CXFunctionIndex, 0),
		CurrentFunction: -1,
	}
}

// ----------------------------------------------------------------
//                             `CXPackage` Getters

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
func (pkg *CXPackage) AddFunction(prgrm *CXProgram, fn *CXFunction) (*CXPackage, CXFunctionIndex) {
	if _, ok := pkg.Functions[fn.Name]; ok {
		println(CompilationError(fn.FileName, fn.FileLine), "function redeclaration")
	}

	fn.Package = CXPackageIndex(pkg.Index)

	fnIdx := prgrm.AddFunctionInArray(fn)
	pkg.Functions[fn.Name] = fnIdx
	pkg.CurrentFunction = fnIdx

	return pkg, fnIdx
}

// RemoveFunction ...
func (pkg *CXPackage) RemoveFunction(fnName string) {
	if _, ok := pkg.Functions[fnName]; !ok {
		return
	}

	delete(pkg.Functions, fnName)
}

// AddStruct ...
func (pkg *CXPackage) AddStruct(prgrm *CXProgram, strct *CXStruct) *CXPackage {
	if _, ok := pkg.Structs[strct.Name]; ok {
		return pkg
	}
	strct.Package = CXPackageIndex(pkg.Index)
	strctIdx := prgrm.AddStructInArray(strct)
	pkg.Structs[strct.Name] = strctIdx

	return pkg
}

// RemoveStruct ...
func (pkg *CXPackage) RemoveStruct(strctName string) {
	if _, ok := pkg.Structs[strctName]; !ok {
		return
	}
	delete(pkg.Structs, strctName)
}

// AddGlobal ...
func (pkg *CXPackage) AddGlobal(prgrm *CXProgram, defIdx CXArgumentIndex) *CXPackage {
	defArg := prgrm.GetCXArgFromArray(defIdx)
	prgrm.CXArgs[defIdx].Package = CXPackageIndex(pkg.Index)

	for _, field := range pkg.Globals.Fields {
		if field.Name == defArg.Name && field.Type == TYPE_CXARGUMENT_DEPRECATE {
			return pkg
		}
	}

	pkg.Globals.AddField_Globals_CXAtomicOps(prgrm, defIdx)

	return pkg
}

func (pkg *CXPackage) AddGlobal_TypeSignature(prgrm *CXProgram, typeSignature *CXTypeSignature) *CXPackage {
	typeSignature.Package = CXPackageIndex(pkg.Index)
	for _, field := range pkg.Globals.Fields {
		if field.Name == typeSignature.Name {
			return pkg
		}
	}

	pkg.Globals.AddField_TypeSignature(prgrm, typeSignature)

	return pkg
}

// RemoveGlobal ...
// func (pkg *CXPackage) RemoveGlobal(prgrm *CXProgram, defName string) {
// 	lenGlobals := len(pkg.Globals)
// 	for i, defIdx := range pkg.Globals {
// 		def := prgrm.GetCXArg(defIdx)
// 		if def.Name == defName {
// 			if i == lenGlobals-1 {
// 				pkg.Globals = pkg.Globals[:len(pkg.Globals)-1]
// 			} else {
// 				pkg.Globals = append(pkg.Globals[:i], pkg.Globals[i+1:]...)
// 			}
// 			break
// 		}
// 	}
// }
