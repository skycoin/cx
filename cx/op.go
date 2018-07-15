package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func GetFinalOffset(stack *CXStack, fp int, arg *CXArgument, opType int) int {
	var elt *CXArgument
	var finalOffset int = arg.Offset
	var fldIdx int
	var memType int

	if opType == MEM_READ {
		memType = arg.MemoryRead
	} else {
		memType = arg.MemoryWrite
	}

	elt = arg

	var dbg bool
	if arg.Name != "" {
		dbg = false
	}
	if dbg {
		fmt.Println("(start", arg.Name, finalOffset, arg.DereferenceOperations, opType, arg.MemoryRead, arg.MemoryWrite)
	}

	var addObjectHeader bool
	_ = addObjectHeader

	for _, op := range arg.DereferenceOperations {
		switch op {
		case DEREF_ARRAY:
			// if addObjectHeader {
			// 	finalOffset += OBJECT_HEADER_SIZE
			// 	addObjectHeader = false
			// }
			
			for i, idxArg := range elt.Indexes {
				var subSize int = 1
				for _, len := range elt.Lengths[i+1:] {
					subSize *= len
				}

				if arg.CustomType != nil {
					finalOffset += int(ReadI32(stack, fp, idxArg)) * subSize * arg.CustomType.Size
				} else {
					finalOffset += int(ReadI32(stack, fp, idxArg)) * subSize * elt.Size
					// fmt.Println("finalOffset", finalOffset)
				}
			}
		case DEREF_FIELD:
			elt = arg.Fields[fldIdx]
			// fmt.Println("offset", elt.Name, elt.Offset)
			finalOffset += elt.Offset
			fldIdx++
		case DEREF_POINTER:
			addObjectHeader = true
			for c := 0; c < elt.DereferenceLevels; c++ {
				var offset int32

				byts := stack.Stack[fp+finalOffset : fp+finalOffset+elt.Size]

				encoder.DeserializeAtomic(byts, &offset)

				if offset != 0 {
					finalOffset = int(offset) + OBJECT_HEADER_SIZE
					// finalOffset = int(offset)
				} else {
					finalOffset = 0
				}
			}
		}
		if dbg {
			fmt.Println("update", arg.Name, finalOffset)
		}
	}

	// if addObjectHeader {
	// 	finalOffset += OBJECT_HEADER_SIZE
	// 	fmt.Println("finalOffset", finalOffset)
	// }

	if memType == MEM_HEAP || memType == MEM_DATA {
		// not sure if arg.MemoryRead or arg.MemoryWrite
		if dbg {
			fmt.Println("result1", finalOffset)
			fmt.Println(")")
		}

		return finalOffset
	} else {
		if dbg {
			fmt.Println("result2", fp+finalOffset)
			fmt.Println(")")
		}

		return fp + finalOffset
	}
}

func ReadMemory(stack *CXStack, offset int, arg *CXArgument) (out []byte) {
	switch arg.MemoryRead {
	case MEM_STACK:
		out = stack.Stack[offset : offset+arg.TotalSize]
	case MEM_DATA:
		out = stack.Program.Data[offset : offset+arg.TotalSize]
	case MEM_HEAP:
		out = stack.Program.Heap.Heap[offset : offset+arg.TotalSize]
	default:
		panic("implement the other mem types")
	}
	
	return
}

// marks all the alive objects in the heap
func Mark(prgrm *CXProgram) {
	fp := 0
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			encoder.DeserializeAtomic(prgrm.Stacks[0].Stack[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)

			prgrm.Heap.Heap[heapOffset] = 1
		}

		fp += op.Size
	}
}

func MarkAndCompact(prgrm *CXProgram) {
	var fp int
	var faddr int32 = NULL_HEAP_ADDRESS_OFFSET

	// marking, setting forward addresses and updating references
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			encoder.DeserializeAtomic(prgrm.Stacks[0].Stack[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)

			if heapOffset == NULL_HEAP_ADDRESS {
				continue
			}

			// marking as alive
			prgrm.Heap.Heap[heapOffset] = 1

			for i, byt := range encoder.SerializeAtomic(faddr) {
				// setting forwarding address
				prgrm.Heap.Heap[int(heapOffset)+MARK_SIZE+i] = byt
				// updating reference
				prgrm.Stacks[0].Stack[fp+ptr.Offset+i] = byt
			}

			var objSize int32
			encoder.DeserializeAtomic(prgrm.Heap.Heap[int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE:int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE+OBJECT_SIZE], &objSize)

			faddr += int32(OBJECT_HEADER_SIZE) + objSize
		}

		fp += op.Size
	}

	// relocation of live objects
	newHeapPointer := NULL_HEAP_ADDRESS_OFFSET
	for c := NULL_HEAP_ADDRESS_OFFSET; c < prgrm.Heap.HeapPointer; {
		var forwardingAddress int32
		encoder.DeserializeAtomic(prgrm.Heap.Heap[c+MARK_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE], &forwardingAddress)

		var objSize int32
		encoder.DeserializeAtomic(prgrm.Heap.Heap[c+MARK_SIZE+FORWARDING_ADDRESS_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE], &objSize)

		if prgrm.Heap.Heap[c] == 1 {
			// setting the mark back to 0
			prgrm.Heap.Heap[c] = 0
			// then it's alive and we'll relocate the object
			for i := int32(0); i < OBJECT_HEADER_SIZE+objSize; i++ {
				prgrm.Heap.Heap[forwardingAddress+i] = prgrm.Heap.Heap[int32(c)+i]
			}
			newHeapPointer += OBJECT_HEADER_SIZE + int(objSize)
		}

		c += OBJECT_HEADER_SIZE + int(objSize)
	}

	prgrm.Heap.HeapPointer = newHeapPointer
}

// allocates memory in the heap
func AllocateSeq(prgrm *CXProgram, size int) (offset int) {
	result := prgrm.Heap.HeapPointer
	newFree := result + size

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

func WriteMemory(stack *CXStack, offset int, arg *CXArgument, byts []byte) {
	switch arg.MemoryWrite {
	case MEM_STACK:
		WriteToStack(stack, offset, byts)
	case MEM_HEAP:
		WriteToHeap(&stack.Program.Heap, offset, byts)
	case MEM_DATA:
		WriteToData(&stack.Program.Data, offset, byts)
	default:
		panic("implement the other mem types")
	}
}

// Utilities

func FromBool(in bool) []byte {
	if in {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func FromByte(in byte) []byte {
	return encoder.SerializeAtomic(in)
}

func FromStr (in string) []byte {
	return encoder.Serialize(in)
}

func FromI8 (in int8) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI32(in int32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromUI32(in uint32) []byte {
	return encoder.SerializeAtomic(in)
}

func FromI64(in int64) []byte {
	return encoder.Serialize(in)
}

func FromF32(in float32) []byte {
	return encoder.Serialize(in)
}

func FromF64(in float64) []byte {
	return encoder.Serialize(in)
}

func ReadArray(stack *CXStack, fp int, inp *CXArgument, indexes []int32) (int, int) {
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

func ReadF32A(stack *CXStack, fp int, inp *CXArgument) (out []float32) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	byts := ReadMemory(stack, offset, inp)
	byts = append(encoder.SerializeAtomic(int32(len(byts)/4)), byts...)
	encoder.DeserializeRaw(byts, &out)
	return
}

func ReadBool(stack *CXStack, fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadByte(stack *CXStack, fp int, inp *CXArgument) (out byte) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadStr(stack *CXStack, fp int, inp *CXArgument) (out string) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	
	if inp.Name == "" {
		offset = inp.HeapOffset+OBJECT_HEADER_SIZE

		var size int32
		sizeB := stack.Program.Heap.Heap[offset : offset + STR_HEADER_SIZE]

		encoder.DeserializeAtomic(sizeB, &size)
		encoder.DeserializeRaw(stack.Program.Heap.Heap[offset : offset+STR_HEADER_SIZE+int(size)], &out)
	} else {
		var off int32
		var size int32
		var byts []byte

		if inp.MemoryRead == MEM_STACK {
			byts = stack.Stack[offset : offset+TYPE_POINTER_SIZE]
			encoder.DeserializeAtomic(byts, &off)
		} else {
			byts = stack.Program.Data[offset : offset+TYPE_POINTER_SIZE]
			encoder.DeserializeAtomic(byts, &off)
		}
		
		sizeB := stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+STR_HEADER_SIZE]
		encoder.DeserializeAtomic(sizeB, &size)

		encoder.DeserializeRaw(stack.Program.Heap.Heap[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
	}
	
	// encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadI8 (stack *CXStack, fp int, inp *CXArgument) (out int8) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadI32(stack *CXStack, fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadI64(stack *CXStack, fp int, inp *CXArgument) (out int64) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadF32(stack *CXStack, fp int, inp *CXArgument) (out float32) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadF64(stack *CXStack, fp int, inp *CXArgument) (out float64) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(stack, offset, inp), &out)
	return
}

func ReadFromStack(stack *CXStack, fp int, inp *CXArgument) (out []byte) {
	offset := GetFinalOffset(stack, fp, inp, MEM_READ)
	out = ReadMemory(stack, offset, inp)
	return
}

func WriteToStack(stack *CXStack, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*stack).Stack[offset+c] = out[c]
	}
}

func WriteToHeap(heap *CXHeap, offset int, out []byte) {
	// size := encoder.Serialize(int32(len(out)))

	// var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	
// for c := 5; c < OBJECT_HEADER_SIZE; c++ {
	// 	header[c] = size[c - 5]
	// }

	// for c := 0; c < OBJECT_HEADER_SIZE; c++ {
	// 	(*heap).Heap[offset + c] = header[c]
	// }

	for c := 0; c < len(out); c++ {
		// (*heap).Heap[offset + OBJECT_HEADER_SIZE + c] = out[c]
		(*heap).Heap[offset+c] = out[c]
	}
}

func WriteToData(data *Data, offset int, out []byte) {
	for c := 0; c < len(out); c++ {
		(*data)[offset+c] = out[c]
	}
}
