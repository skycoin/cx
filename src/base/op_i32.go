package base

import (
	"fmt"
	"math"
	"math/rand"
)

func i32_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

func i32_add (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_sub (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_mul (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) * ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_div (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) / ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_abs (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Abs(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_pow (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(math.Pow(float64(ReadI32(stack, fp, inp1)), float64(ReadI32(stack, fp, inp2)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_gt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_gteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) >= ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_lt (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_lteq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) <= ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_eq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) == ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_uneq (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) != ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_mod (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) % ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_rand (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI32(stack, fp, inp1)
	maximum := ReadI32(stack, fp, inp2)
	
	outB1 := FromI32(int32(rand.Intn(int(maximum - minimum)) + int(minimum)))

	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitand (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) & ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) | ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitxor (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) ^ ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitclear (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) &^ ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitshl (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(stack, fp, inp1)) << uint32(ReadI32(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}

func i32_bitshr (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(stack, fp, inp1)) >> uint32(ReadI32(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1), out1, outB1)
}
