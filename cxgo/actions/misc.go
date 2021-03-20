package actions

import (
	"os"
	"runtime"

	"github.com/skycoin/cx/cx"
)

func SelectProgram(prgrm *cxcore.CXProgram) {
	PRGRM = prgrm
}

func SetCorrectArithmeticOp(expr *cxcore.CXExpression) {
	if expr.Operator == nil || len(expr.Outputs) < 1 {
		return
	}
	op := expr.Operator
	typ := expr.Outputs[0].Type

	if cxcore.CheckArithmeticOp(expr) {
		switch op.OpCode {
		case cxcore.OP_I32_MUL:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_MUL]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_MUL]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_MUL]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_MUL]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_MUL]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_MUL]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_MUL]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_MUL]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_MUL]
			}
		case cxcore.OP_I32_DIV:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_DIV]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_DIV]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_DIV]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_DIV]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_DIV]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_DIV]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_DIV]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_DIV]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_DIV]
			}
		case cxcore.OP_I32_MOD:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_MOD]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_MOD]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_MOD]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_MOD]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_MOD]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_MOD]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_MOD]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_MOD]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_MOD]
			}

		case cxcore.OP_I32_ADD:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_ADD]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_ADD]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_ADD]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_ADD]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_ADD]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_ADD]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_ADD]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_ADD]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_ADD]
			}

		case cxcore.OP_I32_SUB:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_SUB]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_SUB]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_SUB]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_SUB]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_SUB]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_SUB]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_SUB]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_SUB]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_SUB]
			}

		case cxcore.OP_I32_NEG:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_NEG]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_NEG]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_NEG]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_NEG]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_NEG]
			}

		case cxcore.OP_I32_BITSHL:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_BITSHL]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_BITSHL]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_BITSHL]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_BITSHL]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_BITSHL]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_BITSHL]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_BITSHL]
			}

		case cxcore.OP_I32_BITSHR:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_BITSHR]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_BITSHR]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_BITSHR]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_BITSHR]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_BITSHR]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_BITSHR]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_BITSHR]
			}

		case cxcore.OP_I32_LT:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_LT]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_LT]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_LT]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_LT]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_LT]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_LT]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_LT]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_LT]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_LT]
			}

		case cxcore.OP_I32_GT:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_GT]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_GT]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_GT]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_GT]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_GT]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_GT]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_GT]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_GT]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_GT]
			}

		case cxcore.OP_I32_LTEQ:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_LTEQ]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_LTEQ]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_LTEQ]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_LTEQ]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_LTEQ]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_LTEQ]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_LTEQ]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_LTEQ]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_LTEQ]
			}

		case cxcore.OP_I32_GTEQ:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_GTEQ]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_GTEQ]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_GTEQ]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_GTEQ]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_GTEQ]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_GTEQ]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_GTEQ]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_GTEQ]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_GTEQ]
			}

		case cxcore.OP_I32_EQ:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_EQ]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_EQ]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_EQ]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_EQ]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_EQ]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_EQ]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_EQ]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_EQ]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_EQ]
			}

		case cxcore.OP_I32_UNEQ:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_UNEQ]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_UNEQ]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_UNEQ]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_UNEQ]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_UNEQ]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_UNEQ]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_UNEQ]
			case cxcore.TYPE_F32:
				expr.Operator = cxcore.Natives[cxcore.OP_F32_UNEQ]
			case cxcore.TYPE_F64:
				expr.Operator = cxcore.Natives[cxcore.OP_F64_UNEQ]
			}

		case cxcore.OP_I32_BITAND:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_BITAND]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_BITAND]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_BITAND]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_BITAND]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_BITAND]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_BITAND]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_BITAND]
			}

		case cxcore.OP_I32_BITXOR:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_BITXOR]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_BITXOR]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_BITXOR]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_BITXOR]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_BITXOR]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_BITXOR]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_BITXOR]
			}

		case cxcore.OP_I32_BITOR:
			switch typ {
			case cxcore.TYPE_I8:
				expr.Operator = cxcore.Natives[cxcore.OP_I8_BITOR]
			case cxcore.TYPE_I16:
				expr.Operator = cxcore.Natives[cxcore.OP_I16_BITOR]
			case cxcore.TYPE_I32:
			case cxcore.TYPE_I64:
				expr.Operator = cxcore.Natives[cxcore.OP_I64_BITOR]
			case cxcore.TYPE_UI8:
				expr.Operator = cxcore.Natives[cxcore.OP_UI8_BITOR]
			case cxcore.TYPE_UI16:
				expr.Operator = cxcore.Natives[cxcore.OP_UI16_BITOR]
			case cxcore.TYPE_UI32:
				expr.Operator = cxcore.Natives[cxcore.OP_UI32_BITOR]
			case cxcore.TYPE_UI64:
				expr.Operator = cxcore.Natives[cxcore.OP_UI64_BITOR]
			}
		}
	}
}

// hasDeclSpec determines if an argument has certain declaration specifier
func hasDeclSpec(arg *cxcore.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DeclarationSpecifiers {
		if s == spec {
			found = true
		}
	}
	return found
}

// hasDerefOp determines if an argument has certain dereference operation
func hasDerefOp(arg *cxcore.CXArgument, spec int) bool {
	found := false
	for _, s := range arg.DereferenceOperations {
		if s == spec {
			found = true
		}
	}
	return found
}

// This function writes those bytes to PRGRM.Data
func WritePrimary(typ int, byts []byte, isGlobal bool) []*cxcore.CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := cxcore.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(cxcore.TypeNames[typ])
		arg.Package = pkg

		var size int

		size = len(byts)

		arg.Size = cxcore.GetArgSize(typ)
		arg.TotalSize = size
		arg.Offset = DataOffset

		if arg.Type == cxcore.TYPE_STR || arg.Type == cxcore.TYPE_AFF {
			arg.PassBy = cxcore.PASSBY_REFERENCE
			arg.Size = cxcore.TYPE_POINTER_SIZE
			arg.TotalSize = cxcore.TYPE_POINTER_SIZE
		}

		// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
		// after the stack segment. These bytes are used to allocate the data segment
		// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
		// we'll start appending the bytes to PRGRM.Memory.
		// After compilation, we calculate how many bytes we need to add to have a heap segment
		// equal to `minHeapSize()` that is allocated after the data segment.
		if size+DataOffset > len(PRGRM.Memory) {
			var i int
			// First we need to fill the remaining free bytes in
			// the current `PRGRM.Memory` slice.
			for i = 0; i < len(PRGRM.Memory)-DataOffset; i++ {
				PRGRM.Memory[DataOffset+i] = byts[i]
			}
			// Then we append the bytes that didn't fit.
			PRGRM.Memory = append(PRGRM.Memory, byts[i:]...)
		} else {
			for i, byt := range byts {
				PRGRM.Memory[DataOffset+i] = byt
			}
		}
		DataOffset += size

		expr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
		expr.Package = pkg
		expr.Outputs = append(expr.Outputs, arg)
		return []*cxcore.CXExpression{expr}
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

func StructLiteralFields(ident string) *cxcore.CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := cxcore.MakeArgument("", CurrentFile, LineNo)
		arg.AddType(cxcore.TypeNames[cxcore.TYPE_IDENTIFIER])
		arg.Name = ident
		arg.Package = pkg

		expr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*cxcore.CXArgument{arg}
		expr.Package = pkg

		return expr
	} else {
		panic(err)
	}
}

func AffordanceStructs(pkg *cxcore.CXPackage, currentFile string, lineNo int) {
	// Argument type
	argStrct := cxcore.MakeStruct("Argument")
	// argStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	argFldName := cxcore.MakeField("Name", cxcore.TYPE_STR, "", 0)
	argFldName.TotalSize = cxcore.GetArgSize(cxcore.TYPE_STR)
	argFldIndex := cxcore.MakeField("Index", cxcore.TYPE_I32, "", 0)
	argFldIndex.TotalSize = cxcore.GetArgSize(cxcore.TYPE_I32)
	argFldType := cxcore.MakeField("Type", cxcore.TYPE_STR, "", 0)
	argFldType.TotalSize = cxcore.GetArgSize(cxcore.TYPE_STR)

	argStrct.AddField(argFldName)
	argStrct.AddField(argFldIndex)
	argStrct.AddField(argFldType)

	pkg.AddStruct(argStrct)

	// Expression type
	exprStrct := cxcore.MakeStruct("Expression")
	// exprStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	exprFldOperator := cxcore.MakeField("Operator", cxcore.TYPE_STR, "", 0)

	exprStrct.AddField(exprFldOperator)

	pkg.AddStruct(exprStrct)

	// Function type
	fnStrct := cxcore.MakeStruct("Function")
	// fnStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_STR)

	fnFldName := cxcore.MakeField("Name", cxcore.TYPE_STR, "", 0)
	fnFldName.TotalSize = cxcore.GetArgSize(cxcore.TYPE_STR)

	fnFldInpSig := cxcore.MakeField("InputSignature", cxcore.TYPE_STR, "", 0)
	fnFldInpSig.Size = cxcore.GetArgSize(cxcore.TYPE_STR)
	fnFldInpSig = DeclarationSpecifiers(fnFldInpSig, []int{0}, cxcore.DECL_SLICE)

	fnFldOutSig := cxcore.MakeField("OutputSignature", cxcore.TYPE_STR, "", 0)
	fnFldOutSig.Size = cxcore.GetArgSize(cxcore.TYPE_STR)
	fnFldOutSig = DeclarationSpecifiers(fnFldOutSig, []int{0}, cxcore.DECL_SLICE)

	fnStrct.AddField(fnFldName)
	fnStrct.AddField(fnFldInpSig)

	fnStrct.AddField(fnFldOutSig)

	pkg.AddStruct(fnStrct)

	// Structure type
	strctStrct := cxcore.MakeStruct("Structure")
	// strctStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	strctFldName := cxcore.MakeField("Name", cxcore.TYPE_STR, "", 0)
	strctFldName.TotalSize = cxcore.GetArgSize(cxcore.TYPE_STR)

	strctStrct.AddField(strctFldName)

	pkg.AddStruct(strctStrct)

	// Package type
	pkgStrct := cxcore.MakeStruct("Structure")
	// pkgStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR)

	pkgFldName := cxcore.MakeField("Name", cxcore.TYPE_STR, "", 0)

	pkgStrct.AddField(pkgFldName)

	pkg.AddStruct(pkgStrct)

	// Caller type
	callStrct := cxcore.MakeStruct("Caller")
	// callStrct.Size = cxcore.GetArgSize(cxcore.TYPE_STR) + cxcore.GetArgSize(cxcore.TYPE_I32)

	callFldFnName := cxcore.MakeField("FnName", cxcore.TYPE_STR, "", 0)
	callFldFnName.TotalSize = cxcore.GetArgSize(cxcore.TYPE_STR)
	callFldFnSize := cxcore.MakeField("FnSize", cxcore.TYPE_I32, "", 0)
	callFldFnSize.TotalSize = cxcore.GetArgSize(cxcore.TYPE_I32)

	callStrct.AddField(callFldFnName)
	callStrct.AddField(callFldFnSize)

	pkg.AddStruct(callStrct)

	// Program type
	prgrmStrct := cxcore.MakeStruct("Program")
	// prgrmStrct.Size = cxcore.GetArgSize(cxcore.TYPE_I32) + cxcore.GetArgSize(cxcore.TYPE_I64)

	prgrmFldCallCounter := cxcore.MakeField("CallCounter", cxcore.TYPE_I32, "", 0)
	prgrmFldCallCounter.TotalSize = cxcore.GetArgSize(cxcore.TYPE_I32)
	prgrmFldFreeHeap := cxcore.MakeField("HeapUsed", cxcore.TYPE_I64, "", 0)
	prgrmFldFreeHeap.TotalSize = cxcore.GetArgSize(cxcore.TYPE_I64)

	// prgrmFldCaller := cxcore.MakeField("Caller", cxcore.TYPE_CUSTOM, "", 0)
	prgrmFldCaller := DeclarationSpecifiersStruct(callStrct.Name, callStrct.Package.Name, false, currentFile, lineNo)
	prgrmFldCaller.Name = "Caller"

	prgrmStrct.AddField(prgrmFldCallCounter)
	prgrmStrct.AddField(prgrmFldFreeHeap)
	prgrmStrct.AddField(prgrmFldCaller)

	pkg.AddStruct(prgrmStrct)
}

func PrimaryIdentifier(ident string) []*cxcore.CXExpression {
	if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
		arg := cxcore.MakeArgument(ident, CurrentFile, LineNo) // fix: line numbers in errors sometimes report +1 or -1. Issue #195
		arg.AddType(cxcore.TypeNames[cxcore.TYPE_IDENTIFIER])
		// arg.Typ = "ident"
		arg.Name = ident
		arg.Package = pkg

		// expr := &cxcore.CXExpression{Outputs: []*cxcore.CXArgument{arg}}
		expr := cxcore.MakeExpression(nil, CurrentFile, LineNo)
		expr.Outputs = []*cxcore.CXArgument{arg}
		expr.Package = pkg

		return []*cxcore.CXExpression{expr}
	} else {
		panic(err)
	}
}

// DefineNewScope marks the first and last expressions to define the boundaries of a scope.
func DefineNewScope(exprs []*cxcore.CXExpression) {
	if len(exprs) > 1 {
		// initialize new scope
		exprs[0].ScopeOperation = cxcore.SCOPE_NEW
		// remove last scope
		exprs[len(exprs)-1].ScopeOperation = cxcore.SCOPE_REM
	}
}

// IsArgBasicType returns true if `arg`'s type is a basic type, false otherwise.
func IsArgBasicType(arg *cxcore.CXArgument) bool {
	switch arg.Type {
	case cxcore.TYPE_BOOL,
		cxcore.TYPE_STR,
		cxcore.TYPE_F32,
		cxcore.TYPE_F64,
		cxcore.TYPE_I8,
		cxcore.TYPE_I16,
		cxcore.TYPE_I32,
		cxcore.TYPE_I64,
		cxcore.TYPE_UI8,
		cxcore.TYPE_UI16,
		cxcore.TYPE_UI32,
		cxcore.TYPE_UI64:
		return true
	}
	return false
}

// IsAllArgsBasicTypes checks if all the input arguments in an expressions are of basic type.
func IsAllArgsBasicTypes(expr *cxcore.CXExpression) bool {
	for _, inp := range expr.Inputs {
		if !IsArgBasicType(inp) {
			return false
		}
	}
	return true
}

// IsUndOp returns true if the operator receives undefined types as input parameters.
func IsUndOp(fn *cxcore.CXFunction) bool {
	switch fn.OpCode {
	case
		cxcore.OP_UND_EQUAL,
		cxcore.OP_UND_UNEQUAL,
		cxcore.OP_UND_BITAND,
		cxcore.OP_UND_BITXOR,
		cxcore.OP_UND_BITOR,
		cxcore.OP_UND_BITCLEAR,
		cxcore.OP_UND_MUL,
		cxcore.OP_UND_DIV,
		cxcore.OP_UND_MOD,
		cxcore.OP_UND_ADD,
		cxcore.OP_UND_SUB,
		cxcore.OP_UND_BITSHL,
		cxcore.OP_UND_BITSHR,
		cxcore.OP_UND_LT,
		cxcore.OP_UND_GT,
		cxcore.OP_UND_LTEQ,
		cxcore.OP_UND_GTEQ,
		cxcore.OP_UND_LEN,
		cxcore.OP_UND_PRINTF,
		cxcore.OP_UND_SPRINTF,
		cxcore.OP_UND_READ:
		return true
	}
	return false
}

// IsUndOpMimicInput returns true if the operator receives undefined types as input parameters but also an operator that needs to mimic its input's type. For example, == should not return its input type, as it is always going to return a boolean.
func IsUndOpMimicInput(fn *cxcore.CXFunction) bool {
	switch fn.OpCode {
	case
		cxcore.OP_UND_BITAND,
		cxcore.OP_UND_BITXOR,
		cxcore.OP_UND_BITOR,
		cxcore.OP_UND_BITCLEAR,
		cxcore.OP_UND_MUL,
		cxcore.OP_UND_DIV,
		cxcore.OP_UND_MOD,
		cxcore.OP_UND_ADD,
		cxcore.OP_UND_SUB,
		cxcore.OP_UND_NEG,
		cxcore.OP_UND_BITSHL, cxcore.OP_UND_BITSHR:
		return true
	}
	return false
}

// IsUndOp returns true if the operator receives undefined types as input parameters and if it's an operator that only works with basic types. For example, `sa + sb` shouldn't work with struct instances.
func IsUndOpBasicTypes(fn *cxcore.CXFunction) bool {
	switch fn.OpCode {
	case
		cxcore.OP_UND_EQUAL,
		cxcore.OP_UND_UNEQUAL,
		cxcore.OP_UND_BITAND,
		cxcore.OP_UND_BITXOR,
		cxcore.OP_UND_BITOR,
		cxcore.OP_UND_BITCLEAR,
		cxcore.OP_UND_MUL,
		cxcore.OP_UND_DIV,
		cxcore.OP_UND_MOD,
		cxcore.OP_UND_ADD,
		cxcore.OP_UND_SUB,
		cxcore.OP_UND_NEG,
		cxcore.OP_UND_BITSHL,
		cxcore.OP_UND_BITSHR,
		cxcore.OP_UND_LT,
		cxcore.OP_UND_GT,
		cxcore.OP_UND_LTEQ,
		cxcore.OP_UND_GTEQ,
		cxcore.OP_UND_PRINTF,
		cxcore.OP_UND_SPRINTF,
		cxcore.OP_UND_READ:
		return true
	}
	return false
}

// UserHome returns the current user home path. Code taken from fiber-init.
func UserHome() string {
	// os/user relies on cgo which is disabled when cross compiling
	// use fallbacks for various OSes instead
	// usr, err := user.Current()
	// if err == nil {
	// 	return usr.HomeDir
	// }
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}
