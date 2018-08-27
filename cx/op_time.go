package base

import (
	"time"
)

func op_time_UnixMilli(expr *CXExpression, mem []byte, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano() / int64(1000000))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_time_UnixNano(expr *CXExpression, mem []byte, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano())
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_time_Sleep(expr *CXExpression, mem []byte, fp int) {
	inp1 := expr.Inputs[0]
	time.Sleep(time.Duration(ReadI32(mem, fp, inp1)) * time.Millisecond)
}
