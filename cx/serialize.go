package base

import (
        "github.com/skycoin/skycoin/src/cipher/encoder"
)

type sIndex struct {
        ProgramOffset                   int32
        CallsOffset                     int32
        PackagesOffset                  int32
        StructsOffset                   int32
        FunctionsOffset                 int32
        ExpressionsOffset               int32
        ArgumentsOffset                 int32
	IntegersOffset                  int32
        NamesOffset                     int32
        MemoryOffset                    int32
}

type sProgram struct {
        PackagesOffset                  int32
        PackagesSize                    int32
        CurrentPackageOffset            int32
        
        InputsOffset                    int32
        InputsSize                      int32

        OutputsOffset                   int32
        OutputsSize                     int32
        
        CallStackOffset                 int32
        CallStackSize                   int32

        CallCounter                     int32

        MemoryOffset                    int32
        MemorySize                      int32

        HeapPointer                     int32
        StackPointer                    int32
        HeapStartsAt                    int32
        
        Terminated                      int32
}

type sCall struct {
        OperatorOffset                  int32
        Line                            int32
        FramePointer                    int32
}

type sPackage struct {
        NameOffset                      int32
        NameSize                        int32
        ImportsOffset                   int32
        ImportsSize                     int32
        StructsOffset                   int32
        StructsSize                     int32
	GlobalsOffset                   int32
        GlobalsSize                     int32
	FunctionsOffset                 int32
        FunctionsSize                   int32
        CurrentFunctionOffset           int32
        CurrentStructOffset             int32
}

type sStruct struct {
        NameOffset                      int32
        NameSize                        int32
        FieldsOffset                    int32
        FieldsSize                      int32

        Size                            int32

        PackageOffset                   int32
}

type sFunction struct {
        NameOffset                      int32
        NameSize                        int32
        InputsOffset                    int32
        InputsSize                      int32
        OutputsOffset                   int32
        OutputsSize                     int32
        ExpressionsOffset               int32
        ExpressionsSize                 int32
        Size                            int32
        Length                          int32

        ListOfPointersOffset            int32
        ListOfPointersSize              int32

	// We're going to determine this when procesing the expressions. Check sExpression type
        // IsNative                        int32
        // OpCode                          int32

        CurrentExpressionOffset         int32
        PackageOffset                   int32
}

type sExpression struct {
        OperatorOffset                  int32
	// we add these two fields here so we don't add every native sFunction to the serialization
	// the CX runtime already knows about the natives properties. We just need the code if IsNative = true
	IsNative                        int32
	OpCode                          int32
	
        InputsOffset                    int32
        InputsSize                      int32
        OutputsOffset                   int32
        OutputsSize                     int32
        
        LabelOffset                     int32
        LabelSize                       int32
        ThenLines                       int32
        ElseLines                       int32

        IsMethodCall                    int32
        IsStructLiteral                 int32
        IsArrayLiteral                  int32

        FunctionOffset                  int32
        PackageOffset                   int32
}

type sArgument struct {
        NameOffset                      int32
        NameSize                        int32
        Type                            int32
        CustomTypeOffset                int32
        Size                            int32
        TotalSize                       int32

        Offset                          int32
        
        IndirectionLevels               int32
        DereferenceLevels               int32
        DereferenceOperationsOffset     int32
        DereferenceOperationsSize       int32
        DeclarationSpecifiersOffset     int32
        DeclarationSpecifiersSize       int32

        IsSlice                         int32
        IsArray                         int32
        IsArrayFirst                    int32
        IsPointer                       int32
        IsReference                     int32
        
        IsDereferenceFirst              int32
        IsStruct                        int32
        IsRest                          int32
        IsLocalDeclaration              int32
        IsShortDeclaration              int32

        PassBy                          int32
        DoesEscape                      int32

        LengthsOffset                   int32
        LengthsSize                     int32
        IndexesOffset                   int32
        IndexesSize                     int32
        FieldsOffset                    int32
        FieldsSize                      int32

        PackageOffset                   int32
}

type sAll struct {
	Index                           sIndex
	Program                         sProgram
	
	Packages                        []sPackage
	PackagesMap                     map[string]int
	Structs                         []sStruct
	StructsMap                      map[string]int
	Functions                       []sFunction
	FunctionsMap                    map[string]int
	
	Expressions                     []sExpression
	Arguments                       []sArgument
	Calls                           []sCall
	
	Names                           []byte
	NamesMap                        map[string]int
	Integers                        []int32
	
	Memory                          []byte
}

func serializeName (name string, s *sAll) (int32, int32) {
	if name == "" {
		return int32(-1), int32(-1)
	}
	
	var size int
	var err error
	
	size, err = encoder.Size(name)
	
	if err != nil {
		panic(err)
	}
	
	if off, found := s.NamesMap[name]; found {
		return int32(off), int32(size)
	} else {
		off = len(s.Names)
		s.Names = append(s.Names, encoder.Serialize(name)...)
		s.NamesMap[name] = off

		return int32(off), int32(size)
	}
}

func indexPackage (pkg *CXPackage, s *sAll) {
	if _, found := s.PackagesMap[pkg.Name]; !found {
		s.PackagesMap[pkg.Name] = len(s.PackagesMap)
	} else {
		panic("duplicated package in serialization process")
	}
}

func indexStruct (strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if _, found := s.StructsMap[strctName]; !found {
		s.StructsMap[strctName] = len(s.StructsMap)
	} else {
		panic("duplicated struct in serialization process")
	}
}

func indexFunction (fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if _, found := s.FunctionsMap[fnName]; !found {
		s.FunctionsMap[fnName] = len(s.FunctionsMap)
	} else {
		panic("duplicated function in serialization process")
	}
}

func indexExpression (expr *CXExpression, s *sAll) {

}

func serializeBoolean (val bool) int32 {
	if val {
		return 1
	} else {
		return 0
	}
}

func serializeIntegers (ints []int, s *sAll) (int32, int32) {
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

func serializeArgument (arg *CXArgument, s *sAll) int {
	s.Arguments = append(s.Arguments, sArgument{})
	argOff := len(s.Arguments) - 1
	sArg := &s.Arguments[argOff]

	sNil := int32(-1)
	
	sArg.NameOffset, sArg.NameSize = serializeName(arg.Name, s)
	
	sArg.Type = int32(arg.Type)
	
	if arg.CustomType == nil {
		sArg.CustomTypeOffset = sNil
	} else {
		strctName := arg.CustomType.Package.Name + "." + arg.CustomType.Name
		if strctOff, found := s.StructsMap[strctName]; found {
			sArg.CustomTypeOffset = int32(strctOff)
		} else {
			panic("struct reference not found")
		}
	}

	sArg.Size = int32(arg.Size)
	sArg.TotalSize = int32(arg.TotalSize)
	sArg.Offset = int32(arg.Offset)
	sArg.IndirectionLevels = int32(arg.IndirectionLevels)
	sArg.DereferenceLevels = int32(arg.DereferenceLevels)
	sArg.DereferenceLevels = int32(arg.DereferenceLevels)

	sArg.DereferenceOperationsOffset,
	sArg.DereferenceOperationsSize = serializeIntegers(arg.DereferenceOperations, s)

	sArg.DeclarationSpecifiersOffset,
	sArg.DeclarationSpecifiersSize = serializeIntegers(arg.DeclarationSpecifiers, s)

	sArg.IsSlice = serializeBoolean(arg.IsSlice)
	sArg.IsArray = serializeBoolean(arg.IsArray)
	sArg.IsArrayFirst = serializeBoolean(arg.IsArrayFirst)
	sArg.IsPointer = serializeBoolean(arg.IsPointer)
	sArg.IsReference = serializeBoolean(arg.IsReference)

	sArg.IsDereferenceFirst = serializeBoolean(arg.IsDereferenceFirst)
	sArg.IsStruct = serializeBoolean(arg.IsStruct)
	sArg.IsRest = serializeBoolean(arg.IsRest)
	sArg.IsLocalDeclaration = serializeBoolean(arg.IsLocalDeclaration)
	sArg.IsShortDeclaration = serializeBoolean(arg.IsShortDeclaration)

	sArg.PassBy = int32(arg.PassBy)
	sArg.DoesEscape = serializeBoolean(arg.DoesEscape)

	sArg.LengthsOffset,
	sArg.LengthsSize = serializeIntegers(arg.Lengths, s)

	sArg.IndexesOffset, sArg.IndexesSize = serializeSliceOfArguments(arg.Indexes, s)
	sArg.FieldsOffset, sArg.FieldsSize = serializeSliceOfArguments(arg.Fields, s)

	if pkgOff, found := s.PackagesMap[arg.Package.Name]; found {
		sArg.PackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	return argOff
}

func serializeSliceOfArguments (args []*CXArgument, s *sAll) (int32, int32) {
	if len(args) == 0 {
		return int32(-1), int32(-1)
	} else {
		idxs := make([]int, len(args))
		for i, arg := range args {
			idxs[i] = serializeArgument(arg, s)
		}
		return serializeIntegers(idxs, s)
	}
}

func serializeCalls (calls []CXCall, s *sAll) (int32, int32) {
	if len(calls) == 0 {
		return int32(-1), int32(-1)
	} else {
		idxs := make([]int, len(calls))
		for i, call := range calls {
			idxs[i] = serializeCall(&call, s)
		}
		return serializeIntegers(idxs, s)
	}
}

func serializeExpression (expr *CXExpression, s *sAll) int {
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

	sExpr.IsMethodCall = serializeBoolean(expr.IsMethodCall)
	sExpr.IsStructLiteral = serializeBoolean(expr.IsStructLiteral)
	sExpr.IsArrayLiteral = serializeBoolean(expr.IsArrayLiteral)

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

func serializeCall (call *CXCall, s *sAll) int {
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

func serializeProgram (prgrm *CXProgram, s *sAll) {
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
	sPrgrm.MemorySize = int32(len(prgrm.Memory))

	sPrgrm.HeapPointer = int32(prgrm.HeapPointer)
	sPrgrm.StackPointer = int32(prgrm.StackPointer)
	sPrgrm.HeapStartsAt = int32(prgrm.HeapStartsAt)

	sPrgrm.Terminated = serializeBoolean(prgrm.Terminated)
}

func sStructArguments (strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if strctOff, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[strctOff]
		sStrct.FieldsOffset, sStrct.FieldsSize = serializeSliceOfArguments(strct.Fields, s)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionArguments (fn *CXFunction, s *sAll) {
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

func sPackageName (pkg *CXPackage, s *sAll) {
	sPkg := &s.Packages[s.PackagesMap[pkg.Name]]
	sPkg.NameOffset, sPkg.NameSize = serializeName(pkg.Name, s)
}

func sStructName (strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	sStrct := &s.Structs[s.StructsMap[strctName]]
	sStrct.NameOffset, sStrct.NameSize = serializeName(strct.Name, s)
}

func sFunctionName (fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.NameOffset, sFn.NameSize = serializeName(fn.Name, s)
	} else {
		panic("function reference not found")
	}
}

func sPackageGlobals (pkg *CXPackage, s *sAll) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]
		sPkg.GlobalsOffset, sPkg.GlobalsSize = serializeSliceOfArguments(pkg.Globals, s)
	} else {
		panic("package reference not found")
	}
}

func sPackageImports (pkg *CXPackage, s *sAll) {
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

func sStructPackage (strct *CXStruct, s *sAll) {
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

func sFunctionPackage (fn *CXFunction, s *sAll) {
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

func sPackageIntegers (pkg *CXPackage, s *sAll) {
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

func sStructIntegers (strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if off, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[off]
		sStrct.Size = int32(strct.Size)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionIntegers (fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.Size = int32(fn.Size)
		sFn.Length = int32(fn.Length)
	} else {
		panic("function reference not found")
	}
}

func initAll (prgrm *CXProgram, s *sAll) {
	s.PackagesMap = make(map[string]int, 0)
	s.StructsMap = make(map[string]int, 0)
	s.FunctionsMap = make(map[string]int, 0)
	s.NamesMap = make(map[string]int, 0)

	s.Calls = make([]sCall, prgrm.CallCounter)
	s.Packages = make([]sPackage, len(prgrm.Packages))

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

func Serialize (prgrm *CXProgram) (byts []byte) {
	s := sAll{}
	initAll(prgrm, &s)

	// indexing packages and serializing their names
	for _, pkg := range prgrm.Packages {
		indexPackage(pkg, &s)
		sPackageName(pkg, &s)
	}
	// we first needed to populate references to all packages
	// now we add the imports' references
	for _, pkg := range prgrm.Packages {
		sPackageImports(pkg, &s)
	}
	
	// structs
	for _, pkg := range prgrm.Packages {
		for _, strct := range pkg.Structs {
			indexStruct(strct, &s)
			sStructName(strct, &s)
			sStructPackage(strct, &s)
			sStructIntegers(strct, &s)
		}
	}
	// we first needed to populate references to all structs
	// now we add fields
	for _, pkg := range prgrm.Packages {
		for _, strct := range pkg.Structs {
			sStructArguments(strct, &s)
		}
	}

	// globals
	for _, pkg := range prgrm.Packages {
		sPackageGlobals(pkg, &s)
	}

	// functions
	for _, pkg := range prgrm.Packages {
		for _, fn := range pkg.Functions {
			indexFunction(fn, &s)
			sFunctionName(fn, &s)
			sFunctionPackage(fn, &s)
			sFunctionIntegers(fn, &s)
			sFunctionArguments(fn, &s)
		}
	}

	// package elements' offsets and sizes
	var fnCounter int32
	var strctCounter int32
	for _, pkg := range prgrm.Packages {
		if pkgOff, found := s.PackagesMap[pkg.Name]; found {
			sPkg := &s.Packages[pkgOff]

			if len(pkg.Structs) == 0 {
				sPkg.StructsOffset = int32(-1)
				sPkg.StructsSize = int32(-1)
			} else {
				sPkg.StructsOffset = strctCounter
				lenStrcts := int32(len(pkg.Structs))
				sPkg.StructsSize = lenStrcts
				strctCounter += lenStrcts
			}

			if len(pkg.Functions) == 0 {
				sPkg.FunctionsOffset = int32(-1)
				sPkg.FunctionsSize = int32(-1)
			} else {
				sPkg.FunctionsOffset = fnCounter
				lenFns := int32(len(pkg.Functions))
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
		sPackageIntegers(pkg, &s)
	}

	// expressions
	for _, pkg := range prgrm.Packages {
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
						exprIdx := serializeExpression(expr, &s)
						if fn.CurrentExpression == expr {
							sFn.CurrentExpressionOffset = int32(exprIdx)
						}
						exprs[i] = exprIdx
					}

					sFn.ExpressionsOffset, sFn.ExpressionsSize = serializeIntegers(exprs, &s)
				}
			} else {
				panic("function reference not found")
			}
		}
	}

	// program
	serializeProgram(prgrm, &s)

	s.Index = sIndex{}
	sIdx := &s.Index

	// assigning relative offset

	idxSize, _ := encoder.Size(s.Index)
	prgrmSize, _ := encoder.Size(s.Program)
	callSize, _ := encoder.Size(s.Calls)
	pkgSize, _ := encoder.Size(s.Packages)
	strctSize, _ := encoder.Size(s.Structs)
	fnSize, _ := encoder.Size(s.Functions)
	exprSize, _ := encoder.Size(s.Expressions)
	argSize, _ := encoder.Size(s.Arguments)
	intSize, _ := encoder.Size(s.Integers)

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

	Debug("sbyts", s.Calls, encoder.Serialize(s.Calls))
	
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

func dsName (off int32, size int32, s *sAll) string {
	var name string
	encoder.DeserializeRaw(s.Names[off : off + size], &name)
	return name
}

// func dsArgument (off int32, size int32, s *sAll) CXArgument {
// 	var arg CXArgument
// 	encoder.DeserializeRaw(s.Arguments[])
// }


// func dsArguments (off int32, size int32, s *sAll) []*CXArgument {
// 	// s.Arguments[off : off + size]

	
// 	var args []CXArgument
// 	var argsPtr []*CXArgument
// 	encoder.DeserializeRaw(byts, &args)

// 	for _, arg := range args {
// 		argsPtr = append(argsPtr, &arg)
// 	}

// 	return argsPtr
	
// }


func coco (byts []byte) {
	
}

func Deserialize (byts []byte) (prgrm *CXProgram) {
	

	
	idxSize, _ := encoder.Size(sIndex{})

	var s sAll
	
	encoder.DeserializeRaw(byts[:idxSize], &s.Index)
	encoder.DeserializeRaw(byts[s.Index.ProgramOffset : s.Index.CallsOffset], &s.Program)
	encoder.DeserializeRaw(byts[s.Index.CallsOffset : s.Index.PackagesOffset], &s.Calls)
	encoder.DeserializeRaw(byts[s.Index.PackagesOffset : s.Index.StructsOffset], &s.Packages)
	encoder.DeserializeRaw(byts[s.Index.StructsOffset : s.Index.FunctionsOffset], &s.Structs)
	encoder.DeserializeRaw(byts[s.Index.FunctionsOffset : s.Index.ExpressionsOffset], &s.Functions)
	encoder.DeserializeRaw(byts[s.Index.ExpressionsOffset : s.Index.ArgumentsOffset], &s.Expressions)
	encoder.DeserializeRaw(byts[s.Index.ArgumentsOffset : s.Index.IntegersOffset], &s.Arguments)
	encoder.DeserializeRaw(byts[s.Index.IntegersOffset : s.Index.NamesOffset], &s.Integers)
	s.Names = byts[s.Index.NamesOffset : s.Index.MemoryOffset]
	s.Memory = byts[s.Index.MemoryOffset : ]


	// prgrm.CallCounter = s.CallCounter
	// prgrm.HeapPointer = s.HeapPointer
	// prgrm.Inputs = dsArguments()


	Debug("hoho", s.Packages)

	return &CXProgram{}
}
