package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
	"time"
)

// Creates a copy of its argument in the stack. opCopy should be used to copy an argument in a heap.
func identity (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	switch inp1.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp1.Offset : fp + inp1.Offset + inp1.Size]
		WriteToStack(stack, fp, out1.Offset, byts)
	case MEM_DATA:
		byts := inp1.Program.Data[inp1.Offset : inp1.Offset + inp1.Size]
		WriteToStack(stack, fp, out1.Offset, byts)
	default:
		panic("implement the other mem types in readI32")
	}
}

func i32_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i32_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i32_gt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i32_lt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i32_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

func i64_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) + ReadI64(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i64_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) - ReadI64(stack, fp, inp2))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func i64_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI64(stack, fp, inp1))
}

func bool_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(stack, fp, inp1))
}

func goTo (expr *CXExpression, call *CXCall) {
	// inp1 := expr.Inputs[0]
	// call.Line = ReadI32(inp1)
}

func jmp (expr *CXExpression, stack *CXStack, fp int, call *CXCall) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	var predicate bool
	var thenLines int32
	var elseLines int32

	switch inp1.MemoryType {
	case MEM_STACK:
		predicateB := stack.Stack[fp + inp1.Offset : fp + inp1.Offset + inp1.Size]
		encoder.DeserializeAtomic(predicateB, &predicate)
	case MEM_DATA:
		predicateB := inp1.Program.Data[fp + inp1.Offset : fp + inp1.Offset + inp1.Size]
		encoder.DeserializeAtomic(predicateB, &predicate)
	default:
		panic("implement the other mem types in readI32")
	}

	if predicate {
		// thenLines and elseLines will always be in the data segment
		thenLinesB := inp2.Program.Data[inp2.Offset : inp2.Offset + inp2.Size]
		encoder.DeserializeRaw(thenLinesB, &thenLines)
		call.Line = call.Line + int(thenLines)
	} else {
		elseLinesB := inp3.Program.Data[inp3.Offset : inp3.Offset + inp3.Size]
		encoder.DeserializeRaw(elseLinesB, &elseLines)
		call.Line = call.Line + int(elseLines)
	}
}

func time_UnixMilli (expr *CXExpression, stack *CXStack, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano() / int64(1000000))
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI64 (in int64) []byte {
	return encoder.Serialize(in)
}

func FromBool (in bool) []byte {
	return encoder.Serialize(in)
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp.Offset : fp + inp.Offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[inp.Offset : inp.Offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func ReadI64 (stack *CXStack, fp int, inp *CXArgument) (out int64) {
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp.Offset : fp + inp.Offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[inp.Offset : inp.Offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func ReadBool (stack *CXStack, fp int, inp *CXArgument) (out bool) {
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp.Offset : fp + inp.Offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[inp.Offset : inp.Offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
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
