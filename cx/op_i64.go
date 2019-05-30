package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI64ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatInt(ReadI64(fp, expr.Inputs[0]), 10))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type i64 to type i8.
func opI64ToI8(expr *CXExpression, fp int) {
	outB0 := FromI8(int8(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i16 function returns operand 1 casted from type i64 to type i16.
func opI64ToI16(expr *CXExpression, fp int) {
	outB0 := FromI16(int16(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i32 function returns operand 1 casted from type i64 to type i32.
func opI64ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui8 function returns operand 1 casted from type i64 to type ui8.
func opI64ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui16 function returns the operand 1 casted from type i64 to type ui16.
func opI64ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui32 function returns the operand 1 casted from type i64 to type ui32.
func opI64ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in ui64 function returns the operand 1 casted from type i64 to type ui64.
func opI64ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in f32 function returns operand 1 casted from type i64 to type f32.
func opI64ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in f64 function returns operand 1 casted from type i64 to type f64.
func opI64ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The print built-in function formats its arguments and prints them.
func opI64Print(expr *CXExpression, fp int) {
	fmt.Println(ReadI64(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of the two operands.
func opI64Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) + ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference between the two operands.
func opI64Sub(expr *CXExpression, fp int) {
	outB0 := FromI64(ReadI64(fp, expr.Inputs[0]) - ReadI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in neg function returns the opposit of operand 1.
func opI64Neg(expr *CXExpression, fp int) {
	outB0 := FromI64(-ReadI64(fp, expr.Inputs[0]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in mul function returns the product of the two operands.
func opI64Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) * ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient of the two operands.
func opI64Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) / ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mod function returns the remainder of operand 1 divided by operand 2.
func opI64Mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) % ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute value of the operand.
func opI64Abs(expr *CXExpression, fp int) {
	inpV0 := ReadI64(fp, expr.Inputs[0])
	sign := inpV0 >> 63
	outB0 := FromI64((inpV0 ^ sign) - sign)
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI64Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) > ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI64Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) >= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than oeprand 2.
func opI64Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) < ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or
// equal to operand 2.
func opI64Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) <= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI64Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) == ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI64Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) != ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI64Rand(expr *CXExpression, fp int) {
	minimum := ReadI64(fp, expr.Inputs[0])
	maximum := ReadI64(fp, expr.Inputs[1])

	r := int(maximum - minimum)
	val0 := int64(0)
	if r > 0 {
		val0 = int64(rand.Intn(r) + int(minimum))
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(val0))
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI64Bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) & ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI64Bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) | ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI64Bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) ^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI64Bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) &^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI64Bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(ReadI64(fp, inp1) << uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI64Bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(ReadI64(fp, inp1) >> uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the greatest value of the two operands.
func opI64Max(expr *CXExpression, fp int) {
	inp0 := ReadI64(fp, expr.Inputs[0])
	inp1 := ReadI64(fp, expr.Inputs[1])
	if inp1 > inp0 {
		inp0 = inp1
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(inp0))
}

// The built-in min function returns the smallest value of the two operands.
func opI64Min(expr *CXExpression, fp int) {
	inp0 := ReadI64(fp, expr.Inputs[0])
	inp1 := ReadI64(fp, expr.Inputs[1])
	if inp1 < inp0 {
		inp0 = inp1
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(inp0))
}
