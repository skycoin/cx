package base

import (
	"fmt"
)

func bool_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(stack, fp, inp1))
}

func bool_not (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(!ReadBool(stack, fp, inp1))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func bool_and (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(stack, fp, inp1) && ReadBool(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func bool_or (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(stack, fp, inp1) || ReadBool(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}
