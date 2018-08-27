package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func GetFinalOffset(mem []byte, fp int, arg *CXArgument, opType int) int {
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
	
	var dbg bool = true
	if arg.Name != "" {
		dbg = true
	}
	if dbg {
		fmt.Println("(start", arg.Name, finalOffset, arg.DereferenceOperations, opType, arg.MemoryRead, arg.MemoryWrite)
	}

	var addObjectHeader bool
	_ = addObjectHeader

	for _, op := range arg.DereferenceOperations {
		switch op {
		case DEREF_ARRAY:
			for i, idxArg := range elt.Indexes {
				var subSize int = 1
				for _, len := range elt.Lengths[i+1:] {
					subSize *= len
				}

				var sizeToUse int
				if arg.CustomType != nil {
					sizeToUse = arg.CustomType.Size
				} else if elt.IsSlice {
					sizeToUse = elt.TotalSize
				} else {
					sizeToUse = elt.Size
				}
				
				// if arg.CustomType != nil {
				// 	// finalOffset += int(ReadI32(mem, fp, idxArg)) * subSize * arg.CustomType.Size
				// 	finalOffset += int(ReadI32(mem, fp, idxArg)) * subSize * sizeToUse
				// } else {
				// 	// finalOffset += int(ReadI32(mem, fp, idxArg)) * subSize * elt.Size
				// 	finalOffset += int(ReadI32(mem, fp, idxArg)) * subSize * sizeToUse
				// }
				finalOffset += int(ReadI32(mem, fp, idxArg)) * subSize * sizeToUse
			}
		case DEREF_FIELD:
			elt = arg.Fields[fldIdx]
			finalOffset += elt.Offset
			fldIdx++
		case DEREF_POINTER:
			addObjectHeader = true
			for c := 0; c < elt.DereferenceLevels; c++ {
				var offset int32
				var byts []byte

				byts = mem[fp+finalOffset : fp+finalOffset+elt.Size]

				encoder.DeserializeAtomic(byts, &offset)
				
				if offset != 0 {
					if arg.IsSlice {
						if elt.Type == TYPE_STR {
							finalOffset = int(offset) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + STR_HEADER_SIZE
						} else {
							finalOffset = int(offset) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE
						}
					} else {
						finalOffset = int(offset) + OBJECT_HEADER_SIZE
					}
					// finalOffset = int(offset) + OBJECT_HEADER_SIZE
				} else {
					// if elt.Type == TYPE_STR {
						
					// } else {
					// 	finalOffset = 0
					// }
					finalOffset = 0
				}
			}
		}
		if dbg {
			fmt.Println("update", arg.Name, finalOffset)
		}
	}

	if memType == MEM_HEAP || memType == MEM_DATA {
		// not sure if arg.MemoryRead or arg.MemoryWrite
		if dbg {
			fmt.Println("result1", finalOffset, ")")
		}

		return finalOffset
	} else {
		if dbg {
			fmt.Println("result2", fp+finalOffset, ")")
		}

		return fp + finalOffset
	}
}

func ReadMemory(mem []byte, offset int, arg *CXArgument) (out []byte) {
	switch arg.MemoryRead {
	case MEM_STACK:
		out = mem[offset : offset+arg.TotalSize]
	case MEM_DATA:
		out = mem[offset : offset+arg.TotalSize]
	case MEM_HEAP:
		out = mem[offset : offset+arg.TotalSize]
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
			encoder.DeserializeAtomic(prgrm.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)

			prgrm.Memory[heapOffset] = 1
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
			encoder.DeserializeAtomic(prgrm.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)

			if heapOffset == NULL_HEAP_ADDRESS {
				continue
			}

			// marking as alive
			prgrm.Memory[heapOffset] = 1

			for i, byt := range encoder.SerializeAtomic(faddr) {
				// setting forwarding address
				prgrm.Memory[int(heapOffset)+MARK_SIZE+i] = byt
				// updating reference
				prgrm.Memory[fp+ptr.Offset+i] = byt
			}

			var objSize int32
			encoder.DeserializeAtomic(prgrm.Memory[int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE:int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE+OBJECT_SIZE], &objSize)

			faddr += int32(OBJECT_HEADER_SIZE) + objSize
		}

		fp += op.Size
	}

	// relocation of live objects
	newHeapPointer := NULL_HEAP_ADDRESS_OFFSET
	for c := NULL_HEAP_ADDRESS_OFFSET; c < prgrm.HeapPointer; {
		var forwardingAddress int32
		encoder.DeserializeAtomic(prgrm.Memory[c+MARK_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE], &forwardingAddress)

		var objSize int32
		encoder.DeserializeAtomic(prgrm.Memory[c+MARK_SIZE+FORWARDING_ADDRESS_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE], &objSize)

		if prgrm.Memory[c] == 1 {
			// setting the mark back to 0
			prgrm.Memory[c] = 0
			// then it's alive and we'll relocate the object
			for i := int32(0); i < OBJECT_HEADER_SIZE+objSize; i++ {
				prgrm.Memory[forwardingAddress+i] = prgrm.Memory[int32(c)+i]
			}
			newHeapPointer += OBJECT_HEADER_SIZE + int(objSize)
		}

		c += OBJECT_HEADER_SIZE + int(objSize)
	}

	prgrm.HeapPointer = newHeapPointer
}

// allocates memory in the heap
func AllocateSeq(prgrm *CXProgram, size int) (offset int) {
	result := prgrm.HeapPointer
	newFree := result + size

	if newFree > INIT_HEAP_SIZE {
		// call GC
		MarkAndCompact(prgrm)
		result = prgrm.HeapPointer
		newFree = prgrm.HeapPointer + size

		if newFree > INIT_HEAP_SIZE {
			// heap exhausted
			panic("heap exhausted")
		}
	}

	prgrm.HeapPointer = newFree

	return result
}

func WriteMemory(mem []byte, offset int, byts []byte) {
	for c := 0; c < len(byts); c++ {
		mem[offset+c] = byts[c]
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

func ReadArray(mem []byte, fp int, inp *CXArgument, indexes []int32) (int, int) {
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

func ReadF32A(mem []byte, fp int, inp *CXArgument) (out []float32) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	byts := ReadMemory(mem, offset, inp)
	byts = append(encoder.SerializeAtomic(int32(len(byts)/4)), byts...)
	encoder.DeserializeRaw(byts, &out)
	return
}

func ReadBool(mem []byte, fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadByte(mem []byte, fp int, inp *CXArgument) (out byte) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(mem, offset, inp), &out)
	return
}

// maybe delete it
func ReadSlice(mem []byte, fp int, inp *CXArgument) (int, int) {
	return 0, 0
}

func ReadStr(mem []byte, fp int, inp *CXArgument) (out string) {
	// if inp.HeapOffset == 0 {
	// 	return ""
	// }
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)

	if inp.Name == "" {
		offset = inp.HeapOffset+OBJECT_HEADER_SIZE

		var size int32
		sizeB := mem[offset : offset + STR_HEADER_SIZE]

		fmt.Println("memory", mem)
		println("offset", offset, size)
		
		encoder.DeserializeAtomic(sizeB, &size)
		encoder.DeserializeRaw(mem[offset : offset+STR_HEADER_SIZE+int(size)], &out)
	} else {
		var off int32
		var size int32
		var byts []byte

		if inp.MemoryRead == MEM_STACK || inp.IsSlice {
			byts = mem[offset : offset+TYPE_POINTER_SIZE]
			encoder.DeserializeAtomic(byts, &off)
		} else {
			byts = mem[offset : offset+TYPE_POINTER_SIZE]
			encoder.DeserializeAtomic(byts, &off)
		}

		var sizeB []byte
		if inp.IsSlice {
			sizeB = mem[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : off+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+SLICE_HEADER_SIZE]
		} else {
			sizeB = mem[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+STR_HEADER_SIZE]
		}

		// sizeB := mem.Program.Heap.Heap[off : off+STR_HEADER_SIZE]
		encoder.DeserializeAtomic(sizeB, &size)

		encoder.DeserializeRaw(mem[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
	}
	
	// encoder.DeserializeRaw(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadI8 (mem []byte, fp int, inp *CXArgument) (out int8) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadI32(mem []byte, fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeAtomic(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadI64(mem []byte, fp int, inp *CXArgument) (out int64) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadF32(mem []byte, fp int, inp *CXArgument) (out float32) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(mem, offset, inp), &out)
	return
}

func ReadF64(mem []byte, fp int, inp *CXArgument) (out float64) {
	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
	encoder.DeserializeRaw(ReadMemory(mem, offset, inp), &out)
	return
}

// func ReadFromStack(mem []byte, fp int, inp *CXArgument) (out []byte) {
// 	offset := GetFinalOffset(mem, fp, inp, MEM_READ)
// 	out = ReadMemory(mem, offset, inp)
// 	return
// }

// func WriteToStack(mem []byte, offset int, out []byte) {
// 	for c := 0; c < len(out); c++ {
// 		(*mem).Stack[offset+c] = out[c]
// 	}
// }

// func WriteToHeap(heap *CXHeap, offset int, out []byte) {
// 	for c := 0; c < len(out); c++ {
// 		(*heap).Heap[offset+c] = out[c]
// 	}
// }

// func WriteToData(data *Data, offset int, out []byte) {
// 	for c := 0; c < len(out); c++ {
// 		(*data)[offset+c] = out[c]
// 	}
// }
