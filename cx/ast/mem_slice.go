package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// IsValidSliceIndex ...
func IsValidSliceIndex(offset types.Pointer, index types.Pointer, sizeofElement types.Pointer) bool {
	sliceLen := GetSliceLen(offset)
	bytesLen := sliceLen * sizeofElement
	index -= types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + offset
	if index >= 0 && index < bytesLen && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetSliceOffset ...
func GetSliceOffset(fp types.Pointer, arg *CXArgument) types.Pointer {
	element := GetAssignmentElement(arg)
	if element.IsSlice {
		return types.Read_ptr(PROGRAM.Memory, GetFinalOffset(fp, arg))
	}

	return types.InvalidPointer
}

// GetSliceHeader ...
func GetSliceHeader(offset types.Pointer) []byte {
	return PROGRAM.Memory[offset+types.OBJECT_HEADER_SIZE : offset+types.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
}

// GetSliceLen ...
func GetSliceLen(offset types.Pointer) types.Pointer {
	sliceHeader := GetSliceHeader(offset)
	return types.Read_ptr(sliceHeader, types.POINTER_SIZE)
}

// GetSlice ...
func GetSlice(offset types.Pointer, sizeofElement types.Pointer) []byte {
	if offset > 0 && offset.IsValid() {
		sliceLen := GetSliceLen(offset)
		if sliceLen > 0 && sliceLen.IsValid() {
			dataOffset := offset + types.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE - types.POINTER_SIZE
			dataLen := types.POINTER_SIZE + sliceLen*sizeofElement
			return PROGRAM.Memory[dataOffset : dataOffset+dataLen]
		}
	}
	return nil
}

// GetSliceData ...
func GetSliceData(offset types.Pointer, sizeofElement types.Pointer) []byte {
	if slice := GetSlice(offset, sizeofElement); slice != nil {
		return slice[types.POINTER_SIZE:]
	}
	return nil
}

// SliceResizeEx does the logic required by `SliceResize`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceResizeEx(outputSliceOffset types.Pointer, count types.Pointer, sizeofElement types.Pointer) types.Pointer {
	var outputSliceHeader []byte
	var outputSliceCap types.Pointer

	if outputSliceOffset > 0 && outputSliceOffset.IsValid() {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
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
		outputSliceOffset = AllocateSeq(outputObjectSize)
		types.Write_ptr(types.Get_obj_header(PROGRAM.Memory, outputSliceOffset), types.OBJECT_GC_HEADER_SIZE, outputObjectSize)

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		types.Write_ptr(outputSliceHeader, 0, newCap)
		types.Write_ptr(outputSliceHeader, types.POINTER_SIZE, newLen)
	}

	return outputSliceOffset
}

// SliceResize ...
func SliceResize(fp types.Pointer, out *CXArgument, inp *CXArgument, count types.Pointer, sizeofElement types.Pointer) types.Pointer {
	outputSliceOffset := GetSliceOffset(fp, out)

	outputSliceOffset = SliceResizeEx(outputSliceOffset, count, sizeofElement)

	SliceCopy(fp, outputSliceOffset, inp, count, sizeofElement)

	return outputSliceOffset
}

// SliceCopyEx does the logic required by `SliceCopy`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceCopyEx(outputSliceOffset types.Pointer, inputSliceOffset types.Pointer, count types.Pointer, sizeofElement types.Pointer) {
	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if outputSliceOffset > 0 && outputSliceOffset.IsValid() {
		outputSliceHeader := GetSliceHeader(outputSliceOffset)
		types.Write_ptr(outputSliceHeader, types.POINTER_SIZE, count)
		outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
		if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
			copy(outputSliceData, GetSliceData(inputSliceOffset, sizeofElement))
		}
	}
}

// SliceCopy copies the contents from the slice located at `inputSliceOffset` to the slice located at `outputSliceOffset`.
func SliceCopy(fp types.Pointer, outputSliceOffset types.Pointer, inp *CXArgument, count types.Pointer, sizeofElement types.Pointer) {
	inputSliceOffset := GetSliceOffset(fp, inp)
	SliceCopyEx(outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(fp types.Pointer, out *CXArgument, inp *CXArgument, sizeofElement types.Pointer, appendLen types.Pointer) types.Pointer {
	inputSliceOffset := GetSliceOffset(fp, inp)
	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	outputSliceOffset := SliceResize(fp, out, inp, inputSliceLen+appendLen, sizeofElement)
	return outputSliceOffset
}

// SliceAppendWrite writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWrite(outputSliceOffset types.Pointer, object []byte, index types.Pointer) {
	sizeofElement := types.Cast_int_to_ptr(len(object))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[index*sizeofElement:], object)
}

// SliceAppendWriteByte writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWriteByte(outputSliceOffset types.Pointer, object []byte, index types.Pointer) {
	outputSliceData := GetSliceData(outputSliceOffset, 1)
	copy(outputSliceData[index:], object)
}

// SliceInsert ...
func SliceInsert(fp types.Pointer, out *CXArgument, inp *CXArgument, index types.Pointer, object []byte) types.Pointer {
	inputSliceOffset := GetSliceOffset(fp, inp)
	// outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	var newLen = inputSliceLen + 1
	sizeofElement := types.Cast_int_to_ptr(len(object))
	outputSliceOffset := SliceResize(fp, out, inp, newLen, sizeofElement)
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[(index+1)*sizeofElement:], outputSliceData[index*sizeofElement:])
	copy(outputSliceData[index*sizeofElement:], object)
	return outputSliceOffset
}

// SliceRemove ...
func SliceRemove(fp types.Pointer, out *CXArgument, inp *CXArgument, index types.Pointer, sizeofElement types.Pointer) types.Pointer {
	inputSliceOffset := GetSliceOffset(fp, inp)
	outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 && inputSliceOffset.IsValid() {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[index*sizeofElement:], outputSliceData[(index+1)*sizeofElement:])
	outputSliceOffset = SliceResize(fp, out, inp, inputSliceLen-1, sizeofElement)
	return outputSliceOffset
}

// WriteToSlice is used to create slices in the backend, i.e. not by calling `append`
// in a CX program, but rather by the CX code itself. This function is used by
// affordances, serialization and to store OS input arguments.
func WriteToSlice(off types.Pointer, inp []byte) types.Pointer {
	// TODO: Check all these parses from/to int32/int.
	var inputSliceLen types.Pointer
	if off != 0 && off.IsValid() {
		inputSliceLen = GetSliceLen(off)
	}

	inpLen := types.Cast_int_to_ptr(len(inp))
	// We first check if a resize is needed. If a resize occurred
	// the address of the new slice will be stored in `newOff` and will
	// be different to `off`.
	newOff := SliceResizeEx(off, inputSliceLen+1, inpLen)

	// Copy the data from the old slice at `off` to `newOff`.
	SliceCopyEx(newOff, off, inputSliceLen+1, inpLen)

	// Write the new slice element `inp` to the slice located at `newOff`.
	SliceAppendWrite(newOff, inp, inputSliceLen)
	return newOff

}
