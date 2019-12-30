package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI32ToStr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromStr(strconv.FormatInt(int64(ReadI32(fp, expr.Inputs[0])), 10))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// The built-in i8 function returns operand 1 casted from type i32 to type i8.
func opI32ToI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := int8(ReadI32(fp, expr.Inputs[0]))
	WriteI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i16 function returns operand 1 casted from type i32 to type i16.
func opI32ToI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := int16(ReadI32(fp, expr.Inputs[0]))
	WriteI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in i64 function returns operand 1 casted from type i32 to type i64.
func opI32ToI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := int64(ReadI32(fp, expr.Inputs[0]))
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui8 function returns operand 1 casted from type i32 to type ui8.
func opI32ToUI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := uint8(ReadI32(fp, expr.Inputs[0]))
	WriteUI8(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i32 to type ui16.
func opI32ToUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := uint16(ReadI32(fp, expr.Inputs[0]))
	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui32 function returns the operand 1 casted from type i32 to type ui32.
func opI32ToUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := uint32(ReadI32(fp, expr.Inputs[0]))
	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in ui64 function returns the operand 1 casted from type i32 to type ui64.
func opI32ToUI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := uint64(ReadI32(fp, expr.Inputs[0]))
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in f32 function returns operand 1 casted from type i32 to type f32.
func opI32ToF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := float32(ReadI32(fp, expr.Inputs[0]))
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in f64 function returns operand 1 casted from type i32 to type f64.
func opI32ToF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := float64(ReadI32(fp, expr.Inputs[0]))
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The print built-in function formats its arguments and prints them.
func opI32Print(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	fmt.Println(ReadI32(fp, expr.Inputs[0]))
}

// The built-in add function returns the sum of two i32 numbers.
func opI32Add(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) + ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in sub function returns the difference of two i32 numbers.
func opI32Sub(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) - ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opI32Neg(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := -ReadI32(fp, expr.Inputs[0])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mul function returns the product of two i32 numbers.
func opI32Mul(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) * ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in div function returns the quotient of two i32 numbers.
func opI32Div(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) / ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in abs function returns the absolute number of the number.
func opI32Abs(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := ReadI32(fp, expr.Inputs[0])
	sign := inpV0 >> 31
	outV0 := (inpV0 ^ sign) - sign
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI32Gt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) > ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI32Gteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) >= ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI32Lt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) < ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI32Lteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) <= ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI32Eq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) == ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI32Uneq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) != ReadI32(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in mod function returns the remainder of operand 1 divided by operand 2.
func opI32Mod(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) % ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI32Rand(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	minimum := ReadI32(fp, expr.Inputs[0])
	maximum := ReadI32(fp, expr.Inputs[1])

	r := int(maximum - minimum)
	outV0 := int32(0)
	if r > 0 {
		outV0 = int32(rand.Intn(r) + int(minimum))
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI32Bitand(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) & ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI32Bitor(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) | ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI32Bitxor(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) ^ ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI32Bitclear(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) &^ ReadI32(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI32Bitshl(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) << uint32(ReadI32(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI32Bitshr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := ReadI32(fp, expr.Inputs[0]) >> uint32(ReadI32(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

// The built-in max function returns the biggest of the two operands.
func opI32Max(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := ReadI32(fp, expr.Inputs[1])
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opI32Min(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := ReadI32(fp, expr.Inputs[1])
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), inpV0)
}
