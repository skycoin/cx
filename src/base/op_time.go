package base

import (
	"time"
)

func time_UnixMilli (expr *CXExpression, stack *CXStack, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano() / int64(1000000))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}


