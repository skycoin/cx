package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
)

// Creates a copy of its argument in the stack. opCopy should be used to copy an argument in a heap.
func identity (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	WriteToStack(stack, fp, out1.Offset, stack.Stack[fp + inp1.Offset : fp + inp1.Offset + inp1.Size])
}

func i32_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	// fmt.Println("offset", inp2.MemoryType)
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	// fmt.Println("before", stack)
	WriteToStack(stack, fp, out1.Offset, outB1)
	// fmt.Println("after", stack)
	// fmt.Println("fp", fp)
}

func i32_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	fmt.Println("data", inp.MemoryType, inp.Offset, inp.Size)
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp.Offset : fp + inp.Offset + inp.Size]
		fmt.Println("here", fp, inp.Offset, inp.Size, byts, inp.Name)
		encoder.DeserializeAtomic(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[inp.Offset : inp.Offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func WriteToStack (stack *CXStack, fp int, offset int, out []byte) {
	// fmt.Println("before", stack)
	// fmt.Println("fp", fp, offset)
	for c := 0; c < len(out); c++ {
		(*stack).Stack[fp + offset + c] = out[c]
	}
	// fmt.Println("after", stack)
}
