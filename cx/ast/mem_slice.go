package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

// IsValidSliceIndex ...
func IsValidSliceIndex(offset int, index int, sizeofElement int) bool {
	sliceLen := GetSliceLen(int32(offset))
	bytesLen := sliceLen * int32(sizeofElement)
	index -= constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + offset

	if index >= 0 && index < int(bytesLen) && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetSliceOffset ...
//TODO: DANGER, WEIRD INT CAST FROM GetFinalOffset
func GetSliceOffset(fp int, arg *CXArgument) int32 {
	element := GetAssignmentElement(arg)
	if element.IsSlice {
		return GetPointerOffset(int32(GetFinalOffset(fp, arg)))
		// return GetPointerOffset(int32(GetOffset_slice(fp, arg)))
	}

	return -1
}

// GetObjectHeader ...
func GetObjectHeader(offset int32) []byte {
	return PROGRAM.Memory[offset : offset+constants.OBJECT_HEADER_SIZE]
}

// GetSliceHeader ...
func GetSliceHeader(offset int32) []byte {
	return PROGRAM.Memory[offset+constants.OBJECT_HEADER_SIZE : offset+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
}

// GetSliceLen ...
func GetSliceLen(offset int32) int32 {
	sliceHeader := GetSliceHeader(offset)
	return helper.Deserialize_i32(sliceHeader[4:8])
}

// GetSlice ...
func GetSlice(offset int32, sizeofElement int) []byte {
	if offset > 0 {
		sliceLen := GetSliceLen(offset)
		if sliceLen > 0 {
			dataOffset := offset + constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE - 4
			dataLen := 4 + sliceLen*int32(sizeofElement)
			return PROGRAM.Memory[dataOffset : dataOffset+dataLen]
		}
	}
	return nil
}

// GetSliceData ...
func GetSliceData(offset int32, sizeofElement int) []byte {
	if slice := GetSlice(offset, sizeofElement); slice != nil {
		return slice[4:]
	}
	return nil
}

// SliceResizeEx does the logic required by `SliceResize`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceResizeEx(outputSliceOffset int32, count int32, sizeofElement int) int {
	if count < 0 {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var outputSliceHeader []byte
	var outputSliceCap int32

	if outputSliceOffset > 0 {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		outputSliceCap = helper.Deserialize_i32(outputSliceHeader[0:4])
	}

	var newLen = count
	var newCap = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize = constants.OBJECT_HEADER_SIZE + constants.SLICE_HEADER_SIZE + newCap*int32(sizeofElement)
		outputSliceOffset = int32(AllocateSeq(int(outputObjectSize)))
		WriteMemI32(GetObjectHeader(outputSliceOffset)[5:9], 0, outputObjectSize)

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		WriteMemI32(outputSliceHeader[0:4], 0, newCap)
		WriteMemI32(outputSliceHeader[4:8], 0, newLen)
	}

	return int(outputSliceOffset)
}

// SliceResize ...
func SliceResize(fp int, out *CXArgument, inp *CXArgument, count int32, sizeofElement int) int {
	outputSliceOffset := GetSliceOffset(fp, out)

	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, sizeofElement))

	SliceCopy(fp, outputSliceOffset, inp, count, sizeofElement)

	return int(outputSliceOffset)
}

// SliceCopyEx does the logic required by `SliceCopy`. It is separated because some other functions might have access to the offsets of the slices, but not the `CXArgument`s.
func SliceCopyEx(outputSliceOffset int32, inputSliceOffset int32, count int32, sizeofElement int) {
	if count < 0 {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if outputSliceOffset > 0 {
		outputSliceHeader := GetSliceHeader(outputSliceOffset)
		WriteMemI32(outputSliceHeader[4:8], 0, count)
		outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
		if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
			copy(outputSliceData, GetSliceData(inputSliceOffset, sizeofElement))
		}
	}
}

// SliceCopy copies the contents from the slice located at `inputSliceOffset` to the slice located at `outputSliceOffset`.
func SliceCopy(fp int, outputSliceOffset int32, inp *CXArgument, count int32, sizeofElement int) {
	inputSliceOffset := GetSliceOffset(fp, inp)
	SliceCopyEx(outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(fp int, out *CXArgument, inp *CXArgument, sizeofElement int) int32 {
	inputSliceOffset := GetSliceOffset(fp, inp)
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	// TODO: Are we limited then to only one element for now? (because of that +1)
	outputSliceOffset := int32(SliceResize(fp, out, inp, inputSliceLen+1, sizeofElement))
	return outputSliceOffset
}

// SliceAppendWrite writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWrite(outputSliceOffset int32, object []byte, index int32) {
	sizeofElement := len(object)
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(index)*sizeofElement:], object)
}

// SliceAppendWriteByte writes `object` to a slice that is guaranteed to be able to hold `object`, i.e. it had to be checked by `SliceAppendResize` first in case it needed to be resized.
func SliceAppendWriteByte(outputSliceOffset int32, object []byte, index int32) {
	outputSliceData := GetSliceData(outputSliceOffset, 1)
	copy(outputSliceData[int(index):], object)
}

// SliceInsert ...
func SliceInsert(fp int, out *CXArgument, inp *CXArgument, index int32, object []byte) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	// outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	var newLen = inputSliceLen + 1
	sizeofElement := len(object)
	outputSliceOffset := int32(SliceResize(fp, out, inp, newLen, sizeofElement))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(index+1)*sizeofElement:], outputSliceData[int(index)*sizeofElement:])
	copy(outputSliceData[int(index)*sizeofElement:], object)
	return int(outputSliceOffset)
}

// SliceRemove ...
func SliceRemove(fp int, out *CXArgument, inp *CXArgument, index int32, sizeofElement int32) int {
	inputSliceOffset := GetSliceOffset(fp, inp)
	outputSliceOffset := GetSliceOffset(fp, out)

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	outputSliceData := GetSliceData(outputSliceOffset, int(sizeofElement))
	copy(outputSliceData[index*sizeofElement:], outputSliceData[(index+1)*sizeofElement:])
	outputSliceOffset = int32(SliceResize(fp, out, inp, inputSliceLen-1, int(sizeofElement)))
	return int(outputSliceOffset)
}

// WriteToSlice is used to create slices in the backend, i.e. not by calling `append`
// in a CX program, but rather by the CX code itself. This function is used by
// affordances, serialization and to store OS input arguments.
func WriteToSlice(off int, inp []byte) int {
	// TODO: Check all these parses from/to int32/int.
	var inputSliceLen int32
	if off != 0 {
		inputSliceLen = GetSliceLen(int32(off))
	}

	inpLen := len(inp)
	// We first check if a resize is needed. If a resize occurred
	// the address of the new slice will be stored in `newOff` and will
	// be different to `off`.
	newOff := SliceResizeEx(int32(off), inputSliceLen+1, inpLen)

	// Copy the data from the old slice at `off` to `newOff`.
	SliceCopyEx(int32(newOff), int32(off), inputSliceLen+1, inpLen)

	// Write the new slice element `inp` to the slice located at `newOff`.
	SliceAppendWrite(int32(newOff), inp, inputSliceLen)
	return newOff

}
