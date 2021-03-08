package cxcore

//NOTE: Temp file for resolving CalculateDereferences issue
//TODO: What should this function be called?

// CalculateDereferences ... ...
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

// TODO: str, bool
// TODO: f32, f64

// CalculateDereferences_i8 ...
func CalculateDereferences_i8(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i16 ...
func CalculateDereferences_i16(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i32 ...
func CalculateDereferences_i32(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i64 ...
func CalculateDereferences_i64(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui8 ...
func CalculateDereferences_ui8(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui16 ...
func CalculateDereferences_ui16(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui32 ...
func CalculateDereferences_ui32(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui64 ...
func CalculateDereferences_ui64(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_f32 ...
func CalculateDereferences_f32(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_f64 ...
func CalculateDereferences_f64(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_str ...
func CalculateDereferences_str(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_bool ...
func CalculateDereferences_bool(arg *CXArgument, finalOffset *int, fp int) {
	CalculateDereferences(arg, finalOffset, fp)
}
