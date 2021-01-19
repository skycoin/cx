package cxcore

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

type sIndex struct {
	ProgramOffset     int32
	CallsOffset       int32
	PackagesOffset    int32
	StructsOffset     int32
	FunctionsOffset   int32
	ExpressionsOffset int32
	ArgumentsOffset   int32
	IntegersOffset    int32
	NamesOffset       int32
	MemoryOffset      int32
}

type sProgram struct {
	PackagesOffset       int32
	PackagesSize         int32
	CurrentPackageOffset int32

	InputsOffset int32
	InputsSize   int32

	OutputsOffset int32
	OutputsSize   int32

	CallStackOffset int32
	CallStackSize   int32

	CallCounter int32

	MemoryOffset int32
	MemorySize   int32

	HeapPointer  int32
	StackPointer int32
	StackSize    int32
	HeapSize     int32
	HeapStartsAt int32

	Terminated int32

	BCPackageCount int32

	VersionOffset int32
	VersionSize   int32
}

type sCall struct {
	OperatorOffset int32
	Line           int32
	FramePointer   int32
}

type sPackage struct {
	NameOffset            int32
	NameSize              int32
	ImportsOffset         int32
	ImportsSize           int32
	StructsOffset         int32
	StructsSize           int32
	GlobalsOffset         int32
	GlobalsSize           int32
	FunctionsOffset       int32
	FunctionsSize         int32
	CurrentFunctionOffset int32
	CurrentStructOffset   int32
}

type sStruct struct {
	NameOffset   int32
	NameSize     int32
	FieldsOffset int32
	FieldsSize   int32

	Size int32

	PackageOffset int32
}

type sFunction struct {
	NameOffset        int32
	NameSize          int32
	InputsOffset      int32
	InputsSize        int32
	OutputsOffset     int32
	OutputsSize       int32
	ExpressionsOffset int32
	ExpressionsSize   int32
	Size              int32
	Length            int32

	ListOfPointersOffset int32
	ListOfPointersSize   int32

	// We're going to determine this when procesing the expressions. Check sExpression type
	// IsNative                        int32
	// OpCode                          int32

	CurrentExpressionOffset int32
	PackageOffset           int32
}

type sExpression struct {
	OperatorOffset int32
	// we add these two fields here so we don't add every native sFunction to the serialization
	// the CX runtime already knows about the natives properties. We just need the code if IsNative = true
	IsNative int32
	OpCode   int32

	InputsOffset  int32
	InputsSize    int32
	OutputsOffset int32
	OutputsSize   int32

	LabelOffset int32
	LabelSize   int32
	ThenLines   int32
	ElseLines   int32

	ScopeOperation int32

	IsMethodCall    int32
	IsStructLiteral int32
	IsArrayLiteral  int32
	IsUndType       int32
	IsBreak         int32
	IsContinue      int32

	FunctionOffset int32
	PackageOffset  int32
}

type sArgument struct {
	NameOffset       int32
	NameSize         int32
	Type             int32
	CustomTypeOffset int32
	Size             int32
	TotalSize        int32

	Offset int32

	IndirectionLevels           int32
	DereferenceLevels           int32
	DereferenceOperationsOffset int32
	DereferenceOperationsSize   int32
	DeclarationSpecifiersOffset int32
	DeclarationSpecifiersSize   int32

	IsSlice      int32
	IsArray      int32
	IsArrayFirst int32
	IsPointer    int32
	IsReference  int32

	IsDereferenceFirst int32
	IsStruct           int32
	IsRest             int32
	IsLocalDeclaration int32
	IsShortDeclaration int32
	PreviouslyDeclared int32

	PassBy     int32
	DoesEscape int32

	LengthsOffset int32
	LengthsSize   int32
	IndexesOffset int32
	IndexesSize   int32
	FieldsOffset  int32
	FieldsSize    int32
	InputsOffset  int32
	InputsSize    int32
	OutputsOffset int32
	OutputsSize   int32

	PackageOffset int32
}

type sAll struct {
	Index   sIndex
	Program sProgram

	Packages     []sPackage
	PackagesMap  map[string]int
	Structs      []sStruct
	StructsMap   map[string]int
	Functions    []sFunction
	FunctionsMap map[string]int

	Expressions []sExpression
	Arguments   []sArgument
	Calls       []sCall

	Names    []byte
	NamesMap map[string]int
	Integers []int32

	Memory []byte
}

func serializeName(name string, s *sAll) (int32, int32) {
	if name == "" {
		return int32(-1), int32(-1)
	}

	size := encoder.Size(name)

	off, found := s.NamesMap[name]
	if found {
		return int32(off), int32(size)
	}
	off = len(s.Names)
	s.Names = append(s.Names, encoder.Serialize(name)...)
	s.NamesMap[name] = off

	return int32(off), int32(size)
}

func indexPackage(pkg *CXPackage, s *sAll) {
	if _, found := s.PackagesMap[pkg.Name]; !found {
		s.PackagesMap[pkg.Name] = len(s.PackagesMap)
	} else {
		panic("duplicated package in serialization process")
	}
}

func indexStruct(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if _, found := s.StructsMap[strctName]; !found {
		s.StructsMap[strctName] = len(s.StructsMap)
	} else {
		panic("duplicated struct in serialization process")
	}
}

func indexFunction(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if _, found := s.FunctionsMap[fnName]; !found {
		s.FunctionsMap[fnName] = len(s.FunctionsMap)
	} else {
		panic("duplicated function in serialization process")
	}
}

func serializeBoolean(val bool) int32 {
	if val {
		return 1
	}
	return 0
}

func serializeIntegers(ints []int, s *sAll) (int32, int32) {
	if len(ints) == 0 {
		return int32(-1), int32(-1)
	}
	off := len(s.Integers)
	l := len(ints)

	ints32 := make([]int32, l)
	for i, int := range ints {
		ints32[i] = int32(int)
	}

	s.Integers = append(s.Integers, ints32...)

	return int32(off), int32(l)
}

func serializeArgument(arg *CXArgument, s *sAll) int {
	s.Arguments = append(s.Arguments, sArgument{})
	argOff := len(s.Arguments) - 1

	sNil := int32(-1)

	s.Arguments[argOff].NameOffset, s.Arguments[argOff].NameSize = serializeName(arg.Name, s)

	s.Arguments[argOff].Type = int32(arg.Type)

	if arg.CustomType == nil {
		s.Arguments[argOff].CustomTypeOffset = sNil
	} else {
		strctName := arg.CustomType.Package.Name + "." + arg.CustomType.Name
		if strctOff, found := s.StructsMap[strctName]; found {
			s.Arguments[argOff].CustomTypeOffset = int32(strctOff)
		} else {
			panic("struct reference not found")
		}
	}

	s.Arguments[argOff].Size = int32(arg.Size)
	s.Arguments[argOff].TotalSize = int32(arg.TotalSize)
	s.Arguments[argOff].Offset = int32(arg.Offset)
	s.Arguments[argOff].IndirectionLevels = int32(arg.IndirectionLevels)
	s.Arguments[argOff].DereferenceLevels = int32(arg.DereferenceLevels)

	s.Arguments[argOff].DereferenceOperationsOffset,
		s.Arguments[argOff].DereferenceOperationsSize = serializeIntegers(arg.DereferenceOperations, s)

	s.Arguments[argOff].DeclarationSpecifiersOffset,
		s.Arguments[argOff].DeclarationSpecifiersSize = serializeIntegers(arg.DeclarationSpecifiers, s)

	s.Arguments[argOff].IsSlice = serializeBoolean(arg.IsSlice)
	s.Arguments[argOff].IsArray = serializeBoolean(arg.IsArray)
	s.Arguments[argOff].IsArrayFirst = serializeBoolean(arg.IsArrayFirst)
	s.Arguments[argOff].IsPointer = serializeBoolean(arg.IsPointer)
	s.Arguments[argOff].IsReference = serializeBoolean(arg.IsReference)

	s.Arguments[argOff].IsDereferenceFirst = serializeBoolean(arg.IsDereferenceFirst)
	s.Arguments[argOff].IsStruct = serializeBoolean(arg.IsStruct)
	s.Arguments[argOff].IsRest = serializeBoolean(arg.IsRest)
	s.Arguments[argOff].IsLocalDeclaration = serializeBoolean(arg.IsLocalDeclaration)
	s.Arguments[argOff].IsShortDeclaration = serializeBoolean(arg.IsShortDeclaration)
	s.Arguments[argOff].PreviouslyDeclared = serializeBoolean(arg.PreviouslyDeclared)

	s.Arguments[argOff].PassBy = int32(arg.PassBy)
	s.Arguments[argOff].DoesEscape = serializeBoolean(arg.DoesEscape)

	s.Arguments[argOff].LengthsOffset, s.Arguments[argOff].LengthsSize = serializeIntegers(arg.Lengths, s)
	s.Arguments[argOff].IndexesOffset, s.Arguments[argOff].IndexesSize = serializeSliceOfArguments(arg.Indexes, s)
	s.Arguments[argOff].FieldsOffset, s.Arguments[argOff].FieldsSize = serializeSliceOfArguments(arg.Fields, s)
	s.Arguments[argOff].InputsOffset, s.Arguments[argOff].InputsSize = serializeSliceOfArguments(arg.Inputs, s)
	s.Arguments[argOff].OutputsOffset, s.Arguments[argOff].OutputsSize = serializeSliceOfArguments(arg.Outputs, s)

	if pkgOff, found := s.PackagesMap[arg.Package.Name]; found {
		s.Arguments[argOff].PackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	return argOff
}

func serializeSliceOfArguments(args []*CXArgument, s *sAll) (int32, int32) {
	if len(args) == 0 {
		return int32(-1), int32(-1)
	}
	idxs := make([]int, len(args))
	for i, arg := range args {
		idxs[i] = serializeArgument(arg, s)
	}
	return serializeIntegers(idxs, s)
}

func serializeCalls(calls []CXCall, s *sAll) (int32, int32) {
	if len(calls) == 0 {
		return int32(-1), int32(-1)
	}
	idxs := make([]int, len(calls))
	for i, call := range calls {
		idxs[i] = serializeCall(&call, s)
	}
	return serializeIntegers(idxs, s)

}

func serializeExpression(expr *CXExpression, s *sAll) int {
	s.Expressions = append(s.Expressions, sExpression{})
	exprOff := len(s.Expressions) - 1
	sExpr := &s.Expressions[exprOff]

	sNil := int32(-1)

	if expr.Operator == nil {
		// then it's a declaration
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = int32(-1)
	} else if expr.Operator.IsNative {
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(true)
		sExpr.OpCode = int32(expr.Operator.OpCode)
	} else {
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = sNil

		opName := expr.Operator.Package.Name + "." + expr.Operator.Name
		if opOff, found := s.FunctionsMap[opName]; found {
			sExpr.OperatorOffset = int32(opOff)
		}
	}

	sExpr.InputsOffset, sExpr.InputsSize = serializeSliceOfArguments(expr.Inputs, s)
	sExpr.OutputsOffset, sExpr.OutputsSize = serializeSliceOfArguments(expr.Outputs, s)

	sExpr.LabelOffset, sExpr.LabelSize = serializeName(expr.Label, s)
	sExpr.ThenLines = int32(expr.ThenLines)
	sExpr.ElseLines = int32(expr.ElseLines)
	sExpr.ScopeOperation = int32(expr.ScopeOperation)

	sExpr.IsMethodCall = serializeBoolean(expr.IsMethodCall)
	sExpr.IsStructLiteral = serializeBoolean(expr.IsStructLiteral)
	sExpr.IsArrayLiteral = serializeBoolean(expr.IsArrayLiteral)
	sExpr.IsUndType = serializeBoolean(expr.IsUndType)
	sExpr.IsBreak = serializeBoolean(expr.IsBreak)
	sExpr.IsContinue = serializeBoolean(expr.IsContinue)

	fnName := expr.Function.Package.Name + "." + expr.Function.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sExpr.FunctionOffset = int32(fnOff)
	} else {
		panic("function reference not found")
	}

	if pkgOff, found := s.PackagesMap[expr.Package.Name]; found {
		sExpr.PackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	return exprOff
}

func serializeCall(call *CXCall, s *sAll) int {
	s.Calls = append(s.Calls, sCall{})
	callOff := len(s.Calls) - 1
	sCall := &s.Calls[callOff]

	opName := call.Operator.Package.Name + "." + call.Operator.Name
	if opOff, found := s.FunctionsMap[opName]; found {
		sCall.OperatorOffset = int32(opOff)
		sCall.Line = int32(call.Line)
		sCall.FramePointer = int32(call.FramePointer)
	} else {
		panic("function reference not found")
	}

	return callOff
}

func serializeProgram(prgrm *CXProgram, s *sAll) {
	s.Program = sProgram{}
	sPrgrm := &s.Program
	sPrgrm.PackagesOffset = int32(0)
	sPrgrm.PackagesSize = int32(len(prgrm.Packages))

	if pkgOff, found := s.PackagesMap[prgrm.CurrentPackage.Name]; found {
		sPrgrm.CurrentPackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	sPrgrm.InputsOffset, sPrgrm.InputsSize = serializeSliceOfArguments(prgrm.Inputs, s)
	sPrgrm.OutputsOffset, sPrgrm.OutputsSize = serializeSliceOfArguments(prgrm.Outputs, s)

	sPrgrm.CallStackOffset, sPrgrm.CallStackSize = serializeCalls(prgrm.CallStack[:prgrm.CallCounter], s)

	sPrgrm.CallCounter = int32(prgrm.CallCounter)

	sPrgrm.MemoryOffset = int32(0)
	sPrgrm.MemorySize = int32(len(PROGRAM.Memory))

	sPrgrm.HeapPointer = int32(prgrm.HeapPointer)
	sPrgrm.StackPointer = int32(prgrm.StackPointer)
	sPrgrm.StackSize = int32(prgrm.StackSize)
	sPrgrm.HeapSize = int32(prgrm.HeapSize)
	sPrgrm.HeapStartsAt = int32(prgrm.HeapStartsAt)

	sPrgrm.Terminated = serializeBoolean(prgrm.Terminated)
	sPrgrm.BCPackageCount = int32(prgrm.BCPackageCount)
	sPrgrm.VersionOffset, sPrgrm.VersionSize = serializeName(prgrm.Version, s)
}

func sStructArguments(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if strctOff, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[strctOff]
		sStrct.FieldsOffset, sStrct.FieldsSize = serializeSliceOfArguments(strct.Fields, s)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionArguments(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[fnOff]

		sFn.InputsOffset, sFn.InputsSize = serializeSliceOfArguments(fn.Inputs, s)
		sFn.OutputsOffset, sFn.OutputsSize = serializeSliceOfArguments(fn.Outputs, s)
		sFn.ListOfPointersOffset, sFn.ListOfPointersSize = serializeSliceOfArguments(fn.ListOfPointers, s)
	} else {
		panic("function reference not found")
	}
}

func sPackageName(pkg *CXPackage, s *sAll) {
	sPkg := &s.Packages[s.PackagesMap[pkg.Name]]
	sPkg.NameOffset, sPkg.NameSize = serializeName(pkg.Name, s)
}

func sStructName(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	sStrct := &s.Structs[s.StructsMap[strctName]]
	sStrct.NameOffset, sStrct.NameSize = serializeName(strct.Name, s)
}

func sFunctionName(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.NameOffset, sFn.NameSize = serializeName(fn.Name, s)
	} else {
		panic("function reference not found")
	}
}

func sPackageGlobals(pkg *CXPackage, s *sAll) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]
		sPkg.GlobalsOffset, sPkg.GlobalsSize = serializeSliceOfArguments(pkg.Globals, s)
	} else {
		panic("package reference not found")
	}
}

func sPackageImports(pkg *CXPackage, s *sAll) {
	l := len(pkg.Imports)
	if l == 0 {
		s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int32(-1)
		s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int32(-1)
		return
	}
	imps := make([]int32, l)
	for i, imp := range pkg.Imports {
		if idx, found := s.PackagesMap[imp.Name]; found {
			imps[i] = int32(idx)
		} else {
			panic("import package reference not found")
		}
	}

	s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int32(len(s.Integers))
	s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int32(l)
	s.Integers = append(s.Integers, imps...)
}

func sStructPackage(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if pkgOff, found := s.PackagesMap[strct.Package.Name]; found {
		if off, found := s.StructsMap[strctName]; found {
			sStrct := &s.Structs[off]
			sStrct.PackageOffset = int32(pkgOff)
		} else {
			panic("struct reference not found")
		}
	} else {
		panic("struct's package reference not found")
	}
}

func sFunctionPackage(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if pkgOff, found := s.PackagesMap[fn.Package.Name]; found {
		if off, found := s.FunctionsMap[fnName]; found {
			sFn := &s.Functions[off]
			sFn.PackageOffset = int32(pkgOff)
		} else {
			panic("function reference not found")
		}
	} else {
		panic("function's package reference not found")
	}
}

func sPackageIntegers(pkg *CXPackage, s *sAll) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]

		if pkg.CurrentFunction == nil {
			// package has no functions
			sPkg.CurrentFunctionOffset = int32(-1)
		} else {
			currFnName := pkg.CurrentFunction.Package.Name + "." + pkg.CurrentFunction.Name

			if fnOff, found := s.FunctionsMap[currFnName]; found {
				sPkg.CurrentFunctionOffset = int32(fnOff)
			} else {
				panic("function reference not found")
			}
		}

		if pkg.CurrentStruct == nil {
			// package has no structs
			sPkg.CurrentStructOffset = int32(-1)
		} else {
			currStrctName := pkg.CurrentStruct.Package.Name + "." + pkg.CurrentStruct.Name

			if strctOff, found := s.StructsMap[currStrctName]; found {
				sPkg.CurrentStructOffset = int32(strctOff)
			} else {
				panic("struct reference not found")
			}
		}
	} else {
		panic("package reference not found")
	}
}

func sStructIntegers(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if off, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[off]
		sStrct.Size = int32(strct.Size)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionIntegers(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.Size = int32(fn.Size)
		sFn.Length = int32(fn.Length)
	} else {
		panic("function reference not found")
	}
}

func initSerialization(prgrm *CXProgram, s *sAll) {
	s.PackagesMap = make(map[string]int)
	s.StructsMap = make(map[string]int)
	s.FunctionsMap = make(map[string]int)
	s.NamesMap = make(map[string]int)

	s.Calls = make([]sCall, prgrm.CallCounter)
	s.Packages = make([]sPackage, len(prgrm.Packages))

	// s.Memory = prgrm.Memory[:PROGRAM.HeapStartsAt+PROGRAM.HeapPointer]
	s.Memory = prgrm.Memory

	var numStrcts int
	var numFns int

	for _, pkg := range prgrm.Packages {
		numStrcts += len(pkg.Structs)
		numFns += len(pkg.Functions)
	}

	s.Structs = make([]sStruct, numStrcts)
	s.Functions = make([]sFunction, numFns)
	// args and exprs need to be appended as they are found
}

// SplitSerialize ...
func splitSerialize(prgrm *CXProgram, s *sAll, fnCounter, strctCounter *int32, from, to int) {
	// indexing packages and serializing their names
	for _, pkg := range prgrm.Packages[from:to] {
		indexPackage(pkg, s)
		sPackageName(pkg, s)
	}
	// we first needed to populate references to all packages
	// now we add the imports' references
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageImports(pkg, s)
	}

	// structs
	for _, pkg := range prgrm.Packages[from:to] {
		for _, strct := range pkg.Structs {
			indexStruct(strct, s)
			sStructName(strct, s)
			sStructPackage(strct, s)
			sStructIntegers(strct, s)
		}
	}
	// we first needed to populate references to all structs
	// now we add fields
	for _, pkg := range prgrm.Packages[from:to] {
		for _, strct := range pkg.Structs {
			sStructArguments(strct, s)
		}
	}

	// globals
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageGlobals(pkg, s)
	}

	// functions
	for _, pkg := range prgrm.Packages[from:to] {
		for _, fn := range pkg.Functions {
			indexFunction(fn, s)
			sFunctionName(fn, s)
			sFunctionPackage(fn, s)
			sFunctionIntegers(fn, s)
			sFunctionArguments(fn, s)
		}
	}

	// package elements' offsets and sizes
	for _, pkg := range prgrm.Packages[from:to] {
		if pkgOff, found := s.PackagesMap[pkg.Name]; found {
			sPkg := &s.Packages[pkgOff]

			if len(pkg.Structs) == 0 {
				sPkg.StructsOffset = int32(-1)
				sPkg.StructsSize = int32(-1)
			} else {
				sPkg.StructsOffset = *strctCounter
				lenStrcts := int32(len(pkg.Structs))
				sPkg.StructsSize = lenStrcts
				*strctCounter += lenStrcts
			}

			if len(pkg.Functions) == 0 {
				sPkg.FunctionsOffset = int32(-1)
				sPkg.FunctionsSize = int32(-1)
			} else {
				sPkg.FunctionsOffset = *fnCounter
				lenFns := int32(len(pkg.Functions))
				sPkg.FunctionsSize = lenFns
				*fnCounter += lenFns
			}
		} else {
			panic("package reference not found")
		}
	}

	// package integers
	// we needed the references to all functions and structs first
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageIntegers(pkg, s)
	}

	// expressions
	for _, pkg := range prgrm.Packages[from:to] {
		for _, fn := range pkg.Functions {
			fnName := fn.Package.Name + "." + fn.Name
			if fnOff, found := s.FunctionsMap[fnName]; found {
				sFn := &s.Functions[fnOff]

				if len(fn.Expressions) == 0 {
					sFn.ExpressionsOffset = int32(-1)
					sFn.ExpressionsSize = int32(-1)
					sFn.CurrentExpressionOffset = int32(-1)
				} else {
					exprs := make([]int, len(fn.Expressions))
					for i, expr := range fn.Expressions {
						exprIdx := serializeExpression(expr, s)
						if fn.CurrentExpression == expr {
							// sFn.CurrentExpressionOffset = int32(exprIdx)
							sFn.CurrentExpressionOffset = int32(i)
						}
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

// Serialize ...
func Serialize(prgrm *CXProgram, split int) (byts []byte) {
	// prgrm.PrintProgram()

	s := sAll{}
	initSerialization(prgrm, &s)

	var fnCounter int32
	var strctCounter int32
	splitSerialize(prgrm, &s, &fnCounter, &strctCounter, 0, split)
	splitSerialize(prgrm, &s, &fnCounter, &strctCounter, split, len(prgrm.Packages))

	// program
	serializeProgram(prgrm, &s)

	s.Index = sIndex{}
	sIdx := &s.Index

	// assigning relative offset

	idxSize := encoder.Size(s.Index)
	prgrmSize := encoder.Size(s.Program)
	callSize := encoder.Size(s.Calls)
	pkgSize := encoder.Size(s.Packages)
	strctSize := encoder.Size(s.Structs)
	fnSize := encoder.Size(s.Functions)
	exprSize := encoder.Size(s.Expressions)
	argSize := encoder.Size(s.Arguments)
	intSize := encoder.Size(s.Integers)

	// assigning absolute offset
	sIdx.ProgramOffset += int32(idxSize)
	sIdx.CallsOffset += sIdx.ProgramOffset + int32(prgrmSize)
	sIdx.PackagesOffset += sIdx.CallsOffset + int32(callSize)
	sIdx.StructsOffset += sIdx.PackagesOffset + int32(pkgSize)
	sIdx.FunctionsOffset += sIdx.StructsOffset + int32(strctSize)
	sIdx.ExpressionsOffset += sIdx.FunctionsOffset + int32(fnSize)
	sIdx.ArgumentsOffset += sIdx.ExpressionsOffset + int32(exprSize)
	sIdx.IntegersOffset += sIdx.ArgumentsOffset + int32(argSize)
	sIdx.NamesOffset += sIdx.IntegersOffset + int32(intSize)
	sIdx.MemoryOffset += sIdx.NamesOffset + int32(len(s.Names))

	// serializing everything
	byts = append(byts, encoder.Serialize(s.Index)...)
	byts = append(byts, encoder.Serialize(s.Program)...)
	byts = append(byts, encoder.Serialize(s.Calls)...)
	byts = append(byts, encoder.Serialize(s.Packages)...)
	byts = append(byts, encoder.Serialize(s.Structs)...)
	byts = append(byts, encoder.Serialize(s.Functions)...)
	byts = append(byts, encoder.Serialize(s.Expressions)...)
	byts = append(byts, encoder.Serialize(s.Arguments)...)
	byts = append(byts, encoder.Serialize(s.Integers)...)
	byts = append(byts, s.Names...)
	byts = append(byts, s.Memory...)

	return byts
}

func opSerialize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	_ = inp1

	var slcOff int
	byts := Serialize(PROGRAM, 0)
	for _, b := range byts {
		slcOff = WriteToSlice(slcOff, []byte{b})
	}

	WriteI32(out1Offset, int32(slcOff))
}

func opDeserialize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp := expr.Inputs[0]

	inpOffset := GetFinalOffset(fp, inp)

	off := mustDeserializeI32(PROGRAM.Memory[inpOffset : inpOffset+TYPE_POINTER_SIZE])

	_l := PROGRAM.Memory[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
	l := mustDeserializeI32(_l[4:8])

	Deserialize(PROGRAM.Memory[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+l]) // BUG : should be l * elt.TotalSize ?
}

func dsName(off int32, size int32, s *sAll) string {
	if size < 1 {
		return ""
	}

	var name string
	mustDeserializeRaw(s.Names[off:off+size], &name)

	return name
}

func dsPackages(s *sAll, prgrm *CXProgram) {
	var fnCounter int32
	var strctCounter int32

	for i, sPkg := range s.Packages {
		// initializing packages with their names,
		// empty functions, structs, imports and globals
		// and current function and struct
		pkg := CXPackage{}
		prgrm.Packages[i] = &pkg

		pkg.Name = dsName(sPkg.NameOffset, sPkg.NameSize, s)

		if sPkg.ImportsSize > 0 {
			prgrm.Packages[i].Imports = make([]*CXPackage, sPkg.ImportsSize)
		}

		if sPkg.FunctionsSize > 0 {
			prgrm.Packages[i].Functions = make([]*CXFunction, sPkg.FunctionsSize)

			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				var fn CXFunction
				fn.Name = dsName(sFn.NameOffset, sFn.NameSize, s)
				prgrm.Packages[i].Functions[j] = &fn
			}
		}

		if sPkg.StructsSize > 0 {
			prgrm.Packages[i].Structs = make([]*CXStruct, sPkg.StructsSize)

			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				var strct CXStruct
				strct.Name = dsName(sStrct.NameOffset, sStrct.NameSize, s)
				prgrm.Packages[i].Structs[j] = &strct
			}
		}

		if sPkg.GlobalsSize > 0 {
			prgrm.Packages[i].Globals = make([]*CXArgument, sPkg.GlobalsSize)
		}

		// // CurrentFunction
		// if sPkg.FunctionsSize > 0 {
		// 	prgrm.Packages[i].CurrentFunction = prgrm.Packages[i].Functions[sPkg.CurrentFunctionOffset-fnCounter]
		// }

		// CurrentStruct
		if sPkg.StructsSize > 0 {
			prgrm.Packages[i].CurrentStruct = prgrm.Packages[i].Structs[sPkg.CurrentStructOffset-strctCounter]
		}

		fnCounter += sPkg.FunctionsSize
		strctCounter += sPkg.StructsSize
	}

	// imports
	for i, sPkg := range s.Packages {
		if sPkg.ImportsSize > 0 {
			// getting indexes of imports
			idxs := dsIntegers(sPkg.ImportsOffset, sPkg.ImportsSize, s)

			for j, idx := range idxs {
				prgrm.Packages[i].Imports[j] = getImport(&s.Packages[idx], s, prgrm)
			}
		}
	}

	// globals
	for i, sPkg := range s.Packages {
		if sPkg.GlobalsSize > 0 {
			prgrm.Packages[i].Globals = dsArguments(sPkg.GlobalsOffset, sPkg.GlobalsSize, s, prgrm)
		}
	}

	// structs
	for i, sPkg := range s.Packages {
		if sPkg.StructsSize > 0 {
			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				dsStruct(&sStrct, prgrm.Packages[i].Structs[j], s, prgrm)
			}
		}
	}

	// functions
	for i, sPkg := range s.Packages {
		if sPkg.FunctionsSize > 0 {
			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				dsFunction(&sFn, prgrm.Packages[i].Functions[j], s, prgrm)
			}
		}
	}

	// current package
	prgrm.CurrentPackage = prgrm.Packages[s.Program.CurrentPackageOffset]
}

func dsStruct(sStrct *sStruct, strct *CXStruct, s *sAll, prgrm *CXProgram) {
	strct.Name = dsName(sStrct.NameOffset, sStrct.NameSize, s)
	strct.Fields = dsArguments(sStrct.FieldsOffset, sStrct.FieldsSize, s, prgrm)
	strct.Size = int(sStrct.Size)
	strct.Package = prgrm.Packages[sStrct.PackageOffset]
}

func dsArguments(off int32, size int32, s *sAll, prgrm *CXProgram) []*CXArgument {
	if size < 1 {
		return nil
	}

	// getting indexes of arguments
	idxs := dsIntegers(off, size, s)

	// sArgs := s.Arguments[off : off + size]
	args := make([]*CXArgument, size)
	for i, idx := range idxs {
		args[i] = dsArgument(&s.Arguments[idx], s, prgrm)
	}
	return args
}

func getCustomType(sArg *sArgument, s *sAll, prgrm *CXProgram) *CXStruct {
	if sArg.CustomTypeOffset < 0 {
		return nil
	}

	customTypePkg := prgrm.Packages[s.Structs[sArg.CustomTypeOffset].PackageOffset]
	sStrct := s.Structs[sArg.CustomTypeOffset]
	customTypeName := dsName(sStrct.NameOffset, sStrct.NameSize, s)

	for _, strct := range customTypePkg.Structs {
		if strct.Name == customTypeName {
			return strct
		}
	}

	return nil
}

func dsArgument(sArg *sArgument, s *sAll, prgrm *CXProgram) *CXArgument {
	var arg CXArgument
	arg.Name = dsName(sArg.NameOffset, sArg.NameSize, s)
	arg.Type = int(sArg.Type)

	arg.CustomType = getCustomType(sArg, s, prgrm)

	arg.Size = int(sArg.Size)
	arg.TotalSize = int(sArg.TotalSize)
	arg.Offset = int(sArg.Offset)
	arg.IndirectionLevels = int(sArg.IndirectionLevels)
	arg.DereferenceLevels = int(sArg.DereferenceLevels)
	arg.PassBy = int(sArg.PassBy)

	arg.DereferenceOperations = dsIntegers(sArg.DereferenceOperationsOffset, sArg.DereferenceOperationsSize, s)
	arg.DeclarationSpecifiers = dsIntegers(sArg.DeclarationSpecifiersOffset, sArg.DeclarationSpecifiersSize, s)

	arg.IsSlice = dsBool(sArg.IsSlice)
	arg.IsArray = dsBool(sArg.IsArray)
	arg.IsArrayFirst = dsBool(sArg.IsArrayFirst)
	arg.IsPointer = dsBool(sArg.IsPointer)
	arg.IsReference = dsBool(sArg.IsReference)
	arg.IsDereferenceFirst = dsBool(sArg.IsDereferenceFirst)
	arg.IsStruct = dsBool(sArg.IsStruct)
	arg.IsRest = dsBool(sArg.IsRest)
	arg.IsLocalDeclaration = dsBool(sArg.IsLocalDeclaration)
	arg.IsShortDeclaration = dsBool(sArg.IsShortDeclaration)
	arg.PreviouslyDeclared = dsBool(sArg.PreviouslyDeclared)
	arg.DoesEscape = dsBool(sArg.DoesEscape)

	arg.Lengths = dsIntegers(sArg.LengthsOffset, sArg.LengthsSize, s)
	arg.Indexes = dsArguments(sArg.IndexesOffset, sArg.IndexesSize, s, prgrm)
	arg.Fields = dsArguments(sArg.FieldsOffset, sArg.FieldsSize, s, prgrm)
	arg.Inputs = dsArguments(sArg.InputsOffset, sArg.InputsSize, s, prgrm)
	arg.Outputs = dsArguments(sArg.OutputsOffset, sArg.OutputsSize, s, prgrm)

	arg.Package = prgrm.Packages[sArg.PackageOffset]

	return &arg
}

func getOperator(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXFunction {
	if sExpr.OperatorOffset < 0 {
		return nil
	}

	opPkg := prgrm.Packages[s.Functions[sExpr.OperatorOffset].PackageOffset]
	sOp := s.Functions[sExpr.OperatorOffset]
	opName := dsName(sOp.NameOffset, sOp.NameSize, s)

	for _, fn := range opPkg.Functions {
		if fn.Name == opName {
			return fn
		}
	}

	return nil
}

func getImport(sImp *sPackage, s *sAll, prgrm *CXProgram) *CXPackage {
	impName := dsName(sImp.NameOffset, sImp.NameSize, s)

	for _, pkg := range prgrm.Packages {
		if pkg.Name == impName {
			return pkg
		}
	}

	return nil
}

func getFunction(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXFunction {
	if sExpr.FunctionOffset < 0 {
		return nil
	}

	fnPkg := prgrm.Packages[s.Functions[sExpr.FunctionOffset].PackageOffset]
	sFn := s.Functions[sExpr.FunctionOffset]
	fnName := dsName(sFn.NameOffset, sFn.NameSize, s)

	for _, fn := range fnPkg.Functions {
		if fn.Name == fnName {
			return fn
		}
	}

	return nil
}

func dsExpressions(off int32, size int32, s *sAll, prgrm *CXProgram) []*CXExpression {
	if size < 1 {
		return nil
	}

	// getting indexes of expressions
	idxs := dsIntegers(off, size, s)

	// sExprs := s.Expressions[off : off + size]
	exprs := make([]*CXExpression, size)
	for i, idx := range idxs {
		exprs[i] = dsExpression(&s.Expressions[idx], s, prgrm)
	}
	return exprs
}

func dsExpression(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXExpression {
	var expr CXExpression

	if dsBool(sExpr.IsNative) {
		expr.Operator = Natives[int(sExpr.OpCode)]
	} else {
		expr.Operator = getOperator(sExpr, s, prgrm)
	}

	expr.Inputs = dsArguments(sExpr.InputsOffset, sExpr.InputsSize, s, prgrm)
	expr.Outputs = dsArguments(sExpr.OutputsOffset, sExpr.OutputsSize, s, prgrm)

	expr.Label = dsName(sExpr.LabelOffset, sExpr.LabelSize, s)

	expr.ThenLines = int(sExpr.ThenLines)
	expr.ElseLines = int(sExpr.ElseLines)
	expr.ScopeOperation = int(sExpr.ScopeOperation)

	expr.IsMethodCall = dsBool(sExpr.IsMethodCall)
	expr.IsStructLiteral = dsBool(sExpr.IsStructLiteral)
	expr.IsArrayLiteral = dsBool(sExpr.IsArrayLiteral)
	expr.IsUndType = dsBool(sExpr.IsUndType)
	expr.IsBreak = dsBool(sExpr.IsBreak)
	expr.IsContinue = dsBool(sExpr.IsContinue)

	expr.Function = getFunction(sExpr, s, prgrm)
	expr.Package = prgrm.Packages[sExpr.PackageOffset]

	return &expr
}

func dsFunction(sFn *sFunction, fn *CXFunction, s *sAll, prgrm *CXProgram) {
	fn.Name = dsName(sFn.NameOffset, sFn.NameSize, s)
	fn.Inputs = dsArguments(sFn.InputsOffset, sFn.InputsSize, s, prgrm)
	fn.Outputs = dsArguments(sFn.OutputsOffset, sFn.OutputsSize, s, prgrm)
	fn.ListOfPointers = dsArguments(sFn.ListOfPointersOffset, sFn.ListOfPointersSize, s, prgrm)
	fn.Expressions = dsExpressions(sFn.ExpressionsOffset, sFn.ExpressionsSize, s, prgrm)
	fn.Size = int(sFn.Size)
	fn.Length = int(sFn.Length)

	if sFn.CurrentExpressionOffset > 0 {
		fn.CurrentExpression = fn.Expressions[sFn.CurrentExpressionOffset]
	}

	fn.Package = prgrm.Packages[sFn.PackageOffset]
}

func dsBool(val int32) bool {
	return val == 1
}

func dsIntegers(off int32, size int32, s *sAll) []int {
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
func initDeserialization(prgrm *CXProgram, s *sAll) {
	prgrm.Memory = s.Memory
	prgrm.Packages = make([]*CXPackage, len(s.Packages))
	prgrm.CallStack = make([]CXCall, CALLSTACK_SIZE)
	prgrm.HeapStartsAt = int(s.Program.HeapStartsAt)
	prgrm.HeapPointer = int(s.Program.HeapPointer)
	prgrm.StackSize = int(s.Program.StackSize)
	prgrm.HeapSize = int(s.Program.HeapSize)
	prgrm.BCPackageCount = int(s.Program.BCPackageCount)
	prgrm.Version = dsName(s.Program.VersionOffset, s.Program.VersionSize, s)

	dsPackages(s, prgrm)
}

// Deserialize deserializes a serialized CX program back to its golang struct representation.
func Deserialize(byts []byte) (prgrm *CXProgram) {
	prgrm = &CXProgram{}
	idxSize := encoder.Size(sIndex{})

	var s sAll

	mustDeserializeRaw(byts[:idxSize], &s.Index)
	mustDeserializeRaw(byts[s.Index.ProgramOffset:s.Index.CallsOffset], &s.Program)
	mustDeserializeRaw(byts[s.Index.CallsOffset:s.Index.PackagesOffset], &s.Calls)
	mustDeserializeRaw(byts[s.Index.PackagesOffset:s.Index.StructsOffset], &s.Packages)
	mustDeserializeRaw(byts[s.Index.StructsOffset:s.Index.FunctionsOffset], &s.Structs)
	mustDeserializeRaw(byts[s.Index.FunctionsOffset:s.Index.ExpressionsOffset], &s.Functions)
	mustDeserializeRaw(byts[s.Index.ExpressionsOffset:s.Index.ArgumentsOffset], &s.Expressions)
	mustDeserializeRaw(byts[s.Index.ArgumentsOffset:s.Index.IntegersOffset], &s.Arguments)
	mustDeserializeRaw(byts[s.Index.IntegersOffset:s.Index.NamesOffset], &s.Integers)
	s.Names = byts[s.Index.NamesOffset:s.Index.MemoryOffset]
	s.Memory = byts[s.Index.MemoryOffset:]

	initDeserialization(prgrm, &s)

	// prgrm.PrintProgram()

	return prgrm
}

// CopyProgramState copies the program state from `prgrm1` to `prgrm2`.
func CopyProgramState(sPrgrm1, sPrgrm2 *[]byte) {
	idxSize := encoder.Size(sIndex{})

	var index1 sIndex
	var index2 sIndex

	mustDeserializeRaw((*sPrgrm1)[:idxSize], &index1)
	mustDeserializeRaw((*sPrgrm2)[:idxSize], &index2)

	var prgrm1Info sProgram
	mustDeserializeRaw((*sPrgrm1)[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info sProgram
	mustDeserializeRaw((*sPrgrm2)[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

	// the stack segment should be 0 for prgrm1, but just in case
	var prgrmState []byte
	prgrmState = append(prgrmState, make([]byte, prgrm2Info.StackSize)...)
	// We are only interested on extracting the data segment
	prgrmState = append(prgrmState, (*sPrgrm1)[index1.NamesOffset+prgrm1Info.StackSize:index1.NamesOffset+prgrm1Info.StackSize+(prgrm2Info.HeapStartsAt-prgrm2Info.StackSize)]...)

	for i, byt := range prgrmState {
		(*sPrgrm2)[i+int(index2.MemoryOffset)] = byt
	}
}

// updateSerializedSize updates the header of each of the serialized parts of a CX program. For example, if in a full CX program there were 5 packages and after extracting the transaction or blockchain parts of it, there are now 3 packages, updateSerializedSize updates this size in the header of the serialization.
func updateSerializedSize(byts *[]byte, off1, off2 int32, n int) {
	if len((*byts)[off1:off2]) == 0 {
		return
	}
	WriteMemI32(*byts, int(off1), int32((off2-off1-4)/int32(n)))
}

// ExtractBlockchainProgram extracts the blockchain program from `sPrgrm2` by removing the contents of `sPrgrm1` from `sPrgrm2`. TxnPrgrm = sPrgrm2 - sPrgrm1.
func ExtractBlockchainProgram(sPrgrm1, sPrgrm2 []byte) []byte {
	idxSize := encoder.Size(sIndex{})

	var index1 sIndex
	var index2 sIndex

	mustDeserializeRaw(sPrgrm1[:idxSize], &index1)
	mustDeserializeRaw(sPrgrm2[:idxSize], &index2)

	var prgrm1Info sProgram
	mustDeserializeRaw(sPrgrm1[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info sProgram
	mustDeserializeRaw(sPrgrm2[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

	var extracted []byte
	// must match the index from sPrgrm1
	extracted = append(extracted, sPrgrm1[:index1.ProgramOffset]...)

	// Program
	var sPrgrm sProgram
	mustDeserializeRaw(sPrgrm1[index1.ProgramOffset:index1.CallsOffset], &sPrgrm)
	// We need the heap pointer calculated after running the program, which is
	// present in `sPrgrm2`.
	sPrgrm.HeapPointer = prgrm2Info.HeapPointer
	extracted = append(extracted, encoder.Serialize(sPrgrm)...)

	extracted = append(extracted, sPrgrm2[index2.CallsOffset:index2.CallsOffset+(index1.PackagesOffset-index1.CallsOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.PackagesOffset:index2.PackagesOffset+(index1.StructsOffset-index1.PackagesOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.StructsOffset:index2.StructsOffset+(index1.FunctionsOffset-index1.StructsOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.FunctionsOffset:index2.FunctionsOffset+(index1.ExpressionsOffset-index1.FunctionsOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.ExpressionsOffset:index2.ExpressionsOffset+(index1.ArgumentsOffset-index1.ExpressionsOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.ArgumentsOffset:index2.ArgumentsOffset+(index1.IntegersOffset-index1.ArgumentsOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.IntegersOffset:index2.IntegersOffset+(index1.NamesOffset-index1.IntegersOffset)]...)
	extracted = append(extracted, sPrgrm2[index2.NamesOffset:index2.NamesOffset+(index1.MemoryOffset-index1.NamesOffset)]...)

	// We were also simulating an empty stack, but it doesn't make sense now.
	// We'll need to store the stack when we add the ability to pause CX chains and update the program state with the paused state.
	prgrm2DataStart := index2.MemoryOffset + prgrm2Info.StackSize
	prgrm1DataSize := prgrm1Info.HeapStartsAt - prgrm1Info.StackSize
	prgrm2HeapStart := index2.MemoryOffset + prgrm2Info.HeapStartsAt

	// Adding data segment.
	extracted = append(extracted, sPrgrm2[prgrm2DataStart:prgrm2DataStart+prgrm1DataSize]...)
	// Adding heap segment.
	extracted = append(extracted, sPrgrm2[prgrm2HeapStart:prgrm2HeapStart+prgrm2Info.HeapPointer]...)

	// correcting sizes
	updateSerializedSize(&extracted, index1.CallsOffset, index1.PackagesOffset, int(encoder.Size(sCall{})))
	updateSerializedSize(&extracted, index1.PackagesOffset, index1.StructsOffset, int(encoder.Size(sPackage{})))
	updateSerializedSize(&extracted, index1.StructsOffset, index1.FunctionsOffset, int(encoder.Size(sStruct{})))
	updateSerializedSize(&extracted, index1.FunctionsOffset, index1.ExpressionsOffset, int(encoder.Size(sFunction{})))
	updateSerializedSize(&extracted, index1.ExpressionsOffset, index1.ArgumentsOffset, int(encoder.Size(sExpression{})))
	updateSerializedSize(&extracted, index1.ArgumentsOffset, index1.IntegersOffset, int(encoder.Size(sArgument{})))
	updateSerializedSize(&extracted, index1.IntegersOffset, index1.NamesOffset, int(encoder.Size(int32(0))))

	return extracted
}

// ExtractTransactionProgram extracts the transaction code (serialized) from a full CX program.
func ExtractTransactionProgram(sPrgrm1, sPrgrm2 []byte) []byte {
	idxSize := encoder.Size(sIndex{})

	var index1 sIndex
	var index2 sIndex

	mustDeserializeRaw(sPrgrm1[:idxSize], &index1)
	mustDeserializeRaw(sPrgrm2[:idxSize], &index2)

	var prgrm1Info sProgram
	mustDeserializeRaw(sPrgrm1[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info sProgram
	mustDeserializeRaw(sPrgrm2[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

	var extracted []byte
	// must match the index from sPrgrm2
	extracted = append(extracted, sPrgrm2[:index2.ProgramOffset]...)
	extracted = append(extracted, sPrgrm2[index2.ProgramOffset:index2.CallsOffset]...)
	extracted = append(extracted, sPrgrm2[index2.CallsOffset+(index1.PackagesOffset-index1.CallsOffset):index2.PackagesOffset]...)
	extracted = append(extracted, sPrgrm2[index2.PackagesOffset+(index1.StructsOffset-index1.PackagesOffset):index2.StructsOffset]...)
	extracted = append(extracted, sPrgrm2[index2.StructsOffset+(index1.FunctionsOffset-index1.StructsOffset):index2.FunctionsOffset]...)
	extracted = append(extracted, sPrgrm2[index2.FunctionsOffset+(index1.ExpressionsOffset-index1.FunctionsOffset):index2.ExpressionsOffset]...)
	extracted = append(extracted, sPrgrm2[index2.ExpressionsOffset+(index1.ArgumentsOffset-index1.ExpressionsOffset):index2.ArgumentsOffset]...)
	extracted = append(extracted, sPrgrm2[index2.ArgumentsOffset+(index1.IntegersOffset-index1.ArgumentsOffset):index2.IntegersOffset]...)
	extracted = append(extracted, sPrgrm2[index2.IntegersOffset+(index1.NamesOffset-index1.IntegersOffset):index2.NamesOffset]...)
	extracted = append(extracted, sPrgrm2[index2.NamesOffset+(index1.MemoryOffset-index1.NamesOffset):index2.MemoryOffset]...)

	// Calculating where the data segment starts and its sizes in `sPrgrm1` and `sPrgrm2`.
	// In this case, the heap segment of the transaction code should not be appended.
	// The transaction code heap should only be auxiliary in the process of updating
	// the CX chain program state.
	prgrm2DataStart := index2.MemoryOffset + prgrm2Info.StackSize
	prgrm1DataSize := prgrm1Info.HeapStartsAt - prgrm1Info.StackSize
	prgrm2DataSize := prgrm2Info.HeapStartsAt - prgrm2Info.StackSize

	// Adding data segment.
	extracted = append(extracted, sPrgrm2[prgrm2DataStart+prgrm1DataSize:prgrm2DataStart+prgrm1DataSize+(prgrm2DataSize-prgrm1DataSize)]...)

	return extracted
}

// MergeTransactionAndBlockchain merges the serialized CX programs that represent a transaction and the program state stored on the blockchain.
func MergeTransactionAndBlockchain(sPrgrm1, sPrgrm2 []byte) []byte {
	idxSize := encoder.Size(sIndex{})

	var index1 sIndex
	var index2 sIndex

	mustDeserializeRaw(sPrgrm1[:idxSize], &index1)
	mustDeserializeRaw(sPrgrm2[:idxSize], &index2)

	var prgrm1Info sProgram
	mustDeserializeRaw(sPrgrm1[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info sProgram
	mustDeserializeRaw(sPrgrm2[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

	var acc int32
	var s int32
	var merged []byte

	// Index
	merged = append(merged, sPrgrm2[:index2.ProgramOffset]...)
	acc = index2.ProgramOffset

	// Program
	var sPrgrm sProgram
	mustDeserializeRaw(sPrgrm2[index2.ProgramOffset:index2.CallsOffset], &sPrgrm)
	// We need to use the heap pointer from the CX chain program state, which is
	// represented by `sPrgrm1`.
	sPrgrm.HeapPointer = prgrm1Info.HeapPointer
	merged = append(merged, encoder.Serialize(sPrgrm)...)
	acc += index2.CallsOffset - index2.ProgramOffset

	// Calls
	s = (index2.PackagesOffset - index2.CallsOffset) - (index1.PackagesOffset - index1.CallsOffset)
	merged = append(merged, sPrgrm1[index1.CallsOffset:index1.PackagesOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Packages
	s = (index2.StructsOffset - index2.PackagesOffset) - (index1.StructsOffset - index1.PackagesOffset)
	merged = append(merged, sPrgrm1[index1.PackagesOffset:index1.StructsOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Structs
	s = (index2.FunctionsOffset - index2.StructsOffset) - (index1.FunctionsOffset - index1.StructsOffset)
	merged = append(merged, sPrgrm1[index1.StructsOffset:index1.FunctionsOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Functions
	s = (index2.ExpressionsOffset - index2.FunctionsOffset) - (index1.ExpressionsOffset - index1.FunctionsOffset)
	merged = append(merged, sPrgrm1[index1.FunctionsOffset:index1.ExpressionsOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Expressions
	s = (index2.ArgumentsOffset - index2.ExpressionsOffset) - (index1.ArgumentsOffset - index1.ExpressionsOffset)
	merged = append(merged, sPrgrm1[index1.ExpressionsOffset:index1.ArgumentsOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Arguments
	s = (index2.IntegersOffset - index2.ArgumentsOffset) - (index1.IntegersOffset - index1.ArgumentsOffset)
	merged = append(merged, sPrgrm1[index1.ArgumentsOffset:index1.IntegersOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Integers
	s = (index2.NamesOffset - index2.IntegersOffset) - (index1.NamesOffset - index1.IntegersOffset)
	merged = append(merged, sPrgrm1[index1.IntegersOffset:index1.NamesOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Names
	s = (index2.MemoryOffset - index2.NamesOffset) - (index1.MemoryOffset - index1.NamesOffset)
	merged = append(merged, sPrgrm1[index1.NamesOffset:index1.MemoryOffset]...)
	merged = append(merged, sPrgrm2[acc:acc+s]...)
	acc += s

	// Memory
	// For now we need to create an empty stack so we can run the merged blockchain and transaction codes.
	merged = append(merged, make([]byte, prgrm2Info.StackSize)...)
	// We're not incrementing `acc` with stack size because we're ignoring that memory segment for now.
	// acc += prgrm2Info.StackSize

	// prgrm1DataStart := index1.MemoryOffset+prgrm1Info.StackSize
	prgrm1DataStart := index1.MemoryOffset
	prgrm1DataSize := prgrm1Info.HeapStartsAt - prgrm1Info.StackSize

	// prgrm2DataStart := index2.MemoryOffset+prgrm2Info.StackSize
	prgrm2DataSize := prgrm2Info.HeapStartsAt - prgrm2Info.StackSize

	s = prgrm2DataSize - prgrm1DataSize

	bcDataSegment := sPrgrm1[prgrm1DataStart : prgrm1DataStart+prgrm1DataSize]
	txnDataSegment := sPrgrm2[acc : acc+s]

	// Data segments from blockchain and transaction codes.
	merged = append(merged, bcDataSegment...)
	merged = append(merged, txnDataSegment...)

	// Adding heap segment.
	bcHeapSegment := sPrgrm1[index1.MemoryOffset+prgrm1DataSize : index1.MemoryOffset+prgrm1DataSize+prgrm1Info.HeapPointer]
	merged = append(merged, bcHeapSegment...)

	// correcting sizes
	updateSerializedSize(&merged, index2.CallsOffset, index2.PackagesOffset, int(encoder.Size(sCall{})))
	updateSerializedSize(&merged, index2.PackagesOffset, index2.StructsOffset, int(encoder.Size(sPackage{})))
	updateSerializedSize(&merged, index2.StructsOffset, index2.FunctionsOffset, int(encoder.Size(sStruct{})))
	updateSerializedSize(&merged, index2.FunctionsOffset, index2.ExpressionsOffset, int(encoder.Size(sFunction{})))
	updateSerializedSize(&merged, index2.ExpressionsOffset, index2.ArgumentsOffset, int(encoder.Size(sExpression{})))
	updateSerializedSize(&merged, index2.ArgumentsOffset, index2.IntegersOffset, int(encoder.Size(sArgument{})))
	updateSerializedSize(&merged, index2.IntegersOffset, index2.NamesOffset, int(encoder.Size(int32(0))))

	return merged
}

// MergePrograms merges `prgrm1` and `prgrm2`, favoring `prgrm1` (if both have a package with the same name, `prgrm1`'s is used). Note: `prgrm2` is permanently altered.
func MergePrograms(prgrm1, prgrm2 *CXProgram) *CXProgram {
	for _, pkg := range prgrm1.Packages {
		// We're always going to keep prgrm2's main
		if pkg.Name == MAIN_PKG {
			continue
		}
		if dupPkg, err := prgrm2.GetPackage(pkg.Name); err == nil {
			// Then it's duplicated and we need to replace it by prgrm1's
			*dupPkg = *pkg
		} else {
			prgrm2.AddPackage(pkg)
		}
	}

	DataOffset := prgrm1.HeapStartsAt
	for _, pkg := range prgrm2.Packages {
		for _, glbl := range pkg.Globals {
			glbl.Offset += DataOffset
		}

		for _, fn := range pkg.Functions {
			for _, expr := range fn.Expressions {
				for _, inp := range expr.Inputs {
					if inp.Offset > prgrm2.StackSize {
						inp.Offset += DataOffset
					}
				}

				for _, out := range expr.Inputs {
					if out.Offset > prgrm2.StackSize {
						out.Offset += DataOffset
					}
				}
			}
		}
	}

	return prgrm2
}

// GetSerializedMemoryOffset returns the offset at which the memory of a serialized CX program starts.
func GetSerializedMemoryOffset(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	mustDeserializeRaw(sPrgrm[:idxSize], &index)
	return int(index.MemoryOffset)
}

// GetSerializedStackSize returns the stack size of a serialized CX program starts.
func GetSerializedStackSize(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	mustDeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo sProgram
	mustDeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.StackSize)
}

// GetSerializedDataSize returns the size of the data segment of a serialized CX program.
func GetSerializedDataSize(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	mustDeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo sProgram
	mustDeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.HeapStartsAt - prgrmInfo.StackSize)
}
