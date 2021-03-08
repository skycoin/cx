package cxcore

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

// GetOffset_i8 ...
func GetOffset_i8(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_i16 ...
func GetOffset_i16(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_i32 ...
func GetOffset_i32(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_i64 ...
func GetOffset_i64(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_ui8 ...
func GetOffset_ui8(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_ui16 ...
func GetOffset_ui16(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_ui32 ...
func GetOffset_ui32(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_ui64 ...
func GetOffset_ui64(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_f32 ...
func GetOffset_f32(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_f64 ...
func GetOffset_f64(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_bool ...
func GetOffset_bool(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}

// GetOffset_str ...
func GetOffset_str(fp int, arg *CXArgument) int {
	return GetFinalOffset(fp, arg)
}
