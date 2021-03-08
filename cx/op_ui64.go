package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI64ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(ReadUI64(fp, expr.Inputs[0]), 10))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type ui64 to type i8.
func opUI64ToI8(expr *CXExpression, fp int) {
	outV0 := int8(ReadUI64(fp, expr.Outputs[0]))
	WriteI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i16 function retruns operand 1 casted from type ui64 to i16.
func opUI64ToI16(expr *CXExpression, fp int) {
	outV0 := int16(ReadUI64(fp, expr.Outputs[0]))
	WriteI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i32 function return operand 1 casted from type ui64 to type i32.
func opUI64ToI32(expr *CXExpression, fp int) {
	outV0 := int32(ReadUI64(fp, expr.Inputs[0]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i64 function return operand 1 casted from type ui64 to type i64.
func opUI64ToI64(expr *CXExpression, fp int) {
	outV0 := int64(ReadUI64(fp, expr.Inputs[0]))
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui8 function returns operand 1 casted from type ui64 to type ui8.
func opUI64ToUI8(expr *CXExpression, fp int) {
	outV0 := uint8(ReadUI64(fp, expr.Inputs[0]))
	WriteUI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui64 to type ui16.
func opUI64ToUI16(expr *CXExpression, fp int) {
	outV0 := uint16(ReadUI64(fp, expr.Inputs[0]))
	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui64 to type ui32.
func opUI64ToUI32(expr *CXExpression, fp int) {
	outV0 := uint32(ReadUI64(fp, expr.Inputs[0]))
	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in f32 function returns operand 1 casted from type ui64 to type f32.
func opUI64ToF32(expr *CXExpression, fp int) {
	outV0 := float32(ReadUI64(fp, expr.Inputs[0]))
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in f64 function returns operand 1 casted from type ui64 to type f64.
func opUI64ToF64(expr *CXExpression, fp int) {
	outV0 := float64(ReadUI64(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI64Print(expr *CXExpression, fp int) {
	fmt.Println(ReadUI64(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of two ui64 numbers.
func opUI64Add(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) + ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in sub function returns the difference of two ui64 numbers.
func opUI64Sub(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) - ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mul function returns the product of two ui64 numbers.
func opUI64Mul(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) * ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in div function returns the quotient of two ui64 numbers.
func opUI64Div(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) / ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI64Gt(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) > ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI64Gteq(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) >= ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI64Lt(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) < ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI64Lteq(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) <= ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI64Eq(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) == ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI64Uneq(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) != ReadUI64(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI64Mod(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) % ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI64Rand(expr *CXExpression, fp int) {
	outV0 := rand.Uint64()
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI64Bitand(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) & ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI64Bitor(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) | ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI64Bitxor(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) ^ ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI64Bitclear(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) &^ ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI64Bitshl(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) << ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI64Bitshr(expr *CXExpression, fp int) {
	outV0 := ReadUI64(fp, expr.Inputs[0]) >> ReadUI64(fp, expr.Inputs[1])
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI64Max(expr *CXExpression, fp int) {
	inpV0 := ReadUI64(fp, expr.Inputs[0])
	inpV1 := ReadUI64(fp, expr.Inputs[1])
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI64Min(expr *CXExpression, fp int) {
	inpV0 := ReadUI64(fp, expr.Inputs[0])
	inpV1 := ReadUI64(fp, expr.Inputs[1])
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), inpV0)
}
