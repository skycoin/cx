// +build base

package cxcore

import (
	. "github.com/SkycoinProject/cx/cx"
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), makeTimestamp())
}

func opTimeUnixNano(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), time.Now().UnixNano())
}

func opTimeSleep(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	time.Sleep(time.Duration(ReadI32(fp, expr.Inputs[0])) * time.Millisecond)
}
