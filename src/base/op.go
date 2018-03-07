package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func GetFinalOffset (stack *CXStack, fp int, arg *CXArgument) int {
	// // it's from the data segment
	// if arg.MemoryType == MEM_DATA {
	// 	if len(arg.Indexes) > 0 {
	// 		fmt.Println("hmm", arg.Name, arg.Offset, arg.Indexes[0].MemoryType)
	// 	}
		
	// 	return arg.Offset
	// }
	
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

	if arg.IsPointer || arg.MemoryType == MEM_DATA {
		return finalOffset
	} else {
		return fp + finalOffset
	}
}

func ReadMemory (stack *CXStack, offset int, arg *CXArgument) (out []byte) {
	switch arg.MemoryType {
	case MEM_STACK:
		out = stack.Stack[offset : offset + arg.TotalSize]
	case MEM_DATA:
		out = arg.Program.Data[offset : offset + arg.TotalSize]
	case MEM_HEAP:
		out = arg.Program.Heap[offset : offset + arg.TotalSize]
	default:
		panic("implement the other mem types")
	}
	return
}

func WriteMemory (stack *CXStack, offset int, arg *CXArgument, byts []byte) {
	switch arg.MemoryType {
	case MEM_STACK:
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&arg.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&arg.Program.Data, offset, byts)
		}
	case MEM_HEAP:

		// encoder.Serialize(int32(len(byts)))
		
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&arg.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&arg.Program.Data, offset, byts)
		}
	case MEM_DATA:
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&arg.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&arg.Program.Data, offset, byts)
		}
	default:
		panic("implement the other mem types")
	}
}

// Utilities

func FromBool (in bool) []byte {
	if in {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func FromByte (in byte) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI32 (in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromUI32 (in uint32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI64 (in int64) []byte {
	return encoder.Serialize(in)
}

func FromF32 (in float32) []byte {
	return encoder.Serialize(in)
}

func FromF64 (in float64) []byte {
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
}

func ReadF32A (stack *CXStack, fp int, inp *CXArgument) (out []float32) {
	// Only used by native functions (i.e. functions implemented in Golang)
	offset := GetFinalOffset(stack, fp, inp)
	byts := ReadMemory(stack, offset, inp)
	byts = append(encoder.SerializeAtomic(int32(len(byts) / 4)), byts...)
	encoder.DeserializeRaw(byts, &out)
	return
}

func ReadBool (stack *CXStack, fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadByte (stack *CXStack, fp int, inp *CXArgument) (out byte) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeAtomic(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadStr (stack *CXStack, fp int, inp *CXArgument) (out string) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadI32 (stack *CXStack, fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeAtomic(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadI64 (stack *CXStack, fp int, inp *CXArgument) (out int64) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadF32 (stack *CXStack, fp int, inp *CXArgument) (out float32) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadF64 (stack *CXStack, fp int, inp *CXArgument) (out float64) {
	offset := GetFinalOffset(stack, fp, inp)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadFromStack (stack *CXStack, fp int, inp *CXArgument) (out []byte) {
	offset := GetFinalOffset(stack, fp, inp)
	out = ReadMemory(stack, offset, inp)
	return
}

func WriteToStack (stack *CXStack, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*stack).Stack[offset + c] = out[c]
	}
}

func WriteToHeap (heap *Heap, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*heap)[offset + c] = out[c]
	}
}

func WriteToData (data *Data, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*data)[offset + c] = out[c]
	}
}
