package base

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
	"time"
)

// func GetDereferenceOffset (stack *CXStack, fp int, offset int32, offsetOffset int, arg *CXArgument, flag bool) (finalOffset int32) {
// 	for c := 0; c < arg.DereferenceLevels; c++ {
// 		switch arg.MemoryType {
// 		case MEM_STACK:
// 			var byts []byte
// 			if flag {
// 				byts = stack.Stack[fp + int(offset) + offsetOffset : fp + int(offset) + offsetOffset + arg.Size]
// 			} else {
// 				byts = stack.Stack[fp + int(offset) : fp + int(offset) + arg.Size]
// 			}
// 			encoder.DeserializeAtomic(byts, &finalOffset)
// 			return
// 		case MEM_DATA:
// 			var byts []byte
// 			if flag {
// 				byts = arg.Program.Data[int(offset) + offsetOffset : int(offset) + offsetOffset + arg.Size]
// 			} else {
// 				byts = arg.Program.Data[int(offset) : int(offset) + arg.Size]
// 			}
// 			encoder.DeserializeAtomic(byts, &finalOffset)
// 			return
// 		default:
// 			panic("implement the other mem types in readI32")
// 		}
// 	}
// 	return offset
// }

func GetDereferenceOffset (memory *[]byte, fp int, offset int32, size int, levels int) int32 {
	for c := 0; c < levels; c++ {
		var byts []byte
		byts = (*memory)[fp + int(offset) : fp + int(offset) + size]
		encoder.DeserializeAtomic(byts, &offset)
	}
	return offset
}

func GetFinalOffset (stack *CXStack, fp int, arg *CXArgument) int {
	var elt *CXArgument
	var finalOffset int = arg.Offset
	var fldIdx int
	elt = arg

	for _, op := range arg.DereferenceOperations {
			switch op {
			case DEREF_ARRAY:
				for i, idxArg := range elt.Indexes {
					var subSize int = 1
					for _, len := range elt.Lengths[i+1:] {
						subSize *= len
					}
					finalOffset += int(ReadI32(stack, fp, idxArg)) * subSize * elt.Size
				}
			case DEREF_FIELD:
				elt = arg.Fields[fldIdx]
				finalOffset += elt.Offset
				fldIdx++
			case DEREF_POINTER:
				for c := 0; c < elt.DereferenceLevels; c++ {
					switch arg.MemoryType {
					case MEM_STACK:
						var offset int32
						byts := stack.Stack[fp + finalOffset : fp + finalOffset + elt.Size]
						encoder.DeserializeAtomic(byts, &offset)
						finalOffset = int(offset)
					case MEM_DATA:
						var offset int32
						byts := arg.Program.Data[finalOffset : finalOffset + elt.Size]
						encoder.DeserializeAtomic(byts, &offset)
						finalOffset = int(offset)
					default:
						panic("implement the other mem types in readI32")
					}
				}
			}
		}

	return finalOffset
}

// func GetFinalOffset (stack *CXStack, fp int, arg *CXArgument) int {
// 	// checking if needs to dereferenced
// 	var offset int32
// 	offset = int32(arg.Offset)
// 	if len(arg.Indexes) < 1 && len(arg.Fields) < 1 {
// 		if arg.IsPointer {
// 			if arg.DereferenceLevels > 0 {
// 				for c := 0; c < arg.DereferenceLevels; c++ {
// 					switch arg.MemoryType {
// 					case MEM_STACK:
// 						byts := stack.Stack[fp + int(offset) : fp + int(offset) + arg.Size]
// 						encoder.DeserializeAtomic(byts, &offset)
// 					case MEM_DATA:
// 						byts := arg.Program.Data[offset : int(offset) + arg.Size]
// 						encoder.DeserializeAtomic(byts, &offset)
// 					default:
// 						panic("implement the other mem types in readI32")
// 					}
// 				}
// 			}
// 			return int(offset)
// 		} else {
// 			switch arg.MemoryType {
// 			case MEM_STACK:
// 				return fp + int(offset)
// 			case MEM_DATA:
// 				return int(offset)
// 			}
// 		}
// 	}

// 	offsetOffset := 0
// 	if arg.IsDereferenceFirst {
// 		offset = GetDereferenceOffset(stack, fp, offset, offsetOffset, arg, false)
// 		for i, idxArg := range arg.Indexes {
// 			var subSize int = 1
// 			for _, len := range arg.Lengths[i+1:] {
// 				subSize *= len
// 			}
// 			offsetOffset += int(ReadI32(stack, fp, idxArg)) * subSize * arg.PointeeSize
// 		}
// 	} else {
// 		for i, idxArg := range arg.Indexes {
// 			var subSize int = 1
// 			for _, len := range arg.Lengths[i+1:] {
// 				subSize *= len
// 			}
// 			offsetOffset += int(ReadI32(stack, fp, idxArg)) * subSize * arg.Size
// 		}

// 		offset = GetDereferenceOffset(stack, fp, offset, offsetOffset, arg, true)
// 		if len(arg.Fields) > 0 && arg.IsPointer {
// 			offsetOffset = 0
// 		}
// 	}

// 	fieldOffset := 0
// 	if len(arg.Fields) > 0 {
// 		for _, fld := range arg.Fields {
// 			fieldOffset += fld.Offset
// 			for i, idxArg := range fld.Indexes {
// 				var subSize int = 1
// 				for _, len := range fld.Lengths[i+1:] {
// 					subSize *= len
// 				}
// 				fieldOffset += int(ReadI32(stack, fp, idxArg)) * subSize * fld.Size
// 			}
// 		}
// 	}

// 	fmt.Println(arg.Name, offset, offsetOffset, fieldOffset)
	
// 	if arg.IsPointer {
// 		return int(offset) + offsetOffset + fieldOffset
// 	} else {
// 		return fp + int(offset) + offsetOffset + fieldOffset
// 	}
// }

// Creates a copy of its argument in the stack. opCopy should be used to copy an argument in a heap.
func identity (expr *CXExpression, stack *CXStack, fp int) {

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(stack, fp, inp1)
	out1Offset := GetFinalOffset(stack, fp, out1)

	if out1.IsPointer && out1.DereferenceLevels != out1.IndirectionLevels && !inp1.IsPointer {
		switch inp1.MemoryType {
		case MEM_STACK:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToStack(stack, out1Offset, byts)
		case MEM_DATA:
			byts := encoder.SerializeAtomic(int32(inp1Offset))
			WriteToStack(stack, out1Offset, byts)
		default:
			panic("implement the other mem types in readI32")
		}
	} else if inp1.IsReference {
		WriteToStack(stack, out1Offset, FromI32(int32(inp1Offset)))
	} else {
		switch inp1.MemoryType {
		case MEM_STACK:
			byts := stack.Stack[inp1Offset : inp1Offset + inp1.Size]
			WriteToStack(stack, out1Offset, byts)
		case MEM_DATA:
			byts := inp1.Program.Data[inp1Offset : inp1Offset + inp1.Size]
			WriteToStack(stack, out1Offset, byts)
		default:
			panic("implement the other mem types in readI32")
		}
	}
}

func i32_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func i32_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func i32_gt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func i32_lt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func i32_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

func i64_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) + ReadI64(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func i64_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) - ReadI64(stack, fp, inp2))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
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
	// inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	inp1 := expr.Inputs[0]
	var predicate bool
	// var thenLines int32
	// var elseLines int32

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
		// thenLinesB := inp2.Program.Data[inp2.Offset : inp2.Offset + inp2.Size]
		// encoder.DeserializeRaw(thenLinesB, &thenLines)
		// call.Line = call.Line + int(thenLines)

		call.Line = call.Line + expr.ThenLines
	} else {
		// elseLinesB := inp3.Program.Data[inp3.Offset : inp3.Offset + inp3.Size]
		// encoder.DeserializeRaw(elseLinesB, &elseLines)
		// call.Line = call.Line + int(elseLines)

		call.Line = call.Line + expr.ElseLines
	}
}

func time_UnixMilli (expr *CXExpression, stack *CXStack, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano() / int64(1000000))
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
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
	WriteToStack(stack, GetFinalOffset(stack, fp, out1), outB1)
}

func write_array (expr *CXExpression, stack *CXStack, fp int) {
	indexes := make([]int32, len(expr.Inputs[1:]))
	for i, idx := range expr.Inputs[1:] {
		indexes[i] = ReadI32(stack, fp, idx)
	}
}

// Utilities

func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI64 (in int64) []byte {
	return encoder.Serialize(in)
}

func FromBool (in bool) []byte {
	if in {
		return []byte{1}
	} else {
		return []byte{0}
	}
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
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(stack, fp, inp)
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[offset : offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[offset : offset + inp.Size]
		encoder.DeserializeAtomic(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func ReadFromStack (stack *CXStack, fp int, inp *CXArgument) (out []byte) {
	offset := GetFinalOffset(stack, fp, inp)
	switch inp.MemoryType {
	case MEM_STACK:
		return stack.Stack[offset : offset + inp.Size]
	case MEM_DATA:
		return inp.Program.Data[offset : offset + inp.Size]
	default:
		panic("implement the other mem types in readI32")
	}
}

func ReadI64 (stack *CXStack, fp int, inp *CXArgument) (out int64) {
	offset := GetFinalOffset(stack, fp, inp)
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[offset : offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[offset : offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func ReadBool (stack *CXStack, fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(stack, fp, inp)
	switch inp.MemoryType {
	case MEM_STACK:
		byts := stack.Stack[offset : offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	case MEM_DATA:
		byts := inp.Program.Data[offset : offset + inp.Size]
		encoder.DeserializeRaw(byts, &out)
	default:
		panic("implement the other mem types in readI32")
	}
	
	return
}

func WriteToStack (stack *CXStack, offset int, out []byte) {
	// fmt.Println("before", stack)
	for c := 0; c < len(out); c++ {
		(*stack).Stack[offset + c] = out[c]
	}
	// fmt.Println("after", stack)
}
