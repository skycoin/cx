package mem

/*
// DebugHeap prints the symbols that are acting as pointers in a CX program at certain point during the execution of the program along with the addresses they are pointing. Additionally, a list of the objects in the heap is printed, which shows their address in the heap, if they are marked as alive or as dead by the garbage collector, the address where they used to live after a garbage collector call, the full size of the object, the object itself as a slice of bytes and the pointers that are pointing to that object.
func DebugHeap() {
	// symsToAddrs will hold a list of symbols that are pointing to an address.
	symsToAddrs := make(map[int32][]string)

	// Processing global variables. Adding the address they are pointing to.
	for _, pkg := range PROGRAM.Packages {
		for _, glbl := range pkg.Globals {
			if glbl.IsPointer || glbl.IsSlice {
				heapOffset := helper.Deserialize_i32(PROGRAM.Memory[glbl.Offset : glbl.Offset+constants.TYPE_POINTER_SIZE])

				symsToAddrs[heapOffset] = append(symsToAddrs[heapOffset], glbl.Name)
			}
		}
	}

	// Processing local variables in every active function call in the `CallStack`.
	// Adding the address they are pointing to.
	var fp int
	for c := 0; c <= PROGRAM.CallCounter; c++ {
		op := PROGRAM.CallStack[c].Operator

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
			symName := ptr.Name
			if len(ptr.Fields) > 0 {
				fld := ptr.Fields[len(ptr.Fields)-1]
				offset += fld.Offset
				symName += "." + fld.Name
			}

			if ptr.Offset < PROGRAM.StackSize {
				offset += fp
			}

			heapOffset := helper.Deserialize_i32(PROGRAM.Memory[offset : offset+constants.TYPE_POINTER_SIZE])

			symsToAddrs[heapOffset] = append(symsToAddrs[heapOffset], symName)
		}

		fp += op.Size
	}

	// Printing all the details.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, '.', 0)

	for off, symNames := range symsToAddrs {
		var addrB [4]byte
		cxcore.WriteMemI32(addrB[:], 0, off)
		fmt.Fprintln(w, "Addr:\t", addrB, "\tPtr:\t", symNames)
	}

	// Just a newline.
	fmt.Fprintln(w)
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, '.', 0)

	for c := PROGRAM.HeapStartsAt + constants.NULL_HEAP_ADDRESS_OFFSET; c < PROGRAM.HeapStartsAt+PROGRAM.HeapPointer; {
		objSize := helper.Deserialize_i32(PROGRAM.Memory[c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE : c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE+constants.OBJECT_SIZE])

		// Setting a limit size for the object to be printed if the object is too large.
		// We don't want to print obscenely large objects to standard output.
		printObjSize := objSize
		if objSize > 50 {
			printObjSize = 50
		}

		var addrB [4]byte
		cxcore.WriteMemI32(addrB[:], 0, int32(c))

		fmt.Fprintln(w, "Addr:\t", addrB, "\tMark:\t", PROGRAM.Memory[c:c+constants.MARK_SIZE], "\tFwd:\t", PROGRAM.Memory[c+constants.MARK_SIZE:c+constants.MARK_SIZE+constants.FORWARDING_ADDRESS_SIZE], "\tSize:\t", objSize, "\tObj:\t", PROGRAM.Memory[c+constants.OBJECT_HEADER_SIZE:c+int(printObjSize)], "\tPtrs:", symsToAddrs[int32(c)])

		c += int(objSize)
	}

	// Just a newline.
	fmt.Fprintln(w)
	w.Flush()
}
*/