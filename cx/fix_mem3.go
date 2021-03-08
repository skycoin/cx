package cxcore

import ()

//NOTE: Temp file for resolving GetFinalOffset issue
//TODO: What should this function be called?

// GetFinalOffset ...
// Get byte position to write operand result back to?
// Note: CXArgument is bloated, should pass in required though something else

//NEEDS COMMENT. WTF DOES THIS DO?
//TODO:
//GetFinalOffset
//->
//GetFinalOffsetI32
//GetFinalOffsetF32
//GetfinalOffsetI16
//ETC
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
