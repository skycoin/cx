package base

import "fmt"


func op_i8_print (expr *CXExpression, mem []byte, fp int) {
	inp1 :=	expr.Inputs[0]
	fmt.Println(ReadI8(mem, fp, inp1))
}

func op_i8_add (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) +  ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_sub (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) -  ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_mul (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) *  ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_div (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) /  ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_gt (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) > ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_gteq (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) >= ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_lt (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) < ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_lteq (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) <= ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_eq (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) == ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_uneq (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(mem, fp, inp1) != ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_bitand (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) & ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_bitor (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) | ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_bitxor (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) ^ ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}

func op_i8_bitclear (expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(mem, fp, inp1) &^ ReadI8(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}
