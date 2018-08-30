package base

import (
	"fmt"
	"math"
	"strconv"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_f64_f64(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_STR:
		WriteObject(out1Offset, encoder.Serialize(strconv.FormatFloat(ReadF64(fp, inp1), 'f', -1, 64)))
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadF64(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadF64(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(int64(ReadF64(fp, inp1))))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadF64(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(ReadF64(fp, inp1)))
	}
}

// op_f64_print. The print built-in function formats its arguments in an
// implementation-specific

func op_f64_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadF64(fp, inp1))
}

// op_f64_add. The add built-in function returns the add of two numbers

func op_f64_add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) + ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_sub. The sub built-in function returns the substract of two numbers

func op_f64_sub(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) - ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_sub. The mul built-in function returns the multiplication of two numbers

func op_f64_mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) * ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_sub. The div built-in function returns the divides two numbers

func op_f64_div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) / ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_abs. The div built-in function returns the absolute number of the number

func op_f64_abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Abs(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_pow. The div built-in function returns x**n for n>0 otherwise 1

func op_f64_pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Pow(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_gt. The gt built-in function returns true if x number is greater than a y number

func op_f64_gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) > ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_gteq. The gteq built-in function returns true if x number is greater or
// equal than a y number

func op_f64_gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) >= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_lt. The lt built-in function returns true if x number is less then

func op_f64_lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) < ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_lteq. The lteq built-in function returns true if x number is less or
// equal than a y number

func op_f64_lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) <= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_eq. The eq built-in function returns true if x number is equal to the y number

func op_f64_eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) == ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_uneq. The uneq built-in function returns true if x number is diferent to the y number

func op_f64_uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) != ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_cos. The cos built-in function returns the cosine of x number.

func op_f64_cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Cos(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_cos. The cos built-in function returns the sine of x number.

func op_f64_sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sin(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_sqrt. The sqrt built-in function returns the square root of x number

func op_f64_sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sqrt(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log. The log built-in function returns the natural logarithm of x number

func op_f64_log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log2. The log2 built-in function returns the natural logarithm based 2 of x number

func op_f64_log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log2(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_log10. The log10 built-in function returns the natural logarithm based 2 of x number

func op_f64_log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log10(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_max. The max built-in function returns the max value between x and y numbers

func op_f64_max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Max(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_f64_min. The min built-in function returns the min value between x and y numbers

func op_f64_min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Min(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
