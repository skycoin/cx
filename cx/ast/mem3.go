package ast

//TODO: Delete this eventually
func GetFinalOffset(fp int, arg *CXArgument) int {
	finalOffset := arg.Offset

	//Todo: find way to eliminate this check
	if finalOffset < PROGRAM.StackSize {
		// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
		finalOffset += fp
	}

	// elt = arg
	//TODO: Eliminate all op codes with more than one return type
	//TODO: Eliminate this loop
	//Q: How can CalculateDereferences change offset?
	//Why is finalOffset fed in as a pointer?
	finalOffset = CalculateDereferences(arg, finalOffset, fp)
	for _, fld := range arg.Fields {
		// elt = fld
		finalOffset += fld.Offset
		finalOffset = CalculateDereferences(fld, finalOffset, fp)
	}

	return finalOffset
}

