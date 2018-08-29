package base

import (
	"fmt"
	"math"
	"math/rand"
)

func op_i64_i64(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadI64(fp, inp1))))
		// case TYPE_STR: WriteMemory(out1Offset, FromStr(strconv.Itoa(ReadI32(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadI64(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(ReadI64(fp, inp1)))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadI64(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(float64(ReadI64(fp, inp1))))
	}
}

// op_i64_print. The print built-in function formats its arguments in an
// implementation-specific

func op_i64_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI64(fp, inp1))
}

// op_i64_add. The add built-in function returns the add of two numbers

func op_i64_add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) + ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_sub. The sub built-in function returns the substract of two numbers

func op_i64_sub(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) - ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_sub. The mul built-in function returns the multiplication of two numbers

func op_i64_mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) * ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_sub. The div built-in function returns the divides two numbers

func op_i64_div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) / ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_abs. The div built-in function returns the absolute number of the number

func op_i64_abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Abs(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_pow. The div built-in function returns x**n for n>0 otherwise 1

func op_i64_pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Pow(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_gt. The gt built-in function returns true if x number is greater than a y number

func op_i64_gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) > ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_gteq. The gteq built-in function returns true if x number is greater or
// equal than a y number

func op_i64_gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) >= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_lt. The lt built-in function returns true if x number is less then

func op_i64_lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) < ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_lteq. The lteq built-in function returns true if x number is less or
// equal than a y number

func op_i64_lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) <= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_eq. The eq built-in function returns true if x number is equal to the y number

func op_i64_eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) == ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_uneq. The uneq built-in function returns true if x number is diferent to the y number

func op_i64_uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) != ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) % ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_rand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI64(fp, inp1)
	maximum := ReadI64(fp, inp2)

	outB1 := FromI64(int64(rand.Intn(int(maximum-minimum)) + int(minimum)))

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) & ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) | ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) ^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) &^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(fp, inp1)) << uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(fp, inp1)) >> uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Cos(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_i64_sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Sin(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_sqrt. The sqrt built-in function returns the square root of x number

func op_i64_sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Sqrt(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_log. The log built-in function returns the natural logarithm of x number

func op_i64_log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_log2. The log2 built-in function returns the natural logarithm based 2 of x number

func op_i64_log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log2(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_log10. The log10 built-in function returns the natural logarithm based 2 of x number

func op_i64_log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log10(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_max. The max built-in function returns the max value between x and y numbers

func op_i64_max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Max(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// op_i64_min. The min built-in function returns the min value between x and y numbers

func op_i64_min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Min(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
