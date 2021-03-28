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

/*
// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}
*/

// GetDerefSize ...
func GetDerefSize(arg *CXArgument) int {
	if arg.CustomType != nil {
		return arg.CustomType.Size
	}
	return arg.Size
}


func CalculateDereferences(arg *CXArgument, finalOffset *int, fp int) {
	var isPointer bool
	var baseOffset int
	var sizeofElement int

	idxCounter := 0
	for _, op := range arg.DereferenceOperations {
		switch op {
		case constants.DEREF_SLICE:
			if len(arg.Indexes) == 0 {
				continue
			}

			isPointer = false
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+constants.TYPE_POINTER_SIZE]

			offset = helper.Deserialize_i32(byts)

			*finalOffset = int(offset)

			baseOffset = *finalOffset

			*finalOffset += constants.OBJECT_HEADER_SIZE
			*finalOffset += constants.SLICE_HEADER_SIZE

			sizeToUse := GetDerefSize(arg) //GetDerefSize
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeToUse
			if !IsValidSliceIndex(baseOffset, *finalOffset, sizeToUse) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}

			idxCounter++
		case constants.DEREF_ARRAY:
			if len(arg.Indexes) == 0 {
				continue
			}
			var subSize = int(1)
			for _, len := range arg.Lengths[idxCounter+1:] {
				subSize *= len
			}

			sizeToUse := GetDerefSize(arg) //GetDerefSize

			baseOffset = *finalOffset
			sizeofElement = subSize * sizeToUse
			*finalOffset += int(ReadI32(fp, arg.Indexes[idxCounter])) * sizeofElement
			idxCounter++
		case constants.DEREF_POINTER:
			isPointer = true
			var offset int32
			var byts []byte

			byts = PROGRAM.Memory[*finalOffset : *finalOffset+constants.TYPE_POINTER_SIZE]

			offset = helper.Deserialize_i32(byts)
			*finalOffset = int(offset)
		}

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