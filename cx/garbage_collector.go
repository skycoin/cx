package cxcore

import "github.com/skycoin/skycoin/src/cipher/encoder"

// MarkAndCompact ...
func MarkAndCompact(prgrm *CXProgram) {
	var fp int
	var faddr = int32(NULL_HEAP_ADDRESS_OFFSET)

	// marking, setting forward addresses and updating references
	// global variables
	for _, pkg := range prgrm.Packages {
		for _, glbl := range pkg.Globals {
			if (glbl.IsPointer || glbl.IsSlice || glbl.Type == TYPE_STR) && glbl.CustomType == nil {
				// Getting the offset to the object in the heap
				var heapOffset int32
				_, err := encoder.DeserializeAtomic(prgrm.Memory[glbl.Offset:glbl.Offset+TYPE_POINTER_SIZE], &heapOffset)
				if err != nil {
					panic(err)
				}

				if int(heapOffset) < prgrm.HeapStartsAt {
					continue
				}
				MarkObjectsTree(prgrm, glbl.Offset, glbl.Type, glbl.DeclarationSpecifiers[1:])
			}

			// If `ptr` has fields, we need to navigate the heap and mark its fields too.
			if glbl.CustomType != nil {
				for _, fld := range glbl.CustomType.Fields {
					offset := glbl.Offset + fld.Offset
					// Getting the offset to the object in the heap
					var heapOffset int32
					_, err := encoder.DeserializeAtomic(prgrm.Memory[offset:offset+TYPE_POINTER_SIZE], &heapOffset)
					if err != nil {
						panic(err)
					}

					if int(heapOffset) < prgrm.HeapStartsAt {
						continue
					}

					if fld.IsPointer || fld.IsSlice || fld.Type == TYPE_STR {
						MarkObjectsTree(prgrm, offset, fld.Type, fld.DeclarationSpecifiers[1:])
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

			ptrIsPointer := IsPointer(ptr)

			// Checking if we need to mark `ptr`.
			if ptrIsPointer {
				// If `ptr` has fields, we need to navigate the heap and mark its fields too.
				if ptr.CustomType != nil {
					// Getting the offset to the object in the heap
					var heapOffset int32
					_, err := encoder.DeserializeAtomic(prgrm.Memory[offset:offset+TYPE_POINTER_SIZE], &heapOffset)
					if err != nil {
						panic(err)
					}

					if int(heapOffset) >= prgrm.HeapStartsAt {
						for _, fld := range ptr.CustomType.Fields {
							MarkObjectsTree(prgrm, int(heapOffset)+OBJECT_HEADER_SIZE+fld.Offset, fld.Type, fld.DeclarationSpecifiers[1:])
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
	for c := prgrm.HeapStartsAt + NULL_HEAP_ADDRESS_OFFSET; c < prgrm.HeapStartsAt+prgrm.HeapPointer; {
		objSize := mustDeserializeI32(prgrm.Memory[c+MARK_SIZE+FORWARDING_ADDRESS_SIZE : c+MARK_SIZE+FORWARDING_ADDRESS_SIZE+OBJECT_SIZE])

		if prgrm.Memory[c] == 1 {
			forwardingAddress := mustDeserializeI32(prgrm.Memory[c+MARK_SIZE : c+MARK_SIZE+FORWARDING_ADDRESS_SIZE])

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


