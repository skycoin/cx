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

	out1 := expr.Outputs[0]
	outB1 := FromI64(makeTimestamp())
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_time_UnixNano(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano())
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_time_Sleep(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	time.Sleep(time.Duration(ReadI32(fp, inp1)) * time.Millisecond)
}
