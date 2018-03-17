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
				case MEM_HEAP:
					var offset int32
					byts := stack.Stack[fp + finalOffset : fp + finalOffset + elt.Size]
					// byts := stack.Program.Heap.Heap[finalOffset : finalOffset + elt.Size]
					encoder.DeserializeAtomic(byts, &offset)
					finalOffset = int(offset) + OBJECT_HEADER_SIZE
				case MEM_DATA:
					var offset int32
					byts := stack.Program.Data[finalOffset : finalOffset + elt.Size]
					encoder.DeserializeAtomic(byts, &offset)
					finalOffset = int(offset)
				default:
					panic("implement the other mem types in readI32")
				}
			}
		}
	}

	// return fp + finalOffset

	// if arg.IsPointer || arg.MemoryType == MEM_DATA {
	if arg.MemoryType == MEM_DATA {
		// if arg.IsPointer {
		// 	fmt.Println("hi", arg.Name, arg.Offset, arg.MemoryType)
		// }
		
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
		out = stack.Program.Data[offset : offset + arg.TotalSize]
	case MEM_HEAP:
		out = stack.Program.Heap.Heap[offset : offset + arg.TotalSize]
		// out = arg.Program.Heap.Heap[offset + NULL_HEAP_ADDRESS_OFFSET : offset + NULL_HEAP_ADDRESS_OFFSET + arg.TotalSize]
	default:
		panic("implement the other mem types")
	}
	return
}

func GCMarkCompact (stack *CXStack, heap *CXHeap) {
	
}

// marks all the alive objects in the heap
func Mark (prgrm *CXProgram) {
	fp := 0
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			encoder.DeserializeAtomic(prgrm.Stacks[0].Stack[fp + ptr.Offset : fp + ptr.Offset + TYPE_POINTER_SIZE], &heapOffset)
			
			prgrm.Heap.Heap[heapOffset] = 1
		}
		
		fp += op.Size
	}
}

func ComputeLocations (prgrm *CXProgram) {
	// for c := int32(0); c < INIT_HEAP_SIZE; {
	// 	var objSize int32
		
	// 	encoder.DeserializeAtomic(prgrm.Heap.Heap[int(heapOffset) + 1 + TYPE_POINTER_SIZE : int(heapOffset) + 1 + TYPE_POINTER_SIZE + OBJECT_SIZE], &objSize)

	// 	c += objSize
	// }
}

func MarkAndCompact (prgrm *CXProgram) {
	var fp int
	var faddr int32 = NULL_HEAP_ADDRESS_OFFSET

	// marking, setting forward addresses and updating references
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			encoder.DeserializeAtomic(prgrm.Stacks[0].Stack[fp + ptr.Offset : fp + ptr.Offset + TYPE_POINTER_SIZE], &heapOffset)

			if heapOffset == NULL_HEAP_ADDRESS {
				continue
			}

			// marking as alive
			prgrm.Heap.Heap[heapOffset] = 1

			for i, byt := range encoder.SerializeAtomic(faddr) {
				// setting forwarding address
				prgrm.Heap.Heap[int(heapOffset) + MARK_SIZE + i] = byt
				// updating reference
				prgrm.Stacks[0].Stack[fp + ptr.Offset + i] = byt
			}
			
			var objSize int32
			encoder.DeserializeAtomic(prgrm.Heap.Heap[int(heapOffset) + MARK_SIZE + TYPE_POINTER_SIZE : int(heapOffset) + MARK_SIZE + TYPE_POINTER_SIZE + OBJECT_SIZE], &objSize)

			faddr += int32(OBJECT_HEADER_SIZE) + objSize
		}
		
		fp += op.Size
	}

	// relocation of live objects
	newHeapPointer := NULL_HEAP_ADDRESS_OFFSET
	for c := NULL_HEAP_ADDRESS_OFFSET; c < prgrm.Heap.HeapPointer; {
		var forwardingAddress int32
		encoder.DeserializeAtomic(prgrm.Heap.Heap[c + MARK_SIZE : c + MARK_SIZE + FORWARDING_ADDRESS_SIZE], &forwardingAddress)

		var objSize int32
		encoder.DeserializeAtomic(prgrm.Heap.Heap[c + MARK_SIZE + FORWARDING_ADDRESS_SIZE : c + MARK_SIZE + FORWARDING_ADDRESS_SIZE + OBJECT_SIZE], &objSize)


		if prgrm.Heap.Heap[c] == 1 {
			// setting the mark back to 0
			prgrm.Heap.Heap[c] = 0
			// then it's alive and we'll relocate the object
			for i := int32(0); i < OBJECT_HEADER_SIZE + objSize; i++ {
				prgrm.Heap.Heap[forwardingAddress + i] = prgrm.Heap.Heap[int32(c) + i]
			}
			newHeapPointer += OBJECT_HEADER_SIZE + int(objSize)
		}
		
		c += OBJECT_HEADER_SIZE + int(objSize)
	}

	prgrm.Heap.HeapPointer = newHeapPointer
}

// allocates memory in the heap
func AllocateSeq (prgrm *CXProgram, size int) (offset int) {
	// heap := prgrm.Heap
	result := prgrm.Heap.HeapPointer
	newFree := result + size

	// if prgrm.Heap.HeapPointer > 0 {
	// 	MarkAndCompact(prgrm)
	// 	result = prgrm.Heap.HeapPointer
	// 	newFree = prgrm.Heap.HeapPointer + size
	// }

	if newFree > INIT_HEAP_SIZE {
		// call GC
		MarkAndCompact(prgrm)
		result = prgrm.Heap.HeapPointer
		newFree = prgrm.Heap.HeapPointer + size

		if newFree > INIT_HEAP_SIZE {
			// heap exhausted
			panic("heap exhausted")
		}
	}

	prgrm.Heap.HeapPointer = newFree

	return result
}

func WriteMemory (stack *CXStack, offset int, arg *CXArgument, byts []byte) {
	switch arg.MemoryType {
	case MEM_STACK:
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&stack.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&stack.Program.Data, offset, byts)
		}
	case MEM_HEAP:
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&stack.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&stack.Program.Data, offset, byts)
		}
	case MEM_DATA:
		switch arg.MemoryType {
		case MEM_STACK:
			WriteToStack(stack, offset, byts)
		case MEM_HEAP:
			WriteToHeap(&stack.Program.Heap, offset, byts)
		case MEM_DATA:
			WriteToData(&stack.Program.Data, offset, byts)
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

func WriteToHeap (heap *CXHeap, offset int, out []byte) {
	size := encoder.Serialize(int32(len(out)))
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c - 5]
	}
	
	for c := 0; c < OBJECT_HEADER_SIZE; c++ {
		(*heap).Heap[offset + c] = header[c]
	}
	

	for c := 0; c < len(out); c++ {
		(*heap).Heap[offset + OBJECT_HEADER_SIZE + c] = out[c]
	}

	// for c := 0; c < len(out); c++ {
	// 	(*heap)[offset + c] = out[c]
	// }
}

func WriteToData (data *Data, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*data)[offset + c] = out[c]
	}
}
