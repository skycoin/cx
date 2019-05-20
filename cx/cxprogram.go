package cxcore

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amherag/skycoin/src/cipher/encoder"
	. "github.com/satori/go.uuid" // nolint golint
)

/*
 * The CXProgram struct contains a full program.
 *
 * It is the root data structures for all code, variable and data structures
 * declarations.
 */

// CXProgram is used to represent a full CX program.
//
// It is the root data structure for the declarations of all functions,
// variables and data structures.
//
type CXProgram struct {
	// Metadata
	Path      string // Path to the CX project in the filesystem
	ElementID UUID   // Was supposed to be used for blockchain integration. Needs to be removed.

	// Contents
	Packages []*CXPackage // Packages in a CX program

	// Runtime information
	Inputs       []*CXArgument // OS input arguments
	Outputs      []*CXArgument // outputs to the OS
	Memory       []byte        // Used when running the program
	StackSize    int           // This field stores the size of a CX program's stack
	HeapSize     int           // This field stores the size of a CX program's heap
	HeapStartsAt int           // Offset at which the heap starts in a CX program's memory
	StackPointer int           // At what byte the current stack frame is
	CallStack    []CXCall      // Collection of function calls
	CallCounter  int           // What function call is the currently being executed in the CallStack
	HeapPointer  int           // At what offset a CX program can insert a new object to the heap
	Terminated   bool          // Utility field for the runtime. Indicates if a CX program has already finished or not.

	// Used by the REPL and parser
	CurrentPackage *CXPackage // Represents the currently active package in the REPL or when parsing a CX file.
}

// CXCall ...
type CXCall struct {
	Operator     *CXFunction // What CX function will be called when running this CXCall in the runtime
	Line         int         // What line in the CX function is currently being executed
	FramePointer int         // Where in the stack is this function call's local variables stored
}

// MakeProgram ...
func MakeProgram() *CXProgram {
	newPrgrm := &CXProgram{
		ElementID:   MakeElementID(),
		Packages:    make([]*CXPackage, 0),
		CallStack:   make([]CXCall, CALLSTACK_SIZE),
		Memory:      make([]byte, STACK_SIZE+TYPE_POINTER_SIZE+INIT_HEAP_SIZE),
		StackSize:   STACK_SIZE,
		HeapSize:    INIT_HEAP_SIZE,
		HeapPointer: NULL_HEAP_ADDRESS_OFFSET, // We can start adding objects to the heap after the NULL (nil) bytes
	}

	return newPrgrm
}

// ----------------------------------------------------------------
//                             Getters

// GetCurrentPackage ...
func (cxt *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	if cxt.CurrentPackage != nil {
		return cxt.CurrentPackage, nil
	}
	return nil, errors.New("current package is nil")

}

// GetCurrentStruct ...
func (cxt *CXProgram) GetCurrentStruct() (*CXStruct, error) {
	if cxt.CurrentPackage != nil {
		if cxt.CurrentPackage.CurrentStruct != nil {
			return cxt.CurrentPackage.CurrentStruct, nil
		}
		return nil, errors.New("current struct is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentFunction ...
func (cxt *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	if cxt.CurrentPackage != nil {
		if cxt.CurrentPackage.CurrentFunction != nil {
			return cxt.CurrentPackage.CurrentFunction, nil
		}
		return nil, errors.New("current function is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentExpression ...
func (cxt *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	if cxt.CurrentPackage != nil &&
		cxt.CurrentPackage.CurrentFunction != nil &&
		cxt.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return cxt.CurrentPackage.CurrentFunction.CurrentExpression, nil
	}
	return nil, errors.New("current package, function or expression is nil")

}

// GetGlobal ...
func (cxt *CXProgram) GetGlobal(name string) (*CXArgument, error) {
	mod, err := cxt.GetCurrentPackage()
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

// GetPackage ...
func (cxt *CXProgram) GetPackage(modName string) (*CXPackage, error) {
	if cxt.Packages != nil {
		var found *CXPackage
		for _, mod := range cxt.Packages {
			if modName == mod.Name {
				found = mod
				break
			}
		}
		if found != nil {
			return found, nil
		}
		return nil, fmt.Errorf("package '%s' not found", modName)

	}
	return nil, fmt.Errorf("package '%s' not found", modName)

}

// GetStruct ...
func (cxt *CXProgram) GetStruct(strctName string, modName string) (*CXStruct, error) {
	var foundPkg *CXPackage
	for _, mod := range cxt.Packages {
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

		if mod, err := cxt.GetPackage(modName); err == nil {
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
func (cxt *CXProgram) GetFunction(fnName string, pkgName string) (*CXFunction, error) {
	// I need to first look for the function in the current package
	if pkg, err := cxt.GetCurrentPackage(); err == nil {
		for _, fn := range pkg.Functions {
			if fn.Name == fnName {
				return fn, nil
			}
		}
	}

	var foundPkg *CXPackage
	for _, pkg := range cxt.Packages {
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
	}
	return nil, fmt.Errorf("function '%s' not found in package '%s'", fnName, pkgName)

}

// ----------------------------------------------------------------
//                         Package handling

// AddPackage ...
func (cxt *CXProgram) AddPackage(mod *CXPackage) *CXProgram {
	found := false
	for _, md := range cxt.Packages {
		if md.Name == mod.Name {
			cxt.CurrentPackage = md
			found = true
			break
		}
	}
	if !found {
		cxt.Packages = append(cxt.Packages, mod)
		cxt.CurrentPackage = mod
	}
	return cxt
}

// RemovePackage ...
func (cxt *CXProgram) RemovePackage(modName string) {
	lenMods := len(cxt.Packages)
	for i, mod := range cxt.Packages {
		if mod.Name == modName {
			if i == lenMods-1 {
				cxt.Packages = cxt.Packages[:len(cxt.Packages)-1]
			} else {
				cxt.Packages = append(cxt.Packages[:i], cxt.Packages[i+1:]...)
			}
			// This means that we're removing the package set to be the CurrentPackage.
			// If it is removed from the program's list of packages, cxt.CurrentPackage
			// would be pointing to a package meant to be collected by the GC.
			// We fix this by pointing to the last package in the program's list of packages.
			if mod == cxt.CurrentPackage {
				cxt.CurrentPackage = cxt.Packages[len(cxt.Packages)-1]
			}
			break
		}
	}
}

// ----------------------------------------------------------------
//                             Selectors

// SelectPackage ...
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

// SelectFunction ...
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
	}
	return nil, err

}

// SelectStruct ...
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
	}
	return nil, err

}

// SelectExpression ...
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
	}
	return nil, err

}

// ----------------------------------------------------------------
//                             Debugging

// PrintAllObjects prints all objects in a program
//
func (cxt *CXProgram) PrintAllObjects() {
	fp := 0

	for c := 0; c <= cxt.CallCounter; c++ {
		op := cxt.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			_, err := encoder.DeserializeAtomic(cxt.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)
			if err != nil {
				panic(err)
			}

			var byts []byte

			if ptr.CustomType != nil {
				// then it's a pointer to a struct
				// use CustomStruct to match the fields against the bytes
				// for _, fld := range ptr.Fields {

				// }

				byts = cxt.Memory[int(heapOffset)+OBJECT_HEADER_SIZE : int(heapOffset)+OBJECT_HEADER_SIZE+ptr.CustomType.Size]
			}

			// var currLengths []int
			// var currCustom *CXStruct

			// for c := len(ptr.DeclarationSpecifiers) - 1; c >= 0; c-- {
			// 	// we need to go backwards in here

			// 	switch ptr.DeclarationSpecifiers[c] {
			// 	case DECL_POINTER:
			// 		// we might not need to do anything
			// 	case DECL_ARRAY:
			// 		currLengths = ptr.Lengths
			// 	case DECL_SLICE:
			// 	case DECL_STRUCT:
			// 		currCustom = ptr.CustomType
			// 	case DECL_BASIC:
			// 	}
			// }

			// if len(ptr.Lengths) > 0 {
			// 	fmt.Println("ARRAY")
			// }

			// if ptr.CustomType != nil {
			// 	fmt.Println("STRUCT")
			// }

			fmt.Println("declarat", ptr.DeclarationSpecifiers)

			fmt.Println("obj", ptr.Name, ptr.CustomType, cxt.Memory[heapOffset:int(heapOffset)+op.Size], byts)
		}

		fp += op.Size
	}
}
