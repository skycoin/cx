package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// GetSize ...
func GetSize(arg *CXArgument) types.Pointer {
	if len(arg.Fields) > 0 {
		return GetSize(arg.Fields[len(arg.Fields)-1])
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		if deref == constants.DEREF_SLICE || deref == constants.DEREF_ARRAY {
			declCount := len(arg.DeclarationSpecifiers)
			if declCount > 1 {
			}
			return arg.Size
		}
	}

	for _, decl := range arg.DeclarationSpecifiers {
		if decl == constants.DECL_POINTER || decl == constants.DECL_SLICE || decl == constants.DECL_ARRAY {
			return arg.TotalSize
		}
	}

	if arg.StructType != nil {
		return arg.StructType.Size
	}

	return arg.TotalSize
}

// GetDerefSize ...
func GetDerefSize(arg *CXArgument, index int, derefPointer bool, derefArray bool) types.Pointer {
	if !derefArray && len(arg.Lengths) > 1 && ((index + 1) < len(arg.Lengths)) {
		return types.POINTER_SIZE
	}
	if arg.StructType != nil {
		return arg.StructType.Size
	}
	if derefPointer {
		return arg.TotalSize
	}
	return arg.Size
}

// GetDerefSizeSlice ...
func GetDerefSizeSlice(arg *CXArgument) types.Pointer {
	if len(arg.Lengths) > 1 && (len(arg.Lengths)-len(arg.Indexes)) > 1 {
		return types.POINTER_SIZE
	}
	if arg.StructType != nil {
		return arg.StructType.Size
	}
	return arg.Size
}

func GetFinalOffset(fp types.Pointer, arg *CXArgument) types.Pointer {
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

func CalculateDereferences(arg *CXArgument, finalOffset types.Pointer, fp types.Pointer) types.Pointer {
	var isPointer bool

	var baseOffset types.Pointer
	var sizeofElement types.Pointer
	var derefPointer bool
	idxCounter := 0
	for _, op := range arg.DereferenceOperations {
		switch op {
		case constants.DEREF_SLICE: //TODO: Move to CalculateDereference_slice
			if len(arg.Indexes) == 0 {
				continue
			}

			isPointer = false
			finalOffset = types.Read_ptr(PROGRAM.Memory, finalOffset)
			baseOffset = finalOffset

			finalOffset += types.OBJECT_HEADER_SIZE
			finalOffset += constants.SLICE_HEADER_SIZE

			//TODO: delete
			sizeToUse := GetDerefSize(arg, idxCounter, derefPointer, false) //TODO: is always arg.Size unless arg.StructType != nil
			derefPointer = false

			indexOffset := GetFinalOffset(fp, arg.Indexes[idxCounter])
			indexValue := types.Read_i32(PROGRAM.Memory, indexOffset)

			finalOffset += types.Cast_i32_to_ptr(indexValue) * sizeToUse // TODO:PTR Use ptr/Read_ptr, array/slice indexing only works with i32.
			if !IsValidSliceIndex(baseOffset, finalOffset, sizeToUse) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}

			idxCounter++

		case constants.DEREF_ARRAY: //TODO: Move to CalculateDereference_array
			if len(arg.Indexes) == 0 {
				continue
			}
			var subSize = types.Pointer(1)
			for _, len := range arg.Lengths[idxCounter+1:] {
				subSize *= len
			}

			//TODO: Delete
			sizeToUse := GetDerefSize(arg, idxCounter, derefPointer, true) //TODO: is always arg.Size unless arg.StructType != nil
			derefPointer = false
			baseOffset = finalOffset
			sizeofElement = subSize * sizeToUse

			finalOffset += types.Cast_i32_to_ptr(types.Read_i32(PROGRAM.Memory, GetFinalOffset(fp, arg.Indexes[idxCounter]))) * sizeofElement // TODO:PTR Use Read_ptr
			idxCounter++
		case constants.DEREF_POINTER: //TODO: Move to CalculateDereference_ptr
			isPointer = true
			finalOffset = types.Read_ptr(PROGRAM.Memory, finalOffset)
			derefPointer = true
		}
	}

	// if finalOffset >= PROGRAM.HeapStartsAt {
	if finalOffset.IsValid() && finalOffset >= PROGRAM.HeapStartsAt && isPointer {
		// then it's an object
		finalOffset += types.OBJECT_HEADER_SIZE
		if arg.IsSlice {
			finalOffset += constants.SLICE_HEADER_SIZE
			if !IsValidSliceIndex(baseOffset, finalOffset, sizeofElement) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}

	return finalOffset
}
