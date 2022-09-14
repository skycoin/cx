package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

type DereferenceStruct struct {
	isPointer     bool
	baseOffset    types.Pointer
	sizeofElement types.Pointer
	derefPointer  bool
	idxCounter    int

	arg         *CXArgument
	finalOffset types.Pointer
	fp          types.Pointer
}

func TotalLength(lengths []types.Pointer) types.Pointer {
	total := types.Pointer(1)
	for _, i := range lengths {
		total *= i
	}

	return total
}

// GetSize ...
func GetArgSize(prgrm *CXProgram, arg *CXArgument) types.Pointer {
	if len(arg.Fields) > 0 {
		return GetArgSize(prgrm, prgrm.GetCXArgFromArray(arg.Fields[len(arg.Fields)-1]))
	}

	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]

		switch deref {
		case constants.DEREF_SLICE, constants.DEREF_ARRAY:
			return arg.Size
		case constants.DEREF_POINTER:
			return arg.Type.Size()
		}
	}

	for _, decl := range arg.DeclarationSpecifiers {
		switch decl {
		case constants.DECL_ARRAY:
			return arg.Size * TotalLength(arg.Lengths)
		case constants.DECL_POINTER:
			return arg.Size
		case constants.DECL_SLICE:
			return types.POINTER_SIZE
		}
	}

	if arg.StructType != nil {
		return arg.StructType.GetStructSize(prgrm)
	}

	return arg.Size
}

func GetNativeSize(prgrm *CXProgram, arg *CXArgument) types.Pointer {
	derefCount := len(arg.DereferenceOperations)
	if derefCount > 0 {
		deref := arg.DereferenceOperations[derefCount-1]
		switch deref {
		case constants.DEREF_SLICE, constants.DEREF_ARRAY:
			return arg.Size
		case constants.DEREF_POINTER:
			return arg.Type.Size()
		}
	}

	for _, decl := range arg.DeclarationSpecifiers {
		switch decl {
		case constants.DECL_ARRAY:
			return arg.Size * TotalLength(arg.Lengths)
		case constants.DECL_POINTER:
			return arg.Size
		case constants.DECL_SLICE:
			return types.POINTER_SIZE
		}
	}

	if arg.StructType != nil {
		return arg.StructType.GetStructSize(prgrm)
	}

	return arg.Size
}

// GetDerefSize ...
func GetDerefSize(prgrm *CXProgram, arg *CXArgument, index int, derefPointer bool, derefArray bool) types.Pointer {
	if !derefArray && len(arg.Lengths) > 1 && ((index + 1) < len(arg.Lengths)) {
		return types.POINTER_SIZE
	}
	if arg.StructType != nil {
		return arg.StructType.GetStructSize(prgrm)
	}
	if derefPointer {
		return arg.Type.Size()
	}
	return arg.Size
}

// GetDerefSizeSlice ...
func GetDerefSizeSlice(prgrm *CXProgram, arg *CXArgument) types.Pointer {
	if len(arg.Lengths) > 1 && (len(arg.Lengths)-len(arg.Indexes)) > 1 {
		return types.POINTER_SIZE
	}
	if arg.StructType != nil {
		return arg.StructType.GetStructSize(prgrm)
	}
	return arg.Size
}

func GetFinalOffset(prgrm *CXProgram, fp types.Pointer, oldArg *CXArgument, argTypeSig *CXTypeSignature) types.Pointer {
	var finalOffset types.Pointer

	var arg *CXArgument = &CXArgument{}
	if oldArg != nil {
		arg = oldArg
	} else if argTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
		arg = &prgrm.CXArgs[argTypeSig.Meta]
	} else if argTypeSig.Type == TYPE_ATOMIC {
		argTypeSigOffset := argTypeSig.Offset

		//Todo: find way to eliminate this check
		if argTypeSigOffset < prgrm.Stack.Size {
			// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
			argTypeSigOffset += fp
		}

		return argTypeSigOffset
	} else if argTypeSig.Type == TYPE_POINTER_ATOMIC {
		finalOffset = argTypeSig.Offset

		if argTypeSig.IsDeref {
			finalOffset = types.Read_ptr(prgrm.Memory, finalOffset)
		}

		//Todo: find way to eliminate this check
		if finalOffset < prgrm.Stack.Size {
			// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
			finalOffset += fp
		}

		if finalOffset.IsValid() && finalOffset >= prgrm.Heap.StartsAt {
			// then it's an object
			finalOffset += types.OBJECT_HEADER_SIZE
		}

		return finalOffset
	}
	finalOffset = arg.Offset

	//Todo: find way to eliminate this check
	if finalOffset < prgrm.Stack.Size {
		// Then it's in the stack, not in data or heap and we need to consider the frame pointer.
		finalOffset += fp
	}

	// elt = arg
	//TODO: Eliminate all op codes with more than one return type
	//TODO: Eliminate this loop
	//Q: How can CalculateDereferences change offset?
	//Why is finalOffset fed in as a pointer?

	finalOffset = CalculateDereferences(prgrm, arg, finalOffset, fp)
	for _, fldIdx := range arg.Fields {
		fld := prgrm.GetCXArgFromArray(fldIdx)
		// elt = fld
		finalOffset += fld.Offset
		finalOffset = CalculateDereferences(prgrm, fld, finalOffset, fp)
	}

	return finalOffset
}

func CalculateDereference_Slice(prgrm *CXProgram, drfsStruct *DereferenceStruct) {
	drfsStruct.isPointer = false
	drfsStruct.finalOffset = types.Read_ptr(prgrm.Memory, drfsStruct.finalOffset)
	drfsStruct.baseOffset = drfsStruct.finalOffset

	drfsStruct.finalOffset += types.OBJECT_HEADER_SIZE
	drfsStruct.finalOffset += constants.SLICE_HEADER_SIZE

	//TODO: delete
	sizeToUse := GetDerefSize(prgrm, drfsStruct.arg, drfsStruct.idxCounter, drfsStruct.derefPointer, false) //TODO: is always arg.Size unless arg.StructType != nil
	drfsStruct.derefPointer = false

	indexOffset := GetFinalOffset(prgrm, drfsStruct.fp, nil, prgrm.GetCXTypeSignatureFromArray(drfsStruct.arg.Indexes[drfsStruct.idxCounter]))
	indexValue := types.Read_i32(prgrm.Memory, indexOffset)

	drfsStruct.finalOffset += types.Cast_i32_to_ptr(indexValue) * sizeToUse // TODO:PTR Use ptr/Read_ptr, array/slice indexing only works with i32.
	if !IsValidSliceIndex(prgrm, drfsStruct.baseOffset, drfsStruct.finalOffset, sizeToUse) {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	drfsStruct.idxCounter++
}

func CalculateDereference_Array(prgrm *CXProgram, drfsStruct *DereferenceStruct) {
	var subSize = types.Pointer(1)
	for _, len := range drfsStruct.arg.Lengths[drfsStruct.idxCounter+1:] {
		subSize *= len
	}

	//TODO: Delete
	sizeToUse := GetDerefSize(prgrm, drfsStruct.arg, drfsStruct.idxCounter, drfsStruct.derefPointer, true) //TODO: is always arg.Size unless arg.StructType != nil
	drfsStruct.derefPointer = false
	drfsStruct.baseOffset = drfsStruct.finalOffset
	drfsStruct.sizeofElement = subSize * sizeToUse

	drfsStruct.finalOffset += types.Cast_i32_to_ptr(types.Read_i32(prgrm.Memory, GetFinalOffset(prgrm, drfsStruct.fp, nil, prgrm.GetCXTypeSignatureFromArray(drfsStruct.arg.Indexes[drfsStruct.idxCounter])))) * drfsStruct.sizeofElement // TODO:PTR Use Read_ptr
	drfsStruct.idxCounter++
}

func CalculateDereference_Pointer(prgrm *CXProgram, drfsStruct *DereferenceStruct) {
	drfsStruct.isPointer = true
	drfsStruct.finalOffset = types.Read_ptr(prgrm.Memory, drfsStruct.finalOffset)
	drfsStruct.derefPointer = true
}

func CalculateDereferences(prgrm *CXProgram, arg *CXArgument, finalOffset types.Pointer, fp types.Pointer) types.Pointer {
	drfsStruct := &DereferenceStruct{
		idxCounter:  0,
		arg:         arg,
		finalOffset: finalOffset,
		fp:          fp,
	}

	for _, op := range drfsStruct.arg.DereferenceOperations {
		if len(drfsStruct.arg.Indexes) == 0 && op != constants.DEREF_POINTER {
			continue
		}

		switch op {
		case constants.DEREF_SLICE:
			CalculateDereference_Slice(prgrm, drfsStruct)
		case constants.DEREF_ARRAY:
			CalculateDereference_Array(prgrm, drfsStruct)
		case constants.DEREF_POINTER:
			CalculateDereference_Pointer(prgrm, drfsStruct)
		}
	}

	if drfsStruct.finalOffset.IsValid() && drfsStruct.finalOffset >= prgrm.Heap.StartsAt && drfsStruct.isPointer {
		// then it's an object
		drfsStruct.finalOffset += types.OBJECT_HEADER_SIZE
		if drfsStruct.arg.IsSlice {
			drfsStruct.finalOffset += constants.SLICE_HEADER_SIZE
			if !IsValidSliceIndex(prgrm, drfsStruct.baseOffset, drfsStruct.finalOffset, drfsStruct.sizeofElement) {
				panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
			}
		}
	}

	return drfsStruct.finalOffset
}
