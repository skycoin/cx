package base

import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_lt (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) < ReadByte(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) < ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) < ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) < ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) < ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) < ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) > ReadByte(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) > ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) > ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) > ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) > ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) > ReadF64(fp, inp2))
	}
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) <= ReadByte(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) <= ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) <= ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) <= ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) <= ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) <= ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) >= ReadByte(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) >= ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) >= ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) >= ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) >= ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) >= ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_equal(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) == ReadByte(fp, inp2))
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(fp, inp1) == ReadBool(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) == ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) == ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) == ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) == ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) == ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_unequal(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromBool(ReadByte(fp, inp1) != ReadByte(fp, inp2))
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(fp, inp1) != ReadBool(fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(fp, inp1) != ReadStr(fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(fp, inp1) != ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(fp, inp1) != ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(fp, inp1) != ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(fp, inp1) != ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) & ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) & ReadI64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) | ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) | ReadI64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) ^ ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) ^ ReadI64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromByte(ReadByte(fp, inp1) * ReadByte(fp, inp2))
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) * ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) * ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(fp, inp1) * ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(fp, inp1) * ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromByte(ReadByte(fp, inp1) / ReadByte(fp, inp2))
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) / ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) / ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(fp, inp1) / ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(fp, inp1) / ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromByte(ReadByte(fp, inp1) % ReadByte(fp, inp2))
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) % ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) % ReadI64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromByte(ReadByte(fp, inp1) + ReadByte(fp, inp2))
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) + ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) + ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(fp, inp1) + ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(fp, inp1) + ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_sub(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BYTE:
		outB1 = FromByte(ReadByte(fp, inp1) - ReadByte(fp, inp2))
	case TYPE_I32:
		outB1 = FromI32(ReadI32(fp, inp1) - ReadI32(fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(fp, inp1) - ReadI64(fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(fp, inp1) - ReadF32(fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(fp, inp1) - ReadF64(fp, inp2))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(fp, inp1)) << uint32(ReadI32(fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint64(ReadI64(fp, inp1)) << uint64(ReadI64(fp, inp2))))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(fp, inp1)) >> uint32(ReadI32(fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(fp, inp1)) >> uint32(ReadI64(fp, inp2))))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(fp, inp1)) &^ uint32(ReadI32(fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(fp, inp1)) &^ uint32(ReadI64(fp, inp2))))
	}

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_len(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	elt := GetAssignmentElement(inp1)
	
	if elt.IsSlice || elt.Type == TYPE_AFF {
		var sliceOffset int32 = GetSliceOffset(fp, inp1)
		var sliceLen  []byte
		if sliceOffset > 0  {
			sliceLen = GetSliceHeader(sliceOffset)[4:8]
		} else if sliceOffset == 0 {
			sliceLen = FromI32(0)
		} else {
			panic(CX_RUNTIME_ERROR)
		}

		WriteMemory(GetFinalOffset(fp, out1), sliceLen)
	} else if elt.Type == TYPE_STR {
		var strOffset int = GetStrOffset(fp, inp1)
		WriteMemory(GetFinalOffset(fp, out1), PROGRAM.Memory[strOffset:strOffset+STR_HEADER_SIZE])
	} else {
		outB1 := FromI32(int32(elt.Lengths[0]))
		WriteMemory(GetFinalOffset(fp, out1), outB1)
	}
}

func op_resize(expr *CXExpression, fp int) {
	sliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	var newLen int32 = ReadI32(fp, expr.Inputs[1])
	if sliceOffset >= 0 && newLen >= 0 {
		sliceHeader := GetSliceHeader(sliceOffset)
		var oldCap int32
		encoder.DeserializeAtomic(sliceHeader[0:4], &oldCap)
		if newLen <= oldCap {
			copy(GetSliceHeader(sliceOffset)[4:8], encoder.SerializeAtomic(newLen))
			return
		}
	}

	panic(CX_RUNTIME_INVALID_ARGUMENT)
}

func op_copy(expr *CXExpression, fp int) {
	dstInput := expr.Inputs[0]
	srcInput := expr.Inputs[1]
	dstOffset := GetSliceOffset(fp, dstInput)
	srcOffset := GetSliceOffset(fp, srcInput)

	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		copy(GetSlice(dstOffset, GetAssignmentElement(dstInput).TotalSize), GetSlice(srcOffset, GetAssignmentElement(srcInput).TotalSize))
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
}

var FUCK bool = false

func op_append (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != inp2.Type || inp1.Type != out1.Type {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	sizeofElement := inp2.TotalSize

	inputSliceOffset := GetSliceOffset(fp, inp1)

	outputSliceElement := GetAssignmentElement(out1)
	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := GetPointerOffset(outputSliceElement, int32(outputSlicePointer))

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceHeader := GetSliceHeader(inputSliceOffset)
		encoder.DeserializeAtomic(inputSliceHeader[4:8], &inputSliceLen)

		var inputSliceCap int32
		encoder.DeserializeAtomic(inputSliceHeader[0:4], &inputSliceCap)
	}

	var outputSliceHeader []byte
	var outputSliceLen int32
	var outputSliceCap int32

	if outputSliceOffset > 0 {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		encoder.DeserializeAtomic(outputSliceHeader[0:4], &outputSliceCap)
		encoder.DeserializeAtomic(outputSliceHeader[4:8], &outputSliceLen)
	}

	var newLen int32 = inputSliceLen + 1
	var newCap int32 = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize int32 = OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + newCap * int32(sizeofElement)
		outputSliceOffset = int32(AllocateSeq(int(outputObjectSize)))
		copy(PROGRAM.Memory[outputSlicePointer:], encoder.SerializeAtomic(outputSliceOffset))
		copy(GetObjectHeader(outputSliceOffset)[5:9], encoder.SerializeAtomic(outputObjectSize))

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		copy(outputSliceHeader[0:4], encoder.SerializeAtomic(newCap))
	}

	copy(outputSliceHeader[4:8], encoder.SerializeAtomic(newLen))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)

	if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
		copy(outputSliceData, GetSliceData(inputSliceOffset, sizeofElement))
	}

	var obj2[]byte
	if inp2.Type == TYPE_STR || inp2.Type == TYPE_AFF {
		obj2 = encoder.SerializeAtomic(int32(GetStrOffset(fp, inp2)))
	} else {
		obj2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
	}

	copy(outputSliceData[int(inputSliceLen) * sizeofElement:], obj2)
}

func buildString(expr *CXExpression, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ReadStr(fp, inp1)

	var res []byte
	var specifiersCounter int
	var lenStr int = len(fmtStr)
	
	for c := 0; c < len(fmtStr); c++ {
		var nextCh byte
		ch := fmtStr[c]
		if c < lenStr - 1{
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
			if specifiersCounter + 1 == len(expr.Inputs) {
				res = append(res, []byte(fmt.Sprintf("%%!%c(MISSING)", nextCh))...)
				c++
				continue
			}
			
			inp := expr.Inputs[specifiersCounter + 1]
			switch nextCh {
			case 's':
				res = append(res, []byte(checkForEscapedChars(ReadStr(fp, inp)))...)
			case 'd':
				switch inp.Type {
				case TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI32(fp, inp)), 10))...)
				case TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ReadI64(fp, inp), 10))...)
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

	if specifiersCounter != len(expr.Inputs) - 1 {
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

			if c == lInps - 1 {
				extra += fmt.Sprintf("%s=%s", typ, GetPrintableValue(fp, elt))
			} else {
				extra += fmt.Sprintf("%s=%s, ", typ, GetPrintableValue(fp, elt))
			}
			
		}

		extra += ")"
		
		res = append(res, []byte(extra)...)
	}
	
	// if specifiersCounter != len(expr.Inputs) - 1 {
	// 	panic("meow")
	// }

	return res
}

func op_sprintf(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	byts := encoder.Serialize(string(buildString(expr, fp)))
	WriteObject(out1Offset, byts)
}

func op_printf(expr *CXExpression, fp int) {
	fmt.Print(string(buildString(expr, fp)))
}

func op_read (expr *CXExpression, fp int) {
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
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset + OBJECT_HEADER_SIZE))

	WriteMemory(out1Offset, off)
}
