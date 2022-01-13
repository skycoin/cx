package ast

import (
	"errors"
	"fmt"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// CXProgram is used to represent a full CX program.
//
// It is the root data structure for the declarations of all functions,
// variables and data structures.
//
type CXProgram struct {
	// Metadata
	//Remove Path //moved to cx/globals
	//Path string // Path to the CX project in the filesystem

	Stack StackSegmentStruct
	Data  DataSegmentStruct
	Heap  HeapSegmentStruct

	// Contents
	Packages map[string]CXPackageIndex // Packages in a CX program

	// Runtime information
	ProgramInput []*CXArgument // OS input arguments
	//ProgramOutput []*CXArgument // outputs to the OS
	Memory []byte // Used when running the program

	CallStack   []CXCall      // Collection of function calls. Call stack is not in CX Memory.
	CallCounter types.Pointer // What function call is the currently being executed in the CallStack
	Terminated  bool          // Utility field for the runtime. Indicates if a CX program has already finished or not.
	Version     string        // CX version used to build this CX program.

	// Used by the REPL and cxgo
	CurrentPackage CXPackageIndex // Represents the currently active package in the REPL or when parsing a CX file.
	ProgramError   error

	// For new CX AST arrays
	CXAtomicOps []*CXAtomicOperator
	CXArgs      []CXArgument
	CXLines     []CXLine
	CXPackages  []CXPackage
	CXFunctions []*CXFunction
	// Then reference the package of function by CxFunction id

	// For Initializers
	SysInitExprs []*CXExpression
}

type StackSegmentStruct struct {
	Size     types.Pointer // This field stores the size of a CX program's stack
	StartsAt types.Pointer // Offset at which the stack segment starts in a CX program's memory
	Pointer  types.Pointer // At what byte the current stack frame is
}

type DataSegmentStruct struct {
	Size     types.Pointer // This field stores the size of a CX program's data segment size
	StartsAt types.Pointer // Offset at which the data segment starts in a CX program's memory
}

type HeapSegmentStruct struct {
	Size     types.Pointer // This field stores the size of a CX program's heap
	StartsAt types.Pointer // Offset at which the heap starts in a CX program's memory (normally the stack size)
	Pointer  types.Pointer // At what offset a CX program can insert a new object to the heap
}

// MakeProgram ...
func MakeProgram() *CXProgram {
	minHeapSize := minHeapSize()
	newPrgrm := &CXProgram{
		Packages:  make(map[string]CXPackageIndex, 0),
		CallStack: make([]CXCall, constants.CALLSTACK_SIZE),
		Memory:    make([]byte, constants.STACK_SIZE+minHeapSize),
		Stack: StackSegmentStruct{
			Size: constants.STACK_SIZE,
		},
		Data: DataSegmentStruct{
			StartsAt: constants.STACK_SIZE,
		},
		Heap: HeapSegmentStruct{
			Size:    minHeapSize,
			Pointer: constants.NULL_HEAP_ADDRESS_OFFSET, // We can start adding objects to the heap after the NULL (nil) bytes.
		},
	}
	return newPrgrm
}

// ----------------------------------------------------------------
//                         `CXProgram` Package handling

// AddPackage ...
func (cxprogram *CXProgram) AddPackage(mod *CXPackage) {
	// if cxprogram.Packages[mod.Name] != nil {
	// 	return
	// }

	index := cxprogram.AddPackageInArray(mod)

	cxprogram.Packages[mod.Name] = CXPackageIndex(index)
	cxprogram.CurrentPackage = CXPackageIndex(index)
}

func (cxprogram *CXProgram) AddPackageInArray(pkg *CXPackage) int {
	cxprogram.CXPackages = append(cxprogram.CXPackages, *pkg)

	return len(cxprogram.CXPackages) - 1
}

func (cxprogram *CXProgram) GetPackageFromArray(index CXPackageIndex) (*CXPackage, error) {
	if int(index) > (len(cxprogram.CXPackages) - 1) {
		return nil, fmt.Errorf("error: CXPackages[%d]: index out of bounds", index)
	}

	return &cxprogram.CXPackages[index], nil
}

// RemovePackage ...
func (cxprogram *CXProgram) RemovePackage(modName string) {
	// If doesnt exist, return
	// if cxprogram.Packages[modName] == nil {
	// 	return
	// }

	// Check if it is the current pkg so when it
	// is deleted, it will be replaced with new pkg
	isCurrentPkg := cxprogram.Packages[modName] == cxprogram.CurrentPackage

	// Delete package
	delete(cxprogram.Packages, modName)

	// This means that we're removing the package set to be the CurrentPackage.
	// If it is removed from the program's map of packages, cxprogram.CurrentPackage
	// would be pointing to a package meant to be collected by the GC.
	// We fix this by pointing to random package in the program's map of packages.
	if isCurrentPkg {
		for _, pkg := range cxprogram.Packages {
			cxprogram.CurrentPackage = pkg
			break
		}
	}

}

// ----------------------------------------------------------------
//            `CXProgram` CXLines, CXArgs, and CXAtomicOps Handling

func (cxprogram *CXProgram) GetCXLine(index int) (*CXLine, error) {
	if index > (len(cxprogram.CXLines) - 1) {
		return nil, fmt.Errorf("error: CXLines[%d]: index out of bounds", index)
	}

	return &cxprogram.CXLines[index], nil
}

func (cxprogram *CXProgram) GetCXArg(index int) (*CXArgument, error) {
	if index > (len(cxprogram.CXArgs) - 1) {
		return nil, fmt.Errorf("error: CXArgs[%d]: index out of bounds", index)
	}

	return &cxprogram.CXArgs[index], nil
}

func (cxprogram *CXProgram) GetCXAtomicOp(index int) (*CXAtomicOperator, error) {
	if index > (len(cxprogram.CXAtomicOps) - 1) {
		return nil, fmt.Errorf("error: CXAtomicOps[%d]: index out of bounds", index)
	}

	return cxprogram.CXAtomicOps[index], nil
}

func (cxprogram *CXProgram) AddCXLine(CXLine *CXLine) int {
	cxprogram.CXLines = append(cxprogram.CXLines, *CXLine)

	return len(cxprogram.CXLines) - 1
}

func (cxprogram *CXProgram) AddCXArg(CXArg *CXArgument) int {
	cxprogram.CXArgs = append(cxprogram.CXArgs, *CXArg)

	return len(cxprogram.CXArgs) - 1
}

func (cxprogram *CXProgram) AddCXAtomicOp(CXAtomicOp *CXAtomicOperator) int {
	cxprogram.CXAtomicOps = append(cxprogram.CXAtomicOps, CXAtomicOp)

	return len(cxprogram.CXAtomicOps) - 1
}

// ----------------------------------------------------------------
//                             `CXProgram` Getters

func (cxprogram *CXProgram) GetOperation(expr *CXExpression) (*CXAtomicOperator, *CXArgument, *CXLine, error) {
	var cxAtomicOp *CXAtomicOperator
	var err error
	switch expr.Type {
	case CX_ATOMIC_OPERATOR:
		cxAtomicOp, err = cxprogram.GetCXAtomicOp(int(expr.Index))
		if err != nil {
			return &CXAtomicOperator{}, &CXArgument{}, &CXLine{}, err
		}

		return cxAtomicOp, &CXArgument{}, &CXLine{}, nil
	case CX_LINE:
		cxLine, err := cxprogram.GetCXLine(expr.Index)
		if err != nil {
			return &CXAtomicOperator{}, &CXArgument{}, &CXLine{}, err
		}

		return &CXAtomicOperator{}, &CXArgument{}, cxLine, nil
		// case CX_ARGUMENT:
	}

	return &CXAtomicOperator{}, &CXArgument{}, &CXLine{}, fmt.Errorf("operation type is not found.")
}

func (cxprogram *CXProgram) GetPreviousCXLine(exprs []*CXExpression, currIndex int) (*CXLine, error) {
	for i := currIndex; i >= 0; i-- {
		if exprs[i].Type == CX_LINE {
			_, _, cxLine, err := cxprogram.GetOperation(exprs[i])
			if err != nil {
				return &CXLine{}, err
			}

			return cxLine, nil
		}
	}
	return &CXLine{}, fmt.Errorf("CXLine not found.")
}

func (cxprogram *CXProgram) GetCXAtomicOpFromExpressions(exprs []*CXExpression, currIndex int) (*CXAtomicOperator, error) {
	for i := currIndex; i < len(exprs); i++ {
		if exprs[i].Type == CX_ATOMIC_OPERATOR {
			cxAtomicOp, _, _, err := cxprogram.GetOperation(exprs[i])
			if err != nil {
				return &CXAtomicOperator{}, err
			}

			return cxAtomicOp, nil
		}
	}
	return &CXAtomicOperator{}, fmt.Errorf("CXAtomicOperator not found.")
}

func (cxprogram *CXProgram) GetPreviousCXAtomicOpFromExpressions(exprs []*CXExpression, currIndex int) (*CXAtomicOperator, error) {
	for i := currIndex; i >= 0; i-- {
		if exprs[i].Type == CX_ATOMIC_OPERATOR {
			cxAtomicOp, _, _, err := cxprogram.GetOperation(exprs[i])
			if err != nil {
				return &CXAtomicOperator{}, err
			}

			return cxAtomicOp, nil
		}
	}
	return &CXAtomicOperator{}, fmt.Errorf("CXAtomicOperator not found.")
}

// Only two users, both in cx/execute.go
func (cxprogram *CXProgram) SelectPackage(name string) (*CXPackage, error) {
	if _, ok := cxprogram.Packages[name]; !ok {
		return nil, fmt.Errorf("Package '%s' does not exist", name)
	}

	cxprogram.CurrentPackage = cxprogram.Packages[name]

	return cxprogram.GetPackageFromArray(cxprogram.Packages[name])
}

// GetCurrentPackage ...
func (cxprogram *CXProgram) GetCurrentPackage() (*CXPackage, error) {
	// if ok:=cxprogram.CurrentPackage == nil {
	// 	return nil, errors.New("current package is nil")
	// }

	return cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
}

func (cxprogram *CXProgram) GetCurrentPackageIndex() CXPackageIndex {
	// if ok:=cxprogram.CurrentPackage == nil {
	// 	return nil, errors.New("current package is nil")
	// }

	return cxprogram.CurrentPackage
}

// GetCurrentStruct ...
func (cxprogram *CXProgram) GetCurrentStruct() (*CXStruct, error) {
	// if cxprogram.CurrentPackage == nil {
	// 	return nil, errors.New("current package is nil")
	// }

	currentPackage, err := cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
	if err != nil {
		return &CXStruct{}, err
	}

	if currentPackage.CurrentStruct == nil {
		return nil, errors.New("current struct is nil")

	}

	return currentPackage.CurrentStruct, nil
}

// GetCurrentFunction ...
func (cxprogram *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	// if cxprogram.CurrentPackage == nil {
	// 	return nil, errors.New("current package is nil")
	// }

	currentPackage, err := cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
	if err != nil {
		return &CXFunction{}, err
	}

	if currentPackage.CurrentFunction != nil {
		return nil, errors.New("current function is nil")
	}

	return currentPackage.CurrentFunction, nil

}

// GetCurrentExpression ...
func (cxprogram *CXProgram) GetCurrentExpression() (*CXExpression, error) {
	// if cxprogram.CurrentPackage == nil {
	// 	return nil, errors.New("current package is nil")
	// }

	currentPackage, err := cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
	if err != nil {
		return &CXExpression{}, err
	}

	if currentPackage.CurrentFunction == nil {
		return nil, errors.New("current function is nil")
	}

	if currentPackage.CurrentFunction.CurrentExpression == nil {
		return nil, errors.New("current expression is nil")
	}

	return currentPackage.CurrentFunction.CurrentExpression, nil
}

// GetPackage ...
func (cxprogram *CXProgram) GetPackage(pkgName string) (*CXPackage, error) {
	if _, ok := cxprogram.Packages[pkgName]; !ok {
		return nil, fmt.Errorf("package '%s' not found", pkgName)
	}

	return cxprogram.GetPackageFromArray(cxprogram.Packages[pkgName])
}

// GetStruct ...
func (cxprogram *CXProgram) GetStruct(strctName string, pkgName string) (*CXStruct, error) {
	pkg, err := cxprogram.GetPackage(pkgName)
	if err != nil {
		return nil, err
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		return nil, err
	}

	return strct, nil
}

// GetFunction ...
func (cxprogram *CXProgram) GetFunction(fnName string, pkgName string) (*CXFunction, error) {
	pkg, err := cxprogram.GetPackage(pkgName)
	if err != nil {
		return nil, err
	}

	fn, err := pkg.GetFunction(fnName)
	if err != nil {
		return nil, err
	}

	return fn, nil
}

// GetCurrentCall returns the current CXCall
//TODO: What does this do?
//TODO: Only used in OP_JMP
func (cxprogram *CXProgram) GetCurrentCall() *CXCall {
	return &cxprogram.CallStack[cxprogram.CallCounter]
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

// ----------------------------------------------------------------
//                             `CXProgram` Debugging

// PrintAllObjects prints all objects in a program
//
func (cxprogram *CXProgram) PrintAllObjects() {
	fp := types.Pointer(0)

	for c := types.Pointer(0); c <= cxprogram.CallCounter; c++ {
		op := cxprogram.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			heapOffset := types.Read_ptr(cxprogram.Memory, fp+ptr.Offset)

			var byts []byte

			if ptr.StructType != nil {
				// then it's a pointer to a struct
				// use CustomStruct to match the fields against the bytes
				// for _, fld := range ptr.Fields {

				// }

				byts = types.Get_obj_data(cxprogram.Memory, heapOffset, ptr.StructType.Size)
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
			// 		currCustom = ptr.StructType
			// 	case DECL_BASIC:
			// 	}
			// }

			// if len(ptr.Lengths) > 0 {
			// 	fmt.Println("ARRAY")
			// }

			// if ptr.StructType != nil {
			// 	fmt.Println("STRUCT")
			// }

			fmt.Println("declarat", ptr.DeclarationSpecifiers)

			fmt.Println("obj", ptr.Name, ptr.StructType, cxprogram.Memory[heapOffset:heapOffset+op.Size], byts)
		}

		fp += op.Size
	}
}

// PrintProgram prints the abstract syntax tree of a CX program in a
// human-readable format.
func (cxprogram *CXProgram) PrintProgram() {
	fmt.Println(ToString(cxprogram))
}
