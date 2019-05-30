package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI32ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatInt(int64(ReadI32(fp, expr.Inputs[0])), 10))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type i32 to type i8.
func opI32ToI8(expr *CXExpression, fp int) {
	outB0 := FromI8(int8(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i16 function returns operand 1 casted from type i32 to type i16.
func opI32ToI16(expr *CXExpression, fp int) {
	outB0 := FromI16(int16(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i64 function returns operand 1 casted from type i32 to type i64.
func opI32ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui8 function returns operand 1 casted from type i32 to type ui8.
func opI32ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui16 function returns the operand 1 casted from type i32 to type ui16.
func opI32ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui32 function returns the operand 1 casted from type i32 to type ui32.
func opI32ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui64 function returns the operand 1 casted from type i32 to type ui64.
func opI32ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in f32 function returns operand 1 casted from type i32 to type f32.
func opI32ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in f64 function returns operand 1 casted from type i32 to type f64.
func opI32ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The print built-in function formats its arguments and prints them.
func opI32Print(expr *CXExpression, fp int) {
	fmt.Println(ReadI32(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of two i32 numbers.
func opI32Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) + ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference of two i32 numbers.
func opI32Sub(expr *CXExpression, fp int) {
	outB0 := FromI32(ReadI32(fp, expr.Inputs[0]) - ReadI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in neg function returns the opposite of operand 1.
func opI32Neg(expr *CXExpression, fp int) {
	outB0 := FromI32(-ReadI32(fp, expr.Inputs[0]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in mul function returns the product of two i32 numbers.
func opI32Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) * ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient of two i32 numbers.
func opI32Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) / ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute number of the number.
func opI32Abs(expr *CXExpression, fp int) {
	inpV0 := ReadI32(fp, expr.Inputs[0])
	sign := inpV0 >> 31
	outB0 := FromI32((inpV0 ^ sign) - sign)
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI32Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) > ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI32Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) >= ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI32Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) < ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI32Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) <= ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI32Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) == ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI32Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) != ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mod function returns the remainder of operand 1 divided by operand 2.
func opI32Mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) % ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI32Rand(expr *CXExpression, fp int) {
	minimum := ReadI32(fp, expr.Inputs[0])
	maximum := ReadI32(fp, expr.Inputs[1])

	r := int(maximum - minimum)
	val0 := int32(0)
	if r > 0 {
		val0 = int32(rand.Intn(r) + int(minimum))
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(val0))
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI32Bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) & ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI32Bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) | ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI32Bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) ^ ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI32Bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) &^ ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI32Bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) << uint32(ReadI32(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI32Bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) >> uint32(ReadI32(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the biggest of the two operands.
func opI32Max(expr *CXExpression, fp int) {
	inp0 := ReadI32(fp, expr.Inputs[0])
	inp1 := ReadI32(fp, expr.Inputs[1])
	if inp1 > inp0 {
		inp0 = inp1
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(inp0))
}

// The built-in min function returns the smallest of the two operands.
func opI32Min(expr *CXExpression, fp int) {
	inp0 := ReadI32(fp, expr.Inputs[0])
	inp1 := ReadI32(fp, expr.Inputs[1])
	if inp1 < inp0 {
		inp0 = inp1
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(inp0))
}
