// +build base extra full

package cxcore

import (
	"bytes"
	"encoding/binary"
	//"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const (
	OS_RUN_SUCCESS = iota
	OS_RUN_EMPTY_CMD
	OS_RUN_PANIC // 2
	OS_RUN_START_FAILED
	OS_RUN_WAIT_FAILED
	OS_RUN_TIMEOUT
)

const (
	OS_SEEK_SET = iota
	OS_SEEK_CUR
	OS_SEEK_END
)

var openFiles map[string]*os.File = make(map[string]*os.File, 0)

func op_os_ReadAllText(expr *CXExpression, fp int) {
	if byts, err := ioutil.ReadFile(ReadStr(fp, expr.Inputs[0])); err == nil {
		WriteObject(GetFinalOffset(fp, expr.Outputs[0]), encoder.Serialize(string(byts)))
	} else {
		panic(err)
	}
}

func op_os_Open(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, err := os.Open(name); err == nil {
		openFiles[name] = file
	} else {
		panic(err)
	}
}

func op_os_Close(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, ok := openFiles[name]; ok {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
}

func op_os_Seek(expr *CXExpression, fp int) {
	file := openFiles[ReadStr(fp, expr.Inputs[0])]
	file.Seek(ReadI64(fp, expr.Inputs[1]), int(ReadI32(fp, expr.Inputs[2])))
}

func op_os_ReadF32(expr *CXExpression, fp int) {
	file := openFiles[ReadStr(fp, expr.Inputs[0])]
	var value float32
	err := binary.Read(file, binary.LittleEndian, &value)
	if err != nil {
		panic(err)
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF32(value))
}

func op_os_ReadUI32(expr *CXExpression, fp int) {
	file := openFiles[ReadStr(fp, expr.Inputs[0])]
	var value uint32
	err := binary.Read(file, binary.LittleEndian, &value)
	if err != nil {
		panic(err)
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI32(value))
}

func op_os_ReadUI16(expr *CXExpression, fp int) {
	file := openFiles[ReadStr(fp, expr.Inputs[0])]
	var value uint16
	err := binary.Read(file, binary.LittleEndian, &value)
	if err != nil {
		panic(err)
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI16(value))
}

func op_os_GetWorkingDirectory(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	byts := encoder.Serialize(PROGRAM.Path)
	WriteObject(out1Offset, byts)
}

func op_os_Exit(expr *CXExpression, fp int) {
	inp0 := expr.Inputs[0]
	exitCode := ReadI32(fp, inp0)
	os.Exit(int(exitCode))
}

func op_os_Run(expr *CXExpression, fp int) {
	inp0, inp1, inp2, inp3, out0, out1, out2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Outputs[0], expr.Outputs[1], expr.Outputs[2]
	var runError int32 = OS_RUN_SUCCESS

	command := ReadStr(fp, inp0)
	dir := ReadStr(fp, inp3)
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

	//fmt.Println("COMMAND : ", name, " ARGS : ", args)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	var cmdError int32 = 0

	timeoutMs := ReadI32(fp, inp2)
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
	maxSize := ReadI32(fp, inp1)
	if (maxSize > 0) && (len(stdOutBytes) > int(maxSize)) {
		stdOutBytes = stdOutBytes[0:maxSize]
	}

	WriteMemory(GetFinalOffset(fp, out0), FromI32(runError))
	WriteMemory(GetFinalOffset(fp, out1), FromI32(cmdError))
	WriteObject(GetFinalOffset(fp, out2), FromStr(string(stdOutBytes)))
}
