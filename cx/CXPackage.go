package base

import (
	"errors"
	"fmt"
)

/* The CXPackage struct contains information about a CX package.
 */

type CXPackage struct {
	Imports         []*CXPackage
	Functions       []*CXFunction
	Structs         []*CXStruct
	Globals         []*CXArgument
	Name            string
	CurrentFunction *CXFunction
	CurrentStruct   *CXStruct
	ElementID       UUID
}

func MakePackage(name string) *CXPackage {
	return &CXPackage{
		ElementID: MakeElementID(),
		Name:      name,
		Globals:   make([]*CXArgument, 0, 10),
		Imports:   make([]*CXPackage, 0),
		Functions: make([]*CXFunction, 0, 10),
		Structs:   make([]*CXStruct, 0),
	}
}

// ----------------------------------------------------------------
//                             Getters

func (pkg *CXPackage) GetImport(impName string) (*CXPackage, error) {
	for _, imp := range pkg.Imports {
		if imp.Name == impName {
			return imp, nil
		}
	}
	return nil, fmt.Errorf("package '%s' not imported", impName)
}

func (mod *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		return mod.Functions, nil
	} else {
		return nil, fmt.Errorf("package '%s' has no functions", mod.Name)
	}
}

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

func (pkg *CXPackage) GetMethod(fnName string, receiverType string) (*CXFunction, error) {
	for _, fn := range pkg.Functions {
		if fn.Name == fnName && len(fn.Inputs) > 0 && fn.Inputs[0].CustomType != nil && fn.Inputs[0].CustomType.Name == receiverType {
			return fn, nil
		}
	}

	return nil, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
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

func (mod *CXPackage) GetCurrentFunction() (*CXFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	}

	return nil, errors.New("current function is nil")
}

func (mod *CXPackage) GetCurrentStruct() (*CXStruct, error) {
	if mod.CurrentStruct != nil {
		return mod.CurrentStruct, nil
	}

	return nil, errors.New("current struct is nil")
}

// ----------------------------------------------------------------
//                     Member handling

func (mod *CXPackage) AddImport(imp *CXPackage) *CXPackage {
	found := false
	for _, im := range mod.Imports {
		if im.Name == imp.Name {
			found = true
			break
		}
	}
	if !found {
		mod.Imports = append(mod.Imports, imp)
	}

	return mod
}

func (mod *CXPackage) RemoveImport(impName string) {
	lenImps := len(mod.Imports)
	for i, imp := range mod.Imports {
		if imp.Name == impName {
			if i == lenImps-1 {
				mod.Imports = mod.Imports[:len(mod.Imports)-1]
			} else {
				mod.Imports = append(mod.Imports[:i], mod.Imports[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) AddFunction(fn *CXFunction) *CXPackage {
	fn.Package = mod

	found := false
	for i, f := range mod.Functions {
		if f.Name == fn.Name {
			mod.Functions[i].Name = fn.Name
			mod.Functions[i].Inputs = fn.Inputs
			mod.Functions[i].Outputs = fn.Outputs
			mod.Functions[i].Expressions = fn.Expressions
			mod.Functions[i].CurrentExpression = fn.CurrentExpression
			mod.Functions[i].Package = fn.Package
			mod.CurrentFunction = mod.Functions[i]
			found = true
			break
		}
	}
	if !found {
		mod.Functions = append(mod.Functions, fn)
		mod.CurrentFunction = fn
	}

	return mod
}

func (mod *CXPackage) RemoveFunction(fnName string) {
	lenFns := len(mod.Functions)
	for i, fn := range mod.Functions {
		if fn.Name == fnName {
			if i == lenFns-1 {
				mod.Functions = mod.Functions[:len(mod.Functions)-1]
			} else {
				mod.Functions = append(mod.Functions[:i], mod.Functions[i+1:]...)
			}
			break
		}
	}
}

func (pkg *CXPackage) AddStruct(strct *CXStruct) *CXPackage {
	found := false
	for i, s := range pkg.Structs {
		if s.Name == strct.Name {
			pkg.Structs[i] = strct
			found = true
			break
		}
	}
	if !found {
		pkg.Structs = append(pkg.Structs, strct)
	}

	strct.Package = pkg
	pkg.CurrentStruct = strct

	return pkg
}

func (mod *CXPackage) RemoveStruct(strctName string) {
	lenStrcts := len(mod.Structs)
	for i, strct := range mod.Structs {
		if strct.Name == strctName {
			if i == lenStrcts-1 {
				mod.Structs = mod.Structs[:len(mod.Structs)-1]
			} else {
				mod.Structs = append(mod.Structs[:i], mod.Structs[i+1:]...)
			}
			break
		}
	}
}

func (mod *CXPackage) AddGlobal(def *CXArgument) *CXPackage {
	// def.Program = mod.Program
	def.Package = mod
	found := false
	for i, df := range mod.Globals {
		if df.Name == def.Name {
			mod.Globals[i] = def
			found = true
			break
		}
	}
	if !found {
		mod.Globals = append(mod.Globals, def)
	}
	return mod
}

func (mod *CXPackage) RemoveGlobal(defName string) {
	lenGlobals := len(mod.Globals)
	for i, def := range mod.Globals {
		if def.Name == defName {
			if i == lenGlobals-1 {
				mod.Globals = mod.Globals[:len(mod.Globals)-1]
			} else {
				mod.Globals = append(mod.Globals[:i], mod.Globals[i+1:]...)
			}
			break
		}
	}
}

// ----------------------------------------------------------------
//                             Selectors

func (mod *CXPackage) SelectFunction(name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {

	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectFunction(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

	var found *CXFunction
	for _, fn := range mod.Functions {
		if fn.Name == name {
			mod.CurrentFunction = fn
			found = fn
		}
	}

	if found == nil {
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}

	return found, nil
}

func (mod *CXPackage) SelectStruct(name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectStruct(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)

	var found *CXStruct
	for _, strct := range mod.Structs {
		if strct.Name == name {
			mod.CurrentStruct = strct
			found = strct
		}
	}

	if found == nil {
		return nil, errors.New("Desired structure does not exist")
	}

	return found, nil
}

func (mod *CXPackage) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			mod.SelectExpression(line)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, mod.Context)
	fn, err := mod.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	} else {
		return nil, err
	}
}
