package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/helper"
)

// Debug ...
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// GetType ...
func GetType(arg *ast.CXArgument) int {
    fieldCount := len(arg.Fields)
    if fieldCount > 0 {
        return GetType(arg.Fields[fieldCount - 1])
    }

    return arg.Type
}

// ExprOpName ...
func ExprOpName(expr *ast.CXExpression) string {
	if expr.Operator.IsAtomic {
		return globals.OpNames[expr.Operator.OpCode]
	}
	return expr.Operator.Name

}

// func limitString (str string) string {
// 	if len(str) > 3
// }

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

// PrintStack ...
func (cxprogram *ast.CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Callstack===")

	// we're going backwards in the stack
	fp := cxprogram.StackPointer

	for c := cxprogram.CallCounter; c >= 0; c-- {
		op := cxprogram.CallStack[c].Operator
		fp -= op.Size

		var dupNames []string

		fmt.Printf(">>> %s()\n", op.Name)

		for _, inp := range op.Inputs {
			fmt.Println("ProgramInput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), op.Name, ast.GetPrintableValue(fp, inp))

			dupNames = append(dupNames, inp.Package.Name+inp.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("ProgramOutput")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), op.Name, ast.GetPrintableValue(fp, out))

			dupNames = append(dupNames, out.Package.Name+out.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
		for _, expr := range op.Expressions {
			for _, inp := range expr.Inputs {
				if inp.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == inp.Package.Name+inp.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", inp.Name, "\t", ":", "\t", GetPrintableValue(fp, inp))
				// exprs += fmt.Sprintln("\t", stackValueHeader(inp.FileName, inp.FileLine), "\t", ":", "\t", GetPrintableValue(fp, inp))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), ExprOpName(expr), ast.GetPrintableValue(fp, inp))

				dupNames = append(dupNames, inp.Package.Name+inp.Name)
			}

			for _, out := range expr.Outputs {
				if out.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == out.Package.Name+out.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", out.Name, "\t", ":", "\t", GetPrintableValue(fp, out))
				// exprs += fmt.Sprintln("\t", stackValueHeader(out.FileName, out.FileLine), ":", GetPrintableValue(fp, out))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), ExprOpName(expr), ast.GetPrintableValue(fp, out))

				dupNames = append(dupNames, out.Package.Name+out.Name)
			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
	}
}

// IsTempVar ...
func IsTempVar(name string) bool {
	if len(name) >= len(constants.LOCAL_PREFIX) && name[:len(constants.LOCAL_PREFIX)] == constants.LOCAL_PREFIX {
		return true
	}
	return false
}

func checkForEscapedChars(str string) []byte {
	var res []byte
	var lenStr = int(len(str))
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr-1 {
			nextCh = str[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}

		} else {
			res = append(res, ch)
		}
	}

	return res
}

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

// GetPointerOffset ...
func GetPointerOffset(pointer int32) int32 {
	return helper.Deserialize_i32(ast.PROGRAM.Memory[pointer : pointer+constants.TYPE_POINTER_SIZE])
}

// GetSliceOffset ...
func GetSliceOffset(fp int, arg *ast.CXArgument) int32 {
	element := ast.GetAssignmentElement(arg)
	if element.IsSlice {
		return GetPointerOffset(int32(GetFinalOffset(fp, arg)))
	}

	return -1
}

// GetObjectHeader ...
func GetObjectHeader(offset int32) []byte {
	return ast.PROGRAM.Memory[offset : offset+constants.OBJECT_HEADER_SIZE]
}

// GetSliceHeader ...
func GetSliceHeader(offset int32) []byte {
	return ast.PROGRAM.Memory[offset+constants.OBJECT_HEADER_SIZE : offset+constants.OBJECT_HEADER_SIZE+constants.SLICE_HEADER_SIZE]
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
			return ast.PROGRAM.Memory[dataOffset : dataOffset+dataLen]
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
func SliceResize(fp int, out *ast.CXArgument, inp *ast.CXArgument, count int32, sizeofElement int) int {
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
func SliceCopy(fp int, outputSliceOffset int32, inp *ast.CXArgument, count int32, sizeofElement int) {
	inputSliceOffset := GetSliceOffset(fp, inp)
	SliceCopyEx(outputSliceOffset, inputSliceOffset, count, sizeofElement)
}

// SliceAppendResize prepares a slice to be able to store a new object of length `sizeofElement`. It checks if the slice needs to be relocated in memory, and if it is needed it relocates it and a new `outputSliceOffset` is calculated for the new slice.
func SliceAppendResize(fp int, out *ast.CXArgument, inp *ast.CXArgument, sizeofElement int) int32 {
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
func SliceInsert(fp int, out *ast.CXArgument, inp *ast.CXArgument, index int32, object []byte) int {
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
func SliceRemove(fp int, out *ast.CXArgument, inp *ast.CXArgument, index int32, sizeofElement int32) int {
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

