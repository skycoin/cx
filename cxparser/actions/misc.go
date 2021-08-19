package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

func SelectProgram(prgrm *ast.CXProgram) {
	AST = prgrm
}

func SetCorrectArithmeticOp(expr *ast.CXExpression) {
	if expr.Operator == nil || len(expr.Outputs) < 1 {
		return
	}

	code := expr.Operator.OpCode
	if code > constants.START_OF_OPERATORS && code < constants.END_OF_OPERATORS {
		// TODO: argument type are not fully resolved here, should be move elsewhere.
		//expr.Operator = cxcore.GetTypedOperator(cxcore.GetType(expr.ProgramInput[0]), code)
	}
}

// hasDeclSpec determines if an argument has certain declaration specifier
func hasDeclSpec(arg *ast.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DeclarationSpecifiers {
		if s == spec {
			found = true
		}
	}
	return found
}

// hasDerefOp determines if an argument has certain dereference operation
func hasDerefOp(arg *ast.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DereferenceOperations {
		if s == spec {
			found = true
		}
	}
	return found
}

// This function writes those bytes to AST.Data
func WritePrimary(typeCode types.Code, byts []byte, isSlice bool) []*ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(typeCode)
		arg.ArgDetails.Package = pkg

		size := types.Cast_int_to_ptr(len(byts))

		arg.Size = typeCode.Size()
		arg.TotalSize = size
		arg.Offset = AST.DataSegmentSize + AST.DataSegmentStartsAt

		if arg.Type == types.STR || arg.Type == types.AFF {
			arg.PassBy = constants.PASSBY_REFERENCE
			arg.Size = types.POINTER_SIZE
			arg.TotalSize = types.POINTER_SIZE
			if isSlice == false {
				types.Write_ptr(byts, 0, arg.Offset)
			}
		}

		// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
		// after the stack segment. These bytes are used to allocate the data segment
		// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
		// we'll start appending the bytes to AST.Memory.
		// After compilation, we calculate how many bytes we need to add to have a heap segment
		// equal to `minHeapSize()` that is allocated after the data segment.
		memSize := types.Cast_int_to_ptr(len(AST.Memory))
		if (size + AST.DataSegmentSize + AST.DataSegmentStartsAt) > memSize {
			var i types.Pointer
			// First we need to fill the remaining free bytes in
			// the current `AST.Memory` slice.
			for i = types.Pointer(0); i < memSize-AST.DataSegmentSize+AST.DataSegmentStartsAt; i++ {
				AST.Memory[AST.DataSegmentSize+AST.DataSegmentStartsAt+i] = byts[i]
			}
			// Then we append the bytes that didn't fit.
			AST.Memory = append(AST.Memory, byts[i:]...)
		} else {
			for i, byt := range byts {
				AST.Memory[AST.DataSegmentSize+AST.DataSegmentStartsAt+types.Cast_int_to_ptr(i)] = byt
			}
		}
		AST.DataSegmentSize += size
		AST.HeapStartsAt = AST.DataSegmentSize + AST.DataSegmentStartsAt

		expr := ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*ast.CXExpression{expr}
	} else {
		panic(err)
	}
}

func TotalLength(lengths []types.Pointer) types.Pointer {
	total := types.Pointer(1)
	for _, i := range lengths {
		total *= i
	}
	return total
}

func StructLiteralFields(ident string) *ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(types.IDENTIFIER)
		arg.ArgDetails.Name = ident
		arg.ArgDetails.Package = pkg

		expr := ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*ast.CXArgument{arg}
		expr.Package = pkg

		return expr
	} else {
		panic(err)
	}
}

func AffordanceStructs(pkg *ast.CXPackage, currentFile string, lineNo int) {
	// Argument type
	argStrct := ast.MakeStruct("Argument")
	// argStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	argFldName := ast.MakeField("Name", types.STR, "", 0)
	argFldName.TotalSize = types.STR.Size()
	argFldIndex := ast.MakeField("Index", types.I32, "", 0)
	argFldIndex.TotalSize = types.I32.Size()
	argFldType := ast.MakeField("Type", types.STR, "", 0)
	argFldType.TotalSize = types.STR.Size()

	argStrct.AddField(argFldName)
	argStrct.AddField(argFldIndex)
	argStrct.AddField(argFldType)

	pkg.AddStruct(argStrct)

	// Expression type
	exprStrct := ast.MakeStruct("Expression")
	// exprStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	exprFldOperator := ast.MakeField("Operator", types.STR, "", 0)

	exprStrct.AddField(exprFldOperator)

	pkg.AddStruct(exprStrct)

	// Function type
	fnStrct := ast.MakeStruct("Function")
	// fnStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	fnFldName := ast.MakeField("Name", types.STR, "", 0)
	fnFldName.TotalSize = types.STR.Size()

	fnFldInpSig := ast.MakeField("InputSignature", types.STR, "", 0)
	fnFldInpSig.Size = types.STR.Size()
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, []types.Pointer{0}, constants.DECL_SLICE)

	fnFldOutSig := ast.MakeField("OutputSignature", types.STR, "", 0)
	fnFldOutSig.Size = types.STR.Size()
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, []types.Pointer{0}, constants.DECL_SLICE)

	fnStrct.AddField(fnFldName)
	fnStrct.AddField(fnFldInpSig)

	fnStrct.AddField(fnFldOutSig)

	pkg.AddStruct(fnStrct)

	// Structure type
	strctStrct := ast.MakeStruct("Structure")
	// strctStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	strctFldName := ast.MakeField("Name", types.STR, "", 0)
	strctFldName.TotalSize = types.STR.Size()

	strctStrct.AddField(strctFldName)

	pkg.AddStruct(strctStrct)

	// Package type
	pkgStrct := ast.MakeStruct("Structure")
	// pkgStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	pkgFldName := ast.MakeField("Name", types.STR, "", 0)

	pkgStrct.AddField(pkgFldName)

	pkg.AddStruct(pkgStrct)

	// Caller type
	callStrct := ast.MakeStruct("Caller")
	// callStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_I32)

	callFldFnName := ast.MakeField("FnName", types.STR, "", 0)
	callFldFnName.TotalSize = types.STR.Size()
	callFldFnSize := ast.MakeField("FnSize", types.I32, "", 0)
	callFldFnSize.TotalSize = types.I32.Size()

	callStrct.AddField(callFldFnName)
	callStrct.AddField(callFldFnSize)

	pkg.AddStruct(callStrct)

	// Program type
	prgrmStrct := ast.MakeStruct("Program")
	// prgrmStrct.Size = cxcore.GetArgSize(cxcore.TYPE_I32) + cxcore.GetArgSize(cxcore.TYPE_I64)

	prgrmFldCallCounter := ast.MakeField("CallCounter", types.I32, "", 0)
	prgrmFldCallCounter.TotalSize = types.I32.Size()
	prgrmFldFreeHeap := ast.MakeField("HeapUsed", types.I64, "", 0)
	prgrmFldFreeHeap.TotalSize = types.I64.Size()

	// prgrmFldCaller := cxcore.MakeField("Caller", cxcore.TYPE_STRUCT, "", 0)
	prgrmFldCaller := DeclarationSpecifiersStruct(callStrct.Name, callStrct.Package.Name, false, currentFile, lineNo)
	prgrmFldCaller.ArgDetails.Name = "Caller"

	prgrmStrct.AddField(prgrmFldCallCounter)
	prgrmStrct.AddField(prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrmFldCaller)

	pkg.AddStruct(prgrmStrct)
}

func PrimaryIdentifier(ident string) []*ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument(ident, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
		arg.AddType(types.IDENTIFIER)
		// arg.Typ = "ident"
		arg.ArgDetails.Name = ident
		arg.ArgDetails.Package = pkg

		// expr := &cxcore.CXExpression{ProgramOutput: []*cxcore.CXArgument{arg}}
		expr := ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*ast.CXArgument{arg}
		expr.Package = pkg

		return []*ast.CXExpression{expr}
	} else {
		panic(err)
	}
}

// IsAllArgsBasicTypes checks if all the input arguments in an expressions are of basic type.
func IsAllArgsBasicTypes(expr *ast.CXExpression) bool {
	for _, inp := range expr.Inputs {
		if !inp.Type.IsPrimitive() {
			return false
		}
	}
	return true
}
