package cxcore

import (
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

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

			var sizeToUse int
			if arg.CustomType != nil {
				sizeToUse = arg.CustomType.Size
			} else if arg.IsSlice {
				sizeToUse = arg.TotalSize
			} else {
				sizeToUse = arg.Size
			}

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

			var sizeToUse int
			if arg.CustomType != nil {
				sizeToUse = arg.CustomType.Size
			} else if arg.IsSlice {
				sizeToUse = arg.TotalSize
			} else {
				sizeToUse = arg.Size
			}

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
		dbg = true
	}

	if finalOffset < PROGRAM.StackSize {
		// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
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
	return PROGRAM.Memory[offset : offset+arg.TotalSize]
}

// updateDisplaceReference performs the actual addition or subtraction of `plusOff` to the address being pointed by the element at `atOffset`.
func updateDisplaceReference(prgrm *CXProgram, updated *map[int]int, atOffset, plusOff int) {
	// Checking if it was already updated.
	if _, found := (*updated)[atOffset]; found {
		return
	}

	// Extracting the address being pointed by element at `atOffset`
	sCurrAddr := prgrm.Memory[atOffset : atOffset+TYPE_POINTER_SIZE]
	var dsCurrAddr int32
	_, err := encoder.DeserializeAtomic(sCurrAddr, &dsCurrAddr)
	if err != nil {
		panic(err)
	}

	// Adding `plusOff` to the address and updating the address pointed by
	// element at `atOffset`.
	addrB := encoder.SerializeAtomic(int32(int(dsCurrAddr) + plusOff))
	for i := 0; i < TYPE_POINTER_SIZE; i++ {
		prgrm.Memory[atOffset+i] = addrB[i]
	}

	// Keeping a record of this address. We don't want to displace the object twice.
	// We're using a map to speed things up a tiny bit.
	(*updated)[atOffset] = atOffset
}

// doDisplaceReferences checks if the element at `atOffset` is pointing to an object on the heap and, if this is the case, it displaces it by `plusOff`. `updated` keeps a record of all the offsets that have already been updated.
func doDisplaceReferences(prgrm *CXProgram, updated *map[int]int, atOffset int, plusOff int, baseType int, declSpecs []int) {
	var numDeclSpecs = len(declSpecs)

	// Getting the offset to the object in the heap.
	var heapOffset int32
	_, err := encoder.DeserializeAtomic(prgrm.Memory[atOffset:atOffset+TYPE_POINTER_SIZE], &heapOffset)
	if err != nil {
		panic(err)
	}

	// The whole displacement process is needed because the objects on the heap were
	// displaced by additional data segment bytes. These additional bytes need to be
	// considered if we want to read the objects on the heap. We need to check if the
	// displacement is positive or negative; if it is positive then this means a data
	// segment was added; if it is negative it means that we're done with any CX chains
	// process involving the addition of a data segment, and in some code snippets we can ignore the displacement.
	// TODO: Maybe this condition can be avoided by refactoring the code?
	var condPlusOff int
	if plusOff > 0 {
		condPlusOff = plusOff
	}

	// Displace the address pointed by element at `atOffset`.
	updateDisplaceReference(prgrm, updated, atOffset, plusOff)

	// It can't be a tree of objects.
	if numDeclSpecs == 0 || int(heapOffset) <= prgrm.HeapStartsAt+condPlusOff {
		return
	}

	// Checking if it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == DECL_SLICE ||
				declSpecs[1] == DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == TYPE_STR) {
			// Then we need to iterate each of the slice objects
			// and check if we need to update their address.
			var sliceLen int32
			_, err := encoder.DeserializeAtomic(GetSliceHeader(heapOffset + int32(condPlusOff))[4:8], &sliceLen)
			if err != nil {
				panic(err)
			}

			offsetToElements := OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE

			for c := int32(0); c < sliceLen; c++ {
				var cHeapOffset int32
				_, err := encoder.DeserializeAtomic(prgrm.Memory[int(heapOffset)+condPlusOff+offsetToElements+int(c*TYPE_POINTER_SIZE):int(heapOffset)+condPlusOff+offsetToElements+int(c*TYPE_POINTER_SIZE)+4], &cHeapOffset)
				if err != nil {
					panic(err)
				}

				if int(cHeapOffset) <= prgrm.HeapStartsAt+condPlusOff {
					// Then it's pointing to null or data segment
					continue
				}

				// Displacing this child element.
				updateDisplaceReference(prgrm, updated, int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE), plusOff)
			}
		}
	}
}

// DisplaceReferences displaces all the pointer-like variables, slice elements or field structures by `off`. `numPkgs` tells us the number of packages to consider for the reference desplacement (this number should equal to the number of packages that represent the blockchain code in a CX chain).
func DisplaceReferences(prgrm *CXProgram, off int, numPkgs int) {
	// We're going to keep a record of all the references that were already updated.
	updated := make(map[int]int)

	for c := 0; c < numPkgs; c++ {
		pkg := prgrm.Packages[c]

		// In a CX chain we're only interested on considering global variables,
		// as any other object should be destroyed, as the program finished its
		// execution.
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				doDisplaceReferences(prgrm, &updated, glbl.Offset, off, glbl.Type, glbl.DeclarationSpecifiers[1:])
			}

			// If it's a struct instance we need to displace each of its fields.
			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					if fld.IsPointer || fld.IsSlice {
						doDisplaceReferences(prgrm, &updated, glbl.Offset+fld.Offset, off, fld.Type, fld.DeclarationSpecifiers[1:])
					}
				}
			}
		}
	}
}

// Mark marks the object located at `heapOffset` as alive and sets the object's referencing address to `heapOffset`.
func Mark(prgrm *CXProgram, heapOffset int32) {
	// Marking as alive.
	prgrm.Memory[heapOffset] = 1

	for i, byt := range encoder.SerializeAtomic(heapOffset) {
		// Setting forwarding address. This address is used to know where the object used to live on the heap. With it we can know what symbols were pointing to that dead object and then update their address.
		prgrm.Memory[int(heapOffset)+MARK_SIZE+i] = byt
	}
}

// MarkObjectsTree traverses and marks a possible tree of heap objects (slices of slices, slices of pointers, etc.).
func MarkObjectsTree(prgrm *CXProgram, offset int, baseType int, declSpecs []int) {
	lenMem := len(prgrm.Memory)
	// Checking if it's a valid heap address. An invalid address
	// usually occurs in CX chains, with the split of blockchain
	// and transaction codes in a CX chain program state.
	if offset > lenMem || offset+TYPE_POINTER_SIZE > lenMem {
		return
	}

	var numDeclSpecs = len(declSpecs)

	// Getting the offset to the object in the heap
	var heapOffset int32
	_, err := encoder.DeserializeAtomic(prgrm.Memory[offset:offset+TYPE_POINTER_SIZE], &heapOffset)
	if err != nil {
		panic(err)
	}

	// Then it's a pointer to an object in the stack and it should not be marked.
	if heapOffset <= int32(prgrm.HeapStartsAt) {
		return
	}

	// marking the root object
	Mark(prgrm, heapOffset)

	if numDeclSpecs == 0 {
		return
	}

	// Then it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == DECL_SLICE ||
				declSpecs[1] == DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == TYPE_STR) {
			// Then we need to iterate each of the slice objects and mark them as alive
			var sliceLen int32
			_, err := encoder.DeserializeAtomic(GetSliceHeader(heapOffset)[4:8], &sliceLen)
			if err != nil {
				panic(err)
			}

			for c := int32(0); c < sliceLen; c++ {
				offsetToElements := OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE
				var cHeapOffset int32
				_, err := encoder.DeserializeAtomic(prgrm.Memory[int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE):int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE)+4], &cHeapOffset)
				if err != nil {
					panic(err)
				}

				if cHeapOffset <= int32(prgrm.HeapStartsAt) {
					// Then it's pointing to null or data segment
					continue
				}

				MarkObjectsTree(prgrm, int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE), baseType, declSpecs[1:])
			}
		}
	}
}

// updatePointer changes the address of the pointer located at `atOffset` to `newAddress`.
func updatePointer(prgrm *CXProgram, atOffset int, toAddress int32) {
	addrB := encoder.SerializeAtomic(toAddress)
	for i := 0; i < TYPE_POINTER_SIZE; i++ {
		prgrm.Memory[atOffset+i] = addrB[i]
	}
}

// updatePointerTree changes the address of the pointer located at `atOffset` to `newAddress` and checks if it is the
// root of a tree of objects, such as a slice or the instance of a struct where some of its fields are pointers.
func updatePointerTree(prgrm *CXProgram, atOffset int, oldAddr, newAddr int32, baseType int, declSpecs []int) {
	var numDeclSpecs = len(declSpecs)

	// Getting the offset to the object in the heap
	var heapOffset int32
	_, err := encoder.DeserializeAtomic(prgrm.Memory[atOffset:atOffset+TYPE_POINTER_SIZE], &heapOffset)
	if err != nil {
		panic(err)
	}

	if heapOffset == oldAddr {
		// Updating the root pointer.
		updatePointer(prgrm, atOffset, newAddr)
	}

	// It can't be a tree of objects.
	if numDeclSpecs == 0 || int(heapOffset) <= prgrm.HeapStartsAt {
		return
	}

	// Checking if it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == DECL_SLICE ||
				declSpecs[1] == DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == TYPE_STR) {
			// Then we need to iterate each of the slice objects
			// and check if we need to update their address.
			var sliceLen int32
			_, err := encoder.DeserializeAtomic(GetSliceHeader(heapOffset)[4:8], &sliceLen)
			if err != nil {
				panic(err)
			}

			offsetToElements := OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE

			for c := int32(0); c < sliceLen; c++ {
				var cHeapOffset int32
				_, err := encoder.DeserializeAtomic(prgrm.Memory[int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE):int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE)+4], &cHeapOffset)
				if err != nil {
					panic(err)
				}

				// Then it's not pointing to the object moved by the GC or it's pointing to
				// an object in the stack segment or nil.
				if cHeapOffset == oldAddr {
					updatePointerTree(prgrm, int(heapOffset)+offsetToElements+int(c*TYPE_POINTER_SIZE), oldAddr, newAddr, baseType, declSpecs[1:])
				}
			}
		}
	}
}

// updatePointers updates all the references to objects on the heap to their new addresses after calling the garbage collector.
// For example, if `foo` was pointing to an object located at address 5151 and after calling the garbage collector it was
// moved to address 4141, every symbol in a `CXProgram`'s `CallStack` and in its global variables need to be updated to point now to
// 4141 instead of 5151.
func updatePointers(prgrm *CXProgram, oldAddr, newAddr int32) {
	// TODO: `oldAddr` could be received as a slice of bytes that represent the old address of the object,
	// as it needs to be converted to bytes later on anyways. However, I'm sticking to an int32
	// for a bit more of clarity.
	for _, pkg := range prgrm.Packages {
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				updatePointerTree(prgrm, glbl.Offset, oldAddr, newAddr, glbl.Type, glbl.DeclarationSpecifiers[1:])
			}

			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					if fld.IsPointer || fld.IsSlice {
						updatePointerTree(prgrm, glbl.Offset+fld.Offset, oldAddr, newAddr, fld.Type, fld.DeclarationSpecifiers[1:])
					}
				}
			}
		}
	}

	var fp int

	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		// TODO: Some standard library functions "manually" add a function
		// call (callbacks) to `PRGRM.CallStack`. These functions do not have an
		// operator associated to them. This can be considered as a bug or as an
		// undesirable mechanic.
		// [2019-06-24 Mon 22:39] Actually, if the GC is triggered in the middle
		// of a callback, things will certainly break.
		if op == nil {
			continue
		}

		for _, ptr := range op.ListOfPointers {
			offset := ptr.Offset
			offset += fp

			// If we're accessing to a field of that pointer, we need to
			// take into consideration its offset.
			// TODO: I think this could be pre-computed at compile-time so we can just use `ptr.Offset`.
			if len(ptr.Fields) > 0 {
				fld := ptr.Fields[len(ptr.Fields)-1]
				updatePointerTree(prgrm, offset+fld.Offset, oldAddr, newAddr, fld.Type, fld.DeclarationSpecifiers[1:])
			} else {
				updatePointerTree(prgrm, offset, oldAddr, newAddr, ptr.Type, ptr.DeclarationSpecifiers[1:])
			}

		}

		fp += op.Size
	}
}

// MarkAndCompact ...
func MarkAndCompact(prgrm *CXProgram) {
	var fp int
	var faddr = int32(NULL_HEAP_ADDRESS_OFFSET)

	// marking, setting forward addresses and updating references
	// global variables
	for _, pkg := range prgrm.Packages {
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				MarkObjectsTree(prgrm, glbl.Offset, glbl.Type, glbl.DeclarationSpecifiers[1:])
			}
			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					if fld.IsPointer || fld.IsSlice {
						MarkObjectsTree(prgrm, glbl.Offset+fld.Offset, fld.Type, fld.DeclarationSpecifiers[1:])
					}
				}
			}
		}
	}

	// marking, setting forward addresses and updating references
	// local variables
	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		// TODO: Some standard library functions "manually" add a function
		// call (callbacks) to `PRGRM.CallStack`. These functions do not have an
		// operator associated to them. This can be considered as a bug or as an
		// undesirable mechanic.
		// [2019-06-24 Mon 22:39] Actually, if the GC is triggered in the middle
		// of a callback, things will certainly break.
		if op == nil {
			continue
		}

		for _, ptr := range op.ListOfPointers {
			offset := ptr.Offset
			offset += fp

			// If we're accessing to a field of that pointer, we need to
			// take into consideration its offset.
			// TODO: I think this could be pre-computed at compile-time so we can just use `ptr.Offset`.
			if len(ptr.Fields) > 0 {
				fld := ptr.Fields[len(ptr.Fields)-1]
				MarkObjectsTree(prgrm, offset+fld.Offset, fld.Type, fld.DeclarationSpecifiers[1:])
			} else {
				MarkObjectsTree(prgrm, offset, ptr.Type, ptr.DeclarationSpecifiers[1:])
			}
		}

		fp += op.Size
	}

	// relocation of live objects
	for c := prgrm.HeapStartsAt + NULL_HEAP_ADDRESS_OFFSET; c < prgrm.HeapStartsAt+prgrm.HeapPointer; {
		var objSize int32
		_, err := encoder.DeserializeAtomic(prgrm.Memory[c+MARK_SIZE+FORWARDING_ADDRESS_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE], &objSize)
		if err != nil {
			panic(err)
		}

		if prgrm.Memory[c] == 1 {
			var forwardingAddress int32
			_, err := encoder.DeserializeAtomic(prgrm.Memory[c+MARK_SIZE:c+MARK_SIZE+FORWARDING_ADDRESS_SIZE], &forwardingAddress)
			if err != nil {
				panic(err)
			}

			// We update the pointers that are pointing to the just moved object.
			updatePointers(prgrm, forwardingAddress, int32(prgrm.HeapStartsAt)+faddr)

			// setting the mark back to 0
			prgrm.Memory[c] = 0
			// then it's alive and we'll relocate the object
			for i := int32(0); i < objSize; i++ {
				prgrm.Memory[faddr+int32(prgrm.HeapStartsAt)+i] = prgrm.Memory[int32(c)+i]
			}

			faddr += objSize
		}

		c += int(objSize)
	}

	prgrm.HeapPointer = int(faddr)
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
			_, err := encoder.DeserializeRaw(slice, &out)
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
	_, err := encoder.DeserializeRaw(byts, &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) (out bool) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
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

// ReadStrOffset gets the string in heap pointed out by `offset`.
func ReadStrOffset(offset int32) (out string) {
	if offset == 0 {
		// Then it's nil string.
		out = ""
		return
	}

	var size int32
	var err error
	// We need to check if the string lives on the data segment or on the
	// heap to know if we need to take into consideration the object header's size.
	if int(offset) > PROGRAM.HeapStartsAt {
		_, err = encoder.DeserializeAtomic(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE:offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE], &size)

		if err != nil {
			panic(err)
		}

		_, err = encoder.DeserializeRaw(PROGRAM.Memory[offset+OBJECT_HEADER_SIZE:offset+OBJECT_HEADER_SIZE+STR_HEADER_SIZE+size], &out)
		if err != nil {
			panic(err)
		}
	} else {
		_, err = encoder.DeserializeAtomic(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE], &size)
		if err != nil {
			panic(err)
		}

		_, err = encoder.DeserializeRaw(PROGRAM.Memory[offset:offset+STR_HEADER_SIZE+size], &out)
		if err != nil {
			panic(err)
		}
	}

	return out
}

// ReadStr ...
func ReadStr(fp int, inp *CXArgument) (out string) {
	var offset int32
	off := GetFinalOffset(fp, inp)
	if inp.Name == "" {
		// Then it's a literal.
		offset = int32(off)
	} else {
		_, err := encoder.DeserializeAtomic(PROGRAM.Memory[off:off+TYPE_POINTER_SIZE], &offset)
		if err != nil {
			panic(err)
		}
	}

	return ReadStrOffset(offset)
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
	_, err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) (out float32) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) (out float64) {
	offset := GetFinalOffset(fp, inp)
	_, err := encoder.DeserializeRaw(ReadMemory(offset, inp), &out)
	if err != nil {
		panic(err)
	}
	return
}
