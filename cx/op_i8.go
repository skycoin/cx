package base

import "fmt"


func op_i8_print (expr *CXExpression, fp int) {
	inp1 :=	expr.Inputs[0]
	fmt.Println(ReadI8(fp, inp1))
}

func op_i8_add (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) +  ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_sub (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) -  ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_mul (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) *  ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_div (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) /  ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_gt (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) > ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_gteq (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) >= ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_lt (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) < ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_lteq (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) <= ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_eq (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) == ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_uneq (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(fp, inp1) != ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_bitand (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) & ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_bitor (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) | ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_bitxor (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) ^ ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i8_bitclear (expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) &^ ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
