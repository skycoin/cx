package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI8ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatInt(int64(ReadI8(fp, expr.Inputs[0])), 10))
	WriteObject(GetOffsetI8(fp, expr.Outputs[0]), outB0)
}

// The built-in i16 function returns operand 1 casted from type i8 to type i16.
func opI8ToI16(expr *CXExpression, fp int) {
	outV0 := int16(ReadI8(fp, expr.Inputs[0]))
	WriteI16(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in i32 function returns operand 1 casted from type i8 to type i32.
func opI8ToI32(expr *CXExpression, fp int) {
	outV0 := int32(ReadI8(fp, expr.Inputs[0]))
	WriteI32(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in i64 function returns operand 1 casted from type i8 to type i64.
func opI8ToI64(expr *CXExpression, fp int) {
	outV0 := int64(ReadI8(fp, expr.Inputs[0]))
	WriteI64(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui8 function returns operand 1 casted from type i8 to type ui8.
func opI8ToUI8(expr *CXExpression, fp int) {
	outV0 := uint8(ReadI8(fp, expr.Inputs[0]))
	WriteUI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui16 function returns operand 1 casted from type i8 to type ui16.
func opI8ToUI16(expr *CXExpression, fp int) {
	outV0 := uint16(ReadI8(fp, expr.Inputs[0]))
	WriteUI16(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui32 function returns operand 1 casted from type i8 to type ui32.
func opI8ToUI32(expr *CXExpression, fp int) {
	outV0 := uint32(ReadI8(fp, expr.Inputs[0]))
	WriteUI32(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in ui64 function returns operand 1 casted from type i8 to type ui64.
func opI8ToUI64(expr *CXExpression, fp int) {
	outV0 := uint64(ReadI8(fp, expr.Inputs[0]))
	WriteUI64(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in f32 function returns operand 1 casted from type i8 to type f32.
func opI8ToF32(expr *CXExpression, fp int) {
	outV0 := float32(ReadI8(fp, expr.Inputs[0]))
	WriteF32(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in f64 function returns operand 1 casted from type i8 to type uf64.
func opI8ToF64(expr *CXExpression, fp int) {
	outV0 := float64(ReadI8(fp, expr.Inputs[0]))
	WriteF64(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The print built-in function formats its arguments and prints them.
func opI8Print(expr *CXExpression, fp int) {
	fmt.Println(ReadI8(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of two i8 numbers.
func opI8Add(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) + ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in sub function returns the difference of two i8 numbers.
func opI8Sub(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) - ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opI8Neg(expr *CXExpression, fp int) {
	outV0 := -ReadI8(fp, expr.Inputs[0])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in mul function returns the product of two i8 numbers.
func opI8Mul(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) * ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in div function returns the quotient of two i8 numbers.
func opI8Div(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) / ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in abs function returns the absolute number of the number.
func opI8Abs(expr *CXExpression, fp int) {
	inpV0 := ReadI8(fp, expr.Inputs[0])
	sign := inpV0 >> 7
	outV0 := (inpV0 ^ sign) - sign
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI8Gt(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) > ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI8Gteq(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) >= ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI8Lt(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) < ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI8Lteq(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) <= ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI8Eq(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) == ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI8Uneq(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) != ReadI8(fp, expr.Inputs[1])
	WriteBool(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opI8Mod(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) % ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI8Rand(expr *CXExpression, fp int) {
	minimum := ReadI8(fp, expr.Inputs[0])
	maximum := ReadI8(fp, expr.Inputs[1])

	r := int(maximum - minimum)
	outV0 := int8(0)
	if r > 0 {
		outV0 = int8(rand.Intn(r) + int(minimum))
	}

	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI8Bitand(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) & ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI8Bitor(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) | ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI8Bitxor(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) ^ ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI8Bitclear(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) &^ ReadI8(fp, expr.Inputs[1])
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI8Bitshl(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) << uint8(ReadI8(fp, expr.Inputs[1]))
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI8Bitshr(expr *CXExpression, fp int) {
	outV0 := ReadI8(fp, expr.Inputs[0]) >> uint8(ReadI8(fp, expr.Inputs[1]))
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), outV0)
}

// The built-in max function returns the biggest of the two operands.
func opI8Max(expr *CXExpression, fp int) {
	inpV0 := ReadI8(fp, expr.Inputs[0])
	inpV1 := ReadI8(fp, expr.Inputs[1])
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opI8Min(expr *CXExpression, fp int) {
	inpV0 := ReadI8(fp, expr.Inputs[0])
	inpV1 := ReadI8(fp, expr.Inputs[1])
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	WriteI8(GetOffsetI8(fp, expr.Outputs[0]), inpV0)
}
