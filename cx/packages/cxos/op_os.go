// +build cxos

package cxos

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cx/util"
)

const (
	OS_SEEK_SET = iota
	OS_SEEK_CUR
	OS_SEEK_END
)

var openFiles []*os.File
var freeFiles []int32

// helper function used to validate file handle from i32
func ValidFile(handle int32) *os.File {
	if handle >= 0 && handle < int32(len(openFiles)) && openFiles[handle] != nil {
		return openFiles[handle]
	}
	return nil
}

func opOsLogFile(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	util.CXLogFile(inputs[0].Get_bool(prgrm))
}

func opOsReadAllText(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false

	if byts, err := util.CXReadFile(inputs[0].Get_str(prgrm)); err == nil {
		outputs[0].Set_str(prgrm, string(byts))
		success = true
	}

	outputs[1].Set_bool(prgrm, success)
}

func getFileHandle(file *os.File) int32 {
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

func opOsOpen(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	handle := int32(-1)
	if file, err := util.CXOpenFile(inputs[0].Get_str(prgrm)); err == nil {
		handle = getFileHandle(file)
	}

	outputs[0].Set_i32(prgrm, int32(handle))
}

func opOsCreate(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	handle := int32(-1)
	if file, err := util.CXCreateFile(inputs[0].Get_str(prgrm)); err == nil {
		handle = getFileHandle(file)
	}

	outputs[0].Set_i32(prgrm, int32(handle))
}

func opOsClose(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false

	handle := inputs[0].Get_i32(prgrm)
	if file := ValidFile(handle); file != nil {
		if err := file.Close(); err == nil {
			success = true
		}

		openFiles[handle] = nil
		freeFiles = append(freeFiles, handle)
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsSeek(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	offset := int64(-1)
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		var err error
		if offset, err = file.Seek(inputs[1].Get_i64(prgrm), int(inputs[2].Get_i32(prgrm))); err != nil {
			offset = -1
		}
	}
	outputs[0].Set_i64(prgrm, offset)
}

func opOsReadStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var len uint64
	var value string
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &len); err == nil {
			bytes := make([]byte, len)
			if err := binary.Read(file, binary.LittleEndian, &bytes); err == nil {
				value = string(bytes)
				success = true
			}
		}
	}
	outputs[0].Set_str(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value float64
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_f64(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value float32
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_f32(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value uint64
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_ui64(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value uint32
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_ui32(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value uint16
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_ui16(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value uint8
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_ui8(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value int64
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i64(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value int32
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i32(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value int16
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i16(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value int8
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i8(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadBOOL(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var value bool
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, value)
	outputs[1].Set_bool(prgrm, success)
}

func opOsWriteStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		value := inputs[1].Get_str(prgrm)
		len := len(value)
		if err := binary.Write(file, binary.LittleEndian, uint64(len)); err == nil {
			if err := binary.Write(file, binary.LittleEndian, []byte(value)); err == nil {
				success = true
			}
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_f64(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_f32(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui64(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i32(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui16(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui8(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i64(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i32(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i16(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i8(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteBOOL(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_bool(prgrm)); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(prgrm, success)
}

func getSlice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) (outputSlicePointer types.Pointer, outputSliceOffset types.Pointer, sizeofElement types.Pointer, count types.Pointer) {
	var inp1, out0 *ast.CXArgument

	var inpType, outType types.Code
	var inpIsSlice, outIsSlice bool
	var inpSize types.Pointer
	if inputs[1].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		inp1 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inputs[1].TypeSignature.Meta))

		inpType = inp1.Type
		inpIsSlice = (inp1.GetAssignmentElement(prgrm)).IsSlice
		inpSize = (inp1.GetAssignmentElement(prgrm)).Size
	} else if inputs[1].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(inputs[1].TypeSignature.Meta)
		inpType = types.Code(sliceDetails.Type)
		inpIsSlice = true
		inpSize = inpType.Size()
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	if outputs[0].TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
		out0 = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(outputs[0].TypeSignature.Meta))

		outType = out0.Type
		outIsSlice = (out0.GetAssignmentElement(prgrm)).IsSlice
	} else if outputs[0].TypeSignature.Type == ast.TYPE_SLICE_ATOMIC {
		sliceDetails := prgrm.GetCXTypeSignatureArrayFromArray(outputs[0].TypeSignature.Meta)

		outType = types.Code(sliceDetails.Type)
		outIsSlice = true
	} else {
		panic("type is not type cx argument deprecate\n\n")
	}

	if inpType != outType || !inpIsSlice || !outIsSlice {
		panic(constants.CX_RUNTIME_INVALID_ARGUMENT)
	}

	count = inputs[2].Get_ptr(prgrm)
	outputSlicePointer = outputs[0].Offset
	sizeofElement = inpSize
	outputSliceOffset = ast.SliceResize(prgrm, outputs[0].FramePointer, outputs[0].TypeSignature, inputs[1].TypeSignature, count, sizeofElement)
	return
}

func opOsReadF64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]float64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_f64(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadF32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]float32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_f32(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]uint64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_ui64(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]uint32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_ui32(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI16Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]uint16, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_ui16(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadUI8Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]uint8, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_ui8(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]int64, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_i64(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]int32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_i32(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI16Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]int16, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_i16(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsReadI8Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(prgrm, inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
			values := make([]int8, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceData := ast.GetSliceData(prgrm, outputSliceOffset, sizeofElement)
				for i := types.Pointer(0); i < count; i++ {
					types.Write_i8(outputSliceData, i*sizeofElement, values[i])
				}
			}
		}
	}

	outputs[0].Set_ptr(prgrm, outputSliceOffset)
	outputs[1].Set_bool(prgrm, success)
}

func opOsWriteF64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_f64(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteF32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_f32(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_f64(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_ui32(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI16Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_ui16(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteUI8Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_ui8(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI64Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_i64(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI32Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_i32(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI16Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_i16(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsWriteI8Slice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if data := inputs[1].GetSlice_i8(prgrm); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(prgrm, success)
}

func opOsGetWorkingDirectory(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(prgrm, globals.CxProgramPath)
}

func opOsExit(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	exitCode := inputs[0].Get_i32(prgrm)
	os.Exit(int(exitCode))
}

func opOsRun(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var runError int32 = OS_RUN_SUCCESS

	command := inputs[0].Get_str(prgrm)
	dir := inputs[3].Get_str(prgrm)
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

	timeoutMs := inputs[2].Get_i32(prgrm)
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
	maxSize := inputs[1].Get_i32(prgrm)
	if (maxSize > 0) && (len(stdOutBytes) > int(maxSize)) {
		stdOutBytes = stdOutBytes[0:maxSize]
	}

	outputs[0].Set_i32(prgrm, runError)
	outputs[1].Set_i32(prgrm, cmdError)
	outputs[2].Set_str(prgrm, string(stdOutBytes))
}
