package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// creates a copy of its argument in the stack
func identity (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(stack, fp, inp1)
	out1Offset := GetFinalOffset(stack, fp, out1)

	if out1.IsPointer && out1.DereferenceLevels != out1.IndirectionLevels && !inp1.IsPointer {
		switch inp1.MemoryType {
		case MEM_STACK:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToStack(stack, out1Offset, byts)
		case MEM_HEAP:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToHeap(&out1.Program.Heap, out1Offset, byts)
		case MEM_DATA:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToData(&inp1.Program.Data, out1Offset, byts)
		default:
			panic("implement the other mem types in readI32")
		}
	} else if inp1.IsReference {
		// WriteToStack(stack, out1Offset, FromI32(int32(inp1Offset)))
		WriteMemory(stack, out1Offset, out1, FromI32(int32(inp1Offset)))
	} else {
		// fmt.Println("hi", inp1.Offset, inp1Offset)
		byts := ReadMemory(stack, inp1Offset, inp1)
		WriteMemory(stack, out1Offset, out1, byts)
	}
}

func goTo (expr *CXExpression, call *CXCall) {
	// inp1 := expr.Inputs[0]
	// call.Line = ReadI32(inp1)
}

func jmp (expr *CXExpression, stack *CXStack, fp int, call *CXCall) {
	// inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	inp1 := expr.Inputs[0]
	var predicate bool
	// var thenLines int32
	// var elseLines int32

	inp1Offset := GetFinalOffset(stack, fp, inp1)

	switch inp1.MemoryType {
	case MEM_STACK:
		predicateB := stack.Stack[inp1Offset: inp1Offset + inp1.Size]
		encoder.DeserializeAtomic(predicateB, &predicate)
	case MEM_DATA:
		predicateB := inp1.Program.Data[inp1Offset : inp1Offset + inp1.Size]
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

func read_array (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	indexes := make([]int32, len(expr.Inputs[1:]))
	for i, idx := range expr.Inputs[1:] {
		indexes[i] = ReadI32(stack, fp, idx)
	}
	offset, size := ReadArray(stack, fp, inp1, indexes)
	var outB1 []byte

	switch inp1.MemoryType {
	case MEM_STACK:
		outB1 = stack.Stack[fp + offset : fp + offset + size]
	case MEM_DATA:
		outB1 = inp1.Program.Data[offset : offset + size]
	default:
		panic("implement the other mem types in readI32")
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}
