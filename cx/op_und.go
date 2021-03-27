package cxcore

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func opLen(expr *CXExpression, fp int) {
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

func opAppend(expr *CXExpression, fp int) {
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

func opResize(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(SliceResize(fp, out1, inp1, ReadI32(fp, inp2), GetAssignmentElement(inp1).TotalSize))
	outputSlicePointer := GetFinalOffset(fp, out1)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opInsert(expr *CXExpression, fp int) {
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

func opRemove(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out1).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := GetFinalOffset(fp, out1)
	outputSliceOffset := int32(SliceRemove(fp, out1, inp1, ReadI32(fp, inp2), int32(GetAssignmentElement(inp1).TotalSize)))
	WriteI32(outputSlicePointer, outputSliceOffset)
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
		// for _, inp := range expr.ProgramInput[:specifiersCounter] {
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
	WriteString(fp, string(buildString(expr, fp)), expr.Outputs[0])
}

func opPrintf(expr *CXExpression, fp int) {
	fmt.Print(string(buildString(expr, fp)))
}

func opRead(inputs []CXValue, outputs []CXValue) {

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
    if err != nil {
		panic(CX_INTERNAL_ERROR)
	}

	// text = strings.Trim(text, " \n")
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
    outputs[0].Set_str(text)
}
