// +build base

package cxcore

import (
	"bytes"
	"encoding/binary"

	. "github.com/SkycoinProject/cx/cx"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

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

	var success bool
	if byts, err := CXReadFile(ReadStr(fp, expr.Inputs[0])); err == nil {
		WriteObject(GetFinalOffset(fp, expr.Outputs[0]), encoder.Serialize(string(byts)))
		success = true
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsOpen(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	handle := int32(-1)

	if file, err := CXOpenFile(ReadStr(fp, expr.Inputs[0])); err == nil {
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

func opOsReadF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value float32
	var success bool

	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadF32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var success bool

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	count := ReadI32(fp, expr.Inputs[1])
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]float32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, 4))
				outputSliceData := GetSliceData(outputSliceOffset, 4)
				for i := int32(0); i < count; i++ {
					WriteMemF32(outputSliceData, int(i*4), values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)

}

func opOsReadUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint32
	var success bool

	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI32Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var success bool

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	count := ReadI32(fp, expr.Inputs[1])
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint32, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, 4))
				outputSliceData := GetSliceData(outputSliceOffset, 4)
				for i := int32(0); i < count; i++ {
					WriteMemUI32(outputSliceData, int(i*4), values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)

}

func opOsReadUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var value uint16
	var success bool

	if file := validFileFromExpr(expr, fp); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}

	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), value)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

func opOsReadUI16Slice(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var success bool

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	count := ReadI32(fp, expr.Inputs[1])
	if count > 0 {
		if file := validFileFromExpr(expr, fp); file != nil {
			values := make([]uint16, count)
			if err := binary.Read(file, binary.LittleEndian, values); err == nil {
				success = true
				outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, count, 2))
				outputSliceData := GetSliceData(outputSliceOffset, 2)
				for i := int32(0); i < count; i++ {
					WriteMemUI16(outputSliceData, int(i*2), values[i])
				}
			}
		}
	}

	WriteI32(outputSlicePointer, outputSliceOffset)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)

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
