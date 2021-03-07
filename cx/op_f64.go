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
	outV0 := int8(ReadF64(fp, expr.Inputs[0]))
	WriteI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i16 function returns operand 1 casted from type f64 to type i16.
func opF64ToI16(expr *CXExpression, fp int) {
	outV0 := int16(ReadF64(fp, expr.Inputs[0]))
	WriteI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i32 function return operand 1 casted from type f64 to type i32.
func opF64ToI32(expr *CXExpression, fp int) {
	outV0 := int32(ReadF64(fp, expr.Inputs[0]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i64 function returns operand 1 casted from type f64 to type i64.
func opF64ToI64(expr *CXExpression, fp int) {
	outV0 := int64(ReadF64(fp, expr.Inputs[0]))
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui8 function returns operand 1 casted from type f64 to type ui8.
func opF64ToUI8(expr *CXExpression, fp int) {
	outV0 := uint8(ReadF64(fp, expr.Inputs[0]))
	WriteUI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui16 function returns the operand 1 casted from type f64 to type ui16.
func opF64ToUI16(expr *CXExpression, fp int) {
	outV0 := uint16(ReadF64(fp, expr.Inputs[0]))
	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui32 function returns the operand 1 casted from type f64 to type ui32.
func opF64ToUI32(expr *CXExpression, fp int) {
	outV0 := uint32(ReadF64(fp, expr.Inputs[0]))
	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui64 function returns the operand 1 casted from type f64 to type ui64.
func opF64ToUI64(expr *CXExpression, fp int) {
	outV0 := uint64(ReadF64(fp, expr.Inputs[0]))
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in f32 function returns operand 1 casted from type f64 to type f32.
func opF64ToF32(expr *CXExpression, fp int) {
	outV0 := float32(ReadF64(fp, expr.Inputs[0]))
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in isnan function returns true if operand is nan value.
func opF64Isnan(expr *CXExpression, fp int) {
	outV0 := math.IsNaN(ReadF64(fp, expr.Inputs[0]))
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The print built-in function formats its arguments and prints them.
func opF64Print(expr *CXExpression, fp int) {
	fmt.Println(ReadF64(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of the two operands.
func opF64Add(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) + ReadF64(fp, expr.Inputs[1])
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in sub function returns the difference between the two operands.
func opF64Sub(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) - ReadF64(fp, expr.Inputs[1])
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opF64Neg(expr *CXExpression, fp int) {
	outV0 := -ReadF64(fp, expr.Inputs[0])
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mul function returns the product of the two operands.
func opF64Mul(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) * ReadF64(fp, expr.Inputs[1])
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in div function returns the quotient between the two operands.
func opF64Div(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) / ReadF64(fp, expr.Inputs[1])
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mod function return the floating-point remainder of operand 1 divided by operand 2.
func opF64Mod(expr *CXExpression, fp int) {
	outV0 := math.Mod(ReadF64(fp, expr.Inputs[0]), ReadF64(fp, expr.Inputs[1]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in abs function returns the absolute value of the operand.
func opF64Abs(expr *CXExpression, fp int) {
	outV0 := math.Abs(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in pow function returns x**n for n>0 otherwise 1.
func opF64Pow(expr *CXExpression, fp int) {
	outV0 := math.Pow(ReadF64(fp, expr.Inputs[0]), ReadF64(fp, expr.Inputs[1]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gt function returns true if operand 1 is larger than operand 2.
func opF64Gt(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) > ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opF64Gteq(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) >= ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opF64Lt(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) < ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 2.
func opF64Lteq(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) <= ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opF64Eq(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) == ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opF64Uneq(expr *CXExpression, fp int) {
	outV0 := ReadF64(fp, expr.Inputs[0]) != ReadF64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in rand function returns a pseudo-random number in [0.0,1.0) from the default Source.
func opF64Rand(expr *CXExpression, fp int) {
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), rand.Float64())
}

// The built-in acos function returns the arc cosine of the operand.
func opF64Acos(expr *CXExpression, fp int) {
	outV0 := math.Acos(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in cos function returns the cosine of the operand.
func opF64Cos(expr *CXExpression, fp int) {
	outV0 := math.Cos(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in asin function returns the arc sine of the operand.
func opF64Asin(expr *CXExpression, fp int) {
	outV0 := math.Asin(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in sin function returns the sine of the operand.
func opF64Sin(expr *CXExpression, fp int) {
	outV0 := math.Sin(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in sqrt function returns the square root of the operand.
func opF64Sqrt(expr *CXExpression, fp int) {
	outV0 := math.Sqrt(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in log function returns the natural logarithm of the operand.
func opF64Log(expr *CXExpression, fp int) {
	outV0 := math.Log(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in log2 function returns the 2-logarithm of the operand.
func opF64Log2(expr *CXExpression, fp int) {
	outV0 := math.Log2(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in log10 function returns the 10-logarithm of the operand.
func opF64Log10(expr *CXExpression, fp int) {
	outV0 := math.Log10(ReadF64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in max function returns the largest value of the two operands.
func opF64Max(expr *CXExpression, fp int) {
	outV0 := math.Max(ReadF64(fp, expr.Inputs[0]), ReadF64(fp, expr.Inputs[1]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in min function returns the smallest value of the two operands.
func opF64Min(expr *CXExpression, fp int) {
	outV0 := math.Min(ReadF64(fp, expr.Inputs[0]), ReadF64(fp, expr.Inputs[1]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}
