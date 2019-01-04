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
		// if !CheckSameNativeType(expr) {
		//      panic("wrong types")
		// }
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

// This function writes those bytes to PRGRM.Data
func WritePrimary(typ int, byts []byte, isGlobal bool) []*CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := MakeArgument("", CurrentFile, LineNo)
		arg.AddType(TypeNames[typ])
		arg.Package = pkg
		// arg.Program = PRGRM

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

func CompilationError(currentFile string, lineNo int) string {
	FoundCompileErrors = true
	return ErrorHeader(currentFile, lineNo)
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
		arg := MakeArgument(ident, CurrentFile, LineNo)
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
