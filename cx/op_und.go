package cxcore

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opLt(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_STR:
		opStrLt(prgrm)
	case TYPE_I8:
		opI8Lt(prgrm)
	case TYPE_I16:
		opI16Lt(prgrm)
	case TYPE_I32:
		opI32Lt(prgrm)
	case TYPE_I64:
		opI64Lt(prgrm)
	case TYPE_UI8:
		opUI8Lt(prgrm)
	case TYPE_UI16:
		opUI16Lt(prgrm)
	case TYPE_UI32:
		opUI32Lt(prgrm)
	case TYPE_UI64:
		opUI64Lt(prgrm)
	case TYPE_F32:
		opF32Lt(prgrm)
	case TYPE_F64:
		opF64Lt(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opGt(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_STR:
		opStrGt(prgrm)
	case TYPE_I8:
		opI8Gt(prgrm)
	case TYPE_I16:
		opI16Gt(prgrm)
	case TYPE_I32:
		opI32Gt(prgrm)
	case TYPE_I64:
		opI64Gt(prgrm)
	case TYPE_UI8:
		opUI8Gt(prgrm)
	case TYPE_UI16:
		opUI16Gt(prgrm)
	case TYPE_UI32:
		opUI32Gt(prgrm)
	case TYPE_UI64:
		opUI64Gt(prgrm)
	case TYPE_F32:
		opF32Gt(prgrm)
	case TYPE_F64:
		opF64Gt(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opLteq(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_STR:
		opStrLteq(prgrm)
	case TYPE_I8:
		opI8Lteq(prgrm)
	case TYPE_I16:
		opI16Lteq(prgrm)
	case TYPE_I32:
		opI32Lteq(prgrm)
	case TYPE_I64:
		opI64Lteq(prgrm)
	case TYPE_UI8:
		opUI8Lteq(prgrm)
	case TYPE_UI16:
		opUI16Lteq(prgrm)
	case TYPE_UI32:
		opUI32Lteq(prgrm)
	case TYPE_UI64:
		opUI64Lteq(prgrm)
	case TYPE_F32:
		opF32Lteq(prgrm)
	case TYPE_F64:
		opF64Lteq(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opGteq(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_STR:
		opStrGteq(prgrm)
	case TYPE_I8:
		opI8Gteq(prgrm)
	case TYPE_I16:
		opI16Gteq(prgrm)
	case TYPE_I32:
		opI32Gteq(prgrm)
	case TYPE_I64:
		opI64Gteq(prgrm)
	case TYPE_UI8:
		opUI8Gteq(prgrm)
	case TYPE_UI16:
		opUI16Gteq(prgrm)
	case TYPE_UI32:
		opUI32Gteq(prgrm)
	case TYPE_UI64:
		opUI64Gteq(prgrm)
	case TYPE_F32:
		opF32Gteq(prgrm)
	case TYPE_F64:
		opF64Gteq(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opEqual(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_BOOL:
		opBoolEqual(prgrm)
	case TYPE_STR:
		opStrEq(prgrm)
	case TYPE_I8:
		opI8Eq(prgrm)
	case TYPE_I16:
		opI16Eq(prgrm)
	case TYPE_I32:
		opI32Eq(prgrm)
	case TYPE_I64:
		opI64Eq(prgrm)
	case TYPE_UI8:
		opUI8Eq(prgrm)
	case TYPE_UI16:
		opUI16Eq(prgrm)
	case TYPE_UI32:
		opUI32Eq(prgrm)
	case TYPE_UI64:
		opUI64Eq(prgrm)
	case TYPE_F32:
		opF32Eq(prgrm)
	case TYPE_F64:
		opF64Eq(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opUnequal(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_BOOL:
		opBoolUnequal(prgrm)
	case TYPE_STR:
		opStrUneq(prgrm)
	case TYPE_I8:
		opI8Uneq(prgrm)
	case TYPE_I16:
		opI16Uneq(prgrm)
	case TYPE_I32:
		opI32Uneq(prgrm)
	case TYPE_I64:
		opI64Uneq(prgrm)
	case TYPE_UI8:
		opUI8Uneq(prgrm)
	case TYPE_UI16:
		opUI16Uneq(prgrm)
	case TYPE_UI32:
		opUI32Uneq(prgrm)
	case TYPE_UI64:
		opUI64Uneq(prgrm)
	case TYPE_F32:
		opF32Uneq(prgrm)
	case TYPE_F64:
		opF64Uneq(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitand(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitand(prgrm)
	case TYPE_I16:
		opI16Bitand(prgrm)
	case TYPE_I32:
		opI32Bitand(prgrm)
	case TYPE_I64:
		opI64Bitand(prgrm)
	case TYPE_UI8:
		opUI8Bitand(prgrm)
	case TYPE_UI16:
		opUI16Bitand(prgrm)
	case TYPE_UI32:
		opUI32Bitand(prgrm)
	case TYPE_UI64:
		opUI64Bitand(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitor(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitor(prgrm)
	case TYPE_I16:
		opI16Bitor(prgrm)
	case TYPE_I32:
		opI32Bitor(prgrm)
	case TYPE_I64:
		opI64Bitor(prgrm)
	case TYPE_UI8:
		opUI8Bitor(prgrm)
	case TYPE_UI16:
		opUI16Bitor(prgrm)
	case TYPE_UI32:
		opUI32Bitor(prgrm)
	case TYPE_UI64:
		opUI64Bitor(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitxor(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitxor(prgrm)
	case TYPE_I16:
		opI16Bitxor(prgrm)
	case TYPE_I32:
		opI32Bitxor(prgrm)
	case TYPE_I64:
		opI64Bitxor(prgrm)
	case TYPE_UI8:
		opUI8Bitxor(prgrm)
	case TYPE_UI16:
		opUI16Bitxor(prgrm)
	case TYPE_UI32:
		opUI32Bitxor(prgrm)
	case TYPE_UI64:
		opUI64Bitxor(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opMul(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Mul(prgrm)
	case TYPE_I16:
		opI16Mul(prgrm)
	case TYPE_I32:
		opI32Mul(prgrm)
	case TYPE_I64:
		opI64Mul(prgrm)
	case TYPE_UI8:
		opUI8Mul(prgrm)
	case TYPE_UI16:
		opUI16Mul(prgrm)
	case TYPE_UI32:
		opUI32Mul(prgrm)
	case TYPE_UI64:
		opUI64Mul(prgrm)
	case TYPE_F32:
		opF32Mul(prgrm)
	case TYPE_F64:
		opF64Mul(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opDiv(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Div(prgrm)
	case TYPE_I16:
		opI16Div(prgrm)
	case TYPE_I32:
		opI32Div(prgrm)
	case TYPE_I64:
		opI64Div(prgrm)
	case TYPE_UI8:
		opUI8Div(prgrm)
	case TYPE_UI16:
		opUI16Div(prgrm)
	case TYPE_UI32:
		opUI32Div(prgrm)
	case TYPE_UI64:
		opUI64Div(prgrm)
	case TYPE_F32:
		opF32Div(prgrm)
	case TYPE_F64:
		opF64Div(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opMod(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Mod(prgrm)
	case TYPE_I16:
		opI16Mod(prgrm)
	case TYPE_I32:
		opI32Mod(prgrm)
	case TYPE_I64:
		opI64Mod(prgrm)
	case TYPE_UI8:
		opUI8Mod(prgrm)
	case TYPE_UI16:
		opUI16Mod(prgrm)
	case TYPE_UI32:
		opUI32Mod(prgrm)
	case TYPE_UI64:
		opUI64Mod(prgrm)
	case TYPE_F32:
		opF32Mod(prgrm)
	case TYPE_F64:
		opF64Mod(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opAdd(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Add(prgrm)
	case TYPE_I16:
		opI16Add(prgrm)
	case TYPE_I32:
		opI32Add(prgrm)
	case TYPE_I64:
		opI64Add(prgrm)
	case TYPE_UI8:
		opUI8Add(prgrm)
	case TYPE_UI16:
		opUI16Add(prgrm)
	case TYPE_UI32:
		opUI32Add(prgrm)
	case TYPE_UI64:
		opUI64Add(prgrm)
	case TYPE_F32:
		opF32Add(prgrm)
	case TYPE_F64:
		opF64Add(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opSub(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Sub(prgrm)
	case TYPE_I16:
		opI16Sub(prgrm)
	case TYPE_I32:
		opI32Sub(prgrm)
	case TYPE_I64:
		opI64Sub(prgrm)
	case TYPE_UI8:
		opUI8Sub(prgrm)
	case TYPE_UI16:
		opUI16Sub(prgrm)
	case TYPE_UI32:
		opUI32Sub(prgrm)
	case TYPE_UI64:
		opUI64Sub(prgrm)
	case TYPE_F32:
		opF32Sub(prgrm)
	case TYPE_F64:
		opF64Sub(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opNeg(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Neg(prgrm)
	case TYPE_I16:
		opI16Neg(prgrm)
	case TYPE_I32:
		opI32Neg(prgrm)
	case TYPE_I64:
		opI64Neg(prgrm)
	case TYPE_F32:
		opF32Neg(prgrm)
	case TYPE_F64:
		opF64Neg(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitshl(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitshl(prgrm)
	case TYPE_I16:
		opI16Bitshl(prgrm)
	case TYPE_I32:
		opI32Bitshl(prgrm)
	case TYPE_I64:
		opI64Bitshl(prgrm)
	case TYPE_UI8:
		opUI8Bitshl(prgrm)
	case TYPE_UI16:
		opUI16Bitshl(prgrm)
	case TYPE_UI32:
		opUI32Bitshl(prgrm)
	case TYPE_UI64:
		opUI64Bitshl(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitshr(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitshr(prgrm)
	case TYPE_I16:
		opI16Bitshr(prgrm)
	case TYPE_I32:
		opI32Bitshr(prgrm)
	case TYPE_I64:
		opI64Bitshr(prgrm)
	case TYPE_UI8:
		opUI8Bitshr(prgrm)
	case TYPE_UI16:
		opUI16Bitshr(prgrm)
	case TYPE_UI32:
		opUI32Bitshr(prgrm)
	case TYPE_UI64:
		opUI64Bitshr(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opBitclear(prgrm *CXProgram) {
	switch prgrm.GetExpr().Inputs[0].Type {
	case TYPE_I8:
		opI8Bitclear(prgrm)
	case TYPE_I16:
		opI16Bitclear(prgrm)
	case TYPE_I32:
		opI32Bitclear(prgrm)
	case TYPE_I64:
		opI64Bitclear(prgrm)
	case TYPE_UI8:
		opUI8Bitclear(prgrm)
	case TYPE_UI16:
		opUI16Bitclear(prgrm)
	case TYPE_UI32:
		opUI32Bitclear(prgrm)
	case TYPE_UI64:
		opUI64Bitclear(prgrm)
	default:
		panic(CX_INTERNAL_ERROR)
	}
}

func opLen(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	elt := GetAssignmentElement(inp1)

	if elt.IsSlice || elt.Type == TYPE_AFF {
		var sliceOffset = GetSliceOffset(fp, inp1)
		if sliceOffset > 0 {
			sliceLen := GetSliceHeader(sliceOffset)[4:8]
			WriteMemory(GetFinalOffset(fp, out1), sliceLen)
		} else if sliceOffset == 0 {
			WriteI32(GetFinalOffset(fp, out1), 0)
		} else {
			panic(CX_RUNTIME_ERROR)
		}

		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == TYPE_STR && elt.Lengths == nil {
		var strOffset = GetStrOffset(fp, inp1)
		// Checking if the string lives on the heap.
		if strOffset > PROGRAM.HeapStartsAt {
			// Then it's on the heap and we need to consider
			// the object's header.
			strOffset += OBJECT_HEADER_SIZE
		}

		WriteMemory(GetFinalOffset(fp, out1), PROGRAM.Memory[strOffset:strOffset+STR_HEADER_SIZE])
	} else {
		outV0 := int32(elt.Lengths[len(elt.Indexes)])
		WriteI32(GetFinalOffset(fp, out1), outV0)
	}
}

func opAppend(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	eltInp1 := GetAssignmentElement(inp1)
	eltOut1 := GetAssignmentElement(out1)
	if inp1.Type != inp2.Type || inp1.Type != out1.Type || !eltInp1.IsSlice || !eltOut1.IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen int32
	inputSliceOffset := GetSliceOffset(fp, inp1)
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := SliceAppendResize(fp, out1, inp1, inp2.Size)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.
	outputSlicePointer := GetFinalOffset(fp, out1)

	if inp2.Type == TYPE_STR || inp2.Type == TYPE_AFF {
		var obj [4]byte
		WriteMemI32(obj[:], 0, int32(GetStrOffset(fp, inp2)))
		SliceAppendWrite(outputSliceOffset, obj[:], inputSliceLen)
	} else {
		obj := ReadMemory(GetFinalOffset(fp, inp2), inp2)
		SliceAppendWrite(outputSliceOffset, obj, inputSliceLen)
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opResize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(SliceResize(fp, out1, inp1, ReadI32(fp, inp2), GetAssignmentElement(inp1).TotalSize))
	outputSlicePointer := GetFinalOffset(fp, out1)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opInsert(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]

	if inp1.Type != inp3.Type || inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)

	if inp3.Type == TYPE_STR || inp3.Type == TYPE_AFF {
		var obj [4]byte
		WriteMemI32(obj[:], 0, int32(GetStrOffset(fp, inp3)))
		outputSliceOffset := int32(SliceInsert(fp, out1, inp1, ReadI32(fp, inp2), obj[:]))
		WriteI32(outputSlicePointer, outputSliceOffset)
	} else {
		obj := ReadMemory(GetFinalOffset(fp, inp3), inp3)
		outputSliceOffset := int32(SliceInsert(fp, out1, inp1, ReadI32(fp, inp2), obj))
		WriteI32(outputSlicePointer, outputSliceOffset)
	}
}

func opRemove(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := int32(SliceRemove(fp, out1, inp1, ReadI32(fp, inp2), int32(GetAssignmentElement(inp1).TotalSize)))
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opCopy(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(count/dstElem.TotalSize))
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

func opSprintf(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	byts := encoder.Serialize(string(buildString(expr, fp)))
	WriteObject(out1Offset, byts)
}

func opPrintf(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	fmt.Print(string(buildString(expr, fp)))
}

func opRead(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)

	var header = make([]byte, OBJECT_HEADER_SIZE)
	WriteMemI32(header, 5, int32(len(byts)))

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	WriteI32(out1Offset, int32(heapOffset+OBJECT_HEADER_SIZE))
}
