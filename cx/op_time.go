// +build base

package cxcore

import (
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func op_time_UnixMilli(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), makeTimestamp())
}

func op_time_UnixNano(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), time.Now().UnixNano())
}

func op_time_Sleep(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	time.Sleep(time.Duration(ReadI32(fp, expr.Inputs[0])) * time.Millisecond)
}
