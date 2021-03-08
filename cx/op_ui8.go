package cxcore

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI8ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(uint64(ReadUI8(fp, expr.Inputs[0])), 10))
	WriteObject(GetOffsetUI8(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type ui8 to type i8.
func opUI8ToI8(expr *CXExpression, fp int) {
	outV0 := int8(ReadUI8(fp, expr.Inputs[0]))
	WriteI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in i16 function returns operand 1 casted from type ui8 to type i16.
func opUI8ToI16(expr *CXExpression, fp int) {
	outV0 := int16(ReadUI8(fp, expr.Inputs[0]))
	WriteI16(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in i32 function return operand 1 casted from type ui8 to type i32.
func opUI8ToI32(expr *CXExpression, fp int) {
	outV0 := int32(ReadUI8(fp, expr.Inputs[0]))
	WriteI32(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in i64 function returns operand 1 casted from type ui8 to type i64.
func opUI8ToI64(expr *CXExpression, fp int) {
	outV0 := int64(ReadUI8(fp, expr.Inputs[0]))
	WriteI64(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui8 to type ui16.
func opUI8ToUI16(expr *CXExpression, fp int) {
	outV0 := uint16(ReadUI8(fp, expr.Inputs[0]))
	WriteUI16(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui8 to type ui32.
func opUI8ToUI32(expr *CXExpression, fp int) {
	outV0 := uint32(ReadUI8(fp, expr.Inputs[0]))
	WriteUI32(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui64 function returns the operand 1 casted from type ui8 to type ui64.
func opUI8ToUI64(expr *CXExpression, fp int) {
	outV0 := uint64(ReadUI8(fp, expr.Inputs[0]))
	WriteUI64(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in f32 function returns operand 1 casted from type ui8 to type f32.
func opUI8ToF32(expr *CXExpression, fp int) {
	outV0 := float32(ReadUI8(fp, expr.Inputs[0]))
	WriteF32(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in f64 function returns operand 1 casted from type ui8 to type f64.
func opUI8ToF64(expr *CXExpression, fp int) {
	outV0 := float64(ReadUI8(fp, expr.Inputs[0]))
	WriteF64(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI8Print(expr *CXExpression, fp int) {
	fmt.Println(ReadUI8(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of two ui8 numbers.
func opUI8Add(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) + ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in sub function returns the difference of two ui8 numbers.
func opUI8Sub(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) - ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in mul function returns the product of two ui8 numbers.
func opUI8Mul(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) * ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in div function returns the quotient of two ui8 numbers.
func opUI8Div(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) / ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI8Gt(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) > ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI8Gteq(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) >= ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI8Lt(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) < ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI8Lteq(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) <= ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI8Eq(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) == ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI8Uneq(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) != ReadUI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI8Mod(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) % ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI8Rand(expr *CXExpression, fp int) {
	outV0 := uint8(rand.Int31n(int32(math.MaxUint8)))
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI8Bitand(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) & ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI8Bitor(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) | ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI8Bitxor(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) ^ ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI8Bitclear(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) &^ ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI8Bitshl(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) << ReadUI8(fp, expr.Inputs[1])
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI8Bitshr(expr *CXExpression, fp int) {
	outV0 := ReadUI8(fp, expr.Inputs[0]) >> uint32(ReadUI8(fp, expr.Inputs[1]))
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI8Max(expr *CXExpression, fp int) {
	inpV0 := ReadUI8(fp, expr.Inputs[0])
	inpV1 := ReadUI8(fp, expr.Inputs[1])
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI8Min(expr *CXExpression, fp int) {
	inpV0 := ReadUI8(fp, expr.Inputs[0])
	inpV1 := ReadUI8(fp, expr.Inputs[1])
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	WriteUI8(GetOffsetUI8(fp, expr.Outputs[0]), inpV0)
}
