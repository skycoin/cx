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

func op_append(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	
	inp1Offset := GetFinalOffset(stack, fp, inp1, MEM_READ)
	inp2Offset := GetFinalOffset(stack, fp, inp2, MEM_READ)
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_READ)

	var off int32
	// var size int32
	var byts []byte
	
	if out1.MemoryRead == MEM_STACK {
		byts = stack.Stack[out1Offset : out1Offset+TYPE_POINTER_SIZE]
		encoder.DeserializeAtomic(byts, &off)
	} else {
		byts = stack.Program.Data[out1Offset : out1Offset+TYPE_POINTER_SIZE]
		encoder.DeserializeAtomic(byts, &off)
	}

	var heapOffset int

	if off == 0 {
		// then out1 == nil
		var sliceHeader []byte
		var len1 int32

		if inp1Offset != 0 {
			// then we need to reserve for obj1 too
			sliceHeader = stack.Program.Heap.Heap[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
			encoder.DeserializeAtomic(sliceHeader[:4], &len1)
			heapOffset = AllocateSeq(stack.Program, (int(len1)*inp2.TotalSize) + inp2.TotalSize + OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
		} else {
			// then obj1 is nil and zero-sized
			heapOffset = AllocateSeq(stack.Program, inp2.TotalSize+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
		}
		
		// writing address to stack
		WriteMemory(stack, out1Offset, out1, encoder.SerializeAtomic(int32(heapOffset)))

		var obj1 []byte
		var obj2 []byte
		
		obj1 = stack.Program.Heap.Heap[inp1Offset : int32(inp1Offset) + len1*int32(inp2.TotalSize)]

		// obj2 = ReadMemory(stack, inp2Offset, inp2)
		
		if inp2.Type == TYPE_STR {
			// obj2 = []byte(ReadStr(stack, inp2Offset, inp2))
			obj2 = encoder.SerializeAtomic(int32(inp2Offset))
		} else {
			obj2 = ReadMemory(stack, inp2Offset, inp2)
		}
		

		// fmt.Println("obj2", ReadStr(stack, inp2Offset, inp2))

		var size []byte
		if inp1Offset != 0 {
			size = encoder.SerializeAtomic(int32(len(obj1)) + int32(len(obj2) + SLICE_HEADER_SIZE))
		} else {
			size = encoder.SerializeAtomic(int32(len(obj2) + SLICE_HEADER_SIZE))
		}

		var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
		for c := 5; c < OBJECT_HEADER_SIZE; c++ {
			header[c] = size[c-5]
		}

		// length
		lenTotal := encoder.SerializeAtomic(len1 + 1)
		capTotal := lenTotal

		var finalObj []byte

		finalObj = append(header, lenTotal...)
		finalObj = append(finalObj, capTotal...)
		if inp1Offset != 0 {
			// then obj1 is not nil, and we need to append
			finalObj = append(finalObj, obj1...)
		}
		finalObj = append(finalObj, obj2...)

		WriteToHeap(&stack.Program.Heap, heapOffset, finalObj)
	} else {
		// then we have access to a size and capacity
		// sliceHeader := stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
		sliceHeader := stack.Program.Heap.Heap[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
		
		var l int32
		var c int32

		encoder.DeserializeAtomic(sliceHeader[:4], &l)
		encoder.DeserializeAtomic(sliceHeader[4:], &c)

		if l >= c {
			// then we need to increase cap and relocate slice
			// prevObj := stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+l*int32(inp2.TotalSize)]

			var obj1 []byte
			var obj2 []byte
			
			// obj1 = stack.Program.Heap.Heap[inp1Offset+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : int32(inp1Offset)+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE + l*int32(inp2.TotalSize)]

			fmt.Println("huehue", inp1Offset, l, inp1.Name, inp1.Size, inp1.TotalSize)
			
			obj1 = stack.Program.Heap.Heap[inp1Offset : int32(inp1Offset) + l*int32(inp2.TotalSize)]
			
			obj2 = ReadMemory(stack, inp2Offset, inp2)

			l++
			c = c * 2

			heapOffset = AllocateSeq(stack.Program, int(c)*inp2.TotalSize + OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
			
			WriteMemory(stack, out1Offset, out1, encoder.SerializeAtomic(int32(heapOffset)))

			size := encoder.SerializeAtomic(int32(int(c)*inp2.TotalSize + SLICE_HEADER_SIZE))

			var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
			for c := 5; c < OBJECT_HEADER_SIZE; c++ {
				header[c] = size[c-5]
			}

			lB := encoder.SerializeAtomic(l)
			cB := encoder.SerializeAtomic(c)

			var finalObj []byte
			
			finalObj = append(header, lB...)
			finalObj = append(finalObj, cB...)
			finalObj = append(finalObj, obj1...)
			finalObj = append(finalObj, obj2...)

			WriteToHeap(&stack.Program.Heap, heapOffset, finalObj)
		} else {
			// then we can simply write the element

			// updating the length
			newL := encoder.SerializeAtomic(l+int32(1))
			for i, byt := range newL {
				stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE+int32(i)] = byt
			}

			// write the obj
			obj := ReadMemory(stack, inp2Offset, inp2)
			for i, byt := range obj {
				stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+int32(int(l)*inp2.TotalSize+i)] = byt
			}
			

			// WriteToHeap(&stack.Program.Heap, heapOffset, finalObj)
		}
	}
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
