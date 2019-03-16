package actions

import (
	. "github.com/skycoin/cx/cx"
)

func SetCorrectArithmeticOp(expr *CXExpression) {
	if expr.Operator == nil || len(expr.Outputs) < 1 {
		return
	}
	op := expr.Operator
	typ := expr.Outputs[0].Type

	if CheckArithmeticOp(expr) {
		switch op.OpCode {
		case OP_I32_MUL:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_MUL]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_MUL]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_MUL]
			}
		case OP_I32_DIV:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_DIV]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_DIV]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_DIV]
			}
		case OP_I32_MOD:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_MOD]
			}

		case OP_I32_ADD:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_ADD]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_ADD]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_ADD]
			}
		case OP_I32_SUB:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_ADD]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_ADD]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_ADD]
			}

		case OP_I32_BITSHL:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITSHL]
			}
		case OP_I32_BITSHR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITSHR]
			}

		case OP_I32_LT:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_LT]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_LT]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_LT]
			}
		case OP_I32_GT:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_GT]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_GT]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_GT]
			}
		case OP_I32_LTEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_LTEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_LTEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_LTEQ]
			}
		case OP_I32_GTEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_GTEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_GTEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_GTEQ]
			}

		case OP_I32_EQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_EQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_EQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_EQ]
			}
		case OP_I32_UNEQ:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_UNEQ]
			case TYPE_F32:
				expr.Operator = Natives[OP_F32_UNEQ]
			case TYPE_F64:
				expr.Operator = Natives[OP_F64_UNEQ]
			}

		case OP_I32_BITAND:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITAND]
			}

		case OP_I32_BITXOR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITXOR]
			}

		case OP_I32_BITOR:
			switch typ {
			case TYPE_I32:
			case TYPE_I64:
				expr.Operator = Natives[OP_I64_BITOR]
			}
		}
	}
}

// hasDeclSpec determines if an argument has certain declaration specifier
func hasDeclSpec(arg *CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DeclarationSpecifiers {
		if s == spec {
			found = true
		}
	}
	return found
}

// hasDerefOp determines if an argument has certain dereference operation
func hasDerefOp(arg *CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DereferenceOperations {
		if s == spec {
			found = true
		}
	}
	return found
}

// This function writes those bytes to PRGRM.Data
func WritePrimary(typ int, byts []byte, isGlobal bool) []*CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument("", CurrentFile, LineNo)
		arg.AddType(TypeNames[typ])
		arg.Package = pkg

		var size int

		size = len(byts)

		arg.Size = GetArgSize(typ)
		arg.TotalSize = size
		arg.Offset = DataOffset

		if arg.Type == TYPE_STR || arg.Type == TYPE_AFF {
			arg.PassBy = PASSBY_REFERENCE
			arg.Size = TYPE_POINTER_SIZE
			arg.TotalSize = TYPE_POINTER_SIZE
		}

		for i, byt := range byts {
			PRGRM.Memory[DataOffset+i] = byt
		}
		DataOffset += size

		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*CXExpression{expr}
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

func StructLiteralFields(ident string) *CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument("", CurrentFile, LineNo)
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		arg.Name = ident
		arg.Package = pkg

		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*CXArgument{arg}
		expr.Package = pkg

		return expr
	} else {
		panic(err)
	}
}

func AffordanceStructs(pkg *CXPackage, currentFile string, lineNo int) {
	// Argument type
	argStrct := MakeStruct("Argument")
	// argStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR)

	argFldName := MakeField("Name", TYPE_STR, "", 0)
	argFldName.TotalSize = GetArgSize(TYPE_STR)
	argFldIndex := MakeField("Index", TYPE_I32, "", 0)
	argFldIndex.TotalSize = GetArgSize(TYPE_I32)
	argFldType := MakeField("Type", TYPE_STR, "", 0)
	argFldType.TotalSize = GetArgSize(TYPE_STR)

	argStrct.AddField(argFldName)
	argStrct.AddField(argFldIndex)
	argStrct.AddField(argFldType)

	pkg.AddStruct(argStrct)

	// Expression type
	exprStrct := MakeStruct("Expression")
	// exprStrct.Size = GetArgSize(TYPE_STR)

	exprFldOperator := MakeField("Operator", TYPE_STR, "", 0)

	exprStrct.AddField(exprFldOperator)

	pkg.AddStruct(exprStrct)

	// Function type
	fnStrct := MakeStruct("Function")
	// fnStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR) + GetArgSize(TYPE_STR)

	fnFldName := MakeField("Name", TYPE_STR, "", 0)
	fnFldName.TotalSize = GetArgSize(TYPE_STR)

	fnFldInpSig := MakeField("InputSignature", TYPE_STR, "", 0)
	fnFldInpSig.Size = GetArgSize(TYPE_STR)
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, 0, DECL_SLICE)

	fnFldOutSig := MakeField("OutputSignature", TYPE_STR, "", 0)
	fnFldOutSig.Size = GetArgSize(TYPE_STR)
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, 0, DECL_SLICE)

	fnStrct.AddField(fnFldName)
	fnStrct.AddField(fnFldInpSig)

	fnStrct.AddField(fnFldOutSig)

	pkg.AddStruct(fnStrct)

	// Structure type
	strctStrct := MakeStruct("Structure")
	// strctStrct.Size = GetArgSize(TYPE_STR)

	strctFldName := MakeField("Name", TYPE_STR, "", 0)
	strctFldName.TotalSize = GetArgSize(TYPE_STR)

	strctStrct.AddField(strctFldName)

	pkg.AddStruct(strctStrct)

	// Package type
	pkgStrct := MakeStruct("Structure")
	// pkgStrct.Size = GetArgSize(TYPE_STR)

	pkgFldName := MakeField("Name", TYPE_STR, "", 0)

	pkgStrct.AddField(pkgFldName)

	pkg.AddStruct(pkgStrct)

	// Caller type
	callStrct := MakeStruct("Caller")
	// callStrct.Size = GetArgSize(TYPE_STR) + GetArgSize(TYPE_I32)

	callFldFnName := MakeField("FnName", TYPE_STR, "", 0)
	callFldFnName.TotalSize = GetArgSize(TYPE_STR)
	callFldFnSize := MakeField("FnSize", TYPE_I32, "", 0)
	callFldFnSize.TotalSize = GetArgSize(TYPE_I32)

	callStrct.AddField(callFldFnName)
	callStrct.AddField(callFldFnSize)

	pkg.AddStruct(callStrct)

	// Program type
	prgrmStrct := MakeStruct("Program")
	// prgrmStrct.Size = GetArgSize(TYPE_I32) + GetArgSize(TYPE_I64)

	prgrmFldCallCounter := MakeField("CallCounter", TYPE_I32, "", 0)
	prgrmFldCallCounter.TotalSize = GetArgSize(TYPE_I32)
	prgrmFldFreeHeap := MakeField("HeapUsed", TYPE_I64, "", 0)
	prgrmFldFreeHeap.TotalSize = GetArgSize(TYPE_I64)

	// prgrmFldCaller := MakeField("Caller", TYPE_CUSTOM, "", 0)
	prgrmFldCaller := DeclarationSpecifiersStruct(callStrct.Name, callStrct.Package.Name, false, currentFile, lineNo)
	prgrmFldCaller.Name = "Caller"

	prgrmStrct.AddField(prgrmFldCallCounter)
	prgrmStrct.AddField(prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrmFldCaller)

	pkg.AddStruct(prgrmStrct)
}

func PrimaryIdentifier(ident string) []*CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument(ident, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
		arg.AddType(TypeNames[TYPE_IDENTIFIER])
		// arg.Typ = "ident"
		arg.Name = ident
		arg.Package = pkg

		// expr := &CXExpression{Outputs: []*CXArgument{arg}}
		expr := MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*CXArgument{arg}
		expr.Package = pkg

		return []*CXExpression{expr}
	} else {
		panic(err)
	}
}

// DefineNewScope marks the first and last expressions to define the boundaries of a scope.
func DefineNewScope (exprs []*CXExpression) {
	if len(exprs) > 1 {
		// initialize new scope
		exprs[0].ScopeOperation = SCOPE_NEW
		// remove last scope
		exprs[len(exprs)-1].ScopeOperation = SCOPE_REM
	}
}

// IsArgBasicType returns true if `arg`'s type is a basic type, false otherwise.
func IsArgBasicType(arg *CXArgument) bool {
	switch arg.Type {
		case TYPE_BOOL,
		TYPE_BYTE,
		TYPE_STR,
		TYPE_F32,
		TYPE_F64,
		TYPE_I8,
		TYPE_I16,
		TYPE_I32,
		TYPE_I64,
		TYPE_UI8,
		TYPE_UI16,
		TYPE_UI32, TYPE_UI64:
		return true
	}
	return false
}

// IsAllArgsBasicTypes checks if all the input arguments in an expressions are of basic type.
func IsAllArgsBasicTypes(expr *CXExpression) bool {
	for _, inp := range expr.Inputs {
		if !IsArgBasicType(inp) {
			return false
		}
	}
	return true
}

// IsUndOp returns true if the operator receives undefined types as input parameters.
func IsUndOp(fn *CXFunction) bool {
	switch fn.OpCode {
	case
		OP_UND_EQUAL,
		OP_UND_UNEQUAL,
		OP_UND_BITAND,
		OP_UND_BITXOR,
		OP_UND_BITOR,
		OP_UND_BITCLEAR,
		OP_UND_MUL,
		OP_UND_DIV,
		OP_UND_MOD,
		OP_UND_ADD,
		OP_UND_SUB,
		OP_UND_BITSHL,
		OP_UND_BITSHR,
		OP_UND_LT,
		OP_UND_GT,
		OP_UND_LTEQ,
		OP_UND_GTEQ,
		OP_UND_LEN,
		OP_UND_PRINTF,
		OP_UND_SPRINTF,
		OP_UND_READ:
		return true
	}
	return false
}

// IsUndOpMimicInput returns true if the operator receives undefined types as input parameters but also an operator that needs to mimic its input's type. For example, == should not return its input type, as it is always going to return a boolean.
func IsUndOpMimicInput(fn *CXFunction) bool {
	switch fn.OpCode {
	case
		OP_UND_BITAND,
		OP_UND_BITXOR,
		OP_UND_BITOR,
		OP_UND_BITCLEAR,
		OP_UND_MUL,
		OP_UND_DIV,
		OP_UND_MOD,
		OP_UND_ADD,
		OP_UND_SUB,
		OP_UND_NEG,
		OP_UND_BITSHL, OP_UND_BITSHR:
		return true
	}
	return false
}

// IsUndOp returns true if the operator receives undefined types as input parameters and if it's an operator that only works with basic types. For example, `sa + sb` shouldn't work with struct instances.
func IsUndOpBasicTypes(fn *CXFunction) bool {
	switch fn.OpCode {
	case
		OP_UND_EQUAL,
		OP_UND_UNEQUAL,
		OP_UND_BITAND,
		OP_UND_BITXOR,
		OP_UND_BITOR,
		OP_UND_BITCLEAR,
		OP_UND_MUL,
		OP_UND_DIV,
		OP_UND_MOD,
		OP_UND_ADD,
		OP_UND_SUB,
		OP_UND_BITSHL,
		OP_UND_BITSHR,
		OP_UND_LT,
		OP_UND_GT,
		OP_UND_LTEQ,
		OP_UND_GTEQ,
		OP_UND_PRINTF,
		OP_UND_SPRINTF,
		OP_UND_READ:
		return true
	}
	return false
}
