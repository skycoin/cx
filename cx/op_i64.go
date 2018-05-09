package base

import (
	"fmt"
	"math"
	"math/rand"
)

func op_i64_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI64(stack, fp, inp1))
}

func op_i64_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) + ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) - ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_mul (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) * ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_div (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) / ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_abs (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Abs(float64(ReadI64(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_pow (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Pow(float64(ReadI64(stack, fp, inp1)), float64(ReadI64(stack, fp, inp2)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_gt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) > ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_gteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) >= ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_lt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) < ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_lteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) <= ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_eq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) == ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_uneq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(stack, fp, inp1) != ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_mod (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) % ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_rand (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI64(stack, fp, inp1)
	maximum := ReadI64(stack, fp, inp2)
	
	outB1 := FromI64(int64(rand.Intn(int(maximum - minimum)) + int(minimum)))

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitand (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) & ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) | ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitxor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) ^ ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitclear (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(stack, fp, inp1) &^ ReadI64(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitshl (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(stack, fp, inp1)) << uint64(ReadI64(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i64_bitshr (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(stack, fp, inp1)) >> uint64(ReadI64(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
