package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

type serializedCXProgramIndex struct {
	ProgramOffset     int64
	CallsOffset       int64
	PackagesOffset    int64
	StructsOffset     int64
	FunctionsOffset   int64
	ExpressionsOffset int64
	ArgumentsOffset   int64
	IntegersOffset    int64
	StringsOffset     int64
	MemoryOffset      int64
}

type serializedProgram struct {
	PackagesOffset       int64
	PackagesSize         int64
	CurrentPackageOffset int64

	InputsOffset int64
	InputsSize   int64

	OutputsOffset int64
	OutputsSize   int64

	CallStackOffset int64
	CallStackSize   int64

	CallCounter int64

	MemoryOffset int64
	MemorySize   int64

	StackPointer int64
	StackSize    int64

	HeapPointer  int64 //HeapPointer is probably related to HeapStartsAt
	HeapStartsAt int64
	HeapSize     int64

	Terminated int64

	VersionOffset int64
	VersionSize   int64
}

type serializedCall struct {
	OperatorOffset int64
	Line           int64
	FramePointer   int64
}

type serializedPackage struct {
	NameOffset            int64
	NameSize              int64
	ImportsOffset         int64
	ImportsSize           int64
	StructsOffset         int64
	StructsSize           int64
	GlobalsOffset         int64
	GlobalsSize           int64
	FunctionsOffset       int64
	FunctionsSize         int64
	CurrentFunctionOffset int64
	CurrentStructOffset   int64
}

type serializedStruct struct {
	NameOffset   int64
	NameSize     int64
	FieldsOffset int64
	FieldsSize   int64

	Size int64

	PackageOffset int64
}

type serializedFunction struct {
	NameOffset        int64
	NameSize          int64
	InputsOffset      int64
	InputsSize        int64
	OutputsOffset     int64
	OutputsSize       int64
	ExpressionsOffset int64
	ExpressionsSize   int64
	Size              int64
	Length            int64

	ListOfPointersOffset int64
	ListOfPointersSize   int64

	// We're going to determine this when procesing the expressions. Check serializedExpression type
	// IsAtomic                        int64
	// OpCode                          int64

	CurrentExpressionOffset int64
	PackageOffset           int64
}

type serializedExpression struct {
	OperatorOffset int64
	// we add these two fields here so we don't add every native serializedFunction to the serialization
	// the CX runtime already knows about the natives properties. We just need the code if IsNative = true
	IsNative int64
	OpCode   int64

	InputsOffset  int64
	InputsSize    int64
	OutputsOffset int64
	OutputsSize   int64

	LabelOffset int64
	LabelSize   int64
	ThenLines   int64
	ElseLines   int64

	ScopeOperation int64

	IsMethodCall    int64
	IsStructLiteral int64
	IsArrayLiteral  int64
	IsUndType       int64
	IsBreak         int64
	IsContinue      int64

	FunctionOffset int64
	PackageOffset  int64
}

type serializedArgument struct {
	NameOffset       int64
	NameSize         int64
	Type             int64
	StructTypeOffset int64 //WTF IS A CUSTOM TYPE!?
	Size             int64
	TotalSize        int64

	Offset int64

	IndirectionLevels           int64
	DereferenceLevels           int64
	DereferenceOperationsOffset int64
	DereferenceOperationsSize   int64
	DeclarationSpecifiersOffset int64
	DeclarationSpecifiersSize   int64

	IsSlice      int64
	IsArray      int64
	IsArrayFirst int64
	IsPointer    int64
	IsReference  int64

	IsDereferenceFirst int64
	IsStruct           int64
	IsRest             int64
	IsLocalDeclaration int64
	IsShortDeclaration int64
	PreviouslyDeclared int64

	PassBy     int64
	DoesEscape int64

	LengthsOffset int64
	LengthsSize   int64
	IndexesOffset int64
	IndexesSize   int64
	FieldsOffset  int64
	FieldsSize    int64
	InputsOffset  int64
	InputsSize    int64
	OutputsOffset int64
	OutputsSize   int64

	PackageOffset int64
}

type serializedCXProgram struct {
	Index   serializedCXProgramIndex
	Program serializedProgram

	Packages     []serializedPackage
	PackagesMap  map[string]int
	Structs      []serializedStruct
	StructsMap   map[string]int
	Functions    []serializedFunction
	FunctionsMap map[string]int

	Expressions []serializedExpression
	Arguments   []serializedArgument
	Calls       []serializedCall

	Strings    []byte
	StringsMap map[string]int
	Integers   []int64

	Memory []byte
}

type SerializedDataSize struct {
	Program     int `json:"program"`
	Calls       int `json:"calls"`
	Packages    int `json:"packages"`
	Structs     int `json:"structs"`
	Functions   int `json:"functions"`
	Expressions int `json:"expressions"`
	Arguments   int `json:"arguments"`
	Integers    int `json:"integers"`
	Strings     int `json:"strings"`
	Memory      int `json:"memory"`
}

func serializeString(name string, s *serializedCXProgram) (int64, int64) {
	if name == "" {
		return int64(-1), int64(-1)
	}

	size := encoder.Size(name)

	off, found := s.StringsMap[name]
	if found {
		return int64(off), int64(size)
	}
	off = len(s.Strings)
	s.Strings = append(s.Strings, encoder.Serialize(name)...)
	s.StringsMap[name] = off

	return int64(off), int64(size)
}

func indexPackage(pkg *CXPackage, s *serializedCXProgram) {
	if _, found := s.PackagesMap[pkg.Name]; !found {
		s.PackagesMap[pkg.Name] = len(s.PackagesMap)
	} else {
		panic("duplicated package in serialization process")
	}
}

func indexStruct(strct *CXStruct, s *serializedCXProgram) {
	strctName := strct.Package.Name + "." + strct.Name
	if _, found := s.StructsMap[strctName]; !found {
		s.StructsMap[strctName] = len(s.StructsMap)
	} else {
		panic("duplicated struct in serialization process")
	}
}

func indexFunction(fn *CXFunction, s *serializedCXProgram) {
	fnName := fn.Package.Name + "." + fn.Name
	if _, found := s.FunctionsMap[fnName]; !found {
		s.FunctionsMap[fnName] = len(s.FunctionsMap)
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

func serializeIntegers(ints []int, s *serializedCXProgram) (int64, int64) {
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

func serializeArgument(arg *CXArgument, s *serializedCXProgram) int {
	s.Arguments = append(s.Arguments, serializedArgument{})
	argOff := len(s.Arguments) - 1

	sNil := int64(-1)

	s.Arguments[argOff].NameOffset, s.Arguments[argOff].NameSize = serializeString(arg.Name, s)

	s.Arguments[argOff].Type = int64(arg.Type)

	if arg.CustomType == nil {
		s.Arguments[argOff].StructTypeOffset = sNil
	} else {
		strctName := arg.CustomType.Package.Name + "." + arg.CustomType.Name
		if strctOff, found := s.StructsMap[strctName]; found {
			s.Arguments[argOff].StructTypeOffset = int64(strctOff)
		} else {
			panic("struct reference not found")
		}
	}

	s.Arguments[argOff].Size = int64(arg.Size)
	s.Arguments[argOff].TotalSize = int64(arg.TotalSize)
	s.Arguments[argOff].Offset = int64(arg.Offset)
	s.Arguments[argOff].IndirectionLevels = int64(arg.IndirectionLevels)
	s.Arguments[argOff].DereferenceLevels = int64(arg.DereferenceLevels)

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
	s.Arguments[argOff].IsShortDeclaration = serializeBoolean(arg.IsShortAssignmentDeclaration)
	s.Arguments[argOff].PreviouslyDeclared = serializeBoolean(arg.PreviouslyDeclared)

	s.Arguments[argOff].PassBy = int64(arg.PassBy)
	s.Arguments[argOff].DoesEscape = serializeBoolean(arg.DoesEscape)

	s.Arguments[argOff].LengthsOffset, s.Arguments[argOff].LengthsSize = serializeIntegers(arg.Lengths, s)
	s.Arguments[argOff].IndexesOffset, s.Arguments[argOff].IndexesSize = serializeSliceOfArguments(arg.Indexes, s)
	s.Arguments[argOff].FieldsOffset, s.Arguments[argOff].FieldsSize = serializeSliceOfArguments(arg.Fields, s)
	s.Arguments[argOff].InputsOffset, s.Arguments[argOff].InputsSize = serializeSliceOfArguments(arg.Inputs, s)
	s.Arguments[argOff].OutputsOffset, s.Arguments[argOff].OutputsSize = serializeSliceOfArguments(arg.Outputs, s)

	if pkgOff, found := s.PackagesMap[arg.Package.Name]; found {
		s.Arguments[argOff].PackageOffset = int64(pkgOff)
	} else {
		panic("package reference not found")
	}

	return argOff
}

func serializeSliceOfArguments(args []*CXArgument, s *serializedCXProgram) (int64, int64) {
	if len(args) == 0 {
		return int64(-1), int64(-1)
	}
	idxs := make([]int, len(args))
	for i, arg := range args {
		idxs[i] = serializeArgument(arg, s)
	}
	return serializeIntegers(idxs, s)
}

func serializeCalls(calls []CXCall, s *serializedCXProgram) (int64, int64) {
	if len(calls) == 0 {
		return int64(-1), int64(-1)
	}
	idxs := make([]int, len(calls))
	for i, call := range calls {
		idxs[i] = serializeCall(&call, s)
	}
	return serializeIntegers(idxs, s)

}

func serializeExpression(expr *CXExpression, s *serializedCXProgram) int {
	s.Expressions = append(s.Expressions, serializedExpression{})
	exprOff := len(s.Expressions) - 1
	sExpr := &s.Expressions[exprOff]

	sNil := int64(-1)

	if expr.Operator == nil {
		// then it's a declaration
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = int64(-1)
	} else if expr.Operator.IsAtomic {
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(true)
		sExpr.OpCode = int64(expr.Operator.OpCode)
	} else {
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = sNil

		opName := expr.Operator.Package.Name + "." + expr.Operator.Name
		if opOff, found := s.FunctionsMap[opName]; found {
			sExpr.OperatorOffset = int64(opOff)
		}
	}

	sExpr.InputsOffset, sExpr.InputsSize = serializeSliceOfArguments(expr.Inputs, s)
	sExpr.OutputsOffset, sExpr.OutputsSize = serializeSliceOfArguments(expr.Outputs, s)

	sExpr.LabelOffset, sExpr.LabelSize = serializeString(expr.Label, s)
	sExpr.ThenLines = int64(expr.ThenLines)
	sExpr.ElseLines = int64(expr.ElseLines)
	sExpr.ScopeOperation = int64(expr.ScopeOperation)

	sExpr.IsMethodCall = serializeBoolean(expr.IsMethodCall)
	sExpr.IsStructLiteral = serializeBoolean(expr.IsStructLiteral)
	sExpr.IsArrayLiteral = serializeBoolean(expr.IsArrayLiteral)
	sExpr.IsUndType = serializeBoolean(expr.IsUndType)
	sExpr.IsBreak = serializeBoolean(expr.IsBreak)
	sExpr.IsContinue = serializeBoolean(expr.IsContinue)

	fnName := expr.Function.Package.Name + "." + expr.Function.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sExpr.FunctionOffset = int64(fnOff)
	} else {
		panic("function reference not found")
	}

	if pkgOff, found := s.PackagesMap[expr.Package.Name]; found {
		sExpr.PackageOffset = int64(pkgOff)
	} else {
		panic("package reference not found")
	}

	return exprOff
}

func serializeCall(call *CXCall, s *serializedCXProgram) int {
	s.Calls = append(s.Calls, serializedCall{})
	callOff := len(s.Calls) - 1
	serializedCall := &s.Calls[callOff]

	opName := call.Operator.Package.Name + "." + call.Operator.Name
	if opOff, found := s.FunctionsMap[opName]; found {
		serializedCall.OperatorOffset = int64(opOff)
		serializedCall.Line = int64(call.Line)
		serializedCall.FramePointer = int64(call.FramePointer)
	} else {
		panic("function reference not found")
	}

	return callOff
}

func serializeStructArguments(strct *CXStruct, s *serializedCXProgram) {
	strctName := strct.Package.Name + "." + strct.Name
	if strctOff, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[strctOff]
		sStrct.FieldsOffset, sStrct.FieldsSize = serializeSliceOfArguments(strct.Fields, s)
	} else {
		panic("struct reference not found")
	}
}

func serializeFunctionArguments(fn *CXFunction, s *serializedCXProgram) {
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

func serializePackageName(pkg *CXPackage, s *serializedCXProgram) {
	sPkg := &s.Packages[s.PackagesMap[pkg.Name]]
	sPkg.NameOffset, sPkg.NameSize = serializeString(pkg.Name, s) //Change Name to String
}

func serializeStructName(strct *CXStruct, s *serializedCXProgram) {
	strctName := strct.Package.Name + "." + strct.Name
	sStrct := &s.Structs[s.StructsMap[strctName]]
	sStrct.NameOffset, sStrct.NameSize = serializeString(strct.Name, s) //Change Name to String
}

func serializeFunctionName(fn *CXFunction, s *serializedCXProgram) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.NameOffset, sFn.NameSize = serializeString(fn.Name, s) //Change Name to String
	} else {
		panic("function reference not found")
	}
}

func serializePackageGlobals(pkg *CXPackage, s *serializedCXProgram) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]
		sPkg.GlobalsOffset, sPkg.GlobalsSize = serializeSliceOfArguments(pkg.Globals, s)
	} else {
		panic("package reference not found")
	}
}

func serializePackageImports(pkg *CXPackage, s *serializedCXProgram) {
	l := len(pkg.Imports)
	if l == 0 {
		s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int64(-1)
		s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int64(-1)
		return
	}
	imps := make([]int64, l)
	for i, imp := range pkg.Imports {
		if idx, found := s.PackagesMap[imp.Name]; found {
			imps[i] = int64(idx)
		} else {
			panic("import package reference not found")
		}
	}

	s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int64(len(s.Integers))
	s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int64(l)
	s.Integers = append(s.Integers, imps...)
}

func serializeStructPackage(strct *CXStruct, s *serializedCXProgram) {
	strctName := strct.Package.Name + "." + strct.Name
	if pkgOff, found := s.PackagesMap[strct.Package.Name]; found {
		if off, found := s.StructsMap[strctName]; found {
			sStrct := &s.Structs[off]
			sStrct.PackageOffset = int64(pkgOff)
		} else {
			panic("struct reference not found")
		}
	} else {
		panic("struct's package reference not found")
	}
}

func serializeFunctionPackage(fn *CXFunction, s *serializedCXProgram) {
	fnName := fn.Package.Name + "." + fn.Name
	if pkgOff, found := s.PackagesMap[fn.Package.Name]; found {
		if off, found := s.FunctionsMap[fnName]; found {
			sFn := &s.Functions[off]
			sFn.PackageOffset = int64(pkgOff)
		} else {
			panic("function reference not found")
		}
	} else {
		panic("function's package reference not found")
	}
}

func serializePackageIntegers(pkg *CXPackage, s *serializedCXProgram) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]

		if pkg.CurrentFunction == nil {
			// package has no functions
			sPkg.CurrentFunctionOffset = int64(-1)
		} else {
			currFnName := pkg.CurrentFunction.Package.Name + "." + pkg.CurrentFunction.Name

			if fnOff, found := s.FunctionsMap[currFnName]; found {
				sPkg.CurrentFunctionOffset = int64(fnOff)
			} else {
				panic("function reference not found")
			}
		}

		if pkg.CurrentStruct == nil {
			// package has no structs
			sPkg.CurrentStructOffset = int64(-1)
		} else {
			currStrctName := pkg.CurrentStruct.Package.Name + "." + pkg.CurrentStruct.Name

			if strctOff, found := s.StructsMap[currStrctName]; found {
				sPkg.CurrentStructOffset = int64(strctOff)
			} else {
				panic("struct reference not found")
			}
		}
	} else {
		panic("package reference not found")
	}
}

func serializeStructIntegers(strct *CXStruct, s *serializedCXProgram) {
	strctName := strct.Package.Name + "." + strct.Name
	if off, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[off]
		sStrct.Size = int64(strct.Size)
	} else {
		panic("struct reference not found")
	}
}

func serializeFunctionIntegers(fn *CXFunction, s *serializedCXProgram) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.Size = int64(fn.Size)
		sFn.Length = int64(fn.Length)
	} else {
		panic("function reference not found")
	}
}

// initSerialization initializes the
// container for our serialized cx program.
// Program memory is also added here to our container
// if memory is to be included.
func initSerialization(prgrm *CXProgram, s *serializedCXProgram, includeMemory bool) {
	s.PackagesMap = make(map[string]int)
	s.StructsMap = make(map[string]int)
	s.FunctionsMap = make(map[string]int)
	s.StringsMap = make(map[string]int)

	s.Calls = make([]serializedCall, prgrm.CallCounter)
	s.Packages = make([]serializedPackage, len(prgrm.Packages))

	if includeMemory {
		// s.Memory = prgrm.Memory[:PROGRAM.HeapStartsAt+PROGRAM.HeapPointer]
		s.Memory = prgrm.Memory
	}

	var numStrcts int
	var numFns int

	for _, pkg := range prgrm.Packages {
		numStrcts += len(pkg.Structs)
		numFns += len(pkg.Functions)
	}

	s.Structs = make([]serializedStruct, numStrcts)
	s.Functions = make([]serializedFunction, numFns)
	// args and exprs need to be appended as they are found
}

// serializeProgram serializes
// program of cx program.
func serializeProgram(prgrm *CXProgram, s *serializedCXProgram) {
	s.Program = serializedProgram{}
	sPrgrm := &s.Program
	sPrgrm.PackagesOffset = int64(0)
	sPrgrm.PackagesSize = int64(len(prgrm.Packages))

	if pkgOff, found := s.PackagesMap[prgrm.CurrentPackage.Name]; found {
		sPrgrm.CurrentPackageOffset = int64(pkgOff)
	} else {
		panic("package reference not found")
	}

	sPrgrm.InputsOffset, sPrgrm.InputsSize = serializeSliceOfArguments(prgrm.ProgramInput, s)
	sPrgrm.OutputsOffset, sPrgrm.OutputsSize = serializeSliceOfArguments(prgrm.ProgramOutput, s)

	sPrgrm.CallStackOffset, sPrgrm.CallStackSize = serializeCalls(prgrm.CallStack[:prgrm.CallCounter], s)

	sPrgrm.CallCounter = int64(prgrm.CallCounter)

	sPrgrm.MemoryOffset = int64(0)
	sPrgrm.MemorySize = int64(len(PROGRAM.Memory))

	sPrgrm.HeapPointer = int64(prgrm.HeapPointer)
	sPrgrm.StackPointer = int64(prgrm.StackPointer)
	sPrgrm.StackSize = int64(prgrm.StackSize)
	sPrgrm.HeapSize = int64(prgrm.HeapSize)
	sPrgrm.HeapStartsAt = int64(prgrm.HeapStartsAt)

	sPrgrm.Terminated = serializeBoolean(prgrm.Terminated)
	sPrgrm.VersionOffset, sPrgrm.VersionSize = serializeString(prgrm.Version, s)
}

// serializeCXProgramElements is used serializing CX program's
// elements (packages, structs, functions, etc.).
func serializeCXProgramElements(prgrm *CXProgram, s *serializedCXProgram) {
	var fnCounter int64
	var strctCounter int64

	// indexing packages and serializing their names
	for _, pkg := range prgrm.Packages {
		indexPackage(pkg, s)
		serializePackageName(pkg, s)
	}
	// we first needed to populate references to all packages
	// now we add the imports' references
	for _, pkg := range prgrm.Packages {
		serializePackageImports(pkg, s)
	}

	// structs
	for _, pkg := range prgrm.Packages {
		for _, strct := range pkg.Structs {
			indexStruct(strct, s)
			serializeStructName(strct, s)
			serializeStructPackage(strct, s)
			serializeStructIntegers(strct, s)
		}
	}
	// we first needed to populate references to all structs
	// now we add fields
	for _, pkg := range prgrm.Packages {
		for _, strct := range pkg.Structs {
			serializeStructArguments(strct, s)
		}
	}

	// globals
	for _, pkg := range prgrm.Packages {
		serializePackageGlobals(pkg, s)
	}

	// functions
	for _, pkg := range prgrm.Packages {
		for _, fn := range pkg.Functions {
			indexFunction(fn, s)
			serializeFunctionName(fn, s)
			serializeFunctionPackage(fn, s)
			serializeFunctionIntegers(fn, s)
			serializeFunctionArguments(fn, s)
		}
	}

	// package elements' offsets and sizes
	for _, pkg := range prgrm.Packages {
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
	for _, pkg := range prgrm.Packages {
		serializePackageIntegers(pkg, s)
	}

	// expressions
	for _, pkg := range prgrm.Packages {
		for _, fn := range pkg.Functions {
			fnName := fn.Package.Name + "." + fn.Name
			if fnOff, found := s.FunctionsMap[fnName]; found {
				sFn := &s.Functions[fnOff]

				if len(fn.Expressions) == 0 {
					sFn.ExpressionsOffset = int64(-1)
					sFn.ExpressionsSize = int64(-1)
					sFn.CurrentExpressionOffset = int64(-1)
				} else {
					exprs := make([]int, len(fn.Expressions))
					for i, expr := range fn.Expressions {
						exprIdx := serializeExpression(expr, s)
						if fn.CurrentExpression == expr {
							// sFn.CurrentExpressionOffset = int32(exprIdx)
							sFn.CurrentExpressionOffset = int64(i)
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

// SerializeCXProgram translates cx program to slice of bytes that we can save.
// These slice of bytes can then be deserialize in the future and
// be translated back to cx program.
func SerializeCXProgram(prgrm *CXProgram, includeMemory bool) (b []byte) {
	s := serializedCXProgram{}
	initSerialization(prgrm, &s, includeMemory)

	// serialize cx program's packages,
	// structs, functions, etc.
	serializeCXProgramElements(prgrm, &s)

	// serialize cx program's program
	serializeProgram(prgrm, &s)

	s.Index = serializedCXProgramIndex{}
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
	sIdx.ProgramOffset += int64(idxSize)
	sIdx.CallsOffset += sIdx.ProgramOffset + int64(prgrmSize)
	sIdx.PackagesOffset += sIdx.CallsOffset + int64(callSize)
	sIdx.StructsOffset += sIdx.PackagesOffset + int64(pkgSize)
	sIdx.FunctionsOffset += sIdx.StructsOffset + int64(strctSize)
	sIdx.ExpressionsOffset += sIdx.FunctionsOffset + int64(fnSize)
	sIdx.ArgumentsOffset += sIdx.ExpressionsOffset + int64(exprSize)
	sIdx.IntegersOffset += sIdx.ArgumentsOffset + int64(argSize)
	sIdx.StringsOffset += sIdx.IntegersOffset + int64(intSize)
	sIdx.MemoryOffset += sIdx.StringsOffset + int64(len(s.Strings))

	// serializing everything
	b = append(b, encoder.Serialize(s.Index)...)
	b = append(b, encoder.Serialize(s.Program)...)
	b = append(b, encoder.Serialize(s.Calls)...)
	b = append(b, encoder.Serialize(s.Packages)...)
	b = append(b, encoder.Serialize(s.Structs)...)
	b = append(b, encoder.Serialize(s.Functions)...)
	b = append(b, encoder.Serialize(s.Expressions)...)
	b = append(b, encoder.Serialize(s.Arguments)...)
	b = append(b, encoder.Serialize(s.Integers)...)
	b = append(b, s.Strings...)
	b = append(b, s.Memory...)

	return b
}

// SerializeDebugInfo prints the name of the serialized segment and byte size.
func SerializeDebugInfo(prgrm *CXProgram, includeMemory bool) SerializedDataSize {
	idxSize := encoder.Size(serializedCXProgramIndex{})
	var s serializedCXProgram

	bytes := SerializeCXProgram(prgrm, includeMemory)
	helper.DeserializeRaw(bytes[:idxSize], &s.Index)

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

func deserializeString(off int64, size int64, s *serializedCXProgram) string {
	if size < 1 {
		return ""
	}

	var name string
	helper.DeserializeRaw(s.Strings[off:off+size], &name)

	return name
}

func deserializePackages(s *serializedCXProgram, prgrm *CXProgram) {
	var fnCounter int64
	var strctCounter int64

	for i, sPkg := range s.Packages {
		// initializing packages with their names,
		// empty functions, structs, imports and globals
		// and current function and struct
		pkg := CXPackage{}
		prgrm.Packages[i] = &pkg

		pkg.Name = deserializeString(sPkg.NameOffset, sPkg.NameSize, s)

		if sPkg.ImportsSize > 0 {
			prgrm.Packages[i].Imports = make([]*CXPackage, sPkg.ImportsSize)
		}

		if sPkg.FunctionsSize > 0 {
			prgrm.Packages[i].Functions = make([]*CXFunction, sPkg.FunctionsSize)

			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				var fn CXFunction
				fn.Name = deserializeString(sFn.NameOffset, sFn.NameSize, s)
				prgrm.Packages[i].Functions[j] = &fn
			}
		}

		if sPkg.StructsSize > 0 {
			prgrm.Packages[i].Structs = make([]*CXStruct, sPkg.StructsSize)

			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				var strct CXStruct
				strct.Name = deserializeString(sStrct.NameOffset, sStrct.NameSize, s)
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
			idxs := deserializeIntegers(sPkg.ImportsOffset, sPkg.ImportsSize, s)

			for j, idx := range idxs {
				prgrm.Packages[i].Imports[j] = deserializePackageImport(&s.Packages[idx], s, prgrm)
			}
		}
	}

	// globals
	for i, sPkg := range s.Packages {
		if sPkg.GlobalsSize > 0 {
			prgrm.Packages[i].Globals = deserializeArguments(sPkg.GlobalsOffset, sPkg.GlobalsSize, s, prgrm)
		}
	}

	// structs
	for i, sPkg := range s.Packages {
		if sPkg.StructsSize > 0 {
			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				deserializeStruct(&sStrct, prgrm.Packages[i].Structs[j], s, prgrm)
			}
		}
	}

	// functions
	for i, sPkg := range s.Packages {
		if sPkg.FunctionsSize > 0 {
			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				deserializeFunction(&sFn, prgrm.Packages[i].Functions[j], s, prgrm)
			}
		}
	}

	// current package
	prgrm.CurrentPackage = prgrm.Packages[s.Program.CurrentPackageOffset]
}

func deserializeStruct(sStrct *serializedStruct, strct *CXStruct, s *serializedCXProgram, prgrm *CXProgram) {
	strct.Name = deserializeString(sStrct.NameOffset, sStrct.NameSize, s)
	strct.Fields = deserializeArguments(sStrct.FieldsOffset, sStrct.FieldsSize, s, prgrm)
	strct.Size = int(sStrct.Size)
	strct.Package = prgrm.Packages[sStrct.PackageOffset]
}

func deserializeArguments(off int64, size int64, s *serializedCXProgram, prgrm *CXProgram) []*CXArgument {
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

// func getStructType(sArg *serializedArgument, s *serializedCXProgram, prgrm *CXProgram) *CXStruct {
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

func deserializeArgument(sArg *serializedArgument, s *serializedCXProgram, prgrm *CXProgram) *CXArgument {
	var arg CXArgument
	arg.Name = deserializeString(sArg.NameOffset, sArg.NameSize, s)
	arg.Type = int(sArg.Type)

	//arg.CustomType = getStructType(sArg, s, prgrm)

	arg.Size = int(sArg.Size)
	arg.TotalSize = int(sArg.TotalSize)
	arg.Offset = int(sArg.Offset)
	arg.IndirectionLevels = int(sArg.IndirectionLevels)
	arg.DereferenceLevels = int(sArg.DereferenceLevels)
	arg.PassBy = int(sArg.PassBy)

	arg.DereferenceOperations = deserializeIntegers(sArg.DereferenceOperationsOffset, sArg.DereferenceOperationsSize, s)
	arg.DeclarationSpecifiers = deserializeIntegers(sArg.DeclarationSpecifiersOffset, sArg.DeclarationSpecifiersSize, s)

	arg.IsSlice = deserializeBool(sArg.IsSlice)
	arg.IsArray = deserializeBool(sArg.IsArray)
	arg.IsArrayFirst = deserializeBool(sArg.IsArrayFirst)
	arg.IsPointer = deserializeBool(sArg.IsPointer)
	arg.IsReference = deserializeBool(sArg.IsReference)
	arg.IsDereferenceFirst = deserializeBool(sArg.IsDereferenceFirst)
	arg.IsStruct = deserializeBool(sArg.IsStruct)
	arg.IsRest = deserializeBool(sArg.IsRest)
	arg.IsLocalDeclaration = deserializeBool(sArg.IsLocalDeclaration)
	arg.IsShortAssignmentDeclaration = deserializeBool(sArg.IsShortDeclaration)
	arg.PreviouslyDeclared = deserializeBool(sArg.PreviouslyDeclared)
	arg.DoesEscape = deserializeBool(sArg.DoesEscape)

	arg.Lengths = deserializeIntegers(sArg.LengthsOffset, sArg.LengthsSize, s)
	arg.Indexes = deserializeArguments(sArg.IndexesOffset, sArg.IndexesSize, s, prgrm)
	arg.Fields = deserializeArguments(sArg.FieldsOffset, sArg.FieldsSize, s, prgrm)
	arg.Inputs = deserializeArguments(sArg.InputsOffset, sArg.InputsSize, s, prgrm)
	arg.Outputs = deserializeArguments(sArg.OutputsOffset, sArg.OutputsSize, s, prgrm)

	arg.Package = prgrm.Packages[sArg.PackageOffset]

	return &arg
}

func deserializeOperator(sExpr *serializedExpression, s *serializedCXProgram, prgrm *CXProgram) *CXFunction {
	if sExpr.OperatorOffset < 0 {
		return nil
	}

	opPkg := prgrm.Packages[s.Functions[sExpr.OperatorOffset].PackageOffset]
	sOp := s.Functions[sExpr.OperatorOffset]
	opName := deserializeString(sOp.NameOffset, sOp.NameSize, s)

	for _, fn := range opPkg.Functions {
		if fn.Name == opName {
			return fn
		}
	}

	return nil
}

func deserializePackageImport(sImp *serializedPackage, s *serializedCXProgram, prgrm *CXProgram) *CXPackage {
	impName := deserializeString(sImp.NameOffset, sImp.NameSize, s)

	for _, pkg := range prgrm.Packages {
		if pkg.Name == impName {
			return pkg
		}
	}

	return nil
}

func deserializeExpressionFunction(sExpr *serializedExpression, s *serializedCXProgram, prgrm *CXProgram) *CXFunction {
	if sExpr.FunctionOffset < 0 {
		return nil
	}

	fnPkg := prgrm.Packages[s.Functions[sExpr.FunctionOffset].PackageOffset]
	sFn := s.Functions[sExpr.FunctionOffset]
	fnName := deserializeString(sFn.NameOffset, sFn.NameSize, s)

	for _, fn := range fnPkg.Functions {
		if fn.Name == fnName {
			return fn
		}
	}

	return nil
}

func deserializeExpressions(off int64, size int64, s *serializedCXProgram, prgrm *CXProgram) []*CXExpression {
	if size < 1 {
		return nil
	}

	// getting indexes of expressions
	idxs := deserializeIntegers(off, size, s)

	// sExprs := s.Expressions[off : off + size]
	exprs := make([]*CXExpression, size)
	for i, idx := range idxs {
		exprs[i] = deserializeExpression(&s.Expressions[idx], s, prgrm)
	}
	return exprs
}

func deserializeExpression(sExpr *serializedExpression, s *serializedCXProgram, prgrm *CXProgram) *CXExpression {
	var expr CXExpression

	if deserializeBool(sExpr.IsNative) {
		expr.Operator = Natives[int(sExpr.OpCode)]
	} else {
		expr.Operator = deserializeOperator(sExpr, s, prgrm)
	}

	expr.Inputs = deserializeArguments(sExpr.InputsOffset, sExpr.InputsSize, s, prgrm)
	expr.Outputs = deserializeArguments(sExpr.OutputsOffset, sExpr.OutputsSize, s, prgrm)

	expr.Label = deserializeString(sExpr.LabelOffset, sExpr.LabelSize, s)

	expr.ThenLines = int(sExpr.ThenLines)
	expr.ElseLines = int(sExpr.ElseLines)
	expr.ScopeOperation = int(sExpr.ScopeOperation)

	expr.IsMethodCall = deserializeBool(sExpr.IsMethodCall)
	expr.IsStructLiteral = deserializeBool(sExpr.IsStructLiteral)
	expr.IsArrayLiteral = deserializeBool(sExpr.IsArrayLiteral)
	expr.IsUndType = deserializeBool(sExpr.IsUndType)
	expr.IsBreak = deserializeBool(sExpr.IsBreak)
	expr.IsContinue = deserializeBool(sExpr.IsContinue)

	expr.Function = deserializeExpressionFunction(sExpr, s, prgrm)
	expr.Package = prgrm.Packages[sExpr.PackageOffset]

	return &expr
}

func deserializeFunction(sFn *serializedFunction, fn *CXFunction, s *serializedCXProgram, prgrm *CXProgram) {
	fn.Name = deserializeString(sFn.NameOffset, sFn.NameSize, s)
	fn.Inputs = deserializeArguments(sFn.InputsOffset, sFn.InputsSize, s, prgrm)
	fn.Outputs = deserializeArguments(sFn.OutputsOffset, sFn.OutputsSize, s, prgrm)
	fn.ListOfPointers = deserializeArguments(sFn.ListOfPointersOffset, sFn.ListOfPointersSize, s, prgrm)
	fn.Expressions = deserializeExpressions(sFn.ExpressionsOffset, sFn.ExpressionsSize, s, prgrm)
	fn.Size = int(sFn.Size)
	fn.Length = int(sFn.Length)

	if sFn.CurrentExpressionOffset > 0 {
		fn.CurrentExpression = fn.Expressions[sFn.CurrentExpressionOffset]
	}

	fn.Package = prgrm.Packages[sFn.PackageOffset]
}

func deserializeBool(val int64) bool {
	return val == 1
}

func deserializeIntegers(off int64, size int64, s *serializedCXProgram) []int {
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
func initDeserialization(prgrm *CXProgram, s *serializedCXProgram) {
	prgrm.Memory = s.Memory
	prgrm.Packages = make([]*CXPackage, len(s.Packages))
	prgrm.CallStack = make([]CXCall, constants.CALLSTACK_SIZE)
	prgrm.HeapStartsAt = int(s.Program.HeapStartsAt)
	prgrm.HeapPointer = int(s.Program.HeapPointer)
	prgrm.StackSize = int(s.Program.StackSize)
	prgrm.HeapSize = int(s.Program.HeapSize)
	prgrm.Version = deserializeString(s.Program.VersionOffset, s.Program.VersionSize, s)

	deserializePackages(s, prgrm)
}

// Deserialize deserializes a serialized CX program back to its golang struct representation.
func Deserialize(b []byte) (prgrm *CXProgram) {
	prgrm = &CXProgram{}
	idxSize := encoder.Size(serializedCXProgramIndex{})

	var s serializedCXProgram

	helper.DeserializeRaw(b[:idxSize], &s.Index)
	helper.DeserializeRaw(b[s.Index.ProgramOffset:s.Index.CallsOffset], &s.Program)
	helper.DeserializeRaw(b[s.Index.CallsOffset:s.Index.PackagesOffset], &s.Calls)
	helper.DeserializeRaw(b[s.Index.PackagesOffset:s.Index.StructsOffset], &s.Packages)
	helper.DeserializeRaw(b[s.Index.StructsOffset:s.Index.FunctionsOffset], &s.Structs)
	helper.DeserializeRaw(b[s.Index.FunctionsOffset:s.Index.ExpressionsOffset], &s.Functions)
	helper.DeserializeRaw(b[s.Index.ExpressionsOffset:s.Index.ArgumentsOffset], &s.Expressions)
	helper.DeserializeRaw(b[s.Index.ArgumentsOffset:s.Index.IntegersOffset], &s.Arguments)
	helper.DeserializeRaw(b[s.Index.IntegersOffset:s.Index.StringsOffset], &s.Integers)
	s.Strings = b[s.Index.StringsOffset:s.Index.MemoryOffset]
	s.Memory = b[s.Index.MemoryOffset:]

	initDeserialization(prgrm, &s)

	// prgrm.PrintProgram()

	return prgrm
}

// CopyProgramState copies the program state from `prgrm1` to `prgrm2`.
func CopyProgramState(sPrgrm1, sPrgrm2 *[]byte) {
	idxSize := encoder.Size(serializedCXProgramIndex{})

	var index1 serializedCXProgramIndex
	var index2 serializedCXProgramIndex

	helper.DeserializeRaw((*sPrgrm1)[:idxSize], &index1)
	helper.DeserializeRaw((*sPrgrm2)[:idxSize], &index2)

	var prgrm1Info serializedProgram
	helper.DeserializeRaw((*sPrgrm1)[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info serializedProgram
	helper.DeserializeRaw((*sPrgrm2)[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

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
	idxSize := encoder.Size(serializedCXProgramIndex{})
	var index serializedCXProgramIndex
	helper.DeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo serializedProgram
	helper.DeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.StackSize)
}

// GetSerializedDataSize returns the size of the data segment of a serialized CX program.
func GetSerializedDataSize(sPrgrm []byte) int {
	idxSize := encoder.Size(serializedCXProgramIndex{})
	var index serializedCXProgramIndex
	helper.DeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo serializedProgram
	helper.DeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.HeapStartsAt - prgrmInfo.StackSize)
}
