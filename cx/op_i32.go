package base

import (
	"fmt"
	// "strconv"
	"math"
	"math/rand"
)

func op_i32_i32(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)
	switch out1.Type {
	case TYPE_BYTE:
		WriteMemory(stack, out1Offset, out1, FromByte(byte(ReadI32(stack, fp, inp1))))
	// case TYPE_STR: WriteMemory(stack, out1Offset, out1, FromStr(strconv.Itoa(ReadI32(stack, fp, inp1))))
	case TYPE_I32:
		WriteMemory(stack, out1Offset, out1, FromI32(ReadI32(stack, fp, inp1)))
	case TYPE_I64:
		WriteMemory(stack, out1Offset, out1, FromI64(int64(ReadI32(stack, fp, inp1))))
	case TYPE_F32:
		WriteMemory(stack, out1Offset, out1, FromF32(float32(ReadI32(stack, fp, inp1))))
	case TYPE_F64:
		WriteMemory(stack, out1Offset, out1, FromF64(float64(ReadI32(stack, fp, inp1))))
	}
}

func op_i32_print(expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

func op_i32_add(expr *CXExpression, stack *CXStack, fp int) {
	//
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_sub(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_mul(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) * ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_div(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) / ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_abs(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Abs(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_pow(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(math.Pow(float64(ReadI32(stack, fp, inp1)), float64(ReadI32(stack, fp, inp2)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_gt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_gteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) >= ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_lt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_lteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) <= ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_eq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) == ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_uneq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(stack, fp, inp1) != ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_mod(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) % ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_rand(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI32(stack, fp, inp1)
	maximum := ReadI32(stack, fp, inp2)

	outB1 := FromI32(int32(rand.Intn(int(maximum-minimum)) + int(minimum)))

	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitand(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) & ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) | ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitxor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) ^ ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitclear(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(stack, fp, inp1) &^ ReadI32(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitshl(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(stack, fp, inp1)) << uint32(ReadI32(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_i32_bitshr(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(stack, fp, inp1)) >> uint32(ReadI32(stack, fp, inp2))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

//implement more native functions for i32
func op_i32_print2(expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(stack, fp, inp1))
}

/*
sqrt() for i32
*/
func op_i32_sqrt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Sqrt(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

/*
log() for i32. returns the natural logarithm of x.
*/
func op_i32_log(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

/*
log2() for i32. returns the binary logarithm of x.
*/

func op_i32_log2(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log2(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

/*
log10() for i32. returns the decimal logarithm of x.
*/
func op_i32_log10(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log10(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

/*
round() for i32. returns the decimal logarithm of x.
*/
func op_i32_round(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Round(float64(ReadI32(stack, fp, inp1)))))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
