package base

import (
	"strconv"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_lt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) < ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) < ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) < ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_gt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) > ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) > ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) > ReadF64(stack, fp, inp2))
	}
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_lteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) <= ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) <= ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) <= ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) <= ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_gteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) >= ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) >= ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) >= ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) >= ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_equal(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(stack, fp, inp1) == ReadBool(stack, fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(stack, fp, inp1) == ReadStr(stack, fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) == ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) == ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) == ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) == ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_unequal(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(stack, fp, inp1) != ReadBool(stack, fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) != ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) != ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) != ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) != ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitand(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) & ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) & ReadI64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) | ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) | ReadI64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitxor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) ^ ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) ^ ReadI64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_mul(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) * ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) * ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) * ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) * ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_div(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) / ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) / ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) / ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) / ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_mod(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) % ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) % ReadI64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_add(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) + ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) + ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) + ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_sub(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) - ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) - ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) - ReadF64(stack, fp, inp2))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitshl(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(stack, fp, inp1)) << uint32(ReadI32(stack, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint64(ReadI64(stack, fp, inp1)) << uint64(ReadI64(stack, fp, inp2))))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitshr(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(stack, fp, inp1)) >> uint32(ReadI32(stack, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(stack, fp, inp1)) >> uint32(ReadI64(stack, fp, inp2))))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitclear(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(stack, fp, inp1)) &^ uint32(ReadI32(stack, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(stack, fp, inp1)) &^ uint32(ReadI64(stack, fp, inp2))))
	}

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_len(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(inp1.Lengths[len(inp1.Lengths)-1]))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func buildString(expr *CXExpression, stack *CXStack, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ReadStr(stack, fp, inp1)

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
			inp := expr.Inputs[specifiersCounter+1]
			switch nextCh {
			case 's':
				res = append(res, []byte(checkForEscapedChars(ReadStr(stack, fp, inp)))...)
			case 'd':
				switch inp.Type {
				case TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI32(stack, fp, inp)), 10))...)
				case TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ReadI64(stack, fp, inp), 10))...)
				}
			case 'f':
				switch inp.Type {
				case TYPE_F32:
					res = append(res, []byte(strconv.FormatFloat(float64(ReadF32(stack, fp, inp)), 'f', 7, 32))...)
				case TYPE_F64:
					res = append(res, []byte(strconv.FormatFloat(ReadF64(stack, fp, inp), 'f', 16, 64))...)
				}
			}
			c++
			specifiersCounter++
		} else {
			res = append(res, ch)
		}
	}

	return res
}

func op_sprintf(expr *CXExpression, stack *CXStack, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)
	
	// out1 := expr.Outputs[0]
	byts := encoder.Serialize(string(buildString(expr, stack, fp)))
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(stack.Program, len(byts)+OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteToHeap(&stack.Program.Heap, heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteToStack(stack, out1Offset, off)
}

func op_printf(expr *CXExpression, stack *CXStack, fp int) {
	fmt.Print(string(buildString(expr, stack, fp)))
}
