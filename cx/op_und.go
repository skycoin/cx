package cxcore

import (
	"bufio"
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/mem"
	"github.com/skycoin/cx/cx/tostring"
	"github.com/skycoin/cx/cx/util2"
	"os"
	"strconv"
	"strings"
)

func opLen(expr *ast.CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	elt := ast.GetAssignmentElement(inp1)

	if elt.IsSlice || elt.Type == constants.TYPE_AFF {
		var sliceOffset = ast.GetSliceOffset(fp, inp1)
		if sliceOffset > 0 {
			sliceLen := ast.GetSliceHeader(sliceOffset)[4:8]
			mem.WriteMemory(ast.GetFinalOffset(fp, out1), sliceLen)
		} else if sliceOffset == 0 {
			mem.WriteI32(ast.GetFinalOffset(fp, out1), 0)
		} else {
			panic(constants.CX_RUNTIME_ERROR)
		}

		// TODO: Had to add elt.Lengths to avoid doing this for arrays, but not entirely sure why
	} else if elt.Type == constants.TYPE_STR && elt.Lengths == nil {
		var strOffset = ast.GetStrOffset(fp, inp1)
		// Checking if the string lives on the heap.
		if strOffset > ast.PROGRAM.HeapStartsAt {
			// Then it's on the heap and we need to consider
			// the object's header.
			strOffset += constants.OBJECT_HEADER_SIZE
		}

		mem.WriteMemory(ast.GetFinalOffset(fp, out1), ast.PROGRAM.Memory[strOffset:strOffset+constants.STR_HEADER_SIZE])
	} else {
		outV0 := int32(elt.Lengths[len(elt.Indexes)])
		mem.WriteI32(ast.GetFinalOffset(fp, out1), outV0)
	}
}

func opAppend(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	eltInp1 := ast.GetAssignmentElement(inp1)
	eltOut1 := ast.GetAssignmentElement(out1)
	if inp1.Type != inp2.Type || inp1.Type != out1.Type || !eltInp1.IsSlice || !eltOut1.IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var inputSliceLen int32
	inputSliceOffset := ast.GetSliceOffset(fp, inp1)
	if inputSliceOffset != 0 {
		inputSliceLen = ast.GetSliceLen(inputSliceOffset)
	}

	// Preparing slice in case more memory is needed for the new element.
	outputSliceOffset := ast.SliceAppendResize(fp, out1, inp1, inp2.Size)

	// We need to update the address of the output and input, as the final offsets
	// could be on the heap and they could have been moved by the GC.
	outputSlicePointer := ast.GetFinalOffset(fp, out1)

	if inp2.Type == constants.TYPE_STR || inp2.Type == constants.TYPE_AFF {
		var obj [4]byte
		mem.WriteMemI32(obj[:], 0, int32(ast.GetStrOffset(fp, inp2)))
		ast.SliceAppendWrite(outputSliceOffset, obj[:], inputSliceLen)
	} else {
		obj := ast.ReadMemory(ast.GetFinalOffset(fp, inp2), inp2)
		ast.SliceAppendWrite(outputSliceOffset, obj, inputSliceLen)
	}

	mem.WriteI32(outputSlicePointer, outputSliceOffset)
}

func opResize(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSliceOffset := int32(ast.SliceResize(fp, out1, inp1, ReadI32(fp, inp2), ast.GetAssignmentElement(inp1).TotalSize))
	outputSlicePointer := ast.GetFinalOffset(fp, out1)
	mem.WriteI32(outputSlicePointer, outputSliceOffset)
}

func opInsert(expr *ast.CXExpression, fp int) {
	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]

	if inp1.Type != inp3.Type || inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := ast.GetFinalOffset(fp, out1)

	if inp3.Type == constants.TYPE_STR || inp3.Type == constants.TYPE_AFF {
		var obj [4]byte
		mem.WriteMemI32(obj[:], 0, int32(ast.GetStrOffset(fp, inp3)))
		outputSliceOffset := int32(ast.SliceInsert(fp, out1, inp1, ReadI32(fp, inp2), obj[:]))
		mem.WriteI32(outputSlicePointer, outputSliceOffset)
	} else {
		obj := ast.ReadMemory(ast.GetFinalOffset(fp, inp3), inp3)
		outputSliceOffset := int32(ast.SliceInsert(fp, out1, inp1, ReadI32(fp, inp2), obj))
		mem.WriteI32(outputSlicePointer, outputSliceOffset)
	}
}

func opRemove(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	if inp1.Type != out1.Type || !ast.GetAssignmentElement(inp1).IsSlice || !ast.GetAssignmentElement(out1).IsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	outputSlicePointer := ast.GetFinalOffset(fp, out1)
	outputSliceOffset := int32(ast.SliceRemove(fp, out1, inp1, ReadI32(fp, inp2), int32(ast.GetAssignmentElement(inp1).TotalSize)))
	mem.WriteI32(outputSlicePointer, outputSliceOffset)
}

func opCopy(expr *ast.CXExpression, fp int) {
	dstInput := expr.Inputs[0]
	srcInput := expr.Inputs[1]
	dstOffset := ast.GetSliceOffset(fp, dstInput)
	srcOffset := ast.GetSliceOffset(fp, srcInput)

	dstElem := ast.GetAssignmentElement(dstInput)
	srcElem := ast.GetAssignmentElement(srcInput)

	if dstInput.Type != srcInput.Type || !dstElem.IsSlice || !srcElem.IsSlice || dstElem.TotalSize != srcElem.TotalSize {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	var count int
	if dstInput.Type == srcInput.Type && dstOffset >= 0 && srcOffset >= 0 {
		count = copy(ast.GetSliceData(dstOffset, dstElem.TotalSize), ast.GetSliceData(srcOffset, srcElem.TotalSize))
		if count%dstElem.TotalSize != 0 {
			panic(constants.CX_RUNTIME_ERROR)
		}
	} else {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}
	mem.WriteI32(ast.GetFinalOffset(fp, expr.Outputs[0]), int32(count/dstElem.TotalSize))
}

func buildString(expr *ast.CXExpression, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ast.ReadStr(fp, inp1)

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
				res = append(res, []byte(util2.CheckForEscapedChars(ast.ReadStr(fp, inp)))...)
			case 'd':
				switch inp.Type {
				case constants.TYPE_I8:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI8(fp, inp)), 10))...)
				case constants.TYPE_I16:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI16(fp, inp)), 10))...)
				case constants.TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI32(fp, inp)), 10))...)
				case constants.TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ReadI64(fp, inp), 10))...)
				case constants.TYPE_UI8:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI8(fp, inp)), 10))...)
				case constants.TYPE_UI16:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI16(fp, inp)), 10))...)
				case constants.TYPE_UI32:
					res = append(res, []byte(strconv.FormatUint(uint64(ReadUI32(fp, inp)), 10))...)
				case constants.TYPE_UI64:
					res = append(res, []byte(strconv.FormatUint(ReadUI64(fp, inp), 10))...)
				}
			case 'f':
				switch inp.Type {
				case constants.TYPE_F32:
					res = append(res, []byte(strconv.FormatFloat(float64(ReadF32(fp, inp)), 'f', 7, 32))...)
				case constants.TYPE_F64:
					res = append(res, []byte(strconv.FormatFloat(ReadF64(fp, inp), 'f', 16, 64))...)
				}
			case 'v':
				res = append(res, []byte(tostring.GetPrintableValue(fp, inp))...)
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
			elt := ast.GetAssignmentElement(inp)
			typ := ""
			_ = typ
			if elt.CustomType != nil {
				// then it's custom type
				typ = elt.CustomType.Name
			} else {
				// then it's native type
				typ = constants.TypeNames[elt.Type]
			}

			if c == lInps-1 {
				extra += fmt.Sprintf("%s=%s", typ, tostring.GetPrintableValue(fp, elt))
			} else {
				extra += fmt.Sprintf("%s=%s, ", typ, tostring.GetPrintableValue(fp, elt))
			}

		}

		extra += ")"

		res = append(res, []byte(extra)...)
	}

	return res
}

func opSprintf(expr *ast.CXExpression, fp int) {
	mem.WriteString(fp, string(buildString(expr, fp)), expr.Outputs[0])
}

func opPrintf(expr *ast.CXExpression, fp int) {
	fmt.Print(string(buildString(expr, fp)))
}

func opRead(inputs []ast.CXValue, outputs []ast.CXValue) {

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
    if err != nil {
		panic(constants.CX_INTERNAL_ERROR)
	}

	// text = strings.Trim(text, " \n")
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
    outputs[0].Set_str(text)
}
