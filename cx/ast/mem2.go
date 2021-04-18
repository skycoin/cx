package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

//NOTE: Temp file for resolving CalculateDereferences issue
//TODO: What should this function be called?

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

// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size //TODO: WTF is a custom type?
	}
	return arg.Size
}

func CalculateDereferences(arg *CXArgument, finalOffset int, fp int) int {
	var isPointer bool
	var baseOffset int
	var sizeofElement int

	idxCounter := 0
	for _, op := range arg.DereferenceOperations {
		switch op {
		case constants.DEREF_SLICE: //TODO: Move to CalculateDereference_slice
			if len(arg.Indexes) == 0 {
				continue
			}

			isPointer = false
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[finalOffset : finalOffset+constants.TYPE_POINTER_SIZE]

			offset = helper.Deserialize_i32(byts)

			finalOffset = int(offset)

			baseOffset = finalOffset

			finalOffset += constants.OBJECT_HEADER_SIZE
			finalOffset += constants.SLICE_HEADER_SIZE

			//TODO: delete
			sizeToUse := GetDerefSize(arg) //TODO: is always arg.Size unless arg.CustomType != nil
			finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeToUse
			if !IsValidSliceIndex(baseOffset, finalOffset, sizeToUse) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}

			idxCounter++

		case constants.DEREF_ARRAY: //TODO: Move to CalculateDereference_array
			if len(arg.Indexes) == 0 {
				continue
			}
			var subSize = int(1)
			for _, len := range arg.Lengths[idxCounter+1:] {
				subSize *= len
			}

			//TODO: Delete
			sizeToUse := GetDerefSize(arg) //TODO: is always arg.Size unless arg.CustomType != nil

			baseOffset = finalOffset
			sizeofElement = subSize * sizeToUse
			// finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeofElement //TODO: FIX INTEGER CAST
			finalOffset += int(ReadArray(fp, arg.Indexes[idxCounter])) * sizeofElement //TODO: FIX INTEGER CAST
			idxCounter++
		case constants.DEREF_POINTER: //TODO: Move to CalculateDereference_ptr
			isPointer = true
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[finalOffset : finalOffset+constants.TYPE_POINTER_SIZE]

			offset = helper.Deserialize_i32(byts)
			finalOffset = int(offset) //TODO: FIX INTEGER CAST
		}

	}

	// if finalOffset >= PROGRAM.HeapStartsAt {
	if finalOffset >= PROGRAM.HeapStartsAt && isPointer {
		// then it's an object
		finalOffset += constants.OBJECT_HEADER_SIZE
		if arg.IsSlice {
			finalOffset += constants.SLICE_HEADER_SIZE
			if !IsValidSliceIndex(baseOffset, finalOffset, sizeofElement) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}

	return finalOffset
}

// CalculateDereferences_array ...
func CalculateDereferences_array(arg *CXArgument, finalOffset *int, fp int) {
	var sizeofElement int

	idxCounter := 0
	for _, _ = range arg.DereferenceOperations {
		if len(arg.Indexes) == 0 {
			continue
		}
		var subSize = int(1)
		for _, len := range arg.Lengths[idxCounter+1:] {
			subSize *= len
		}

		//TODO: Delete
		sizeToUse := GetDerefSize(arg) //TODO: is always arg.Size unless arg.CustomType != nil

		sizeofElement = subSize * sizeToUse
		// *finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeofElement //TODO: FIX INTEGER CAST
		*finalOffset += int(ReadArray(fp, arg.Indexes[idxCounter])) * sizeofElement //TODO: FIX INTEGER CAST
		idxCounter++
	}
}

// CalculateDereferences_slice
func CalculateDereferences_slice(arg *CXArgument, finalOffset *int, fp int) {

	// remove this check
	if !arg.IsSlice {
		panic("not slice")
	}
	var baseOffset int

	idxCounter := 0
	for _, _ = range arg.DereferenceOperations {
		if len(arg.Indexes) == 0 {
			continue
		}

		var offset int32
		var byts []byte

		byts = PROGRAM.Memory[*finalOffset : *finalOffset+constants.TYPE_POINTER_SIZE]

		offset = helper.Deserialize_i32(byts)

		*finalOffset = int(offset)

		baseOffset = *finalOffset

		*finalOffset += constants.OBJECT_HEADER_SIZE
		*finalOffset += constants.SLICE_HEADER_SIZE

		//TODO: delete
		sizeToUse := GetDerefSize(arg) //TODO: is always arg.Size unless arg.CustomType != nil
		// *finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeToUse
		*finalOffset += int(ReadSlice(fp, arg.Indexes[idxCounter])) * sizeToUse
		if !IsValidSliceIndex(baseOffset, *finalOffset, sizeToUse) {
			panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
		}

		idxCounter++
	}

}

// CalculateDereferences_ptr
func CalculateDereferences_ptr(arg *CXArgument, finalOffset *int, fp int) {
	// remove this check
	if !arg.IsPointer && !arg.IsSlice {
		panic("not pointer")
	}
	var isPointer bool
	var baseOffset int
	var sizeofElement int

	for _, _ = range arg.DereferenceOperations {

		isPointer = true
		var offset int32
		var byts []byte

		byts = PROGRAM.Memory[*finalOffset : *finalOffset+constants.TYPE_POINTER_SIZE]

		offset = helper.Deserialize_i32(byts)
		*finalOffset = int(offset) //TODO: FIX INTEGER CAST

	}

	// if *finalOffset >= PROGRAM.HeapStartsAt {
	if *finalOffset >= PROGRAM.HeapStartsAt && isPointer {
		// then it's an object
		*finalOffset += constants.OBJECT_HEADER_SIZE
		if arg.IsSlice {
			*finalOffset += constants.SLICE_HEADER_SIZE
			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeofElement) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}
}

// CalculateDereferences_i8 ...
func CalculateDereferences_i8(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i16 ...
func CalculateDereferences_i16(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i32 ...
func CalculateDereferences_i32(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_i64 ...
func CalculateDereferences_i64(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui8 ...
func CalculateDereferences_ui8(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui16 ...
func CalculateDereferences_ui16(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui32 ...
func CalculateDereferences_ui32(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_ui64 ...
func CalculateDereferences_ui64(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_f32 ...
func CalculateDereferences_f32(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_f64 ...
func CalculateDereferences_f64(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_str ...
func CalculateDereferences_str(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}

// CalculateDereferences_bool ...
func CalculateDereferences_bool(arg *CXArgument, finalOffset int, fp int) int {
	if len(arg.DereferenceOperations) == 0 {
		panic("0 dereference operations")
	}
	return CalculateDereferences(arg, finalOffset, fp)
}
