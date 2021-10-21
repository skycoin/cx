package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// MarkAndCompact ...
func MarkAndCompact(prgrm *CXProgram) {
	var fp types.Pointer
	var faddr = constants.NULL_HEAP_ADDRESS_OFFSET

	// marking, setting forward addresses and updating references
	// global variables
	for _, pkg := range prgrm.Packages {
		for _, glbl := range pkg.Globals {
			if (glbl.IsPointer || glbl.IsSlice || glbl.Type == types.STR) && glbl.CustomType == nil {
				// Getting the offset to the object in the heap
				heapOffset := types.Read_ptr(prgrm.Memory, glbl.Offset)
				if heapOffset < prgrm.Heap.StartsAt {
					continue
				}
				MarkObjectsTree(prgrm, glbl.Offset, glbl.Type, glbl.DeclarationSpecifiers[1:])
			}

			// If `ptr` has fields, we need to navigate the heap and mark its fields too.
			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					offset := glbl.Offset + fld.Offset
					// Getting the offset to the object in the heap
					heapOffset := types.Read_ptr(prgrm.Memory, offset)
					if heapOffset < prgrm.Heap.StartsAt {
						continue
					}

					if fld.IsPointer || fld.IsSlice || fld.Type == types.STR {
						MarkObjectsTree(prgrm, offset, fld.Type, fld.DeclarationSpecifiers[1:])
					}
				}
			}
		}
	}

	// marking, setting forward addresses and updating references
	// local variables
	for c := types.Pointer(0); c <= prgrm.CallCounter; c++ {
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

			ptrIsPointer := IsPointer(ptr)

			// Checking if we need to mark `ptr`.
			if ptrIsPointer {
				// If `ptr` has fields, we need to navigate the heap and mark its fields too.
				if ptr.CustomType != nil {
					// Getting the offset to the object in the heap
					heapOffset := types.Read_ptr(prgrm.Memory, offset)
					if heapOffset >= prgrm.Heap.StartsAt {
						for _, fld := range ptr.CustomType.Fields {
							MarkObjectsTree(prgrm, heapOffset+types.OBJECT_HEADER_SIZE+fld.Offset, fld.Type, fld.DeclarationSpecifiers[1:])
						}
					}
				}

				MarkObjectsTree(prgrm, offset, ptr.Type, ptr.DeclarationSpecifiers[1:])
			}

			// Checking if the field being accessed needs to be marked.
			// If the root (`ptr`) is a pointer, this step is unnecessary.
			if len(ptr.Fields) > 0 && !ptrIsPointer && IsPointer(ptr.Fields[len(ptr.Fields)-1]) {
				fld := ptr.Fields[len(ptr.Fields)-1]
				MarkObjectsTree(prgrm, offset+fld.Offset, fld.Type, fld.DeclarationSpecifiers[1:])
			}
		}

		fp += op.Size
	}

	// Relocation of live objects.
	for c := prgrm.Heap.StartsAt + constants.NULL_HEAP_ADDRESS_OFFSET; c < prgrm.Heap.StartsAt+prgrm.Heap.Pointer; {
		objSize := types.Read_obj_size(prgrm.Memory, c)

		if prgrm.Memory[c] == 1 {
			forwardingAddress := types.Read_obj_forwarding_address(prgrm.Memory, c)

			// We update the pointers that are pointing to the just moved object.
			updatePointers(prgrm, forwardingAddress, prgrm.Heap.StartsAt+faddr)

			// setting the mark back to 0
			prgrm.Memory[c] = 0
			// then it's alive and we'll relocate the object
			for i := types.Pointer(0); i < objSize; i++ {
				prgrm.Memory[faddr+prgrm.Heap.StartsAt+i] = prgrm.Memory[c+i]
			}

			faddr += objSize
		}

		c += objSize
	}

	prgrm.Heap.Pointer = faddr
}

// updateDisplaceReference performs the actual addition or subtraction of `plusOff` to the address being pointed by the element at `atOffset`.
func updateDisplaceReference(prgrm *CXProgram, updated *map[types.Pointer]types.Pointer, atOffset types.Pointer, plusOff types.Pointer) {
	// Checking if it was already updated.
	if _, found := (*updated)[atOffset]; found {
		return
	}

	// Extracting the address being pointed by element at `atOffset`
	dsCurrAddr := types.Read_ptr(prgrm.Memory, atOffset)

	// Adding `plusOff` to the address and updating the address pointed by
	// element at `atOffset`.
	types.Write_ptr(prgrm.Memory, atOffset, dsCurrAddr+plusOff)

	// Keeping a record of this address. We don't want to displace the object twice.
	// We're using a map to speed things up a tiny bit.
	(*updated)[atOffset] = atOffset
}

// doDisplaceReferences checks if the element at `atOffset` is pointing to an object on the heap and, if this is the case, it displaces it by `plusOff`. `updated` keeps a record of all the offsets that have already been updated.
func doDisplaceReferences(prgrm *CXProgram, updated *map[types.Pointer]types.Pointer, atOffset types.Pointer, plusOff types.Pointer, baseType types.Code, declSpecs []int) {
	var numDeclSpecs = len(declSpecs)

	// Getting the offset to the object in the heap.
	heapOffset := types.Read_ptr(prgrm.Memory, atOffset)

	// The whole displacement process is needed because the objects on the heap were
	// displaced by additional data segment bytes. These additional bytes need to be
	// considered if we want to read the objects on the heap. We need to check if the
	// displacement is positive or negative; if it is positive then this means a data
	// segment was added; if it is negative it means that we're done with any CX chains
	// process involving the addition of a data segment, and in some code snippets we can ignore the displacement.
	// TODO: Maybe this condition can be avoided by refactoring the code?
	var condPlusOff types.Pointer
	if plusOff > 0 {
		condPlusOff = plusOff
	}

	// Displace the address pointed by element at `atOffset`.
	updateDisplaceReference(prgrm, updated, atOffset, plusOff)

	// It can't be a tree of objects.
	if numDeclSpecs == 0 || heapOffset <= prgrm.Heap.StartsAt+condPlusOff {
		return
	}

	// Checking if it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == constants.DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == constants.DECL_SLICE ||
				declSpecs[1] == constants.DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == types.STR) {
			// Then we need to iterate each of the slice objects
			// and check if we need to update their address.
			sliceLen := GetSliceLen(heapOffset + condPlusOff)

			offsetToElements := types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE

			for c := types.Pointer(0); c < sliceLen; c++ {
				cHeapOffset := types.Read_ptr(prgrm.Memory, heapOffset+condPlusOff+offsetToElements+c*types.POINTER_SIZE)
				if cHeapOffset <= prgrm.Heap.StartsAt+condPlusOff {
					// Then it's pointing to null or data segment
					continue
				}

				// Displacing this child element.
				updateDisplaceReference(prgrm, updated, heapOffset+offsetToElements+c*types.POINTER_SIZE, plusOff)
			}
		}
	}
}

// DisplaceReferences displaces all the pointer-like variables, slice elements or field structures by `off`. `numPkgs` tells us the number of packages to consider for the reference desplacement (this number should equal to the number of packages that represent the blockchain code in a CX chain).
func DisplaceReferences(prgrm *CXProgram, off types.Pointer, numPkgs int) {
	// We're going to keep a record of all the references that were already updated.
	updated := make(map[types.Pointer]types.Pointer)

	count := 0
	for _, pkg := range prgrm.Packages {
		if count > numPkgs {
			break
		}

		// In a CX chain we're only interested on considering global variables,
		// as any other object should be destroyed, as the program finished its
		// execution.
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				doDisplaceReferences(prgrm, &updated, glbl.Offset, off, glbl.Type, glbl.DeclarationSpecifiers[1:]) // TODO:PTR remove hardcoded offsets
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

		count++
	}
}

// Mark marks the object located at `heapOffset` as alive and sets the object's referencing address to `heapOffset`.
func Mark(prgrm *CXProgram, heapOffset types.Pointer) {
	// Marking as alive.
	types.Write_obj_mark(prgrm.Memory, heapOffset, 1)

	// Setting forwarding address. This address is used to know where the object used to live on the heap. With it we can know what symbols were pointing to that dead object and then update their address.
	types.Write_obj_forwarding_address(prgrm.Memory, heapOffset, heapOffset)
}

// MarkObjectsTree traverses and marks a possible tree of heap objects (slices of slices, slices of pointers, etc.).
func MarkObjectsTree(prgrm *CXProgram, offset types.Pointer, baseType types.Code, declSpecs []int) {
	lenMem := types.Cast_int_to_ptr(len(prgrm.Memory))
	// Checking if it's a valid heap address. An invalid address
	// usually occurs in CX chains, with the split of blockchain
	// and transaction codes in a CX chain program state.
	if offset > lenMem || offset+types.POINTER_SIZE > lenMem {
		return
	}

	// Getting the offset to the object in the heap
	heapOffset := types.Read_ptr(prgrm.Memory, offset)

	// Then it's a pointer to an object in the stack and it should not be marked.
	if heapOffset <= prgrm.Heap.StartsAt {
		return
	}

	// marking the root object
	Mark(prgrm, heapOffset)

	var numDeclSpecs = len(declSpecs)
	if numDeclSpecs == 0 {
		return
	}

	// Then it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == constants.DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == constants.DECL_SLICE ||
				declSpecs[1] == constants.DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == types.STR) {
			// Then we need to iterate each of the slice objects and mark them as alive
			sliceLen := GetSliceLen(heapOffset)

			for c := types.Pointer(0); c < sliceLen; c++ {
				offsetToElements := types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE
				cHeapOffset := types.Read_ptr(prgrm.Memory, heapOffset+offsetToElements+c*types.POINTER_SIZE)
				if cHeapOffset <= prgrm.Heap.StartsAt {
					// Then it's pointing to null or data segment
					continue
				}

				MarkObjectsTree(prgrm, heapOffset+offsetToElements+c*types.POINTER_SIZE, baseType, declSpecs[1:])
			}
		}
	}
}

// updatePointer changes the address of the pointer located at `atOffset` to `newAddress`.
func updatePointer(prgrm *CXProgram, atOffset types.Pointer, toAddress types.Pointer) {
	types.Write_ptr(prgrm.Memory, atOffset, toAddress)
}

// updatePointerTree changes the address of the pointer located at `atOffset` to `newAddress` and checks if it is the
// root of a tree of objects, such as a slice or the instance of a struct where some of its fields are pointers.
func updatePointerTree(prgrm *CXProgram, atOffset types.Pointer, oldAddr, newAddr types.Pointer, baseType types.Code, declSpecs []int) {
	var numDeclSpecs = len(declSpecs)

	// Getting the offset to the object in the heap
	heapOffset := types.Read_ptr(prgrm.Memory, atOffset)

	if heapOffset == oldAddr {
		// Updating the root pointer.
		updatePointer(prgrm, atOffset, newAddr)
	}

	// It can't be a tree of objects.
	if numDeclSpecs == 0 || heapOffset <= prgrm.Heap.StartsAt {
		return
	}

	// Checking if it's a tree of objects.
	// TODO: We're not considering struct instances with pointer fields.
	if declSpecs[0] == constants.DECL_SLICE {
		if (numDeclSpecs > 1 &&
			(declSpecs[1] == constants.DECL_SLICE ||
				declSpecs[1] == constants.DECL_POINTER)) ||
			(numDeclSpecs == 1 && baseType == types.STR) {
			// Then we need to iterate each of the slice objects
			// and check if we need to update their address.
			sliceLen := GetSliceLen(heapOffset)
			offsetToElements := types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE

			for c := types.Pointer(0); c < sliceLen; c++ {
				cHeapOffset := types.Read_ptr(prgrm.Memory, heapOffset+offsetToElements+c*types.POINTER_SIZE)
				if cHeapOffset <= prgrm.Heap.StartsAt {
					// Then it's pointing to null or data segment
					continue
				}

				// Then it's not pointing to the object moved by the GC or it's pointing to
				// an object in the stack segment or nil.
				if cHeapOffset == oldAddr {
					updatePointerTree(prgrm, heapOffset+offsetToElements+c*types.POINTER_SIZE, oldAddr, newAddr, baseType, declSpecs[1:]) // TODO:PTR remove hardcode
				}
			}
		}
	}
}

// updatePointers updates all the references to objects on the heap to their new addresses after calling the garbage collector.
// For example, if `foo` was pointing to an object located at address 5151 and after calling the garbage collector it was
// moved to address 4141, every symbol in a `CXProgram`'s `CallStack` and in its global variables need to be updated to point now to
// 4141 instead of 5151.
func updatePointers(prgrm *CXProgram, oldAddr, newAddr types.Pointer) {
	if oldAddr == newAddr {
		return
	}
	// TODO: `oldAddr` could be received as a slice of bytes that represent the old address of the object,
	// as it needs to be converted to bytes later on anyways. However, I'm sticking to an int32
	// for a bit more of clarity.
	for _, pkg := range prgrm.Packages {
		for _, glbl := range pkg.Globals {
			if (glbl.IsPointer || glbl.IsSlice || glbl.Type == types.STR) && glbl.CustomType == nil {
				// Getting the offset to the object in the heap
				heapOffset := types.Read_ptr(prgrm.Memory, glbl.Offset)
				if heapOffset < prgrm.Heap.StartsAt {
					continue
				}

				updatePointerTree(prgrm, glbl.Offset, oldAddr, newAddr, glbl.Type, glbl.DeclarationSpecifiers[1:]) // TODO:PTR remove hardcoded 1 value
			}

			// If `ptr` has fields, we need to navigate the heap and mark its fields too.
			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					if !IsPointer(fld) {
						continue
					}
					offset := glbl.Offset + fld.Offset

					// Getting the offset to the object in the heap
					heapOffset := types.Read_ptr(prgrm.Memory, offset)
					if heapOffset < prgrm.Heap.StartsAt {
						continue
					}

					if fld.IsPointer || fld.IsSlice || fld.Type == types.STR {
						updatePointerTree(prgrm, offset, oldAddr, newAddr, fld.Type, fld.DeclarationSpecifiers[1:])
					}
				}
			}
		}
	}

	var fp types.Pointer
	for c := types.Pointer(0); c <= prgrm.CallCounter; c++ {
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

			ptrIsPointer := IsPointer(ptr)

			// Checking if we need to mark `ptr`.
			if ptrIsPointer {
				// Getting the offset to the object in the heap
				heapOffset := types.Read_ptr(prgrm.Memory, offset)
				if heapOffset > prgrm.Heap.StartsAt {
					updatePointerTree(prgrm, offset, oldAddr, newAddr, ptr.Type, ptr.DeclarationSpecifiers[1:])

					// If `ptr` has fields, we need to navigate the heap and mark its fields too.
					if ptr.CustomType != nil {
						if heapOffset >= prgrm.Heap.StartsAt {
							for _, fld := range ptr.CustomType.Fields {
								updatePointerTree(prgrm, heapOffset+types.OBJECT_HEADER_SIZE+fld.Offset, oldAddr, newAddr, fld.Type, fld.DeclarationSpecifiers[1:])
							}
						}
					}
				}
			}

			// Checking if the field being accessed needs to be marked.
			// If the root (`ptr`) is a pointer, this step is unnecessary.
			if len(ptr.Fields) > 0 && !ptrIsPointer && IsPointer(ptr.Fields[len(ptr.Fields)-1]) {
				fld := ptr.Fields[len(ptr.Fields)-1]

				// Getting the offset to the object in the heap
				heapOffset := types.Read_ptr(prgrm.Memory, offset+fld.Offset)
				if heapOffset > prgrm.Heap.StartsAt {
					updatePointerTree(prgrm, offset+fld.Offset, oldAddr, newAddr, fld.Type, fld.DeclarationSpecifiers[1:])
				}
			}

		}

		fp += op.Size
	}
}
