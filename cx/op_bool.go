package base

import (
	"fmt"
)

func op_bool_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(fp, inp1))
}

func op_bool_equal(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) == ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bool_unequal(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) != ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bool_not(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(!ReadBool(fp, inp1))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bool_and(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) && ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_bool_or(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) || ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
