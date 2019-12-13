// +build base

package cxcore

import (
	"bytes"
	//"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

var openFiles map[string]*os.File = make(map[string]*os.File, 0)

func op_os_ReadFile(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	_ = out1

	if byts, err := ioutil.ReadFile(ReadStr(fp, inp1)); err == nil {
		_ = byts
		// sByts := encoder.Serialize(byts)
		// assignOutput(0, sByts, "[]byte", expr, call)
	} else {
		panic(err)
	}
}

func op_os_Open(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, err := os.Open(name); err == nil {
		openFiles[name] = file
	} else {
		panic(err)
	}
}

func op_os_Close(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, ok := openFiles[name]; ok {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
}

func op_os_GetWorkingDirectory(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	byts := encoder.Serialize(PROGRAM.Path)
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), byts)
}

func op_os_Exit(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	exitCode := ReadI32(fp, expr.Inputs[0])
	os.Exit(int(exitCode))
}

func op_os_Run(prgrm *CXProgram) {
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

	//fmt.Println("COMMAND : ", name, " ARGS : ", args)
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
