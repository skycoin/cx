package cxcore

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opF64ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatFloat(ReadF64(fp, expr.Inputs[0]), 'f', -1, 64))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type f64 to type i8.
func opF64ToI8(expr *CXExpression, fp int) {
	outB0 := FromI8(int8(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i16 function returns operand 1 casted from type f64 to type i16.
func opF64ToI16(expr *CXExpression, fp int) {
	outB0 := FromI16(int16(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i32 function return operand 1 casted from type f64 to type i32.
func opF64ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i64 function returns operand 1 casted from type f64 to type i64.
func opF64ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui8 function returns operand 1 casted from type f64 to type ui8.
func opF64ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui16 function returns the operand 1 casted from type f64 to type ui16.
func opF64ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui32 function returns the operand 1 casted from type f64 to type ui32.
func opF64ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui64 function returns the operand 1 casted from type f64 to type ui64.
func opF64ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in f32 function returns operand 1 casted from type f64 to type f32.
func opF64ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in isnan function returns true if operand is nan value.
func opF64Isnan(expr *CXExpression, fp int) {
	outB0 := FromBool(math.IsNaN(ReadF64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The print built-in function formats its arguments and prints them.
func opF64Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadF64(fp, inp1))
}

// The built-in add function returns the sum of the two operands.
func opF64Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) + ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference between the two operands.
func opF64Sub(expr *CXExpression, fp int) {
	outB0 := FromF64(ReadF64(fp, expr.Inputs[0]) - ReadF64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in neg function returns the opposite of operand 1.
func opF64Neg(expr *CXExpression, fp int) {
	outB0 := FromF64(-ReadF64(fp, expr.Inputs[0]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in mul function returns the product of the two operands.
func opF64Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) * ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient between the two operands.
func opF64Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) / ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mod function return the floating-point remainder of operand 1 divided by operand 2.
func opF64Mod(expr *CXExpression, fp int) {
	outB0 := FromF64(math.Mod(ReadF64(fp, expr.Inputs[0]), ReadF64(fp, expr.Inputs[1])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in abs function returns the absolute value of the operand.
func opF64Abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Abs(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in pow function returns x**n for n>0 otherwise 1.
func opF64Pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Pow(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gt function returns true if operand 1 is larger than operand 2.
func opF64Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) > ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opF64Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) >= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opF64Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) < ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 2.
func opF64Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) <= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opF64Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) == ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opF64Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) != ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in rand function returns a pseudo-random number in [0.0,1.0) from the default Source.
func opF64Rand(expr *CXExpression, fp int) {
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF64(rand.Float64()))
}

// The built-in acos function returns the arc cosine of the operand.
func opF64Acos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Acos(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in cos function returns the cosine of the operand.
func opF64Cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Cos(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in asin function returns the arc sine of the operand.
func opF64Asin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Asin(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sin function returns the sine of the operand.
func opF64Sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sin(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sqrt function returns the square root of the operand.
func opF64Sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sqrt(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log function returns the natural logarithm of the operand.
func opF64Log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log2 function returns the 2-logarithm of the operand.
func opF64Log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log2(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log10 function returns the 10-logarithm of the operand.
func opF64Log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log10(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the largest value of the two operands.
func opF64Max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Max(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in min function returns the smallest value of the two operands.
func opF64Min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Min(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
