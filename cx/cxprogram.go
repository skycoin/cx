package base

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/satori/go.uuid"
)

/*
 * The CXProgram struct contains a full program.
 *
 * It is the root data structures for all code, variable and data structures
 * declarations.
 */

type CXProgram struct {
	Packages       []*CXPackage
	Memory         []byte
	Inputs         []*CXArgument
	Outputs        []*CXArgument
	CallStack      []CXCall
	Path           string
	CurrentPackage *CXPackage
	CallCounter    int
	HeapPointer    int
	StackPointer   int
	HeapStartsAt   int
	ElementID      UUID
	Terminated     bool
}

func MakeProgram() *CXProgram {
	newPrgrm := &CXProgram{
		ElementID: MakeElementID(),
		Packages:  make([]*CXPackage, 0),
		CallStack: make([]CXCall, CALLSTACK_SIZE),
		Memory:    make([]byte, STACK_SIZE+TYPE_POINTER_SIZE+INIT_HEAP_SIZE),
	}

	return newPrgrm
}

// ----------------------------------------------------------------
//                             Getters

func (prgrm *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	if prgrm.CurrentPackage != nil {
		return prgrm.CurrentPackage, nil
	} else {
		return nil, errors.New("current package is nil")
	}
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

func (prgrm *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	if prgrm.CurrentPackage != nil &&
		prgrm.CurrentPackage.CurrentFunction != nil &&
		prgrm.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return prgrm.CurrentPackage.CurrentFunction.CurrentExpression, nil
	} else {
		return nil, errors.New("current package, function or expression is nil")
	}
}

func (prgrm *CXProgram) GetGlobal(name string) (*CXArgument, error) {
	mod, err := prgrm.GetCurrentPackage()
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
	} else {
		return foundArgument, nil
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
			return nil, fmt.Errorf("package '%s' not found", modName)
		}
	} else {
		return nil, fmt.Errorf("package '%s' not found", modName)
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
		return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, modName)
	}
}

func (prgrm *CXProgram) GetFunction(fnName string, pkgName string) (*CXFunction, error) {
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
		return nil, fmt.Errorf("package '%s' not found", pkgName)
	}

	if foundPkg != nil && foundFn != nil {
		return foundFn, nil
	} else {
		return nil, fmt.Errorf("function '%s' not found in package '%s'", fnName, pkgName)
	}
}

// ----------------------------------------------------------------
//                         Package handling

func (prgrm *CXProgram) AddPackage(mod *CXPackage) *CXProgram {
	found := false
	for _, md := range prgrm.Packages {
		if md.Name == mod.Name {
			prgrm.CurrentPackage = md
			found = true
			break
		}
	}
	if !found {
		prgrm.Packages = append(prgrm.Packages, mod)
		prgrm.CurrentPackage = mod
	}
	return prgrm
}

func (prgrm *CXProgram) RemovePackage(modName string) {
	lenMods := len(prgrm.Packages)
	for i, mod := range prgrm.Packages {
		if mod.Name == modName {
			if i == lenMods-1 {
				prgrm.Packages = prgrm.Packages[:len(prgrm.Packages)-1]
			} else {
				prgrm.Packages = append(prgrm.Packages[:i], prgrm.Packages[i+1:]...)
			}
			break
		}
	}
}

// ----------------------------------------------------------------
//                             Selectors

func (cxt *CXProgram) SelectPackage(name string) (*CXPackage, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectPackage(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)

	var found *CXPackage
	for _, mod := range cxt.Packages {
		if mod.Name == name {
			cxt.CurrentPackage = mod
			found = mod
		}
	}

	if found == nil {
		return nil, fmt.Errorf("Package '%s' does not exist", name)
	}

	return found, nil
}

func (cxt *CXProgram) SelectFunction(name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectFunction(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)

	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectFunction(name)
	} else {
		return nil, err
	}
}

func (cxt *CXProgram) SelectStruct(name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectStruct(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)

	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectStruct(name)
	} else {
		return nil, err
	}
}

func (cxt *CXProgram) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		cxt.SelectExpression(line)
	// 	},
	// }
	// saveProgramStep(prgrmStep, cxt)

	mod, err := cxt.GetCurrentPackage()
	if err == nil {
		return mod.SelectExpression(line)
	} else {
		return nil, err
	}
}
