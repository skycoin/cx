package base

import (
	"strconv"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_lt(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) < ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) < ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) < ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) < ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_gt(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) > ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) > ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) > ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) > ReadF64(mem, fp, inp2))
	}
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_lteq(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) <= ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) <= ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) <= ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) <= ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_gteq(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) >= ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) >= ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) >= ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) >= ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_equal(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(mem, fp, inp1) == ReadBool(mem, fp, inp2))
	case TYPE_STR:
		outB1 = FromBool(ReadStr(mem, fp, inp1) == ReadStr(mem, fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) == ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) == ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) == ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) == ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_unequal(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_BOOL:
		outB1 = FromBool(ReadBool(mem, fp, inp1) != ReadBool(mem, fp, inp2))
	case TYPE_I32:
		outB1 = FromBool(ReadI32(mem, fp, inp1) != ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(mem, fp, inp1) != ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(mem, fp, inp1) != ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(mem, fp, inp1) != ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitand(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) & ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) & ReadI64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitor(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) | ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) | ReadI64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitxor(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) ^ ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) ^ ReadI64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_mul(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) * ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) * ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(mem, fp, inp1) * ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(mem, fp, inp1) * ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_div(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) / ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) / ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(mem, fp, inp1) / ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(mem, fp, inp1) / ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_mod(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) % ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) % ReadI64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_add(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) + ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) + ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(mem, fp, inp1) + ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(mem, fp, inp1) + ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_sub(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(mem, fp, inp1) - ReadI32(mem, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(mem, fp, inp1) - ReadI64(mem, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(mem, fp, inp1) - ReadF32(mem, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(mem, fp, inp1) - ReadF64(mem, fp, inp2))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitshl(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(mem, fp, inp1)) << uint32(ReadI32(mem, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint64(ReadI64(mem, fp, inp1)) << uint64(ReadI64(mem, fp, inp2))))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitshr(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(mem, fp, inp1)) >> uint32(ReadI32(mem, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(mem, fp, inp1)) >> uint32(ReadI64(mem, fp, inp2))))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bitclear(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(mem, fp, inp1)) &^ uint32(ReadI32(mem, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(mem, fp, inp1)) &^ uint32(ReadI64(mem, fp, inp2))))
	}

	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_len(expr *CXExpression, mem []byte, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(inp1.Lengths[len(inp1.Lengths)-1]))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_append(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	
	inp1Offset := GetFinalOffset(mem, fp, inp1, MEM_READ)
	inp2Offset := GetFinalOffset(mem, fp, inp2, MEM_READ)
	out1Offset := GetFinalOffset(mem, fp, out1, MEM_READ)

	var off int32
	// var size int32
	var byts []byte

	byts = mem[out1Offset : out1Offset+TYPE_POINTER_SIZE]
	encoder.DeserializeAtomic(byts, &off)

	var heapOffset int

	if off == 0 {
		// then out1 == nil
		var sliceHeader []byte
		var len1 int32

		if inp1Offset != 0 {
			// then we need to reserve for obj1 too
			sliceHeader = mem[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]

			encoder.DeserializeAtomic(sliceHeader[:4], &len1)
			heapOffset = AllocateSeq(expr.Program, (int(len1)*inp2.TotalSize) + inp2.TotalSize + OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
		} else {
			// then obj1 is nil and zero-sized
			heapOffset = AllocateSeq(expr.Program, inp2.TotalSize+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
		}
		
		// writing address to mem
		// WriteMemory(mem, out1Offset, out1, encoder.SerializeAtomic(int32(heapOffset)))
		// WriteMemory(mem, out1Offset, out1, encoder.SerializeAtomic(int32(heapOffset)))
		WriteMemory(mem, out1Offset, encoder.SerializeAtomic(int32(heapOffset)))

		var obj1 []byte
		var obj2 []byte
		
		obj1 = mem[inp1Offset : int32(inp1Offset) + len1*int32(inp2.TotalSize)]

		// obj2 = ReadMemory(mem, inp2Offset, inp2)
		
		if inp2.Type == TYPE_STR {
			// obj2 = []byte(ReadStr(mem, inp2Offset, inp2))
			obj2 = encoder.SerializeAtomic(int32(inp2Offset))
		} else {
			obj2 = ReadMemory(mem, inp2Offset, inp2)
		}
		

		// fmt.Println("obj2", ReadStr(mem, inp2Offset, inp2))

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

		WriteMemory(mem, heapOffset, finalObj)
	} else {
		// then we have access to a size and capacity
		// sliceHeader := mem.Program.Heap.Heap[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
		var sliceHeader []byte
		if inp1.Type == TYPE_STR {
			sliceHeader = mem[off + OBJECT_HEADER_SIZE : off + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE]
		} else {
			sliceHeader = mem[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
		}
		// sliceHeader = mem.Program.Heap.Heap[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
		// fmt.Println("houhou", off, inp1Offset-SLICE_HEADER_SIZE, inp1Offset)
		
		// sliceHeader = mem.Program.Heap.Heap[off-SLICE_HEADER_SIZE : off]
		
		var l int32
		var c int32

		encoder.DeserializeAtomic(sliceHeader[:4], &l)
		encoder.DeserializeAtomic(sliceHeader[4:], &c)

		if l >= c {
			// then we need to increase cap and relocate slice
			// prevObj := mem.Program.Heap.Heap[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+l*int32(inp2.TotalSize)]

			var obj1 []byte
			var obj2 []byte
			
			// obj1 = mem.Program.Heap.Heap[inp1Offset+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : int32(inp1Offset)+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE + l*int32(inp2.TotalSize)]

			if inp1.Type == TYPE_STR {
				obj1 = mem[off + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE : off + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + l*int32(inp2.TotalSize)]
			} else {
				obj1 = mem[inp1Offset : int32(inp1Offset) + l*int32(inp2.TotalSize)]
			}

			if inp2.Type == TYPE_STR {
				obj2 = encoder.SerializeAtomic(int32(inp2Offset))
			} else {
				obj2 = ReadMemory(mem, inp2Offset, inp2)
			}
			// obj2 = ReadMemory(mem, inp2Offset, inp2)

			l++
			c = c * 2

			heapOffset = AllocateSeq(expr.Program, int(c)*inp2.TotalSize + OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE)
			
			// WriteMemory(mem, out1Offset, out1, encoder.SerializeAtomic(int32(heapOffset)))
			WriteMemory(mem, out1Offset, encoder.SerializeAtomic(int32(heapOffset)))

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

			WriteMemory(mem, heapOffset, finalObj)
		} else {
			// then we can simply write the element

			// updating the length
			newL := encoder.SerializeAtomic(l+int32(1))
			for i, byt := range newL {
				mem[off+OBJECT_HEADER_SIZE+int32(i)] = byt
			}

			// write the obj
			obj := ReadMemory(mem, inp2Offset, inp2)
			for i, byt := range obj {
				mem[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+int32(int(l)*inp2.TotalSize+i)] = byt
			}

			// WriteMemory(&mem.Program.Heap, heapOffset, finalObj)
		}
	}
}

func buildString(expr *CXExpression, mem []byte, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ReadStr(mem, fp, inp1)

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
				res = append(res, []byte(checkForEscapedChars(ReadStr(mem, fp, inp)))...)
			case 'd':
				switch inp.Type {
				case TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ReadI32(mem, fp, inp)), 10))...)
				case TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ReadI64(mem, fp, inp), 10))...)
				}
			case 'f':
				switch inp.Type {
				case TYPE_F32:
					res = append(res, []byte(strconv.FormatFloat(float64(ReadF32(mem, fp, inp)), 'f', 7, 32))...)
				case TYPE_F64:
					res = append(res, []byte(strconv.FormatFloat(ReadF64(mem, fp, inp), 'f', 16, 64))...)
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

func op_sprintf(expr *CXExpression, mem []byte, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(mem, fp, out1, MEM_WRITE)

	byts := encoder.Serialize(string(buildString(expr, mem, fp)))
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(expr.Program, len(byts)+OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(mem, heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(mem, out1Offset, off)
}

func op_printf(expr *CXExpression, mem []byte, fp int) {
	fmt.Print(string(buildString(expr, mem, fp)))
}
