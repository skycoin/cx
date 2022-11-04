package ast

import (
	"errors"
	"fmt"

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
	ProgramInput []CXArgumentIndex // OS input arguments
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
	CXAtomicOps             []CXAtomicOperator
	CXArgs                  []CXArgument
	CXLines                 []CXLine
	CXPackages              []CXPackage
	CXFunctions             []CXFunction
	CXStructs               []CXStruct
	CXTypeSignatures        []CXTypeSignature
	TypeSignatureForArrays  []CXTypeSignature_Array
	TypeSignatureForStructs []CXTypeSignature_Struct
	// Then reference the package of function by CxFunction id

	// For Initializers
	SysInitExprs []CXExpression
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

// ----------------------------------------------------------------
//                         `CXProgram` Package handling

// AddPackage ...
func (cxprogram *CXProgram) AddPackage(mod *CXPackage) CXPackageIndex {
	if idx, ok := cxprogram.Packages[mod.Name]; ok {
		cxprogram.CurrentPackage = idx
		return idx
	}

	index := cxprogram.AddPackageInArray(mod)

	cxprogram.Packages[mod.Name] = CXPackageIndex(index)
	cxprogram.CurrentPackage = CXPackageIndex(index)

	return index
}

func (cxprogram *CXProgram) AddPackageInArray(pkg *CXPackage) CXPackageIndex {
	// The index of pkg after it will be added in the array
	pkg.Index = len(cxprogram.CXPackages)
	cxprogram.CXPackages = append(cxprogram.CXPackages, *pkg)

	return CXPackageIndex(pkg.Index)
}

func (cxprogram *CXProgram) GetPackageFromArray(index CXPackageIndex) (*CXPackage, error) {
	if int(index) > (len(cxprogram.CXPackages) - 1) {
		return nil, fmt.Errorf("error: CXPackages[%d]: index out of bounds", index)
	}

	return &cxprogram.CXPackages[index], nil
}

// RemovePackage ...
// func (cxprogram *CXProgram) RemovePackage(modName string) {
// 	// If doesnt exist, return
// 	// if cxprogram.Packages[modName] == nil {
// 	// 	return
// 	// }

// 	// Check if it is the current pkg so when it
// 	// is deleted, it will be replaced with new pkg
// 	isCurrentPkg := cxprogram.Packages[modName] == cxprogram.CurrentPackage

// 	// Delete package
// 	delete(cxprogram.Packages, modName)

// 	// This means that we're removing the package set to be the CurrentPackage.
// 	// If it is removed from the program's map of packages, cxprogram.CurrentPackage
// 	// would be pointing to a package meant to be collected by the GC.
// 	// We fix this by pointing to random package in the program's map of packages.
// 	if isCurrentPkg {
// 		for _, pkg := range cxprogram.Packages {
// 			cxprogram.CurrentPackage = pkg
// 			break
// 		}
// 	}

// }

// ----------------------------------------------------------------
//                         `CXProgram` Function handling
func (cxprogram *CXProgram) AddFunctionInArray(fn *CXFunction) CXFunctionIndex {
	// The index of fn after it will be added in the array
	fn.Index = len(cxprogram.CXFunctions)

	cxprogram.CXFunctions = append(cxprogram.CXFunctions, *fn)

	return CXFunctionIndex(fn.Index)
}

func (cxprogram *CXProgram) AddNativeFunctionInArray(fn *CXNativeFunction) CXFunctionIndex {
	fnNative := &CXFunction{
		Index:          len(cxprogram.CXFunctions),
		AtomicOPCode:   fn.AtomicOPCode,
		Inputs:         &CXStruct{},
		Outputs:        &CXStruct{},
		LocalVariables: []string{},
	}

	// Add inputs to cx arg array
	for _, argIn := range fn.Inputs {
		err := fnNative.AddLocalVariableName(argIn.Name)
		if err != nil {
			// TODO: improve error handling
			panic("error adding local variable name")
		}

		newField := GetCXTypeSignatureRepresentationOfCXArg(cxprogram, argIn)

		newFieldIdx := cxprogram.AddCXTypeSignatureInArray(newField)
		fnNative.Inputs.Fields = append(fnNative.Inputs.Fields, newFieldIdx)
	}

	// Add outputs to cx arg array
	for _, argOut := range fn.Outputs {
		err := fnNative.AddLocalVariableName(argOut.Name)
		if err != nil {
			// TODO: improve error handling
			panic("error adding local variable name")
		}

		newField := GetCXTypeSignatureRepresentationOfCXArg(cxprogram, argOut)
		newFieldIdx := cxprogram.AddCXTypeSignatureInArray(newField)
		fnNative.Outputs.Fields = append(fnNative.Outputs.Fields, newFieldIdx)
	}

	cxprogram.CXFunctions = append(cxprogram.CXFunctions, *fnNative)

	return CXFunctionIndex(fnNative.Index)
}

func (cxprogram *CXProgram) GetFunctionFromArray(index CXFunctionIndex) *CXFunction {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.CXFunctions) - 1) {
		panic(fmt.Errorf("error: CXFunctions[%d]: index out of bounds", index))
	}

	return &cxprogram.CXFunctions[index]
}

// ----------------------------------------------------------------
//                         `CXProgram` Structs handling
func (cxprogram *CXProgram) AddStructInArray(strct *CXStruct) CXStructIndex {
	// The index of fn after it will be added in the array
	strct.Index = len(cxprogram.CXStructs)

	cxprogram.CXStructs = append(cxprogram.CXStructs, *strct)

	return CXStructIndex(strct.Index)
}

func (cxprogram *CXProgram) GetStructFromArray(index CXStructIndex) *CXStruct {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.CXStructs) - 1) {
		panic(fmt.Errorf("error: CXStructs[%d]: index out of bounds", index))
	}

	return &cxprogram.CXStructs[index]
}

// ----------------------------------------------------------------
//            `CXProgram` CXLines, CXArgs, and CXAtomicOps Handling

func (cxprogram *CXProgram) GetCXLine(index int) (*CXLine, error) {
	if index > (len(cxprogram.CXLines) - 1) {
		return nil, fmt.Errorf("error: CXLines[%d]: index out of bounds", index)
	}

	return &cxprogram.CXLines[index], nil
}

func (cxprogram *CXProgram) GetCXArg(index CXArgumentIndex) *CXArgument {
	if int(index) > (len(cxprogram.CXArgs) - 1) {
		panic(fmt.Errorf("error: CXArgs[%d]: index out of bounds", index))
	}

	return &cxprogram.CXArgs[index]
}

func (cxprogram *CXProgram) GetCXAtomicOp(index int) (*CXAtomicOperator, error) {
	if index > (len(cxprogram.CXAtomicOps) - 1) {
		return nil, fmt.Errorf("error: CXAtomicOps[%d]: index out of bounds", index)
	}

	return &cxprogram.CXAtomicOps[index], nil
}

func (cxprogram *CXProgram) AddCXLine(CXLine *CXLine) int {
	cxprogram.CXLines = append(cxprogram.CXLines, *CXLine)

	return len(cxprogram.CXLines) - 1
}

// func (cxprogram *CXProgram) AddCXArg(CXArg *CXArgument) int {
// 	cxprogram.CXArgs = append(cxprogram.CXArgs, *CXArg)

// 	return len(cxprogram.CXArgs) - 1
// }

func (cxprogram *CXProgram) AddCXAtomicOp(CXAtomicOp *CXAtomicOperator) int {
	cxprogram.CXAtomicOps = append(cxprogram.CXAtomicOps, *CXAtomicOp)

	return len(cxprogram.CXAtomicOps) - 1
}

// ----------------------------------------------------------------
//                         `CXProgram` CXTypeSignatures handling
func (cxprogram *CXProgram) AddCXTypeSignatureInArray(typeSignature *CXTypeSignature) CXTypeSignatureIndex {
	typeSignature.Index = CXTypeSignatureIndex(len(cxprogram.CXTypeSignatures))
	cxprogram.CXTypeSignatures = append(cxprogram.CXTypeSignatures, *typeSignature)

	return typeSignature.Index
}

func (cxprogram *CXProgram) GetCXTypeSignatureFromArray(index CXTypeSignatureIndex) *CXTypeSignature {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.CXTypeSignatures) - 1) {
		panic(fmt.Errorf("error: CXTypeSignature[%d]: index out of bounds", index))
	}

	return &cxprogram.CXTypeSignatures[index]
}

// ----------------------------------------------------------------
//                         `CXProgram` TypeSignature_Array handling
func (cxprogram *CXProgram) AddCXTypeSignatureArrayInArray(typeSignatureArray *CXTypeSignature_Array) int {
	cxprogram.TypeSignatureForArrays = append(cxprogram.TypeSignatureForArrays, *typeSignatureArray)

	return len(cxprogram.TypeSignatureForArrays) - 1
}

func (cxprogram *CXProgram) GetCXTypeSignatureArrayFromArray(index int) *CXTypeSignature_Array {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.TypeSignatureForArrays) - 1) {
		panic(fmt.Errorf("error: TypeSignatureForArrays[%d]: index out of bounds", index))
	}

	return &cxprogram.TypeSignatureForArrays[index]
}

// ----------------------------------------------------------------
//                         `CXProgram` TypeSignature_Struct handling
func (cxprogram *CXProgram) AddCXTypeSignatureStructInArray(typeSignatureStruct *CXTypeSignature_Struct) int {
	cxprogram.TypeSignatureForStructs = append(cxprogram.TypeSignatureForStructs, *typeSignatureStruct)

	return len(cxprogram.TypeSignatureForStructs) - 1
}

func (cxprogram *CXProgram) GetCXTypeSignatureStructFromArray(index int) *CXTypeSignature_Struct {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.TypeSignatureForStructs) - 1) {
		panic(fmt.Errorf("error: TypeSignatureForArrays[%d]: index out of bounds", index))
	}

	return &cxprogram.TypeSignatureForStructs[index]
}

// ----------------------------------------------------------------
//                         `CXProgram` CXArgument handling
func (cxprogram *CXProgram) AddCXArgInArray(cxArg *CXArgument) CXArgumentIndex {
	// The index of fn after it will be added in the array
	cxArg.Index = len(cxprogram.CXArgs)
	cxprogram.CXArgs = append(cxprogram.CXArgs, *cxArg)

	return CXArgumentIndex(cxArg.Index)
}

func (cxprogram *CXProgram) GetCXArgFromArray(index CXArgumentIndex) *CXArgument {
	if index == -1 {
		return nil
	}

	if int(index) > (len(cxprogram.CXArgs) - 1) {
		panic(fmt.Errorf("error: CXArgument[%d]: index out of bounds", index))
	}

	return &cxprogram.CXArgs[index]
}

func (cxprogram *CXProgram) ConvertIndexArgsToPointerArgs(idxs []CXArgumentIndex) []*CXArgument {
	var cxArgs []*CXArgument
	for _, idx := range idxs {
		arg := cxprogram.GetCXArgFromArray(idx)
		cxArgs = append(cxArgs, arg)
	}
	return cxArgs
}

func (cxprogram *CXProgram) AddPointerArgsToCXArgsArray(cxArgs []*CXArgument) []CXArgumentIndex {
	var cxArgsIdxs []CXArgumentIndex
	for _, cxArg := range cxArgs {
		cxArgIdx := cxprogram.AddCXArgInArray(cxArg)
		cxArgsIdxs = append(cxArgsIdxs, cxArgIdx)
	}
	return cxArgsIdxs
}

// Temporary only for intitial new struct def implementation.
func (cxprogram *CXProgram) ConvertIndexTypeSignaturesToPointerArgs(idxs []CXTypeSignatureIndex) []*CXArgument {
	var cxArgs []*CXArgument
	for _, typeSignatureIdx := range idxs {
		typeSignature := cxprogram.GetCXTypeSignatureFromArray(typeSignatureIdx)
		if typeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
			idx := typeSignature.Meta
			arg := cxprogram.GetCXArgFromArray(CXArgumentIndex(idx))
			cxArgs = append(cxArgs, arg)
		}
	}
	return cxArgs
}

// Temporary only for intitial new struct def implementation.
func (cxprogram *CXProgram) AddPointerArgsToTypeSignaturesArray(cxArgs []*CXArgument) []CXTypeSignatureIndex {
	var cxTypeSignaturesIdxs []CXTypeSignatureIndex
	for _, cxArg := range cxArgs {
		cxArgIdx := cxprogram.AddCXArgInArray(cxArg)

		newCXTypeSignature := &CXTypeSignature{
			Name:       cxArg.Name,
			Package:    cxArg.Package,
			Offset:     cxArg.Offset,
			Type:       TYPE_CXARGUMENT_DEPRECATE,
			Meta:       int(cxArgIdx),
			ArgDetails: cxArg.ArgDetails,
		}

		newCXTypeSignatureIdx := cxprogram.AddCXTypeSignatureInArray(newCXTypeSignature)
		cxTypeSignaturesIdxs = append(cxTypeSignaturesIdxs, newCXTypeSignatureIdx)
	}
	return cxTypeSignaturesIdxs
}

// ----------------------------------------------------------------
//                             `CXProgram` Getters

func (cxprogram *CXProgram) GetOperation(expr *CXExpression) (*CXAtomicOperator, *CXArgument, *CXLine, error) {
	switch expr.Type {
	case CX_ATOMIC_OPERATOR:
		cxAtomicOp, err := cxprogram.GetCXAtomicOp(expr.Index)
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

func (cxprogram *CXProgram) GetPreviousCXLine(exprs []CXExpression, currIndex int) (*CXLine, error) {
	for i := currIndex; i >= 0; i-- {
		if exprs[i].Type == CX_LINE {
			_, _, cxLine, err := cxprogram.GetOperation(&exprs[i])
			if err != nil {
				return &CXLine{}, err
			}

			return cxLine, nil
		}
	}
	return &CXLine{}, fmt.Errorf("CXLine not found.")
}

func (cxprogram *CXProgram) GetCXAtomicOpFromExpressions(exprs []CXExpression, currIndex int) (*CXAtomicOperator, error) {
	for i := currIndex; i < len(exprs); i++ {
		if exprs[i].Type == CX_ATOMIC_OPERATOR {
			cxAtomicOp, err := cxprogram.GetCXAtomicOp(exprs[i].Index)
			if err != nil {
				return &CXAtomicOperator{}, err
			}

			return cxAtomicOp, nil
		}
	}
	return &CXAtomicOperator{}, fmt.Errorf("CXAtomicOperator not found.")
}

func (cxprogram *CXProgram) GetPreviousCXAtomicOpFromExpressions(exprs []CXExpression, currIndex int) (*CXAtomicOperator, error) {
	for i := currIndex; i >= 0; i-- {
		if exprs[i].Type == CX_ATOMIC_OPERATOR {
			cxAtomicOp, err := cxprogram.GetCXAtomicOp(exprs[i].Index)
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
	if cxprogram.CurrentPackage == -1 {
		return nil, errors.New("current package is nil")
	}

	return cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
}

// GetCurrentFunction ...
func (cxprogram *CXProgram) GetCurrentFunction() (*CXFunction, error) {
	if cxprogram.CurrentPackage == -1 {
		return nil, errors.New("current package is nil")
	}

	currentPackage, err := cxprogram.GetPackageFromArray(cxprogram.CurrentPackage)
	if err != nil {
		return &CXFunction{}, err
	}

	if currentPackage.CurrentFunction == -1 {
		return nil, errors.New("current function is nil")
	}

	return cxprogram.GetFunctionFromArray(currentPackage.CurrentFunction), nil

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

	strct, err := pkg.GetStruct(cxprogram, strctName)
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

	fn, err := pkg.GetFunction(cxprogram, fnName)
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

				byts = types.Get_obj_data(cxprogram.Memory, heapOffset, ptr.StructType.GetStructSize(cxprogram))
			}

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
