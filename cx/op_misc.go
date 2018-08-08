package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func EscapeAnalysis (stack *CXStack, fp int, inpOffset, outOffset int, arg *CXArgument) {
	var heapOffset int
	if arg.HeapOffset > 0 {
		// then it's a reference to the symbol
		var off int32
		encoder.DeserializeAtomic(stack.Stack[fp+arg.HeapOffset:fp+arg.HeapOffset+TYPE_POINTER_SIZE], &off)

		if off > 0 {
			// non-nil, i.e. object is already allocated
			heapOffset = int(off)
		} else {
			// nil, needs to be allocated
			heapOffset = AllocateSeq(stack.Program, arg.TotalSize+OBJECT_HEADER_SIZE)
			o := GetFinalOffset(stack, fp, arg, MEM_WRITE)
			WriteMemory(stack, o, arg, encoder.SerializeAtomic(int32(heapOffset)))
		}
	}

	byts := ReadMemory(stack, inpOffset, arg)
	// creating a header for this object
	size := encoder.SerializeAtomic(int32(len(byts)))

	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteToHeap(&stack.Program.Heap, heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteToStack(stack, outOffset, off)
	// WriteMemory(stack, outOffset, arg, off)
}

func op_identity(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(stack, fp, inp1, MEM_READ)
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)

	if out1.DoesEscape {
		EscapeAnalysis(stack, fp, inp1Offset, out1Offset, inp1)
	} else {
		switch out1.PassBy {
		case PASSBY_VALUE:
			WriteMemory(stack, out1Offset, out1, ReadMemory(stack, inp1Offset, inp1))
		case PASSBY_REFERENCE:
			WriteMemory(stack, out1Offset, out1, encoder.SerializeAtomic(int32(inp1Offset)))
		}
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
