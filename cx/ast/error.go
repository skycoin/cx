package ast

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
)

// ErrorHeader ...
func ErrorHeader(currentFile string, lineNo int) string {
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

// CompilationError is a helper function that concatenates the `currentFile` and `lineNo` data to a error header and returns the full error string.
func CompilationError(currentFile string, lineNo int) string {
	globals.FoundCompileErrors = true
	return ErrorHeader(currentFile, lineNo)
}

// ErrorString ...
func ErrorString(code int) string {
	if str, found := constants.ErrorStrings[code]; found {
		return str
	}
	return constants.ErrorStrings[constants.CX_RUNTIME_ERROR]
}

func errorCode(r interface{}) int {
	switch v := r.(type) {
	case int:
		return int(v)
	default:
		return constants.CX_RUNTIME_ERROR
	}
}

func RuntimeErrorInfo(prgrm *CXProgram, r interface{}, printStack bool, defaultError int) {
	call := prgrm.CallStack[prgrm.CallCounter]
	expr := call.Operator.Expressions[call.Line]
	code := errorCode(r)
	if code == constants.CX_RUNTIME_ERROR {
		code = defaultError
	}

	fmt.Printf("%s, %s, %v", ErrorHeader(expr.FileName, expr.FileLine), ErrorString(code), r)

	if printStack {
		prgrm.PrintStack()
	}

	if globals.DBG_GOLANG_STACK_TRACE {
		debug.PrintStack()
	}

	os.Exit(code)
}

// RuntimeError ...
func RuntimeError(prgrm *CXProgram) {
	if r := recover(); r != nil {
		switch r {
		case constants.STACK_OVERFLOW_ERROR:
			call := prgrm.CallStack[prgrm.CallCounter]
			if prgrm.CallCounter > 0 {
				prgrm.CallCounter--
				prgrm.Stack.Pointer = call.FramePointer
				RuntimeErrorInfo(prgrm, r, true, constants.CX_RUNTIME_STACK_OVERFLOW_ERROR)
			} else {
				// error at entry point
				RuntimeErrorInfo(prgrm, r, false, constants.CX_RUNTIME_STACK_OVERFLOW_ERROR)
			}
		case constants.HEAP_EXHAUSTED_ERROR:
			RuntimeErrorInfo(prgrm, r, true, constants.CX_RUNTIME_HEAP_EXHAUSTED_ERROR)
		default:
			RuntimeErrorInfo(prgrm, r, true, constants.CX_RUNTIME_ERROR)
		}
	}
}
