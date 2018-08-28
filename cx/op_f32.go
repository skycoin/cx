package base

import (
	"fmt"
	"math"
)

func op_f32_f32(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadF32(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadF32(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(int64(ReadF32(fp, inp1))))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadF32(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(float64(ReadF32(fp, inp1))))
	}
}

func op_f32_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadF32(fp, inp1))
}

// op_f32_add. The add built-in function returns the add of two numbers

func op_f32_add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) + ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f32_sub. The sub built-in function returns the substract of two numbers

func op_f32_sub(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) - ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f32_sub. The mul built-in function returns the multiplication of two numbers

func op_f32_mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) * ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f32_sub. The div built-in function returns the divides two numbers

func op_f32_div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) / ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f32_abs. The div built-in function returns the absolute number of the number

func op_f32_abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Abs(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_pow. The div built-in function returns x**n for n>0 otherwise 1

func op_f32_pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Pow(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_gt. The gt built-in function returns true if x number is greater than a y number

func op_f32_gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) > ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f32_gteq. The gteq built-in function returns true if x number is greater or
// equal than a y number

func op_f32_gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) >= ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_lt. The lt built-in function returns true if x number is less then

func op_f32_lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) < ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_lteq. The lteq built-in function returns true if x number is less or
// equal than a y number

func op_f32_lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) <= ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_eq. The eq built-in function returns true if x number is equal to the y number

func op_f32_eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) == ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_uneq. The uneq built-in function returns true if x number is diferent to the y number

func op_f32_uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) != ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_cos. The cos built-in function returns the cosine of x number.

func op_f32_cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Cos(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_cos. The cos built-in function returns the sine of x number.

func op_f32_sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Sin(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_sqrt. The sqrt built-in function returns the square root of x number

func op_f32_sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Sqrt(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log. The log built-in function returns the natural logarithm of x number

func op_f32_log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log2. The log2 built-in function returns the natural logarithm based 2 of x number

func op_f32_log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log2(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log10. The log10 built-in function returns the natural logarithm based 2 of x number

func op_f32_log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log10(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_max. The max built-in function returns the max value between x and y numbers

func op_f32_max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Max(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_min. The min built-in function returns the min value between x and y numbers

func op_f32_min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Min(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
