package base

import "fmt"


func op_i8_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 :=	expr.Inputs[0]
	fmt.Println(ReadI8(stack, fp, inp1))
}

func op_i8_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) +  ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) -  ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_mul (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) *  ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_div (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) /  ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_gt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) > ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_gteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) >= ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_lt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) < ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_lteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) <= ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_eq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) == ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_uneq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1:= FromBool(ReadI8(stack, fp, inp1) != ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_bitand (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) & ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_bitor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) | ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_bitxor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) ^ ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i8_bitclear (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(stack, fp, inp1) &^ ReadI8(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
