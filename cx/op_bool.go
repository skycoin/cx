package base

import (
	"fmt"
)

func op_bool_print(expr *CXExpression, mem []byte, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(mem, fp, inp1))
}

func op_bool_equal(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(mem, fp, inp1) == ReadBool(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bool_unequal(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(mem, fp, inp1) != ReadBool(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bool_not(expr *CXExpression, mem []byte, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(!ReadBool(mem, fp, inp1))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bool_and(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(mem, fp, inp1) && ReadBool(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_bool_or(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(mem, fp, inp1) || ReadBool(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}
