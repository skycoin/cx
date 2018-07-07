package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// creates a copy of its argument in the stack
func op_identity(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(stack, fp, inp1, MEM_READ)
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)

	// fmt.Println("huh", out1.Name, out1.Fields, out1.)
	// fmt.Println("isStruct", inp1.IsStruct)

	// fmt.Println("stats", inp1.MemoryRead, inp1.MemoryWrite, out1.MemoryRead, out1.MemoryWrite)

	if out1.IsPointer && out1.DereferenceLevels != out1.IndirectionLevels && !inp1.IsPointer || (inp1.MemoryRead == MEM_HEAP && inp1.MemoryWrite == MEM_HEAP) {
		switch inp1.MemoryWrite {
		case MEM_STACK:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToStack(stack, out1Offset, byts)
		case MEM_HEAP:
			// if heapoffset > 0 look in here and don't allocate

			if inp1.MemoryRead == MEM_HEAP {
				WriteToStack(stack, out1Offset, encoder.SerializeAtomic(int32(inp1.Offset)))
			} else if inp1.MemoryRead == MEM_DATA {
				WriteToData(&stack.Program.Data, out1Offset, encoder.SerializeAtomic(int32(inp1.Offset)))
			} else {
				
				var heapOffset int
				if inp1.HeapOffset > 0 {
					// then it's a reference to the symbol
					var off int32
					encoder.DeserializeAtomic(stack.Stack[fp+inp1.HeapOffset:fp+inp1.HeapOffset+TYPE_POINTER_SIZE], &off)

					if off > 0 {
						// non-nil, i.e. object is already allocated
						heapOffset = int(off)
					} else {
						// nil, needs to be allocated
						heapOffset = AllocateSeq(stack.Program, inp1.TotalSize+OBJECT_HEADER_SIZE)
						WriteToStack(stack, fp+inp1.HeapOffset, encoder.SerializeAtomic(int32(heapOffset)))
					}
				}

				byts := ReadMemory(stack, inp1Offset, inp1)
				// creating a header for this object
				size := encoder.Serialize(int32(len(byts)))

				var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
				for c := 5; c < OBJECT_HEADER_SIZE; c++ {
					header[c] = size[c-5]
				}

				obj := append(header, byts...)

				// WriteToHeap(&stack.Program.Heap, heapOffset, byts)
				WriteToHeap(&stack.Program.Heap, heapOffset, obj)

				offset := encoder.SerializeAtomic(int32(heapOffset))

				// WriteToStack(stack, fp + out1Offset, offset)
				WriteToStack(stack, out1Offset, offset)
			}
		case MEM_DATA:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToData(&stack.Program.Data, out1Offset, byts)
		default:
			panic("implement the other mem types")
		}
	} else if (inp1.IsReference && out1.IsPointer) || (inp1.IsPointer && out1.IsPointer) {
		// if inp1.Type == TYPE_STR {
		// 	fmt.Println("count", inp1.Name, ReadMemory(stack, inp1Offset, inp1), inp1.Type, inp1.TotalSize)
		// 	WriteToStack(stack, out1Offset, ReadMemory(stack, inp1Offset, inp1))
		// } else {
		// 	WriteMemory(stack, out1Offset, out1, ReadMemory(stack, inp1Offset, inp1))
		// }
		WriteMemory(stack, out1Offset, out1, ReadMemory(stack, inp1Offset, inp1))
	} else {
		WriteMemory(stack, out1Offset, out1, ReadMemory(stack, inp1Offset, inp1))
	}
}

func op_goTo(expr *CXExpression, call *CXCall) {
	// inp1 := expr.Inputs[0]
	// call.Line = ReadI32(inp1)
}

func op_jmp(expr *CXExpression, stack *CXStack, fp int, call *CXCall) {
	// inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	inp1 := expr.Inputs[0]
	var predicate bool
	// var thenLines int32
	// var elseLines int32

	if expr.Label != "" {
		// then it's a goto
		call.Line = call.Line + expr.ThenLines
	} else {
		inp1Offset := GetFinalOffset(stack, fp, inp1, MEM_READ)

		switch inp1.MemoryRead {
		case MEM_STACK:
			predicateB := stack.Stack[inp1Offset : inp1Offset+inp1.Size]
			encoder.DeserializeAtomic(predicateB, &predicate)
		case MEM_DATA:
			predicateB := inp1.Program.Data[inp1Offset : inp1Offset+inp1.Size]
			encoder.DeserializeAtomic(predicateB, &predicate)
		default:
			panic("implement the other mem types in readI32")
		}

		if predicate {
			call.Line = call.Line + expr.ThenLines
		} else {
			call.Line = call.Line + expr.ElseLines
		}
	}

}
