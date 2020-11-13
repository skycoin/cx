// +build base

package cxcore

import (
	"bytes"
	"encoding/binary"

	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	. "github.com/SkycoinProject/cx/cx"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

const (
	OS_SEEK_SET = iota
	OS_SEEK_CUR
	OS_SEEK_END
)

var openFiles []*os.File
var freeFiles []int32

// helper function used to validate json handle from expr
func validFileFromExpr(expr *CXExpression, fp int) *os.File {
	handle := ReadI32(fp, expr.Inputs[0])
	return ValidFile(handle)
}

// helper function used to validate file handle from i32
func ValidFile(handle int32) *os.File {
	if handle >= 0 && handle < int32(len(openFiles)) && openFiles[handle] != nil {
		return openFiles[handle]
	}
	return nil
}

func opOsLogFile(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	CXLogFile(ReadBool(fp, expr.Inputs[0]))
}

func opOsReadAllText(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false

	if byts, err := CXReadFile(ReadStr(fp, expr.Inputs[0])); err == nil {
		WriteObject(GetFinalOffset(fp, expr.Outputs[0]), encoder.Serialize(string(byts)))
		success = true
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func getFileHandle(file *os.File) (int32) {
	handle := int32(-1)
	freeCount := len(freeFiles)
	if freeCount > 0 {
		freeCount--
		handle = int32(freeFiles[freeCount])
		freeFiles = freeFiles[:freeCount]
	} else {
		handle = int32(len(openFiles))
		openFiles = append(openFiles, nil)
	}

	if handle < 0 || handle >= int32(len(openFiles)) {
		panic("internal error")
	}

	openFiles[handle] = file
	return handle
}

func opOsOpen(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	handle := int32(-1)
	if file, err := CXOpenFile(ReadStr(fp, expr.Inputs[0])); err == nil {
		handle = getFileHandle(file)
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(handle))
}

func opOsCreate(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	handle := int32(-1)
	if file, err := CXCreateFile(ReadStr(fp, expr.Inputs[0])); err == nil {
		handle = getFileHandle(file)
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(handle))
}

func opOsClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false

	handle := ReadI32(fp, expr.Inputs[0])
	if file := ValidFile(handle); file != nil {
		if err := file.Close(); err == nil {
			success = true
		}

		openFiles[handle] = nil
		freeFiles = append(freeFiles, handle)
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsSeek(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	offset := int64(-1)
	if file := validFileFromExpr(expr, fp); file != nil {
		var err error
		if offset, err = file.Seek(ReadI64(fp, expr.Inputs[1]), int(ReadI32(fp, expr.Inputs[2]))); err != nil {
			offset = -1
		}
	}
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), offset)
}

func opOsReadStr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var len uint64
	var value string
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &len); err == nil {
			bytes := make([]byte, len)
			if err := binary.Read(file, binary.LittleEndian, &bytes); err == nil {
				value = string(bytes)
				success = true
			}
		}
	}
	WriteString(fp, value, expr.Outputs[0])
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value float64
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value float32
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint64
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint32
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint16
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint8
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI8(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value int64
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value int32
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value int16
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteI16(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value int8
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteI8(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadBOOL(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value bool
	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)	
}

func opOsWriteStr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		value := ReadStr(fp, expr.Inputs[1])
		len := len(value)
		if err := binary.Write(file, binary.LittleEndian, uint64(len)); err == nil {
			if err := binary.Write(file, binary.LittleEndian, []byte(value)); err == nil {
				success = true
			}
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadF64(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadF32(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadUI64(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadUI32(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadUI16(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadUI8(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadI64(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadI32(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadI16(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadI8(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteBOOL(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Write(file, binary.LittleEndian, ReadBool(fp, expr.Inputs[1])); err == nil {
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func getSlice(expr *CXExpression, fp int) (outputSlicePointer int, outputSliceOffset int32, sizeofElement int, count uint64) {
	inp1, out0 := expr.Inputs[1], expr.Outputs[0]
	if inp1.Type != out0.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out0).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
	count = ReadUI64(fp, expr.Inputs[2])
	outputSlicePointer = GetFinalOffset(fp, out0)
	sizeofElement = GetAssignmentElement(inp1).Size
	outputSliceOffset = int32(SliceResize(fp, out0, inp1, int32(count), sizeofElement))
	return
}

func opOsReadF64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]float64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemF64(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadF32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]float32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemF32(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemUI64(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemUI32(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI16Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint16, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemUI16(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI8Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint8, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemUI8(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]int64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemI64(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]int32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemI32(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI16Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]int16, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemI16(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadI8Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	outputSlicePointer, outputSliceOffset, sizeofElement, count := getSlice(expr, fp)
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]int8, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
				for i := uint64(0); i < count; i++ {
					WriteMemI8(outputSliceData, int(i)*sizeofElement, values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsWriteF64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_F64); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteF32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_F32); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_UI64); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_UI32); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI16Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_UI16); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteUI8Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_UI8); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI64Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_I64); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_I32); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI16Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_I16); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsWriteI8Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false
	if file := validFileFromExpr(expr, fp); file != nil {
		if data := ReadData(fp, expr.Inputs[1], TYPE_I8); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

func opOsGetWorkingDirectory(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	byts := encoder.Serialize(PROGRAM.Path)
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), byts)
}

func opOsExit(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	exitCode := ReadI32(fp, expr.Inputs[0])
	os.Exit(int(exitCode))
}

func opOsRun(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var runError int32 = OS_RUN_SUCCESS

	command := ReadStr(fp, expr.Inputs[0])
	dir := ReadStr(fp, expr.Inputs[3])
	args := strings.Split(command, " ")
	if len(args) <= 0 {
		runError = OS_RUN_EMPTY_CMD
	}

	name := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	var cmdError int32 = 0

	timeoutMs := ReadI32(fp, expr.Inputs[2])
	timeout := time.Duration(math.MaxInt64)
	if timeoutMs > 0 {
		timeout = time.Duration(timeoutMs) * time.Millisecond
	}

	if err := cmd.Start(); err != nil {
		runError = OS_RUN_START_FAILED
	} else {
		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		select {
		case <-time.After(timeout):
			cmd.Process.Kill()
			runError = OS_RUN_TIMEOUT
		case err := <-done:
			if err != nil {
				if exiterr, ok := err.(*exec.ExitError); ok {
					// from stackoverflow
					// The program has exited with an exit code != 0
					// This works on both Unix and Windows. Although package
					// syscall is generally platform dependent, WaitStatus is
					// defined for both Unix and Windows and in both cases has
					// an ExitStatus() method with the same signature.
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						cmdError = int32(status.ExitStatus())
					}
				} else {
					runError = OS_RUN_WAIT_FAILED
				}
			}
		}
	}

	stdOutBytes := out.Bytes()
	maxSize := ReadI32(fp, expr.Inputs[1])
	if (maxSize > 0) && (len(stdOutBytes) > int(maxSize)) {
		stdOutBytes = stdOutBytes[0:maxSize]
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), runError)
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), cmdError)
	WriteObject(GetFinalOffset(fp, expr.Outputs[2]), FromStr(string(stdOutBytes)))
}
