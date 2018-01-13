package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Creates a copy of its argument in the stack. opCopy should be used to copy an argument in a heap.
func opIdentity (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	WriteToStack(stack, fp, out1.Offset, stack.Stack[fp + inp1.Offset : fp + inp1.Offset + inp1.Size])
}

func opAdd (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}



func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	byts := stack.Stack[fp + inp.Offset : fp + inp.Offset + inp.Size]
	encoder.DeserializeAtomic(byts, &out)
	return
}

func WriteToStack (stack *CXStack, fp int, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		stack.Stack[fp + offset + c] = out[c]
	}
}
