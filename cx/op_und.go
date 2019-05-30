package cxcore

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

func opLt(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_STR:
		opStrLt(expr, fp)
	case TYPE_I8:
		opI8Lt(expr, fp)
	case TYPE_I16:
		opI16Lt(expr, fp)
	case TYPE_I32:
		opI32Lt(expr, fp)
	case TYPE_I64:
		opI64Lt(expr, fp)
	case TYPE_UI8:
		opUI8Lt(expr, fp)
	case TYPE_UI16:
		opUI16Lt(expr, fp)
	case TYPE_UI32:
		opUI32Lt(expr, fp)
	case TYPE_UI64:
		opUI64Lt(expr, fp)
	case TYPE_F32:
		opF32Lt(expr, fp)
	case TYPE_F64:
		opF64Lt(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opGt(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_STR:
		opStrGt(expr, fp)
	case TYPE_I8:
		opI8Gt(expr, fp)
	case TYPE_I16:
		opI16Gt(expr, fp)
	case TYPE_I32:
		opI32Gt(expr, fp)
	case TYPE_I64:
		opI64Gt(expr, fp)
	case TYPE_UI8:
		opUI8Gt(expr, fp)
	case TYPE_UI16:
		opUI16Gt(expr, fp)
	case TYPE_UI32:
		opUI32Gt(expr, fp)
	case TYPE_UI64:
		opUI64Gt(expr, fp)
	case TYPE_F32:
		opF32Gt(expr, fp)
	case TYPE_F64:
		opF64Gt(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opLteq(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_STR:
		opStrLteq(expr, fp)
	case TYPE_I8:
		opI8Lteq(expr, fp)
	case TYPE_I16:
		opI16Lteq(expr, fp)
	case TYPE_I32:
		opI32Lteq(expr, fp)
	case TYPE_I64:
		opI64Lteq(expr, fp)
	case TYPE_UI8:
		opUI8Lteq(expr, fp)
	case TYPE_UI16:
		opUI16Lteq(expr, fp)
	case TYPE_UI32:
		opUI32Lteq(expr, fp)
	case TYPE_UI64:
		opUI64Lteq(expr, fp)
	case TYPE_F32:
		opF32Lteq(expr, fp)
	case TYPE_F64:
		opF64Lteq(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opGteq(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_STR:
		opStrGteq(expr, fp)
	case TYPE_I8:
		opI8Gteq(expr, fp)
	case TYPE_I16:
		opI16Gteq(expr, fp)
	case TYPE_I32:
		opI32Gteq(expr, fp)
	case TYPE_I64:
		opI64Gteq(expr, fp)
	case TYPE_UI8:
		opUI8Gteq(expr, fp)
	case TYPE_UI16:
		opUI16Gteq(expr, fp)
	case TYPE_UI32:
		opUI32Gteq(expr, fp)
	case TYPE_UI64:
		opUI64Gteq(expr, fp)
	case TYPE_F32:
		opF32Gteq(expr, fp)
	case TYPE_F64:
		opF64Gteq(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opEqual(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_BOOL:
		opBoolEqual(expr, fp)
	case TYPE_STR:
		opStrEq(expr, fp)
	case TYPE_I8:
		opI8Eq(expr, fp)
	case TYPE_I16:
		opI16Eq(expr, fp)
	case TYPE_I32:
		opI32Eq(expr, fp)
	case TYPE_I64:
		opI64Eq(expr, fp)
	case TYPE_UI8:
		opUI8Eq(expr, fp)
	case TYPE_UI16:
		opUI16Eq(expr, fp)
	case TYPE_UI32:
		opUI32Eq(expr, fp)
	case TYPE_UI64:
		opUI64Eq(expr, fp)
	case TYPE_F32:
		opF32Eq(expr, fp)
	case TYPE_F64:
		opF64Eq(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opUnequal(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_BOOL:
		opBoolUnequal(expr, fp)
	case TYPE_STR:
		opStrUneq(expr, fp)
	case TYPE_I8:
		opI8Uneq(expr, fp)
	case TYPE_I16:
		opI16Uneq(expr, fp)
	case TYPE_I32:
		opI32Uneq(expr, fp)
	case TYPE_I64:
		opI64Uneq(expr, fp)
	case TYPE_UI8:
		opUI8Uneq(expr, fp)
	case TYPE_UI16:
		opUI16Uneq(expr, fp)
	case TYPE_UI32:
		opUI32Uneq(expr, fp)
	case TYPE_UI64:
		opUI64Uneq(expr, fp)
	case TYPE_F32:
		opF32Uneq(expr, fp)
	case TYPE_F64:
		opF64Uneq(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitand(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitand(expr, fp)
	case TYPE_I16:
		opI16Bitand(expr, fp)
	case TYPE_I32:
		opI32Bitand(expr, fp)
	case TYPE_I64:
		opI64Bitand(expr, fp)
	case TYPE_UI8:
		opUI8Bitand(expr, fp)
	case TYPE_UI16:
		opUI16Bitand(expr, fp)
	case TYPE_UI32:
		opUI32Bitand(expr, fp)
	case TYPE_UI64:
		opUI64Bitand(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitor(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitor(expr, fp)
	case TYPE_I16:
		opI16Bitor(expr, fp)
	case TYPE_I32:
		opI32Bitor(expr, fp)
	case TYPE_I64:
		opI64Bitor(expr, fp)
	case TYPE_UI8:
		opUI8Bitor(expr, fp)
	case TYPE_UI16:
		opUI16Bitor(expr, fp)
	case TYPE_UI32:
		opUI32Bitor(expr, fp)
	case TYPE_UI64:
		opUI64Bitor(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitxor(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitxor(expr, fp)
	case TYPE_I16:
		opI16Bitxor(expr, fp)
	case TYPE_I32:
		opI32Bitxor(expr, fp)
	case TYPE_I64:
		opI64Bitxor(expr, fp)
	case TYPE_UI8:
		opUI8Bitxor(expr, fp)
	case TYPE_UI16:
		opUI16Bitxor(expr, fp)
	case TYPE_UI32:
		opUI32Bitxor(expr, fp)
	case TYPE_UI64:
		opUI64Bitxor(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opMul(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Mul(expr, fp)
	case TYPE_I16:
		opI16Mul(expr, fp)
	case TYPE_I32:
		opI32Mul(expr, fp)
	case TYPE_I64:
		opI64Mul(expr, fp)
	case TYPE_UI8:
		opUI8Mul(expr, fp)
	case TYPE_UI16:
		opUI16Mul(expr, fp)
	case TYPE_UI32:
		opUI32Mul(expr, fp)
	case TYPE_UI64:
		opUI64Mul(expr, fp)
	case TYPE_F32:
		opF32Mul(expr, fp)
	case TYPE_F64:
		opF64Mul(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opDiv(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Div(expr, fp)
	case TYPE_I16:
		opI16Div(expr, fp)
	case TYPE_I32:
		opI32Div(expr, fp)
	case TYPE_I64:
		opI64Div(expr, fp)
	case TYPE_UI8:
		opUI8Div(expr, fp)
	case TYPE_UI16:
		opUI16Div(expr, fp)
	case TYPE_UI32:
		opUI32Div(expr, fp)
	case TYPE_UI64:
		opUI64Div(expr, fp)
	case TYPE_F32:
		opF32Div(expr, fp)
	case TYPE_F64:
		opF64Div(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opMod(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Mod(expr, fp)
	case TYPE_I16:
		opI16Mod(expr, fp)
	case TYPE_I32:
		opI32Mod(expr, fp)
	case TYPE_I64:
		opI64Mod(expr, fp)
	case TYPE_UI8:
		opUI8Mod(expr, fp)
	case TYPE_UI16:
		opUI16Mod(expr, fp)
	case TYPE_UI32:
		opUI32Mod(expr, fp)
	case TYPE_UI64:
		opUI64Mod(expr, fp)
	case TYPE_F32:
		opF32Mod(expr, fp)
	case TYPE_F64:
		opF64Mod(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opAdd(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Add(expr, fp)
	case TYPE_I16:
		opI16Add(expr, fp)
	case TYPE_I32:
		opI32Add(expr, fp)
	case TYPE_I64:
		opI64Add(expr, fp)
	case TYPE_UI8:
		opUI8Add(expr, fp)
	case TYPE_UI16:
		opUI16Add(expr, fp)
	case TYPE_UI32:
		opUI32Add(expr, fp)
	case TYPE_UI64:
		opUI64Add(expr, fp)
	case TYPE_F32:
		opF32Add(expr, fp)
	case TYPE_F64:
		opF64Add(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opSub(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Sub(expr, fp)
	case TYPE_I16:
		opI16Sub(expr, fp)
	case TYPE_I32:
		opI32Sub(expr, fp)
	case TYPE_I64:
		opI64Sub(expr, fp)
	case TYPE_UI8:
		opUI8Sub(expr, fp)
	case TYPE_UI16:
		opUI16Sub(expr, fp)
	case TYPE_UI32:
		opUI32Sub(expr, fp)
	case TYPE_UI64:
		opUI64Sub(expr, fp)
	case TYPE_F32:
		opF32Sub(expr, fp)
	case TYPE_F64:
		opF64Sub(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opNeg(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Neg(expr, fp)
	case TYPE_I16:
		opI16Neg(expr, fp)
	case TYPE_I32:
		opI32Neg(expr, fp)
	case TYPE_I64:
		opI64Neg(expr, fp)
	case TYPE_F32:
		opF32Neg(expr, fp)
	case TYPE_F64:
		opF64Neg(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitshl(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitshl(expr, fp)
	case TYPE_I16:
		opI16Bitshl(expr, fp)
	case TYPE_I32:
		opI32Bitshl(expr, fp)
	case TYPE_I64:
		opI64Bitshl(expr, fp)
	case TYPE_UI8:
		opUI8Bitshl(expr, fp)
	case TYPE_UI16:
		opUI16Bitshl(expr, fp)
	case TYPE_UI32:
		opUI32Bitshl(expr, fp)
	case TYPE_UI64:
		opUI64Bitshl(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitshr(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitshr(expr, fp)
	case TYPE_I16:
		opI16Bitshr(expr, fp)
	case TYPE_I32:
		opI32Bitshr(expr, fp)
	case TYPE_I64:
		opI64Bitshr(expr, fp)
	case TYPE_UI8:
		opUI8Bitshr(expr, fp)
	case TYPE_UI16:
		opUI16Bitshr(expr, fp)
	case TYPE_UI32:
		opUI32Bitshr(expr, fp)
	case TYPE_UI64:
		opUI64Bitshr(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitclear(expr *CXExpression, fp int) {
	switch expr.Inputs[0].Type {
	case TYPE_I8:
		opI8Bitclear(expr, fp)
	case TYPE_I16:
		opI16Bitclear(expr, fp)
	case TYPE_I32:
		opI32Bitclear(expr, fp)
	case TYPE_I64:
		opI64Bitclear(expr, fp)
	case TYPE_UI8:
		opUI8Bitclear(expr, fp)
	case TYPE_UI16:
		opUI16Bitclear(expr, fp)
	case TYPE_UI32:
		opUI32Bitclear(expr, fp)
	case TYPE_UI64:
		opUI64Bitclear(expr, fp)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opLen(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	elt := GetAssignmentElement(inp1)

	if elt.IsSlice || elt.Type == TYPE_AFF {
		var sliceOffset = GetSliceOffset(fp, inp1)
		var sliceLen []byte
		if sliceOffset > 0 {
			sliceLen = GetSliceHeader(sliceOffset)[4:8]
		} else if sliceOffset == 0 {
			sliceLen = FromI32(0)
		} else {
			panic(CX_RUNTIME_ERROR)
		}

		WriteMemory(GetFinalOffset(fp, out1), sliceLen)
		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == TYPE_STR && elt.Lengths == nil {
		var strOffset = GetStrOffset(fp, inp1)
		WriteMemory(GetFinalOffset(fp, out1), PROGRAM.Memory[strOffset:strOffset+STR_HEADER_SIZE])
	} else {
		outB1 := FromI32(int32(elt.Lengths[0]))
		WriteMemory(GetFinalOffset(fp, out1), outB1)
	}
}

func opAppend(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != inp2.Type || inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	inputSliceOffset := GetSliceOffset(fp, inp1)

	var obj []byte
	if inp2.Type == TYPE_STR || inp2.Type == TYPE_AFF {
		obj = encoder.SerializeAtomic(int32(GetStrOffset(fp, inp2)))
	} else {
		obj = ReadMemory(GetFinalOffset(fp, inp2), inp2)
	}

	outputSliceOffset = int32(SliceAppend(outputSliceOffset, inputSliceOffset, obj))
	copy(PROGRAM.Memory[outputSlicePointer:], encoder.SerializeAtomic(outputSliceOffset))
}

func opResize(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))

	inputSliceOffset := GetSliceOffset(fp, inp1)

	outputSliceOffset = int32(SliceResize(outputSliceOffset, inputSliceOffset, ReadI32(fp, inp2), GetAssignmentElement(inp1).TotalSize))
	copy(PROGRAM.Memory[outputSlicePointer:], encoder.SerializeAtomic(outputSliceOffset))
}

func opInsert(expr *CXExpression, fp int) {
	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]

	if inp1.Type != inp3.Type || inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))

	inputSliceOffset := GetSliceOffset(fp, inp1)

	var obj []byte
	if inp3.Type == TYPE_STR || inp3.Type == TYPE_AFF {
		obj = encoder.SerializeAtomic(int32(GetStrOffset(fp, inp3)))
	} else {
		obj = ReadMemory(GetFinalOffset(fp, inp3), inp3)
	}

	outputSliceOffset = int32(SliceInsert(outputSliceOffset, inputSliceOffset, ReadI32(fp, inp2), obj))
	copy(PROGRAM.Memory[outputSlicePointer:], encoder.SerializeAtomic(outputSliceOffset))
}

func opRemove(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))

	inputSliceOffset := GetSliceOffset(fp, inp1)

	outputSliceOffset = int32(SliceRemove(outputSliceOffset, inputSliceOffset, ReadI32(fp, inp2), int32(GetAssignmentElement(inp1).TotalSize)))
	copy(PROGRAM.Memory[outputSlicePointer:], encoder.SerializeAtomic(outputSliceOffset))
}

func opCopy(expr *CXExpression, fp int) {
	dstInput := expr.Inputs[0]
	srcInput := expr.Inputs[1]
	dstOffset := GetSliceOffset(fp, dstInput)
	srcOffset := GetSliceOffset(fp, srcInput)

	dstElem := GetAssignmentElement(dstInput)
	srcElem := GetAssignmentElement(srcInput)

	if dstInput.Type != srcInput.Type || !dstElem.IsSlice || !srcElem.IsSlice || dstElem.TotalSize != srcElem.TotalSize {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	var count int
	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		count = copy(GetSliceData(dstOffset, dstElem.TotalSize), GetSliceData(srcOffset, srcElem.TotalSize))
		if count%dstElem.TotalSize != 0 {
			panic(CX_RUNTIME_ERROR)
		}
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(count/dstElem.TotalSize)))
}

func buildString(expr *CXExpression, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ReadStr(fp, inp1)

	var res []byte
	var specifiersCounter int
	var lenStr = int(len(fmtStr))

	for c := 0; c < len(fmtStr); c++ {
		var nextCh byte
		ch := fmtStr[c]
		if c < lenStr-1 {
			nextCh = fmtStr[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}
		}
		if ch == '%' {
			if specifiersCounter+1 == len(expr.Inputs) {
				res = append(res, []byte(fmt.Sprintf("%%!%c(MISSING)", nextCh))...)
				c++
				continue
			}

			inp := expr.Inputs[specifiersCounter+1]
			switch nextCh {
			case 's':
				res = append(res, []byte(checkForEscapedChars(ReadStr(fp, inp)))...)
			case 'd':
				switch inp.Type {
				case TYPE_I8:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI8(fp, inp)), 10))...)
				case TYPE_I16:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI16(fp, inp)), 10))...)
				case TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI32(fp, inp)), 10))...)
				case TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ReadI64(fp, inp), 10))...)
				case TYPE_UI8:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI8(fp, inp)), 10))...)
				case TYPE_UI16:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI16(fp, inp)), 10))...)
				case TYPE_UI32:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI32(fp, inp)), 10))...)
				case TYPE_UI64:
					res = append(res, []byte(strconv.FormatUint(ReadUI64(fp, inp), 10))...)
				}
			case 'f':
				switch inp.Type {
				case TYPE_F32:
					res = append(res, []byte(strconv.FormatFloat(float64(ReadF32(fp, inp)), 'f', 7, 32))...)
				case TYPE_F64:
					res = append(res, []byte(strconv.FormatFloat(ReadF64(fp, inp), 'f', 16, 64))...)
				}
			case 'v':
				res = append(res, []byte(GetPrintableValue(fp, inp))...)
			}
			c++
			specifiersCounter++
		} else {
			res = append(res, ch)
		}
	}

	if specifiersCounter != len(expr.Inputs)-1 {
		extra := "%!(EXTRA "
		// for _, inp := range expr.Inputs[:specifiersCounter] {
		lInps := len(expr.Inputs[specifiersCounter+1:])
		for c := 0; c < lInps; c++ {
			inp := expr.Inputs[specifiersCounter+1+c]
			elt := GetAssignmentElement(inp)
			typ := ""
			_ = typ
			if elt.CustomType != nil {
				// then it's custom type
				typ = elt.CustomType.Name
			} else {
				// then it's native type
				typ = TypeNames[elt.Type]
			}

			if c == lInps-1 {
				extra += fmt.Sprintf("%s=%s", typ, GetPrintableValue(fp, elt))
			} else {
				extra += fmt.Sprintf("%s=%s, ", typ, GetPrintableValue(fp, elt))
			}

		}

		extra += ")"

		res = append(res, []byte(extra)...)
	}

	return res
}

func opSprintf(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	byts := encoder.Serialize(string(buildString(expr, fp)))
	WriteObject(out1Offset, byts)
}

func opPrintf(expr *CXExpression, fp int) {
	fmt.Print(string(buildString(expr, fp)))
}

func opRead(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	// text = strings.Trim(text, " \n")
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	if err != nil {
		panic("")
	}
	byts := encoder.Serialize(text)
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)

	var header = make([]byte, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset + OBJECT_HEADER_SIZE))

	WriteMemory(out1Offset, off)
}
