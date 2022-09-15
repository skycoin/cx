package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// TODO: rewrite slice api.
// TODO: remove all offset > 0 tests.
// TODO: panic(constants.CX_RUNTIME_*) instead of silent failure.

// IsValidSliceIndex ...
func IsValidSliceIndex(prgrm *CXProgram, offset types.Pointer, index types.Pointer, sizeofElement types.Pointer) bool {
	sliceLen := GetSliceLen(prgrm, offset)
	bytesLen := sliceLen * sizeofElement
	index -= types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + offset
	if index >= 0 && index < bytesLen && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetSliceOffset ...
func GetSliceOffset(prgrm *CXProgram, fp types.Pointer, argTypeSig *CXTypeSignature) types.Pointer {
	var element *CXArgument
	if argTypeSig.Type == TYPE_CXARGUMENT_DEPRECATE {
		arg := prgrm.GetCXArgFromArray(CXArgumentIndex(argTypeSig.Meta))
		element = arg.GetAssignmentElement(prgrm)
	} else if argTypeSig.Type == TYPE_ATOMIC || argTypeSig.Type == TYPE_POINTER_ATOMIC {
		return types.InvalidPointer
	} else {
		panic("type is not known")
	}

	if element.IsSlice {
		return types.Read_ptr(prgrm.Memory, GetFinalOffset(prgrm, fp, nil, argTypeSig))
	}

	return types.InvalidPointer
}

// GetSliceHeader ...
func GetSliceHeader(prgrm *CXProgram, offset types.Pointer) []byte {
	if offset > 0 && offset.IsValid() {
		return prgrm.Memory[offset+types.OBJECT_HEADER_SIZE : offset+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
	}
	return nil
}

// GetSliceLen ...
func GetSliceLen(prgrm *CXProgram, offset types.Pointer) types.Pointer {
	if offset > 0 && offset.IsValid() {
		sliceHeader := GetSliceHeader(prgrm, offset)
		return types.Read_ptr(sliceHeader, types.POINTER_SIZE)
	}
	return 0
}

// GetSlice ...
func GetSlice(prgrm *CXProgram, offset types.Pointer, sizeofElement types.Pointer) []byte {
	if offset > 0 && offset.IsValid() {
		sliceLen := GetSliceLen(prgrm, offset)
		if sliceLen > 0 && sliceLen.IsValid() {
			dataOffset := offset + types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE - types.POINTER_SIZE
			dataLen := types.POINTER_SIZE + sliceLen*sizeofElement
			return prgrm.Memory[dataOffset : dataOffset+dataLen]
		}
	}
	return nil
}

// GetSliceData ...
func GetSliceData(prgrm *CXProgram, offset types.Pointer, sizeofElement types.Pointer) []byte {
	if slice := GetSlice(prgrm, offset, sizeofElement); slice != nil {
		return slice[types.POINTER_SIZE:]
	}
	return nil
}

// SliceResizeEx does the logic required by `SliceResize`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceResizeEx(prgrm *CXProgram, outputSliceOffset types.Pointer, count types.Pointer, sizeofElement types.Pointer) types.Pointer {
	var outputSliceHeader []byte
	var outputSliceCap types.Pointer

	if outputSliceOffset > 0 && outputSliceOffset.IsValid() {
		outputSliceHeader = GetSliceHeader(prgrm, outputSliceOffset)
		outputSliceCap = types.Read_ptr(outputSliceHeader, 0)
	}

	var newLen = count
	var newCap = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize = types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + newCap*sizeofElement
		outputSliceOffset = AllocateSeq(prgrm, outputObjectSize)
		types.Write_ptr(types.Get_obj_header(prgrm.Memory, outputSliceOffset), types.OBJECT_GC_HEADER_SIZE, outputObjectSize)

		outputSliceHeader = GetSliceHeader(prgrm, outputSliceOffset)
		types.Write_ptr(outputSliceHeader, 0, newCap)
	}

	if outputSliceHeader != nil {
		types.Write_ptr(outputSliceHeader, types.POINTER_SIZE, newLen)
	}
	return outputSliceOffset
}

// SliceResize ...
func SliceResize(prgrm *CXProgram, fp types.Pointer, out *CXTypeSignature, inp *CXTypeSignature, count types.Pointer, sizeofElement types.Pointer) types.Pointer {
	inputSliceOffset := GetSliceOffset(prgrm, fp, inp)
	outputSliceOffset := SliceResizeEx(prgrm, inputSliceOffset, count, sizeofElement)

	if outputSliceOffset != inputSliceOffset {
		inputSliceLen := GetSliceLen(prgrm, inputSliceOffset)
		if inputSliceLen > 0 {
			SliceCopy(prgrm, fp, outputSliceOffset, inp, inputSliceLen, sizeofElement)
		}
	}

	return outputSliceOffset
}

// SliceCopyEx does the logic required by `SliceCopy`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceCopyEx(prgrm *CXProgram, outputSliceOffset types.Pointer, inputSliceOffset types.Pointer, count types.Pointer, sizeofElement types.Pointer) {
	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(prgrm, inputSliceOffset)
	}

	if outputSliceOffset > 0 && outputSliceOffset.IsValid() {
		outputSliceData := GetSliceData(prgrm, outputSliceOffset, sizeofElement)
		if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
			copy(outputSliceData, GetSliceData(prgrm, inputSliceOffset, sizeofElement))
		}
	}
}

// SliceCopy copies the contents from the slice located at `inputSliceOffset` to the slice located at `outputSliceOffset`.
func SliceCopy(prgrm *CXProgram, fp types.Pointer, outputSliceOffset types.Pointer, inp *CXTypeSignature, count types.Pointer, sizeofElement types.Pointer) {
	inputSliceOffset := GetSliceOffset(prgrm, fp, inp)
	SliceCopyEx(prgrm, outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(prgrm *CXProgram, fp types.Pointer, out *CXTypeSignature, inp *CXTypeSignature, sizeofElement types.Pointer, appendLen types.Pointer) types.Pointer {
	inputSliceOffset := GetSliceOffset(prgrm, fp, inp)
	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(prgrm, inputSliceOffset)
	}

	outputSliceOffset := SliceResize(prgrm, fp, out, inp, inputSliceLen+appendLen, sizeofElement)
	return outputSliceOffset
}

// SliceAppendWrite writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWrite(prgrm *CXProgram, outputSliceOffset types.Pointer, object []byte, index types.Pointer) {
	sizeofElement := types.Cast_int_to_ptr(len(object))
	outputSliceData := GetSliceData(prgrm, outputSliceOffset, sizeofElement)
	copy(outputSliceData[index*sizeofElement:], object)
}

// SliceAppendWriteByte writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWriteByte(prgrm *CXProgram, outputSliceOffset types.Pointer, object []byte, index types.Pointer) {
	outputSliceData := GetSliceData(prgrm, outputSliceOffset, 1)
	copy(outputSliceData[index:], object)
}

// SliceInsert ...
func SliceInsert(prgrm *CXProgram, fp types.Pointer, out *CXTypeSignature, inp *CXTypeSignature, index types.Pointer, object []byte) types.Pointer {
	inputSliceOffset := GetSliceOffset(prgrm, fp, inp)

	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(prgrm, inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	var newLen = inputSliceLen + 1
	sizeofElement := types.Cast_int_to_ptr(len(object))
	outputSliceOffset := SliceResize(prgrm, fp, out, inp, newLen, sizeofElement)
	outputSliceData := GetSliceData(prgrm, outputSliceOffset, sizeofElement)
	copy(outputSliceData[(index+1)*sizeofElement:], outputSliceData[index*sizeofElement:])
	copy(outputSliceData[index*sizeofElement:], object)
	return outputSliceOffset
}

// SliceRemove ...
func SliceRemove(prgrm *CXProgram, fp types.Pointer, out *CXTypeSignature, inp *CXTypeSignature, index types.Pointer, sizeofElement types.Pointer) types.Pointer {
	inputSliceOffset := GetSliceOffset(prgrm, fp, inp)
	outputSliceOffset := GetSliceOffset(prgrm, fp, out)

	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(prgrm, inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	outputSliceData := GetSliceData(prgrm, outputSliceOffset, sizeofElement)
	copy(outputSliceData[index*sizeofElement:], outputSliceData[(index+1)*sizeofElement:])
	outputSliceHeader := GetSliceHeader(prgrm, outputSliceOffset)
	types.Write_ptr(outputSliceHeader, types.POINTER_SIZE, inputSliceLen-1)
	return outputSliceOffset
}

// WriteToSlice is used to create slices in the backend, i.e. not by calling `append`
// in a CX program, but rather by the CX code itself. This function is used by
// affordances, serialization and to store OS input arguments.
func WriteToSlice(prgrm *CXProgram, off types.Pointer, inp []byte) types.Pointer {
	// TODO: Check all these parses from/to int32/int.
	var inputSliceLen types.Pointer
	if off != 0 && off.IsValid() {
		inputSliceLen = GetSliceLen(prgrm, off)
	}

	inpLen := types.Cast_int_to_ptr(len(inp))
	// We first check if a resize is needed. If a resize occurred
	// the address of the new slice will be stored in `newOff` and will
	// be different to `off`.
	newOff := SliceResizeEx(prgrm, off, inputSliceLen+1, inpLen)

	// Copy the data from the old slice at `off` to `newOff`.
	SliceCopyEx(prgrm, newOff, off, inputSliceLen, inpLen)

	// Write the new slice element `inp` to the slice located at `newOff`.
	SliceAppendWrite(prgrm, newOff, inp, inputSliceLen)
	return newOff

}
