package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
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
func WritePrimary(typ int, byts []byte, isGlobal bool) []*ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(constants.TypeNames[typ])
		arg.Package = pkg

		var size = len(byts)

		arg.Size = constants.GetArgSize(typ)
		arg.TotalSize = size
		arg.Offset = DataOffset

		if arg.Type == constants.TYPE_STR || arg.Type == constants.TYPE_AFF {
			arg.PassBy = constants.PASSBY_REFERENCE
			arg.Size = constants.TYPE_POINTER_SIZE
			arg.TotalSize = constants.TYPE_POINTER_SIZE
		}

		// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
		// after the stack segment. These bytes are used to allocate the data segment
		// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
		// we'll start appending the bytes to AST.Memory.
		// After compilation, we calculate how many bytes we need to add to have a heap segment
		// equal to `minHeapSize()` that is allocated after the data segment.
		if size+DataOffset > len(AST.Memory) {
			var i int
			// First we need to fill the remaining free bytes in
			// the current `AST.Memory` slice.
			for i = 0; i < len(AST.Memory)-DataOffset; i++ {
				AST.Memory[DataOffset+i] = byts[i]
			}
			// Then we append the bytes that didn't fit.
			AST.Memory = append(AST.Memory, byts[i:]...)
		} else {
			for i, byt := range byts {
				AST.Memory[DataOffset+i] = byt
			}
		}
		DataOffset += size

		expr := ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*ast.CXExpression{expr}
	} else {
		panic(err)
	}
}

func TotalLength(lengths []int) int {
	var total int = 1
	for _, i := range lengths {
		total *= i
	}
	return total
}

func StructLiteralFields(ident string) *ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(constants.TypeNames[constants.TYPE_IDENTIFIER])
		arg.Name = ident
		arg.Package = pkg

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

	argFldName := ast.MakeField("Name", constants.TYPE_STR, "", 0)
	argFldName.TotalSize = constants.GetArgSize(constants.TYPE_STR)
	argFldIndex := ast.MakeField("Index", constants.TYPE_I32, "", 0)
	argFldIndex.TotalSize = constants.GetArgSize(constants.TYPE_I32)
	argFldType := ast.MakeField("Type", constants.TYPE_STR, "", 0)
	argFldType.TotalSize = constants.GetArgSize(constants.TYPE_STR)

	argStrct.AddField(argFldName)
	argStrct.AddField(argFldIndex)
	argStrct.AddField(argFldType)

	pkg.AddStruct(argStrct)

	// Expression type
	exprStrct := ast.MakeStruct("Expression")
	// exprStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	exprFldOperator := ast.MakeField("Operator", constants.TYPE_STR, "", 0)

	exprStrct.AddField(exprFldOperator)

	pkg.AddStruct(exprStrct)

	// Function type
	fnStrct := ast.MakeStruct("Function")
	// fnStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	fnFldName := ast.MakeField("Name", constants.TYPE_STR, "", 0)
	fnFldName.TotalSize = constants.GetArgSize(constants.TYPE_STR)

	fnFldInpSig := ast.MakeField("InputSignature", constants.TYPE_STR, "", 0)
	fnFldInpSig.Size = constants.GetArgSize(constants.TYPE_STR)
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, []int{0}, constants.DECL_SLICE)

	fnFldOutSig := ast.MakeField("OutputSignature", constants.TYPE_STR, "", 0)
	fnFldOutSig.Size = constants.GetArgSize(constants.TYPE_STR)
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, []int{0}, constants.DECL_SLICE)

	fnStrct.AddField(fnFldName)
	fnStrct.AddField(fnFldInpSig)

	fnStrct.AddField(fnFldOutSig)

	pkg.AddStruct(fnStrct)

	// Structure type
	strctStrct := ast.MakeStruct("Structure")
	// strctStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	strctFldName := ast.MakeField("Name", constants.TYPE_STR, "", 0)
	strctFldName.TotalSize = constants.GetArgSize(constants.TYPE_STR)

	strctStrct.AddField(strctFldName)

	pkg.AddStruct(strctStrct)

	// Package type
	pkgStrct := ast.MakeStruct("Structure")
	// pkgStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	pkgFldName := ast.MakeField("Name", constants.TYPE_STR, "", 0)

	pkgStrct.AddField(pkgFldName)

	pkg.AddStruct(pkgStrct)

	// Caller type
	callStrct := ast.MakeStruct("Caller")
	// callStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_I32)

	callFldFnName := ast.MakeField("FnName", constants.TYPE_STR, "", 0)
	callFldFnName.TotalSize = constants.GetArgSize(constants.TYPE_STR)
	callFldFnSize := ast.MakeField("FnSize", constants.TYPE_I32, "", 0)
	callFldFnSize.TotalSize = constants.GetArgSize(constants.TYPE_I32)

	callStrct.AddField(callFldFnName)
	callStrct.AddField(callFldFnSize)

	pkg.AddStruct(callStrct)

	// Program type
	prgrmStrct := ast.MakeStruct("Program")
	// prgrmStrct.Size = cxcore.GetArgSize(cxcore.TYPE_I32) + cxcore.GetArgSize(cxcore.TYPE_I64)

	prgrmFldCallCounter := ast.MakeField("CallCounter", constants.TYPE_I32, "", 0)
	prgrmFldCallCounter.TotalSize = constants.GetArgSize(constants.TYPE_I32)
	prgrmFldFreeHeap := ast.MakeField("HeapUsed", constants.TYPE_I64, "", 0)
	prgrmFldFreeHeap.TotalSize = constants.GetArgSize(constants.TYPE_I64)

	// prgrmFldCaller := cxcore.MakeField("Caller", cxcore.TYPE_CUSTOM, "", 0)
	prgrmFldCaller := DeclarationSpecifiersStruct(callStrct.Name, callStrct.Package.Name, false, currentFile, lineNo)
	prgrmFldCaller.Name = "Caller"

	prgrmStrct.AddField(prgrmFldCallCounter)
	prgrmStrct.AddField(prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrmFldCaller)

	pkg.AddStruct(prgrmStrct)
}

func PrimaryIdentifier(ident string) []*ast.CXExpression {
	if pkg, err := AST.GetCurrentPackage(); err == nil {
		arg := ast.MakeArgument(ident, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
		arg.AddType(constants.TypeNames[constants.TYPE_IDENTIFIER])
		// arg.Typ = "ident"
		arg.Name = ident
		arg.Package = pkg

		// expr := &cxcore.CXExpression{ProgramOutput: []*cxcore.CXArgument{arg}}
		expr := ast.MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*ast.CXArgument{arg}
		expr.Package = pkg

		return []*ast.CXExpression{expr}
	} else {
		panic(err)
	}
}

// IsArgBasicType returns true if `arg`'s type is a basic type, false otherwise.
func IsArgBasicType(arg *ast.CXArgument) bool {
	switch arg.Type {
	case constants.TYPE_BOOL,
		constants.TYPE_STR, //A STRING IS NOT AN ATOMIC TYPE
		constants.TYPE_F32,
		constants.TYPE_F64,
		constants.TYPE_I8,
		constants.TYPE_I16,
		constants.TYPE_I32,
		constants.TYPE_I64,
		constants.TYPE_UI8,
		constants.TYPE_UI16,
		constants.TYPE_UI32,
		constants.TYPE_UI64:
		return true
	}
	return false
}

// IsAllArgsBasicTypes checks if all the input arguments in an expressions are of basic type.
func IsAllArgsBasicTypes(expr *ast.CXExpression) bool {
	for _, inp := range expr.Inputs {
		if !IsArgBasicType(inp) {
			return false
		}
	}
	return true
}