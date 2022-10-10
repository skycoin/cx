package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func deserializeRaw(byts []byte, offset types.Pointer, size types.Pointer, item interface{}) {
	_, err := encoder.DeserializeRaw(byts[offset:offset+size], item)
	if err != nil {
		panic(err)
	}
}

func serializeString(name string, s *SerializedCXProgram) (int64, int64) {
	if name == "" {
		return int64(-1), int64(-1)
	}

	size := encoder.Size(name)

	off, found := s.StringsMap[name]
	if found {
		return int64(off), int64(size)
	}
	off = int64(len(s.Strings))
	s.Strings = append(s.Strings, encoder.Serialize(name)...)
	s.StringsMap[name] = off

	return int64(off), int64(size)
}

func indexPackage(pkg *CXPackage, s *SerializedCXProgram) {
	if _, found := s.PackagesMap[pkg.Name]; !found {
		s.PackagesMap[pkg.Name] = int64(len(s.PackagesMap))
	} else {
		panic("duplicated package in serialization process")
	}
}

func indexStruct(prgrm *CXProgram, strct *CXStruct, s *SerializedCXProgram) {
	strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
	if err != nil {
		panic(err)
	}

	strctName := strctPkg.Name + "." + strct.Name
	if _, found := s.StructsMap[strctName]; !found {
		s.StructsMap[strctName] = int64(len(s.StructsMap))
	} else {
		panic("duplicated struct in serialization process")
	}
}

func indexFunction(prgrm *CXProgram, fn *CXFunction, s *SerializedCXProgram) {
	fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
	if err != nil {
		panic(err)
	}
	fnName := fnPkg.Name + "." + fn.Name
	if _, found := s.FunctionsMap[fnName]; !found {
		s.FunctionsMap[fnName] = int64(len(s.FunctionsMap))
	} else {
		panic("duplicated function in serialization process")
	}
}

func serializeBoolean(val bool) int64 {
	if val {
		return 1
	}
	return 0
}

func serializePointers(pointers []types.Pointer, s *SerializedCXProgram) (int64, int64) {
	if len(pointers) == 0 {
		return int64(-1), int64(-1)
	}
	off := len(s.Integers)
	l := len(pointers)

	ints := make([]int64, l)
	for i, pointer := range pointers {
		ints[i] = int64(pointer)
	}

	s.Integers = append(s.Integers, ints...)

	return int64(off), int64(l)
}

func serializeIntegers(ints []int, s *SerializedCXProgram) (int64, int64) {
	if len(ints) == 0 {
		return int64(-1), int64(-1)
	}
	off := len(s.Integers)
	l := len(ints)

	ints32 := make([]int64, l)
	for i, int := range ints {
		ints32[i] = int64(int)
	}

	s.Integers = append(s.Integers, ints32...)

	return int64(off), int64(l)
}

func serializeArgument(prgrm *CXProgram, arg *CXArgument, s *SerializedCXProgram) int {
	s.Arguments = append(s.Arguments, serializedArgument{})
	argOff := len(s.Arguments) - 1

	sNil := int64(-1)

	s.Arguments[argOff].NameOffset, s.Arguments[argOff].NameSize = serializeString(arg.Name, s)

	s.Arguments[argOff].Type = int64(arg.Type)

	if arg.StructType == nil {
		s.Arguments[argOff].StructTypeOffset = sNil
	} else {
		strctPkg, err := prgrm.GetPackageFromArray(arg.StructType.Package)
		if err != nil {
			panic(err)
		}

		strctName := strctPkg.Name + "." + arg.StructType.Name
		if strctOff, found := s.StructsMap[strctName]; found {
			s.Arguments[argOff].StructTypeOffset = int64(strctOff)
		} else {
			panic("struct reference not found")
		}
	}

	s.Arguments[argOff].Size = int64(arg.Size)
	s.Arguments[argOff].Offset = int64(arg.Offset)

	s.Arguments[argOff].DeclarationSpecifiersOffset,
		s.Arguments[argOff].DeclarationSpecifiersSize = serializeIntegers(arg.DeclarationSpecifiers, s)

	s.Arguments[argOff].IsSlice = serializeBoolean(arg.IsSlice)
	s.Arguments[argOff].PreviouslyDeclared = serializeBoolean(arg.PreviouslyDeclared)

	s.Arguments[argOff].PassBy = int64(arg.PassBy)

	s.Arguments[argOff].LengthsOffset, s.Arguments[argOff].LengthsSize = serializePointers(arg.Lengths, s)
	// TODO: include Indexes of type CTypeSignature
	// s.Arguments[argOff].IndexesOffset, s.Arguments[argOff].IndexesSize = serializeSliceOfArguments(prgrm, prgrm.ConvertIndexArgsToPointerArgs(arg.Indexes), s)
	s.Arguments[argOff].FieldsOffset, s.Arguments[argOff].FieldsSize = serializeSliceOfArguments(prgrm, prgrm.ConvertIndexArgsToPointerArgs(arg.Fields), s)
	s.Arguments[argOff].InputsOffset, s.Arguments[argOff].InputsSize = serializeSliceOfArguments(prgrm, prgrm.ConvertIndexArgsToPointerArgs(arg.Inputs), s)
	s.Arguments[argOff].OutputsOffset, s.Arguments[argOff].OutputsSize = serializeSliceOfArguments(prgrm, prgrm.ConvertIndexArgsToPointerArgs(arg.Outputs), s)

	argPkg, err := prgrm.GetPackageFromArray(arg.Package)
	if err != nil {
		panic(err)
	}

	if _, found := s.PackagesMap[argPkg.Name]; found {
		s.Arguments[argOff].PackageName = argPkg.Name
	} else {
		panic("package reference not found")
	}

	return argOff
}

func serializeSliceOfArguments(prgrm *CXProgram, args []*CXArgument, s *SerializedCXProgram) (int64, int64) {
	if len(args) == 0 {
		return int64(-1), int64(-1)
	}
	idxs := make([]int, len(args))
	for i, arg := range args {
		idxs[i] = serializeArgument(prgrm, arg, s)
	}
	return serializeIntegers(idxs, s)
}

func serializeCalls(prgrm *CXProgram, calls []CXCall, s *SerializedCXProgram) (int64, int64) {
	if len(calls) == 0 {
		return int64(-1), int64(-1)
	}
	idxs := make([]int, len(calls))
	for i, call := range calls {
		idxs[i] = serializeCall(prgrm, &call, s)
	}
	return serializeIntegers(idxs, s)

}

func serializeExpression(prgrm *CXProgram, expr *CXExpression, s *SerializedCXProgram) int {
	s.Expressions = append(s.Expressions, serializedExpression{})
	exprOff := len(s.Expressions) - 1
	sExpr := &s.Expressions[exprOff]

	sNil := int64(-1)

	sExpr.Type = int64(expr.Type)
	switch expr.Type {
	case CX_LINE:
		_, _, cxLine, err := prgrm.GetOperation(expr)
		if err != nil {
			panic(err)
		}

		sExpr.ExpressionType = int64(expr.ExpressionType)
		sExpr.FileName = cxLine.FileName
		sExpr.LineNumber = int64(cxLine.LineNumber)
		sExpr.LineStr = cxLine.LineStr

	case CX_ATOMIC_OPERATOR:
		cxAtomicOp, err := prgrm.GetCXAtomicOp(expr.Index)
		if err != nil {
			panic(err)
		}

		cxAtomicOpOperator := prgrm.GetFunctionFromArray(cxAtomicOp.Operator)

		if cxAtomicOpOperator == nil {
			// then it's a declaration
			sExpr.OperatorOffset = sNil
			sExpr.IsNative = serializeBoolean(false)
			sExpr.OpCode = int64(-1)
		} else if cxAtomicOpOperator.IsBuiltIn() {
			sExpr.OperatorOffset = sNil
			sExpr.IsNative = serializeBoolean(true)
			sExpr.OpCode = int64(cxAtomicOpOperator.AtomicOPCode)
		} else {
			sExpr.IsNative = serializeBoolean(false)
			sExpr.OpCode = sNil

			opPkg, err := prgrm.GetPackageFromArray(cxAtomicOpOperator.Package)
			if err != nil {
				panic(err)
			}
			opName := opPkg.Name + "." + cxAtomicOpOperator.Name
			if opOff, found := s.FunctionsMap[opName]; found {
				sExpr.OperatorOffset = int64(opOff)
			}
		}

		inputCXArgs := prgrm.ConvertIndexTypeSignaturesToPointerArgs(cxAtomicOp.GetInputs(prgrm))
		sExpr.InputsOffset, sExpr.InputsSize = serializeSliceOfArguments(prgrm, inputCXArgs, s)
		outputCXArgs := prgrm.ConvertIndexTypeSignaturesToPointerArgs(cxAtomicOp.GetOutputs(prgrm))
		sExpr.OutputsOffset, sExpr.OutputsSize = serializeSliceOfArguments(prgrm, outputCXArgs, s)

		sExpr.LabelOffset, sExpr.LabelSize = serializeString(cxAtomicOp.Label, s)
		sExpr.ThenLines = int64(cxAtomicOp.ThenLines)
		sExpr.ElseLines = int64(cxAtomicOp.ElseLines)

		sExpr.ExpressionType = int64(expr.ExpressionType)

		cxAtomicOpFunction := prgrm.GetFunctionFromArray(cxAtomicOp.Function)

		fnPkg, err := prgrm.GetPackageFromArray(cxAtomicOpFunction.Package)
		if err != nil {
			panic(err)
		}
		fnName := fnPkg.Name + "." + cxAtomicOpFunction.Name
		if fnOff, found := s.FunctionsMap[fnName]; found {
			sExpr.FunctionOffset = int64(fnOff)
		} else {
			panic("function reference not found")
		}

		cxAtomicOpPkg, err := prgrm.GetPackageFromArray(cxAtomicOp.Package)
		if err != nil {
			panic(err)
		}

		if _, found := s.PackagesMap[cxAtomicOpPkg.Name]; found {
			sExpr.PackageName = cxAtomicOpPkg.Name
		} else {
			panic("package reference not found")
		}

	}
	return exprOff
}

func serializeCall(prgrm *CXProgram, call *CXCall, s *SerializedCXProgram) int {
	s.Calls = append(s.Calls, serializedCall{})
	callOff := len(s.Calls) - 1
	serializedCall := &s.Calls[callOff]

	opPkg, err := prgrm.GetPackageFromArray(call.Operator.Package)
	if err != nil {
		panic(err)
	}

	opName := opPkg.Name + "." + call.Operator.Name
	if opOff, found := s.FunctionsMap[opName]; found {
		serializedCall.OperatorOffset = int64(opOff)
		serializedCall.Line = int64(call.Line)
		serializedCall.FramePointer = int64(call.FramePointer)
	} else {
		panic("function reference not found")
	}

	return callOff
}

func serializeStructArguments(prgrm *CXProgram, strct *CXStruct, s *SerializedCXProgram) {
	strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
	if err != nil {
		panic(err)
	}

	strctName := strctPkg.Name + "." + strct.Name
	if strctOff, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[strctOff]
		sStrct.FieldsOffset, sStrct.FieldsSize = serializeSliceOfArguments(prgrm, prgrm.ConvertIndexTypeSignaturesToPointerArgs(strct.Fields), s)
	} else {
		panic("struct reference not found")
	}
}

func serializeFunctionArguments(prgrm *CXProgram, fn *CXFunction, s *SerializedCXProgram) {
	fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
	if err != nil {
		panic(err)
	}
	fnName := fnPkg.Name + "." + fn.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[fnOff]

		arrCXArgs := prgrm.ConvertIndexTypeSignaturesToPointerArgs(fn.GetInputs(prgrm))
		sFn.InputsOffset, sFn.InputsSize = serializeSliceOfArguments(prgrm, arrCXArgs, s)
		outputCXArgs := prgrm.ConvertIndexTypeSignaturesToPointerArgs(fn.GetOutputs(prgrm))
		sFn.OutputsOffset, sFn.OutputsSize = serializeSliceOfArguments(prgrm, outputCXArgs, s)
		sFn.ListOfPointersOffset, sFn.ListOfPointersSize = serializeSliceOfArguments(prgrm, fn.ListOfPointers, s)
	} else {
		panic("function reference not found")
	}
}

func serializePackageName(pkg *CXPackage, s *SerializedCXProgram) {
	sPkg := &s.Packages[s.PackagesMap[pkg.Name]]
	sPkg.NameOffset, sPkg.NameSize = serializeString(pkg.Name, s) //Change Name to String
}

func serializeStructName(prgrm *CXProgram, strct *CXStruct, s *SerializedCXProgram) {
	strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
	if err != nil {
		panic(err)
	}

	strctName := strctPkg.Name + "." + strct.Name
	sStrct := &s.Structs[s.StructsMap[strctName]]
	sStrct.NameOffset, sStrct.NameSize = serializeString(strct.Name, s) //Change Name to String
}

func serializeFunctionName(prgrm *CXProgram, fn *CXFunction, s *SerializedCXProgram) {
	fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
	if err != nil {
		panic(err)
	}
	fnName := fnPkg.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.NameOffset, sFn.NameSize = serializeString(fn.Name, s) //Change Name to String
	} else {
		panic("function reference not found")
	}
}

func serializePackageGlobals(prgrm *CXProgram, pkg *CXPackage, s *SerializedCXProgram) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]

		var glblArgs []*CXArgument
		for _, glblFldIdx := range pkg.Globals.Fields {
			glblFld := prgrm.GetCXTypeSignatureFromArray(glblFldIdx)
			// Assuming only all are TYPE_CXARGUMENT_DEPRECATE
			// TODO: To be replaced
			glbl := prgrm.GetCXArg(CXArgumentIndex(glblFld.Meta))
			glblArgs = append(glblArgs, glbl)
		}
		sPkg.GlobalsOffset, sPkg.GlobalsSize = serializeSliceOfArguments(prgrm, glblArgs, s)
	} else {
		panic("package reference not found")
	}
}

func serializePackageImports(prgrm *CXProgram, pkg *CXPackage, s *SerializedCXProgram) {
	l := len(pkg.Imports)
	if l == 0 {
		s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int64(-1)
		s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int64(-1)
		return
	}
	imps := make([]int64, l)
	count := 0
	for _, impIdx := range pkg.Imports {
		impPkg, err := prgrm.GetPackageFromArray(impIdx)
		if err != nil {
			panic(err)
		}
		if idx, found := s.PackagesMap[impPkg.Name]; found {
			imps[count] = int64(idx)
		} else {
			panic("import package reference not found")
		}

		count++
	}

	s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int64(len(s.Integers))
	s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int64(l)
	s.Integers = append(s.Integers, imps...)
}

func serializeStructPackage(prgrm *CXProgram, strct *CXStruct, s *SerializedCXProgram) {
	strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
	if err != nil {
		panic(err)
	}

	strctName := strctPkg.Name + "." + strct.Name
	if _, found := s.PackagesMap[strctPkg.Name]; found {
		if off, found := s.StructsMap[strctName]; found {
			sStrct := &s.Structs[off]
			sStrct.PackageName = strctPkg.Name
		} else {
			panic("struct reference not found")
		}
	} else {
		panic("struct's package reference not found")
	}
}

func serializeFunctionPackage(prgrm *CXProgram, fn *CXFunction, s *SerializedCXProgram) {
	fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
	if err != nil {
		panic(err)
	}
	fnName := fnPkg.Name + "." + fn.Name
	if _, found := s.PackagesMap[fnPkg.Name]; found {
		if off, found := s.FunctionsMap[fnName]; found {
			sFn := &s.Functions[off]
			sFn.PackageName = fnPkg.Name
		} else {
			panic("function reference not found")
		}
	} else {
		panic("function's package reference not found")
	}
}

func serializePackageIntegers(prgrm *CXProgram, pkg *CXPackage, s *SerializedCXProgram) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]

		if pkg.CurrentFunction == -1 {
			// package has no functions
			sPkg.CurrentFunctionName = ""
		} else {
			currFn := prgrm.GetFunctionFromArray(pkg.CurrentFunction)

			currFnPkg, err := prgrm.GetPackageFromArray(currFn.Package)
			if err != nil {
				panic(err)
			}
			currFnName := currFnPkg.Name + "." + currFn.Name

			if _, found := s.FunctionsMap[currFnName]; found {
				sPkg.CurrentFunctionName = currFnName
			} else {
				panic("function reference not found")
			}
		}

	} else {
		panic("package reference not found")
	}
}

// func serializeStructIntegers(prgrm *CXProgram, strct *CXStruct, s *SerializedCXProgram) {
// 	strctPkg, err := prgrm.GetPackageFromArray(strct.Package)
// 	if err != nil {
// 		panic(err)
// 	}

// 	strctName := strctPkg.Name + "." + strct.Name
// 	if off, found := s.StructsMap[strctName]; found {
// 		sStrct := &s.Structs[off]
// 		sStrct.Size = int64(strct.Size)
// 	} else {
// 		panic("struct reference not found")
// 	}
// }

func serializeFunctionIntegers(prgrm *CXProgram, fn *CXFunction, s *SerializedCXProgram) {
	fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
	if err != nil {
		panic(err)
	}
	fnName := fnPkg.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.Size = int64(fn.Size)
		sFn.Length = int64(fn.LineCount)
	} else {
		panic("function reference not found")
	}
}

// initSerialization initializes the
// container for our serialized cx program.
// Program memory is also added here to our container
// if memory is to be included.
func initSerialization(prgrm *CXProgram, s *SerializedCXProgram, includeDataMemory, useCompression bool) {
	s.PackagesMap = make(map[string]int64)
	s.StructsMap = make(map[string]int64)
	s.FunctionsMap = make(map[string]int64)
	s.StringsMap = make(map[string]int64)

	s.Calls = make([]serializedCall, prgrm.CallCounter)
	s.Packages = make([]serializedPackage, len(prgrm.Packages))

	// If use compression, whole memory will be included
	// If not and if includeDataMemory, only data segment memory will be included
	if useCompression {
		s.Memory = prgrm.Memory
	} else if includeDataMemory && len(prgrm.Memory) != 0 {
		s.DataSegmentMemory = prgrm.Memory[prgrm.Data.StartsAt : prgrm.Data.StartsAt+prgrm.Data.Size]
	}

	var numStrcts int
	var numFns int

	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		numStrcts += len(pkg.Structs)
		numFns += len(pkg.Functions)
	}

	s.Structs = make([]serializedStruct, numStrcts)
	s.Functions = make([]serializedFunction, numFns)
	// args and exprs need to be appended as they are found
}

// serializeProgram serializes
// program of cx program.
func serializeProgram(prgrm *CXProgram, s *SerializedCXProgram) {
	s.Program = serializedProgram{}
	sPrgrm := &s.Program
	sPrgrm.PackagesOffset = int64(0)
	sPrgrm.PackagesSize = int64(len(prgrm.Packages))

	currPkg, err := prgrm.GetPackageFromArray(prgrm.CurrentPackage)
	if err != nil {
		panic(err)
	}
	if _, found := s.PackagesMap[currPkg.Name]; found {
		sPrgrm.CurrentPackageName = currPkg.Name
	} else {
		panic("package reference not found")
	}

	args := []*CXArgument{}
	for _, argIdx := range prgrm.ProgramInput {
		arg := prgrm.GetCXArg(argIdx)
		args = append(args, arg)
	}
	sPrgrm.InputsOffset, sPrgrm.InputsSize = serializeSliceOfArguments(prgrm, args, s)
	//sPrgrm.OutputsOffset, sPrgrm.OutputsSize = serializeSliceOfArguments(prgrm.ProgramOutput, s)

	sPrgrm.CallStackOffset, sPrgrm.CallStackSize = serializeCalls(prgrm, prgrm.CallStack[:prgrm.CallCounter], s)

	sPrgrm.CallCounter = int64(prgrm.CallCounter)

	sPrgrm.MemoryOffset = int64(0)
	sPrgrm.MemorySize = int64(len(prgrm.Memory))

	sPrgrm.HeapPointer = int64(prgrm.Heap.Pointer)
	sPrgrm.StackPointer = int64(prgrm.Stack.Pointer)
	sPrgrm.StackSize = int64(prgrm.Stack.Size)
	sPrgrm.DataSegmentSize = int64(prgrm.Data.Size)
	sPrgrm.DataSegmentStartsAt = int64(prgrm.Data.StartsAt)
	sPrgrm.HeapSize = int64(prgrm.Heap.Size)
	sPrgrm.HeapStartsAt = int64(prgrm.Heap.StartsAt)

	sPrgrm.Terminated = serializeBoolean(prgrm.Terminated)
	sPrgrm.VersionOffset, sPrgrm.VersionSize = serializeString(prgrm.Version, s)
}

// serializeCXProgramElements is used serializing CX program's
// elements (packages, structs, functions, etc.).
func serializeCXProgramElements(prgrm *CXProgram, s *SerializedCXProgram) {
	var fnCounter int64
	var strctCounter int64

	// indexing packages and serializing their names
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		indexPackage(pkg, s)
		serializePackageName(pkg, s)
	}
	// we first needed to populate references to all packages
	// now we add the imports' references
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		serializePackageImports(prgrm, pkg, s)
	}

	// structs
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		for _, strctIdx := range pkg.Structs {
			strct := &prgrm.CXStructs[strctIdx]
			indexStruct(prgrm, strct, s)
			serializeStructName(prgrm, strct, s)
			serializeStructPackage(prgrm, strct, s)
			// serializeStructIntegers(prgrm, strct, s)
		}
	}
	// we first needed to populate references to all structs
	// now we add fields
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		for _, strctIdx := range pkg.Structs {
			strct := &prgrm.CXStructs[strctIdx]
			serializeStructArguments(prgrm, strct, s)
		}
	}

	// globals
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		serializePackageGlobals(prgrm, pkg, s)
	}

	// functions
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		for _, fnIdx := range pkg.Functions {
			fn := prgrm.GetFunctionFromArray(fnIdx)

			indexFunction(prgrm, fn, s)
			serializeFunctionName(prgrm, fn, s)
			serializeFunctionPackage(prgrm, fn, s)
			serializeFunctionIntegers(prgrm, fn, s)
			serializeFunctionArguments(prgrm, fn, s)
		}
	}

	// package elements' offsets and sizes
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		if pkgOff, found := s.PackagesMap[pkg.Name]; found {
			sPkg := &s.Packages[pkgOff]

			if len(pkg.Structs) == 0 {
				sPkg.StructsOffset = int64(-1)
				sPkg.StructsSize = int64(-1)
			} else {
				sPkg.StructsOffset = strctCounter
				lenStrcts := int64(len(pkg.Structs))
				sPkg.StructsSize = lenStrcts
				strctCounter += lenStrcts
			}

			if len(pkg.Functions) == 0 {
				sPkg.FunctionsOffset = int64(-1)
				sPkg.FunctionsSize = int64(-1)
			} else {
				sPkg.FunctionsOffset = fnCounter
				lenFns := int64(len(pkg.Functions))
				sPkg.FunctionsSize = lenFns
				fnCounter += lenFns
			}
		} else {
			panic("package reference not found")
		}
	}

	// package integers
	// we needed the references to all functions and structs first
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		serializePackageIntegers(prgrm, pkg, s)
	}

	// expressions
	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		for _, fnIdx := range pkg.Functions {
			fn := prgrm.GetFunctionFromArray(fnIdx)

			fnPkg, err := prgrm.GetPackageFromArray(fn.Package)
			if err != nil {
				panic(err)
			}
			fnName := fnPkg.Name + "." + fn.Name
			if fnOff, found := s.FunctionsMap[fnName]; found {
				sFn := &s.Functions[fnOff]

				if len(fn.Expressions) == 0 {
					sFn.ExpressionsOffset = int64(-1)
					sFn.ExpressionsSize = int64(-1)
					sFn.CurrentExpressionOffset = int64(-1)
				} else {
					exprs := make([]int, len(fn.Expressions))
					for i, expr := range fn.Expressions {
						exprIdx := serializeExpression(prgrm, &expr, s)
						exprs[i] = exprIdx
					}

					sFn.ExpressionsOffset, sFn.ExpressionsSize = serializeIntegers(exprs, s)
				}
			} else {
				panic("function reference not found")
			}
		}
	}
}

// SerializeCXProgram translates cx program to slice of bytes that we can save.
// These slice of bytes can then be deserialize in the future and
// be translated back to cx program.
func SerializeCXProgram(prgrm *CXProgram, includeDataMemory, useCompression bool) (b []byte) {
	s := SerializedCXProgram{}
	initSerialization(prgrm, &s, includeDataMemory, useCompression)

	// serialize cx program's packages,
	// structs, functions, etc.
	serializeCXProgramElements(prgrm, &s)

	// serialize cx program's program
	serializeProgram(prgrm, &s)

	// serializing everything
	b = encoder.Serialize(s)

	if useCompression {
		// Compress using LZ4
		CompressBytesLZ4(&b)
	}
	return b
}

// SerializeDebugInfo prints the name of the serialized segment and byte size.
func SerializeDebugInfo(prgrm *CXProgram, includeMemory, useCompression bool) SerializedDataSize {
	idxSize := encoder.Size(serializedCXProgramIndex{})
	var s SerializedCXProgram

	bytes := SerializeCXProgram(prgrm, includeMemory, useCompression)
	deserializeRaw(bytes, 0, types.Cast_ui64_to_ptr(idxSize), &s.Index)

	data := &SerializedDataSize{
		Program:     len(bytes[s.Index.ProgramOffset:s.Index.CallsOffset]),
		Calls:       len(bytes[s.Index.CallsOffset:s.Index.PackagesOffset]),
		Packages:    len(bytes[s.Index.PackagesOffset:s.Index.StructsOffset]),
		Structs:     len(bytes[s.Index.StructsOffset:s.Index.FunctionsOffset]),
		Functions:   len(bytes[s.Index.FunctionsOffset:s.Index.ExpressionsOffset]),
		Expressions: len(bytes[s.Index.ExpressionsOffset:s.Index.ArgumentsOffset]),
		Arguments:   len(bytes[s.Index.ArgumentsOffset:s.Index.IntegersOffset]),
		Integers:    len(bytes[s.Index.IntegersOffset:s.Index.StringsOffset]),
		Strings:     len(bytes[s.Index.StringsOffset:s.Index.MemoryOffset]),
		Memory:      len(bytes[s.Index.MemoryOffset:]),
	}

	return *data
}

func deserializeString(off int64, size int64, s *SerializedCXProgram) string {
	if size < 1 {
		return ""
	}

	var name string
	deserializeRaw(s.Strings, types.Cast_i64_to_ptr(off), types.Cast_i64_to_ptr(size), &name)

	return name
}

func deserializePackages(s *SerializedCXProgram, prgrm *CXProgram) {
	var fnCounter int64
	var strctCounter int64

	for _, sPkg := range s.Packages {
		// initializing packages with their names,
		// empty functions, structs, imports and globals
		// and current function and struct
		pkg := &CXPackage{}
		pkg.Name = deserializeString(sPkg.NameOffset, sPkg.NameSize, s)
		pkgIdx := prgrm.AddPackage(pkg)
		pkg, _ = prgrm.GetPackageFromArray(pkgIdx)

		if sPkg.ImportsSize > 0 {
			pkg.Imports = make(map[string]CXPackageIndex, sPkg.ImportsSize)
		}

		if sPkg.FunctionsSize > 0 {
			pkg.Functions = make(map[string]CXFunctionIndex, sPkg.FunctionsSize)

			for _, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				var fn CXFunction
				fn.Name = deserializeString(sFn.NameOffset, sFn.NameSize, s)

				fnIdx := prgrm.AddFunctionInArray(&fn)
				pkg.Functions[fn.Name] = fnIdx
			}
		}

		if sPkg.StructsSize > 0 {
			pkg.Structs = make(map[string]CXStructIndex, sPkg.StructsSize)

			for _, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				var strct CXStruct
				strct.Name = deserializeString(sStrct.NameOffset, sStrct.NameSize, s)
				strctIdx := prgrm.AddStructInArray(&strct)
				pkg.Structs[strct.Name] = strctIdx
			}
		}

		// if sPkg.GlobalsSize > 0 {
		// 	pkg.Globals = make([]CXArgumentIndex, sPkg.GlobalsSize)
		// }

		// CurrentFunction
		if sPkg.FunctionsSize > 0 {
			pkg.CurrentFunction = pkg.Functions[sPkg.CurrentFunctionName]
		}

		fnCounter += sPkg.FunctionsSize
		strctCounter += sPkg.StructsSize

		// imports
		if sPkg.ImportsSize > 0 {
			// getting indexes of imports
			idxs := deserializeIntegers(sPkg.ImportsOffset, sPkg.ImportsSize, s)

			for _, idx := range idxs {
				impPkg := deserializePackageImport(&s.Packages[idx], s, prgrm)
				pkg.Imports[impPkg.Name] = CXPackageIndex(impPkg.Index)
			}
		}

		// globals
		if sPkg.GlobalsSize > 0 {
			glblArgs := deserializeArguments(sPkg.GlobalsOffset, sPkg.GlobalsSize, s, prgrm)
			var glblArgsIdxs []CXArgumentIndex
			for _, glbl := range glblArgs {
				glblIdx := prgrm.AddCXArgInArray(glbl)
				glblArgsIdxs = append(glblArgsIdxs, CXArgumentIndex(glblIdx))
			}
			// pkg.Globals = glblArgsIdxs
			for _, glblIdx := range glblArgsIdxs {
				pkg.Globals.AddField_Globals_CXAtomicOps(prgrm, glblIdx)
			}
		}

		// structs
		if sPkg.StructsSize > 0 {
			for _, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				strctName := deserializeString(sStrct.NameOffset, sStrct.NameSize, s)
				strct := &prgrm.CXStructs[pkg.Structs[strctName]]
				deserializeStruct(&sStrct, strct, s, prgrm)
			}
		}

		// functions
		if sPkg.FunctionsSize > 0 {
			for _, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				fnName := deserializeString(sFn.NameOffset, sFn.NameSize, s)
				pkgFn := prgrm.GetFunctionFromArray(pkg.Functions[fnName])

				deserializeFunction(&sFn, pkgFn, s, prgrm)
			}
		}
		prgrm.CXPackages[pkgIdx] = *pkg
	}

	// current package
	prgrm.CurrentPackage = prgrm.Packages[s.Program.CurrentPackageName]
}

func deserializeStruct(sStrct *serializedStruct, strct *CXStruct, s *SerializedCXProgram, prgrm *CXProgram) {
	strct.Name = deserializeString(sStrct.NameOffset, sStrct.NameSize, s)
	strct.Fields = prgrm.AddPointerArgsToTypeSignaturesArray(deserializeArguments(sStrct.FieldsOffset, sStrct.FieldsSize, s, prgrm))

	strct.Package = prgrm.Packages[sStrct.PackageName]
}

func deserializeArguments(off int64, size int64, s *SerializedCXProgram, prgrm *CXProgram) []*CXArgument {
	if size < 1 {
		return nil
	}

	// getting indexes of arguments
	idxs := deserializeIntegers(off, size, s)

	// sArgs := s.Arguments[off : off + size]
	args := make([]*CXArgument, size)
	for i, idx := range idxs {
		args[i] = deserializeArgument(&s.Arguments[idx], s, prgrm)
	}
	return args
}

// func getStructType(sArg *serializedArgument, s *SerializedCXProgram, prgrm *CXProgram) *CXStruct {
// 	if sArg.StructTypeOffset < 0 {
// 		return nil
// 	}

// 	//structTypePkg := prgrm.Packages[s.Structs[sArg.StructTypeOffset].PackageOffset]
// 	sStrct := s.Structs[sArg.StructTypeOffset]
// 	structTypeName := deserializeString(sStrct.NameOffset, sStrct.NameSize, s)

// 	for _, strct := range structTypePkg.Structs {
// 		if strct.Name == structTypeName {
// 			return strct
// 		}
// 	}

// 	return nil
// }

func deserializeArgument(sArg *serializedArgument, s *SerializedCXProgram, prgrm *CXProgram) *CXArgument {
	var arg CXArgument
	arg.ArgDetails = &CXArgumentDebug{}
	arg.Name = deserializeString(sArg.NameOffset, sArg.NameSize, s)
	arg.Type = types.Code(sArg.Type)

	//arg.StructType = getStructType(sArg, s, prgrm)

	arg.Size = types.Cast_i64_to_ptr(sArg.Size)
	arg.Offset = types.Cast_i64_to_ptr(sArg.Offset)
	arg.PassBy = int(sArg.PassBy)

	arg.DeclarationSpecifiers = deserializeIntegers(sArg.DeclarationSpecifiersOffset, sArg.DeclarationSpecifiersSize, s)

	arg.IsSlice = deserializeBool(sArg.IsSlice)
	// arg.IsPointer = deserializeBool(sArg.IsPointer)
	// arg.IsReference = deserializeBool(sArg.IsReference)

	arg.PreviouslyDeclared = deserializeBool(sArg.PreviouslyDeclared)

	arg.Lengths = deserializePointers(sArg.LengthsOffset, sArg.LengthsSize, s)
	// TODO: include serializing of CXTypeSignature for indexes
	// arg.Indexes = prgrm.AddPointerArgsToCXArgsArray(deserializeArguments(sArg.IndexesOffset, sArg.IndexesSize, s, prgrm))
	arg.Fields = prgrm.AddPointerArgsToCXArgsArray(deserializeArguments(sArg.FieldsOffset, sArg.FieldsSize, s, prgrm))
	arg.Inputs = prgrm.AddPointerArgsToCXArgsArray(deserializeArguments(sArg.InputsOffset, sArg.InputsSize, s, prgrm))
	arg.Outputs = prgrm.AddPointerArgsToCXArgsArray(deserializeArguments(sArg.OutputsOffset, sArg.OutputsSize, s, prgrm))
	arg.Package = -1
	if _, ok := prgrm.Packages[sArg.PackageName]; ok {
		arg.Package = prgrm.Packages[sArg.PackageName]
	}

	return &arg
}

func deserializeOperator(sExpr *serializedExpression, s *SerializedCXProgram, prgrm *CXProgram) *CXFunction {
	if sExpr.OperatorOffset < 0 {
		return nil
	}

	opPkgIdx := prgrm.Packages[s.Functions[sExpr.OperatorOffset].PackageName]
	sOp := s.Functions[sExpr.OperatorOffset]
	opName := deserializeString(sOp.NameOffset, sOp.NameSize, s)

	opPkg, err := prgrm.GetPackageFromArray(opPkgIdx)
	if err != nil {
		panic(err)
	}

	for _, fnIdx := range opPkg.Functions {
		fn := prgrm.GetFunctionFromArray(fnIdx)

		if fn.Name == opName {
			return fn
		}
	}

	return nil
}

func deserializePackageImport(sImp *serializedPackage, s *SerializedCXProgram, prgrm *CXProgram) *CXPackage {
	impName := deserializeString(sImp.NameOffset, sImp.NameSize, s)

	for _, pkgIdx := range prgrm.Packages {
		pkg, err := prgrm.GetPackageFromArray(pkgIdx)
		if err != nil {
			panic(err)
		}

		if pkg.Name == impName {
			return pkg
		}
	}

	return nil
}

func deserializeExpressionFunction(sExpr *serializedExpression, s *SerializedCXProgram, prgrm *CXProgram) CXFunctionIndex {
	if sExpr.FunctionOffset < 0 {
		return -1
	}

	fnPkgIdx := prgrm.Packages[s.Functions[sExpr.FunctionOffset].PackageName]
	sFn := s.Functions[sExpr.FunctionOffset]
	fnName := deserializeString(sFn.NameOffset, sFn.NameSize, s)

	fnPkg, err := prgrm.GetPackageFromArray(fnPkgIdx)
	if err != nil {
		panic(err)
	}
	for _, fnIdx := range fnPkg.Functions {
		fn := prgrm.GetFunctionFromArray(fnIdx)

		if fn.Name == fnName {
			return fnIdx
		}
	}

	return -1
}

func deserializeExpressions(off int64, size int64, s *SerializedCXProgram, prgrm *CXProgram) []CXExpression {
	if size < 1 {
		return nil
	}

	// getting indexes of expressions
	idxs := deserializeIntegers(off, size, s)

	// sExprs := s.Expressions[off : off + size]
	exprs := make([]CXExpression, size)
	for i, idx := range idxs {
		exprs[i] = deserializeExpression(&s.Expressions[idx], s, prgrm)
	}
	return exprs
}

func deserializeExpression(sExpr *serializedExpression, s *SerializedCXProgram, prgrm *CXProgram) CXExpression {
	var expr CXExpression

	expr.ExpressionType = CXEXPR_TYPE(sExpr.ExpressionType)
	switch sExpr.Type {
	case int64(CX_LINE):
		cxLine := &CXLine{}

		cxLine.FileName = sExpr.FileName
		cxLine.LineNumber = int(sExpr.LineNumber)
		cxLine.LineStr = sExpr.LineStr
		index := prgrm.AddCXLine(cxLine)
		expr.Index = index
		expr.Type = CX_LINE
	case int64(CX_ATOMIC_OPERATOR):
		cxAtomicOpOperatorIdx := prgrm.AddFunctionInArray(&CXFunction{})
		cxAtomicOp := &CXAtomicOperator{Operator: cxAtomicOpOperatorIdx}

		if deserializeBool(sExpr.IsNative) {
			opIdx := prgrm.AddNativeFunctionInArray(Natives[int(sExpr.OpCode)])
			cxAtomicOp.Operator = opIdx
		} else {
			opIdx := prgrm.AddFunctionInArray(deserializeOperator(sExpr, s, prgrm))
			cxAtomicOp.Operator = opIdx
		}

		inputCXArgsArray := deserializeArguments(sExpr.InputsOffset, sExpr.InputsSize, s, prgrm)
		for _, inputCXArg := range inputCXArgsArray {
			typeSignature := GetCXTypeSignatureRepresentationOfCXArg(prgrm, inputCXArg)
			typeSignatureIdx := prgrm.AddCXTypeSignatureInArray(typeSignature)
			cxAtomicOp.AddInput(prgrm, typeSignatureIdx)
		}

		outputCXArgsArray := deserializeArguments(sExpr.OutputsOffset, sExpr.OutputsSize, s, prgrm)
		for _, outputCXArg := range outputCXArgsArray {
			typeSignature := GetCXTypeSignatureRepresentationOfCXArg(prgrm, outputCXArg)
			typeSignatureIdx := prgrm.AddCXTypeSignatureInArray(typeSignature)
			cxAtomicOp.AddOutput(prgrm, typeSignatureIdx)
		}

		cxAtomicOp.Label = deserializeString(sExpr.LabelOffset, sExpr.LabelSize, s)

		cxAtomicOp.ThenLines = int(sExpr.ThenLines)
		cxAtomicOp.ElseLines = int(sExpr.ElseLines)
		cxAtomicOp.Package = prgrm.Packages[sExpr.PackageName]
		cxAtomicOp.Function = deserializeExpressionFunction(sExpr, s, prgrm)

		index := prgrm.AddCXAtomicOp(cxAtomicOp)
		expr.Index = index
		expr.Type = CX_ATOMIC_OPERATOR
	}

	return expr
}

func deserializeFunction(sFn *serializedFunction, fn *CXFunction, s *SerializedCXProgram, prgrm *CXProgram) {
	fn.Name = deserializeString(sFn.NameOffset, sFn.NameSize, s)

	inputCXArgsArray := deserializeArguments(sFn.InputsOffset, sFn.InputsSize, s, prgrm)
	for _, inputCXArg := range inputCXArgsArray {
		fn.AddInput(prgrm, inputCXArg)
	}

	outputCXArgsArray := deserializeArguments(sFn.OutputsOffset, sFn.OutputsSize, s, prgrm)
	for _, outputCXArg := range outputCXArgsArray {
		fn.AddOutput(prgrm, outputCXArg)
	}

	fn.ListOfPointers = deserializeArguments(sFn.ListOfPointersOffset, sFn.ListOfPointersSize, s, prgrm)
	fn.Package = prgrm.Packages[sFn.PackageName]
	fn.Expressions = deserializeExpressions(sFn.ExpressionsOffset, sFn.ExpressionsSize, s, prgrm)
	fn.Size = types.Cast_i64_to_ptr(sFn.Size)
	fn.LineCount = int(sFn.Length)
	prgrm.CXFunctions[fn.Index] = *fn
}

func deserializeBool(val int64) bool {
	return val == 1
}

func deserializePointers(off int64, size int64, s *SerializedCXProgram) []types.Pointer {
	if size < 1 {
		return nil
	}
	ints := s.Integers[off : off+size]
	res := make([]types.Pointer, len(ints))
	for i, in := range ints {
		res[i] = types.Cast_i64_to_ptr(in)
	}

	return res
}

func deserializeIntegers(off int64, size int64, s *SerializedCXProgram) []int {
	if size < 1 {
		return nil
	}
	ints := s.Integers[off : off+size]
	res := make([]int, len(ints))
	for i, in := range ints {
		res[i] = int(in)
	}

	return res
}

// initDeserialization initializes the CXProgram fields that represent a CX program. This should be refactored, as the names Deserialize and initDeserialization create some naming conflict.
func initDeserialization(prgrm *CXProgram, s *SerializedCXProgram) {
	prgrm.Memory = s.Memory
	prgrm.Packages = make(map[string]CXPackageIndex, len(s.Packages))
	prgrm.CallStack = make([]CXCall, constants.CALLSTACK_SIZE)
	prgrm.Heap.StartsAt = types.Cast_i64_to_ptr(s.Program.HeapStartsAt)
	prgrm.Heap.Pointer = types.Cast_i64_to_ptr(s.Program.HeapPointer)
	prgrm.Stack.Size = types.Cast_i64_to_ptr(s.Program.StackSize)
	prgrm.Data.Size = types.Cast_i64_to_ptr(s.Program.DataSegmentSize)
	prgrm.Data.StartsAt = types.Cast_i64_to_ptr(s.Program.DataSegmentStartsAt)
	prgrm.Heap.Size = types.Cast_i64_to_ptr(s.Program.HeapSize)
	prgrm.Version = deserializeString(s.Program.VersionOffset, s.Program.VersionSize, s)

	// This means reinstantiate memory and add DataSegmentMemory
	if len(s.DataSegmentMemory) > 0 && len(s.Memory) == 0 {
		minHeapSize := MinHeapSize()
		prgrm.Memory = make([]byte, constants.STACK_SIZE+minHeapSize)
		y := 0
		for i := prgrm.Data.StartsAt; i < prgrm.Data.StartsAt+prgrm.Data.Size; i++ {
			prgrm.Memory[i] = s.DataSegmentMemory[y]
			y++
		}
	}
	deserializePackages(s, prgrm)
}

// Deserialize deserializes a serialized CX program back to its golang struct representation.
func Deserialize(b []byte, useCompression bool) (prgrm *CXProgram) {
	prgrm = &CXProgram{}
	var s SerializedCXProgram

	if useCompression {
		// Uncompress using LZ4
		UncompressBytesLZ4(&b)
	}

	deserializeRaw(b, 0, types.Cast_int_to_ptr(len(b)), &s)
	initDeserialization(prgrm, &s)

	return prgrm
}

// CopyProgramState copies the program state from `prgrm1` to `prgrm2`.
func CopyProgramState(sPrgrm1, sPrgrm2 *[]byte) {
	idxSize := types.Cast_ui64_to_ptr(encoder.Size(serializedCXProgramIndex{}))

	var index1 serializedCXProgramIndex
	var index2 serializedCXProgramIndex

	deserializeRaw((*sPrgrm1), 0, idxSize, &index1)
	deserializeRaw((*sPrgrm2), 0, idxSize, &index2)

	var prgrm1Info serializedProgram
	deserializeRaw((*sPrgrm1),
		types.Cast_i64_to_ptr(index1.ProgramOffset),
		types.Cast_i64_to_ptr(index1.CallsOffset-index1.ProgramOffset),
		&prgrm1Info)

	var prgrm2Info serializedProgram
	deserializeRaw((*sPrgrm2),
		types.Cast_i64_to_ptr(index2.ProgramOffset),
		types.Cast_i64_to_ptr(index2.CallsOffset-index2.ProgramOffset),
		&prgrm2Info)

	// the stack segment should be 0 for prgrm1, but just in case
	var prgrmState []byte
	prgrmState = append(prgrmState, make([]byte, prgrm2Info.StackSize)...)
	// We are only interested on extracting the data segment
	prgrmState = append(prgrmState, (*sPrgrm1)[index1.StringsOffset+prgrm1Info.StackSize:index1.StringsOffset+prgrm1Info.StackSize+(prgrm2Info.HeapStartsAt-prgrm2Info.StackSize)]...)

	for i, byt := range prgrmState {
		(*sPrgrm2)[i+int(index2.MemoryOffset)] = byt
	}
}

// GetSerializedStackSize returns the stack size of a serialized CX program starts.
func GetSerializedStackSize(sPrgrm []byte) int {
	idxSize := types.Cast_ui64_to_ptr(encoder.Size(serializedCXProgramIndex{}))
	var index serializedCXProgramIndex
	deserializeRaw(sPrgrm, 0, idxSize, &index)

	var prgrmInfo serializedProgram
	deserializeRaw(sPrgrm,
		types.Cast_i64_to_ptr(index.ProgramOffset),
		types.Cast_i64_to_ptr(index.CallsOffset-index.ProgramOffset),
		&prgrmInfo)

	return int(prgrmInfo.StackSize)
}

// GetSerializedDataSize returns the size of the data segment of a serialized CX program.
func GetSerializedDataSize(sPrgrm []byte) int {
	idxSize := types.Cast_ui64_to_ptr(encoder.Size(serializedCXProgramIndex{}))
	var index serializedCXProgramIndex
	deserializeRaw(sPrgrm, 0, idxSize, &index)

	var prgrmInfo serializedProgram
	deserializeRaw(sPrgrm,
		types.Cast_i64_to_ptr(index.ProgramOffset),
		types.Cast_i64_to_ptr(index.CallsOffset-index.ProgramOffset),
		&prgrmInfo)

	return int(prgrmInfo.HeapStartsAt - prgrmInfo.StackSize)
}
