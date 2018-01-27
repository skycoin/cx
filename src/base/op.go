package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
	"time"
)

// Creates a copy of its argument in the stack. opCopy should be used to copy an argument in a heap.
func identity (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(stack, fp, inp1)
	out1Offset := GetFinalOffset(stack, fp, out1)
	switch inp1.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + inp1Offset : fp + inp1Offset + inp1.Size]
		WriteToStack(stack, fp, out1Offset, byts)
	case MEM_DATA:
		byts := inp1.Program.Data[inp1Offset : inp1Offset + inp1.Size]
		WriteToStack(stack, fp, out1Offset, byts)
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
	
	// outB1 := ReadI32(stack, fp, )
	WriteToStack(stack, fp, out1.Offset, outB1)
}

func write_array (expr *CXExpression, stack *CXStack, fp int) {
	// inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	indexes := make([]int32, len(expr.Inputs[1:]))
	for i, idx := range expr.Inputs[1:] {
		indexes[i] = ReadI32(stack, fp, idx)
	}
	// outB1 := ReadArray(stack, fp, inp1, indexes)
	// WriteToStack(stack, fp, out1.Offset, outB1)
}

// Utilities

func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI64 (in int64) []byte {
	return encoder.Serialize(in)
}

func FromBool (in bool) []byte {
	return encoder.Serialize(in)
}

func ReadArray (stack *CXStack, fp int, inp *CXArgument, indexes []int32) (int, int) {
	var offset int
	var size int = inp.Size
	for i, idx := range indexes {
		offset += int(idx) * inp.Lengths[i]
	}
	for _, len := range indexes {
		size *= int(len)
	}

	return offset, size

	// switch inp.MemoryType {
	// case MEM_STACK:
	// 	out = stack.Stack[fp + inp.Offset : fp + inp.Offset + size]
	// case MEM_DATA:
	// 	out = inp.Program.Data[inp.Offset : inp.Offset + size]
	// default:
	// 	panic("implement the other mem types in readI32")
	// }

	// return
}

func GetFinalOffset (stack *CXStack, fp int, arg *CXArgument) int {
	offsetOffset := 0
	for i, idxArg := range arg.Indexes {
		var subSize int = 1
		for _, len := range arg.Lengths[i+1:] {
			subSize *= len
		}
		offsetOffset += int(ReadI32(stack, fp, idxArg)) * subSize * arg.Size
	}
	return arg.Offset + offsetOffset
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(stack, fp, inp)
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[fp + offset : fp + offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[offset : offset + inp.Size]
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
