package cxcore

import (
	"errors"
	"fmt"
	"strings"
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
	Path string // Path to the CX project in the filesystem

	// Contents
	Packages []*CXPackage // Packages in a CX program

	// Runtime information
	Inputs         []*CXArgument // OS input arguments
	Outputs        []*CXArgument // outputs to the OS
	Memory         []byte        // Used when running the program
	StackSize      int           // This field stores the size of a CX program's stack
	HeapSize       int           // This field stores the size of a CX program's heap
	HeapStartsAt   int           // Offset at which the heap starts in a CX program's memory
	StackPointer   int           // At what byte the current stack frame is
	CallStack      []CXCall      // Collection of function calls
	CallCounter    int           // What function call is the currently being executed in the CallStack
	HeapPointer    int           // At what offset a CX program can insert a new object to the heap
	Terminated     bool          // Utility field for the runtime. Indicates if a CX program has already finished or not.
	BCPackageCount int           // In case of a CX chain, how many packages of this program are part of blockchain code.
	Version        string        // CX version used to build this CX program.

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
	minHeapSize := minHeapSize()
	newPrgrm := &CXProgram{
		Packages:    make([]*CXPackage, 0),
		CallStack:   make([]CXCall, CALLSTACK_SIZE),
		Memory:      make([]byte, STACK_SIZE+minHeapSize),
		StackSize:   STACK_SIZE,
		HeapSize:    minHeapSize,
		HeapPointer: NULL_HEAP_ADDRESS_OFFSET, // We can start adding objects to the heap after the NULL (nil) bytes.
	}

	return newPrgrm
}

// ----------------------------------------------------------------
//                             Getters

// GetCurrentPackage ...
func (prgrm *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	if prgrm.CurrentPackage != nil {
		return prgrm.CurrentPackage, nil
	}
	return nil, errors.New("current package is nil")

}

// GetCurrentStruct ...
func (prgrm *CXProgram) GetCurrentStruct() (*CXStruct, error) {
	if prgrm.CurrentPackage != nil {
		if prgrm.CurrentPackage.CurrentStruct != nil {
			return prgrm.CurrentPackage.CurrentStruct, nil
		}
		return nil, errors.New("current struct is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentFunction ...
func (prgrm *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	if prgrm.CurrentPackage != nil {
		if prgrm.CurrentPackage.CurrentFunction != nil {
			return prgrm.CurrentPackage.CurrentFunction, nil
		}
		return nil, errors.New("current function is nil")

	}
	return nil, errors.New("current package is nil")

}

// GetCurrentExpression ...
func (prgrm *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	if prgrm.CurrentPackage != nil &&
		prgrm.CurrentPackage.CurrentFunction != nil &&
		prgrm.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return prgrm.CurrentPackage.CurrentFunction.CurrentExpression, nil
	}
	return nil, errors.New("current package, function or expression is nil")

}

// GetGlobal ...
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
	}
	return foundArgument, nil
}

// GetPackage ...
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
		}
		return nil, fmt.Errorf("package '%s' not found", modName)

	}
	return nil, fmt.Errorf("package '%s' not found", modName)

}

// GetStruct ...
func (prgrm *CXProgram) GetStruct(strctName string, modName string) (*CXStruct, error) {
	var foundPkg *CXPackage
	for _, mod := range prgrm.Packages {
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
	}
	return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, modName)

}

// GetFunction ...
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
		if foundPkg != nil {
			for _, fn := range foundPkg.Functions {
				if fn.Name == fnName {
					foundFn = fn
					break
				}
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

// GetCall returns the current CXCall
func (prgrm *CXProgram) GetCall() *CXCall {
	return &prgrm.CallStack[prgrm.CallCounter]
}

// GetExpr returns the current CXExpression
func (prgrm *CXProgram) GetExpr() *CXExpression {
	call := prgrm.GetCall()
	return call.Operator.Expressions[call.Line]
}

// GetOpCode returns the current OpCode
func (prgrm *CXProgram) GetOpCode() int {
	return prgrm.GetExpr().Operator.OpCode
}

// GetFramePointer returns the current frame pointer
func (prgrm *CXProgram) GetFramePointer() int {
	return prgrm.GetCall().FramePointer
}

// ----------------------------------------------------------------
//                         Package handling

// AddPackage ...
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

// RemovePackage ...
func (prgrm *CXProgram) RemovePackage(modName string) {
	lenMods := len(prgrm.Packages)
	for i, mod := range prgrm.Packages {
		if mod.Name == modName {
			if i == lenMods-1 {
				prgrm.Packages = prgrm.Packages[:len(prgrm.Packages)-1]
			} else {
				prgrm.Packages = append(prgrm.Packages[:i], prgrm.Packages[i+1:]...)
			}
			// This means that we're removing the package set to be the CurrentPackage.
			// If it is removed from the program's list of packages, prgrm.CurrentPackage
			// would be pointing to a package meant to be collected by the GC.
			// We fix this by pointing to the last package in the program's list of packages.
			if mod == prgrm.CurrentPackage {
				prgrm.CurrentPackage = prgrm.Packages[len(prgrm.Packages)-1]
			}
			break
		}
	}
}

// ----------------------------------------------------------------
//                             Selectors

// SelectProgram sets `PROGRAM` to the the receiver `prgrm`. This is a utility function used mainly
// by CX chains. `PROGRAM` is used in multiple parts of the CX runtime as a convenience; instead of having
// to pass around a parameter of type CXProgram, the CX program currently being run is accessible through
// `PROGRAM`.
func (prgrm *CXProgram) SelectProgram() (*CXProgram, error) {
	PROGRAM = prgrm

	return PROGRAM, nil
}

// GetProgram returns the CX program assigned to global variable `PROGRAM`.
// This function is mainly used for CX chains.
func GetProgram() (*CXProgram, error) {
	if PROGRAM == nil {
		return nil, fmt.Errorf("a CX program has not been loaded")
	}
	return PROGRAM, nil
}

// SelectPackage ...
func (prgrm *CXProgram) SelectPackage(name string) (*CXPackage, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(prgrm *CXProgram) {
	// 		prgrm.SelectPackage(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, prgrm)

	var found *CXPackage
	for _, mod := range prgrm.Packages {
		if mod.Name == name {
			prgrm.CurrentPackage = mod
			found = mod
		}
	}

	if found == nil {
		return nil, fmt.Errorf("Package '%s' does not exist", name)
	}

	return found, nil
}

// SelectFunction ...
func (prgrm *CXProgram) SelectFunction(name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(prgrm *CXProgram) {
	// 		prgrm.SelectFunction(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, prgrm)

	mod, err := prgrm.GetCurrentPackage()
	if err == nil {
		return mod.SelectFunction(name)
	}
	return nil, err

}

// SelectStruct ...
func (prgrm *CXProgram) SelectStruct(name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(prgrm *CXProgram) {
	// 		prgrm.SelectStruct(name)
	// 	},
	// }
	// saveProgramStep(prgrmStep, prgrm)

	mod, err := prgrm.GetCurrentPackage()
	if err == nil {
		return mod.SelectStruct(name)
	}
	return nil, err

}

// SelectExpression ...
func (prgrm *CXProgram) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(prgrm *CXProgram) {
	// 		prgrm.SelectExpression(line)
	// 	},
	// }
	// saveProgramStep(prgrmStep, prgrm)

	mod, err := prgrm.GetCurrentPackage()
	if err == nil {
		return mod.SelectExpression(line)
	}
	return nil, err

}

// ----------------------------------------------------------------
//                             Debugging

// PrintAllObjects prints all objects in a program
//
func (prgrm *CXProgram) PrintAllObjects() {
	fp := 0

	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			heapOffset := mustDeserializeI32(prgrm.Memory[fp+ptr.Offset : fp+ptr.Offset+TYPE_POINTER_SIZE])

			var byts []byte

			if ptr.CustomType != nil {
				// then it's a pointer to a struct
				// use CustomStruct to match the fields against the bytes
				// for _, fld := range ptr.Fields {

				// }

				byts = prgrm.Memory[int(heapOffset)+OBJECT_HEADER_SIZE : int(heapOffset)+OBJECT_HEADER_SIZE+ptr.CustomType.Size]
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

			fmt.Println("obj", ptr.Name, ptr.CustomType, prgrm.Memory[heapOffset:int(heapOffset)+op.Size], byts)
		}

		fp += op.Size
	}
}
