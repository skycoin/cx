package base

import (
	"fmt"
)

func op_str_print(expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(stack, fp, inp1))
}

func op_str_eq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(stack, fp, inp1) == ReadStr(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
