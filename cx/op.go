package cxcore

import (
	"fmt"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

// GetSize ...
func GetSize(arg *CXArgument) int {
	if len(arg.Fields) > 0 {
		return GetSize(arg.Fields[len(arg.Fields)-1])
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		if deref == DEREF_SLICE || deref == DEREF_ARRAY {
			return arg.Size
		}
	}

	for decl := range arg.DeclarationSpecifiers {
		if decl == DECL_POINTER {
			return arg.TotalSize
		}
	}

	if arg.CustomType != nil {
		return arg.CustomType.Size
	}

	return arg.TotalSize
}

// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}

// CalculateDereferences ...
func CalculateDereferences(arg *CXArgument, finalOffset *int, fp int, dbg bool) {
	var isPointer bool
	var baseOffset int
	var sizeofElement int

	idxCounter := 0
	for _, op := range arg.DereferenceOperations {
		switch op {
		case DEREF_SLICE:
			if len(arg.Indexes) == 0 {
				continue
			}

			isPointer = false
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+TYPE_POINTER_SIZE]

			_, err := encoder.DeserializeAtomic(byts, &offset)
			if err != nil {
				panic(err)
			}
			*finalOffset = int(offset)

			baseOffset = *finalOffset

			*finalOffset += OBJECT_HEADER_SIZE
			*finalOffset += SLICE_HEADER_SIZE

			sizeToUse := GetDerefSize(arg)
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeToUse

			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeToUse) {
				panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}

			idxCounter++
		case DEREF_ARRAY:
			if len(arg.Indexes) == 0 {
				continue
			}
			var subSize = int(1)
			for _, len := range arg.Lengths[idxCounter+1:] {
				subSize *= len
			}

			sizeToUse := GetDerefSize(arg)

			baseOffset = *finalOffset
			sizeofElement = subSize * sizeToUse
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeofElement
			idxCounter++
		case DEREF_POINTER:
			isPointer = true
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+TYPE_POINTER_SIZE]

			_, err := encoder.DeserializeAtomic(byts, &offset)
			if err != nil {
				panic(err)
			}
			*finalOffset = int(offset)
		}
		if dbg {
			fmt.Println("\tupdate", arg.Name, arg.DereferenceOperations, *finalOffset, PROGRAM.Memory[*finalOffset:*finalOffset+10])
		}
	}
	if dbg {
		fmt.Println("\tupdate", arg.Name, arg.DereferenceOperations, *finalOffset, PROGRAM.Memory[*finalOffset:*finalOffset+10])
	}

	// if *finalOffset >= PROGRAM.HeapStartsAt {
	if *finalOffset >= PROGRAM.HeapStartsAt && isPointer {
		// then it's an object
		*finalOffset += OBJECT_HEADER_SIZE
		if arg.IsSlice {
			*finalOffset += SLICE_HEADER_SIZE
			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeofElement) {
				panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}
}

// GetStrOffset ...
func GetStrOffset(fp int, arg *CXArgument) int {
	strOffset := GetFinalOffset(fp, arg)
	if arg.Name != "" {
		// then it's not a literal
		var offset = int32(0)
		_, err := encoder.DeserializeAtomic(PROGRAM.Memory[strOffset:strOffset+TYPE_POINTER_SIZE], &offset)
		if err != nil {
			panic(err)
		}
		strOffset = int(offset)
	}
	return strOffset
}

// GetFinalOffset ...
func GetFinalOffset(fp int, arg *CXArgument) int {
	// defer RuntimeError(PROGRAM)
	// var elt *CXArgument
	var finalOffset = int(arg.Offset)
	// var fldIdx int

	// elt = arg

	var dbg bool
	if arg.Name != "" {
		dbg = false
	}

	if finalOffset < STACK_SIZE {
		// then it's in the stack, not in data or heap
		finalOffset += fp
	}

	if dbg {
		fmt.Println("(start", arg.Name, fmt.Sprintf("%s:%d", arg.FileName, arg.FileLine), arg.DereferenceOperations, finalOffset, PROGRAM.Memory[finalOffset:finalOffset+10])
	}

	// elt = arg
	CalculateDereferences(arg, &finalOffset, fp, dbg)
	for _, fld := range arg.Fields {
		// elt = fld
		finalOffset += fld.Offset
		CalculateDereferences(fld, &finalOffset, fp, dbg)
	}

	if dbg {
		fmt.Println("\t\tresult", finalOffset, PROGRAM.Memory[finalOffset:finalOffset+10], "...)")
	}

	return finalOffset
}

// ReadMemory ...
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}

// Mark marks all the alive objects in the heap
func Mark(prgrm *CXProgram) {
	fp := 0
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			_, err := encoder.DeserializeAtomic(prgrm.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)
			if err != nil {
				panic(err)
			}

			prgrm.Memory[heapOffset] = 1
		}

		fp += op.Size
	}
}

// MarkAndCompact ...
func MarkAndCompact() {
	var fp int
	var faddr = int32(NULL_HEAP_ADDRESS_OFFSET)

	// marking, setting forward addresses and updating references
	for c := 0; c <= PROGRAM.CallCounter; c++ {
		op := PROGRAM.CallStack[c].Operator

		if op == nil {
			continue
		}

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			_, err := encoder.DeserializeAtomic(PROGRAM.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)
			if err != nil {
				panic(err)
			}

			if heapOffset == NULL_HEAP_ADDRESS {
				continue
			}

			// marking as alive
			PROGRAM.Memory[heapOffset] = 1

			for i, byt := range encoder.SerializeAtomic(faddr) {
				// setting forwarding address
				PROGRAM.Memory[int(heapOffset)+MARK_SIZE+i] = byt
				// updating reference
				PROGRAM.Memory[fp+ptr.Offset+i] = byt
			}

			var objSize int32
			_, err = encoder.DeserializeAtomic(PROGRAM.Memory[int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE:int(heapOffset)+MARK_SIZE+TYPE_POINTER_SIZE+OBJECT_SIZE], &objSize)
			if err != nil {
				panic(err)
			}

			faddr += int32(OBJECT_HEADER_SIZE) + objSize
		}

		fp += op.Size
	}

	// relocation of live objects
	newHeapPointer := NULL_HEAP_ADDRESS_OFFSET
	for c := NULL_HEAP_ADDRESS_OFFSET; c < PROGRAM.HeapPointer; {
		var forwardingAddress int32
		_, err := encoder.DeserializeAtomic(PROGRAM.Memory[PROGRAM.HeapStartsAt+c+MARK_SIZE:PROGRAM.HeapStartsAt+c+MARK_SIZE+FORWARDING_ADDRESS_SIZE], &forwardingAddress)
		if err != nil {
			panic(err)
		}

		var objSize int32
		_, err = encoder.DeserializeAtomic(PROGRAM.Memory[PROGRAM.HeapStartsAt+c+MARK_SIZE+FORWARDING_ADDRESS_SIZE:PROGRAM.HeapStartsAt+c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE], &objSize)
		if err != nil {
			panic(err)
		}

		if PROGRAM.Memory[c] == 1 {
			// setting the mark back to 0
			PROGRAM.Memory[c] = 0
			// then it's alive and we'll relocate the object
			for i := int32(0); i < OBJECT_HEADER_SIZE+objSize; i++ {
				PROGRAM.Memory[forwardingAddress+i] = PROGRAM.Memory[int32(c)+i]
			}
			newHeapPointer += OBJECT_HEADER_SIZE + int(objSize)
		}

		c += OBJECT_HEADER_SIZE + int(objSize)
	}

	PROGRAM.HeapPointer = newHeapPointer
}

// ResizeMemory ...
func ResizeMemory(newMemSize int, isExpand bool) {
	if newMemSize > MAX_HEAP_SIZE {
		// heap exhausted
		panic(HEAP_EXHAUSTED_ERROR)
	}

	if isExpand {
		PROGRAM.Memory = append(PROGRAM.Memory, make([]byte, MEMORY_SIZE-newMemSize)...)
		MEMORY_SIZE = newMemSize
	} else {
		PROGRAM.Memory = PROGRAM.Memory[:newMemSize]
		MEMORY_SIZE = newMemSize
	}
}

// AllocateSeq allocates memory in the heap
func AllocateSeq(size int) (offset int) {
	result := PROGRAM.HeapStartsAt + PROGRAM.HeapPointer
	newFree := PROGRAM.HeapPointer + size

	// if newFree > MEMORY_SIZE {
	if result+size > MEMORY_SIZE {
		// call GC
		MarkAndCompact()
		result = PROGRAM.HeapStartsAt + PROGRAM.HeapPointer
		newFree = PROGRAM.HeapPointer + size

		freeMemPerc := 1.0 - float32(newFree)/float32(MEMORY_SIZE-PROGRAM.HeapStartsAt)

		if freeMemPerc < float32(MIN_HEAP_FREE_RATIO)/100.0 {
			// then we have less than MIN_HEAP_FREE_RATIO memory left. expand!
			ResizeMemory(int(float32(MIN_HEAP_FREE_RATIO*(MEMORY_SIZE-PROGRAM.HeapStartsAt))/freeMemPerc), true)
		}

		if freeMemPerc > float32(MAX_HEAP_FREE_RATIO)/100.0 {
			// then we have more than MAX_HEAP_FREE_RATIO memory left. shrink!
			ResizeMemory(int(float32(MAX_HEAP_FREE_RATIO*(MEMORY_SIZE-PROGRAM.HeapStartsAt))/freeMemPerc), false)
		}
	}

	PROGRAM.HeapPointer = newFree

	return result
}

// WriteMemory ...
func WriteMemory(offset int, byts []byte) {
	for c := 0; c < len(byts); c++ {
		PROGRAM.Memory[offset+c] = byts[c]
	}
}

// Utilities

// FromBool ...
func FromBool(in bool) []byte {
	if in {
		return []byte{1}
	}
	return []byte{0}

}

// FromByte ...
func FromByte(in byte) []byte {
	return encoder.SerializeAtomic(in)
}

// FromStr ...
func FromStr(in string) []byte {
	return encoder.Serialize(in)
}

// FromI8 ...
func FromI8(in int8) []byte {
	return encoder.SerializeAtomic(in)
}

// FromI32 ...
func FromI32(in int32) []byte {
	return encoder.SerializeAtomic(in)
}

// FromUI32 ...
func FromUI32(in uint32) []byte {
	return encoder.SerializeAtomic(in)
}

// FromI64 ...
func FromI64(in int64) []byte {
	return encoder.Serialize(in)
}

// FromF32 ...
func FromF32(in float32) []byte {
	return encoder.Serialize(in)
}

// FromF64 ...
func FromF64(in float64) []byte {
	return encoder.Serialize(in)
}

// func ReadArray(mem []byte, fp int, inp *CXArgument, indexes []int32) (int, int) {
// 	var offset int
// 	var size int = inp.Size
// 	for i, idx := range indexes {
// 		offset += int(idx) * inp.Lengths[i]
// 	}
// 	for _, len := range indexes {
// 		size *= int(len)
// 	}

// 	return offset, size
// }

// ReadF32Data ...
func ReadF32Data(fp int, inp *CXArgument) interface{} {
	var data interface{}
	elt := GetAssignmentElement(inp)
	var dataF32 []float32
	if elt.IsSlice {
		dataF32 = ReadF32Slice(fp, inp)
	} else if elt.IsArray {
		dataF32 = ReadF32A(fp, inp)
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
	if len(dataF32) > 0 {
		data = dataF32
	}
	return data
}

// ReadF32Slice ...
func ReadF32Slice(fp int, inp *CXArgument) (out []float32) {
	sliceOffset := GetSliceOffset(fp, inp)
	if sliceOffset >= 0 && inp.Type == TYPE_F32 {
		slice := GetSlice(sliceOffset, GetAssignmentElement(inp).TotalSize)
		if slice != nil {
			err := encoder.DeserializeRaw(slice, &out)
			if err != nil {
				panic(err)
			}
		}
	} else {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
	return
}

// ReadF32A ...
func ReadF32A(fp int, inp *CXArgument) (out []float32) {
	offset := GetFinalOffset(fp, inp)
	byts := ReadMemory(offset, inp)
	byts = append(encoder.SerializeAtomic(int32(len(byts)/4)), byts...)
	err := encoder.DeserializeRaw(byts, &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(fp, inp)
	err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadByte ...
func ReadByte(fp int, inp *CXArgument) (out byte) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeAtomic(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	var offset int32
	off := GetFinalOffset(fp, inp)
	if inp.Name == "" {
		// then it's a literal
		offset = int32(off)
	} else {
		_, err := encoder.DeserializeAtomic(PROGRAM.Memory[off:off+TYPE_POINTER_SIZE], &offset)
		if err != nil {
			panic(err)
		}
	}

	if offset == 0 {
		// then it's nil string
		out = ""
		return
	}

	var size int32
	sizeB := PROGRAM.Memory[offset : offset+STR_HEADER_SIZE]

	_, err := encoder.DeserializeAtomic(sizeB, &size)
	if err != nil {
		panic(err)
	}
	err = encoder.DeserializeRaw(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE+size], &out)
	if err != nil {
		panic(err)
	}

	return out
}

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) (out int8) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeAtomic(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) (out int32) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeAtomic(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) (out int64) {
	offset := GetFinalOffset(fp, inp)
	err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) (out float32) {
	offset := GetFinalOffset(fp, inp)
	err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) (out float64) {
	offset := GetFinalOffset(fp, inp)
	err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}
