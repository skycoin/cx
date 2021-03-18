// +build os

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

    . "github.com/skycoin/cx/cx"
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

func opOsLogFile(inputs []CXValue, outputs []CXValue) {
	CXLogFile(inputs[0].Get_bool())
}

func opOsReadAllText(inputs []CXValue, outputs []CXValue) {
	success := false

	if byts, err := CXReadFile(inputs[0].Get_str()); err == nil {
        outputs[0].Set_str(string(byts))
		success = true
	}

	outputs[1].Set_bool(success)
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

func opOsOpen(inputs []CXValue, outputs []CXValue) {
	handle := int32(-1)
	if file, err := CXOpenFile(inputs[0].Get_str()); err == nil {
		handle = getFileHandle(file)
	}

	outputs[0].Set_i32(int32(handle))
}

func opOsCreate(inputs []CXValue, outputs []CXValue) {
	handle := int32(-1)
	if file, err := CXCreateFile(inputs[0].Get_str()); err == nil {
		handle = getFileHandle(file)
	}

	outputs[0].Set_i32(int32(handle))
}

func opOsClose(inputs []CXValue, outputs []CXValue) {
	success := false

	handle := inputs[0].Get_i32()
	if file := ValidFile(handle); file != nil {
		if err := file.Close(); err == nil {
			success = true
		}

		openFiles[handle] = nil
		freeFiles = append(freeFiles, handle)
	}

	outputs[0].Set_bool(success)
}

func opOsSeek(inputs []CXValue, outputs []CXValue) {
	offset := int64(-1)
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		var err error
		if offset, err = file.Seek(inputs[1].Get_i64(), int(inputs[2].Get_i32())); err != nil {
			offset = -1
		}
	}
	outputs[0].Set_i64(offset)
}

func opOsReadStr(inputs []CXValue, outputs []CXValue) {
	var len uint64
	var value string
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &len); err == nil {
			bytes := make([]byte, len)
			if err := binary.Read(file, binary.LittleEndian, &bytes); err == nil {
				value = string(bytes)
				success = true
			}
		}
	}
    outputs[0].Set_str(value)
	outputs[1].Set_bool(success)
}

func opOsReadF64(inputs []CXValue, outputs []CXValue) {
	var value float64
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_f64(value)
	outputs[1].Set_bool(success)
}

func opOsReadF32(inputs []CXValue, outputs []CXValue) {
	var value float32
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_f32(value)
	outputs[1].Set_bool(success)
}

func opOsReadUI64(inputs []CXValue, outputs []CXValue) {
	var value uint64
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_ui64(value)
	outputs[1].Set_bool(success)
}

func opOsReadUI32(inputs []CXValue, outputs []CXValue) {
	var value uint32
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_ui32(value)
	outputs[1].Set_bool(success)
}

func opOsReadUI16(inputs []CXValue, outputs []CXValue) {
	var value uint16
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_ui16(value)
	outputs[1].Set_bool(success)
}

func opOsReadUI8(inputs []CXValue, outputs []CXValue) {
	var value uint8
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_ui8(value)
	outputs[1].Set_bool(success)
}

func opOsReadI64(inputs []CXValue, outputs []CXValue) {
	var value int64
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i64(value)
	outputs[1].Set_bool(success)
}

func opOsReadI32(inputs []CXValue, outputs []CXValue) {
	var value int32
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_i32(value)
	outputs[1].Set_bool(success)
}

func opOsReadI16(inputs []CXValue, outputs []CXValue) {
	var value int16
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_i16(value)
	outputs[1].Set_bool(success)
}

func opOsReadI8(inputs []CXValue, outputs []CXValue) {
	var value int8
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

    outputs[0].Set_i8(value)
	outputs[1].Set_bool(success)
}

func opOsReadBOOL(inputs []CXValue, outputs []CXValue) {
	var value bool
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(value)
	outputs[1].Set_bool(success)
}

func opOsWriteStr(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		value := inputs[1].Get_str()
		len := len(value)
		if err := binary.Write(file, binary.LittleEndian, uint64(len)); err == nil {
			if err := binary.Write(file, binary.LittleEndian, []byte(value)); err == nil {
				success = true
			}
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteF64(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_f64()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteF32(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_f32()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteUI64(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui64()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteUI32(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i32()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteUI16(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui16()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteUI8(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_ui8()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteI64(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i64()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteI32(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i32()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteI16(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i16()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteI8(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_i8()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func opOsWriteBOOL(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if err := binary.Write(file, binary.LittleEndian, inputs[1].Get_bool()); err == nil {
			success = true
		}
	}

	outputs[0].Set_bool(success)
}

func getSlice(inputs []CXValue, outputs []CXValue) (outputSlicePointer int, outputSliceOffset int32, sizeofElement int, count uint64) {
	inp1, out0 := inputs[1].Arg, outputs[0].Arg
    
	if inp1.Type != out0.Type || !GetAssignmentElement(inp1).IsSlice || !GetAssignmentElement(out0).IsSlice {
		panic(CX_RUNTIME_INVALID_ARGUMENT)
	}
    inputs[1].Used = int8(inp1.Type)
	count = inputs[2].Get_ui64()
	outputSlicePointer = outputs[0].Offset
	sizeofElement = GetAssignmentElement(inp1).Size
	outputSliceOffset = int32(SliceResize(outputs[0].FramePointer, out0, inp1, int32(count), sizeofElement))
	return
}

func opOsReadF64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
    _, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadF32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

    outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadUI64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadUI32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadUI16Slice(inputs []CXValue, outputs []CXValue) {
	success := false
    _, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadUI8Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadI64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
    _, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadI32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadI16Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsReadI8Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	_, outputSliceOffset, sizeofElement, count := getSlice(inputs, outputs)
	if count > 0 {
		if file := ValidFile(inputs[0].Get_i32()); file != nil {
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

	outputs[0].SetSlice(outputSliceOffset)
	outputs[1].Set_bool(success)
}

func opOsWriteF64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_f64(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteF32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_f32(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteUI64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_f64(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteUI32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_ui32(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteUI16Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_ui16(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteUI8Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_ui8(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteI64Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_i64(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteI32Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_i32(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteI16Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_i16(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsWriteI8Slice(inputs []CXValue, outputs []CXValue) {
	success := false
	if file := ValidFile(inputs[0].Get_i32()); file != nil {
		if data := inputs[1].GetSlice_i8(); data != nil {
			if err := binary.Write(file, binary.LittleEndian, data); err == nil {
				success = true
			}
		}
	}
	outputs[0].Set_bool(success)
}

func opOsGetWorkingDirectory(inputs []CXValue, outputs []CXValue) {
    outputs[0].Set_str(PROGRAM.Path)
}

func opOsExit(inputs []CXValue, outputs []CXValue) {
	exitCode := inputs[0].Get_i32()
	os.Exit(int(exitCode))
}

func opOsRun(inputs []CXValue, outputs []CXValue) {
	var runError int32 = OS_RUN_SUCCESS

	command := inputs[0].Get_str()
	dir := inputs[3].Get_str()
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

	timeoutMs := inputs[2].Get_i32()
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
	maxSize := inputs[1].Get_i32()
	if (maxSize > 0) && (len(stdOutBytes) > int(maxSize)) {
		stdOutBytes = stdOutBytes[0:maxSize]
	}

	outputs[0].Set_i32(runError)
	outputs[1].Set_i32(cmdError)
    outputs[2].Set_str(string(stdOutBytes))
}
