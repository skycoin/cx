package cxcore

import (
	//"fmt"
	//"math"

	//"github.com/skycoin/skycoin/src/cipher/encoder"
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
//Todo: This function needs comments? What does it do?
//Todo: Can this function be specialized?
//CalculateDeference
// ->
//CalculateDeferenceSlice
//CalculateDeferenceArray
//CalculateDeferencePointer
//CalculateDeferenceInt32, etc (FIXED)
//TODO: Why are we calling this function for fixed data types in flow path
//TODO: For int32, f32, etc, this function should not be called at all
//reduce loops and switches in op code execution flow path
func CalculateDereferences(arg *CXArgument, finalOffset *int, fp int) {
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

			offset = mustDeserializeI32(byts)

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

			offset = mustDeserializeI32(byts)
			*finalOffset = int(offset)
		}

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
		offset := mustDeserializeI32(PROGRAM.Memory[strOffset : strOffset+TYPE_POINTER_SIZE])
		strOffset = int(offset)
	}
	return strOffset
}

// GetFinalOffset ...
// Get byte position to write operand result back to?
// Note: CXArgument is bloated, should pass in required though something else

//TODO:
//CalculateDeference
// ->
//CalculateDeferenceInt8,
//CalculateDeferenceInt16,
//CalculateDeferenceInt32, etc
//CalculateDeferenceInt64, etc (FIXED)
//CalculateDeferenceSlice
//CalculateDeferenceArray
//CalculateDeferencePointer
//Then use them for all fixed types
//InlineCalculateDeferences

func GetFinalOffset(fp int, arg *CXArgument) int {
	// defer RuntimeError(PROGRAM)
	// var elt *CXArgument
	finalOffset := arg.Offset
	// var fldIdx int

	// elt = arg

	//Todo: find way to eliminate this check
	if finalOffset < PROGRAM.StackSize {
		// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
		finalOffset += fp
	}

	// elt = arg
	//TODO: Eliminate all op codes with more than one return type
	//TODO: Eliminate this loop
	CalculateDereferences(arg, &finalOffset, fp)
	for _, fld := range arg.Fields {
		// elt = fld
		finalOffset += fld.Offset
		CalculateDereferences(fld, &finalOffset, fp)
	}

	return finalOffset
}

// ReadMemory ...
//TODO: Avoid all read memory commands for fixed width types (i32,f32,etc)
//TODO: Make "ReadMemoryI32", "ReadMemoryI16", etc
func ReadMemory(offset int, arg *CXArgument) []byte {
	size := GetSize(arg)
	return PROGRAM.Memory[offset : offset+size]
}

// ResizeMemory ...
func ResizeMemory(prgrm *CXProgram, newMemSize int, isExpand bool) {
	// We can't expand memory to a value greater than `memLimit`.
	if newMemSize > MAX_HEAP_SIZE {
		newMemSize = MAX_HEAP_SIZE
	}

	if newMemSize == prgrm.HeapSize {
		// Then we're at the limit; we can't expand anymore.
		// We can only hope that the free memory is enough for the CX program to continue running.
		return
	}

	if isExpand {
		// Adding bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append(prgrm.Memory, make([]byte, newMemSize-prgrm.HeapSize)...)
		prgrm.HeapSize = newMemSize
	} else {
		// Removing bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append([]byte(nil), prgrm.Memory[:prgrm.HeapStartsAt+newMemSize]...)
		prgrm.HeapSize = newMemSize
	}
}

// AllocateSeq allocates memory in the heap
func AllocateSeq(size int) (offset int) {
	// Current object trying to be allocated would use this address.
	addr := PROGRAM.HeapPointer
	// Next object to be allocated will use this address.
	newFree := addr + size

	// Checking if we can allocate the entirety of the object in the current heap.
	if newFree > PROGRAM.HeapSize {
		// It does not fit, so calling garbage collector.
		MarkAndCompact(PROGRAM)
		// Heap pointer got moved by GC and recalculate these variables based on the new pointer.
		addr = PROGRAM.HeapPointer
		newFree = addr + size

		// If the new heap pointer exceeds `MAX_HEAP_SIZE`, there's nothing left to do.
		if newFree > MAX_HEAP_SIZE {
			panic(HEAP_EXHAUSTED_ERROR)
		}

		// According to MIN_HEAP_FREE_RATIO and MAX_HEAP_FREE_RATION we can either shrink
		// or expand the heap to maintain "healthy" heap sizes. The idea is that we don't want
		// to have an absurdly amount of free heap memory, as we would be wasting resources, and we
		// don't want to have a small amount of heap memory left as we'd be calling the garbage collector
		// too frequently.

		// Calculating free heap memory percentage.
		usedPerc := float32(newFree) / float32(PROGRAM.HeapSize)
		freeMemPerc := 1.0 - usedPerc

		// Then we have less than MIN_HEAP_FREE_RATIO memory left. Expand!
		if freeMemPerc < MIN_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MIN_HEAP_FREE_RATIO.
			newMemSize := int(float32(newFree) / (1.0 - MIN_HEAP_FREE_RATIO))
			ResizeMemory(PROGRAM, newMemSize, true)
		}

		// Then we have more than MAX_HEAP_FREE_RATIO memory left. Shrink!
		if freeMemPerc > MAX_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MAX_HEAP_FREE_RATIO.
			newMemSize := int(float32(newFree) / (1.0 - MAX_HEAP_FREE_RATIO))

			// This check guarantees that the CX program has always at least INIT_HEAP_SIZE bytes to work with.
			// A flag could be added later to remove this, as in some cases this mechanism could not be desired.
			if newMemSize > INIT_HEAP_SIZE {
				ResizeMemory(PROGRAM, newMemSize, false)
			}
		}
	}

	PROGRAM.HeapPointer = newFree

	// Returning absolute memory address (not relative to where heap starts at).
	// Above this point we were performing all operations taking into
	// consideration only heap offsets.
	return addr + PROGRAM.HeapStartsAt
}


