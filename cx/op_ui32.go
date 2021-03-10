package cxcore

import (
	"fmt"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI32ToStr(inputs []CXValue, outputs []CXValue) {
	outV0 := strconv.FormatUint(uint64(inputs[0].Get_ui32()), 10)
	outputs[0].Set_str(outV0)
}

// The built-in i8 function returns operand 1 casted from type ui32 to type i8.
func opUI32ToI8(inputs []CXValue, outputs []CXValue) {
	outV0 := int8(inputs[0].Get_ui32())
    outputs[0].Set_i8(outV0)
}

// The built-in i16 function returns operand 1 casted from type ui32 to type i16.
func opUI32ToI16(inputs []CXValue, outputs []CXValue) {
	outV0 := int16(inputs[0].Get_ui32())
    outputs[0].Set_i16(outV0)
}

// The built-in i32 function return operand 1 casted from type ui32 to type i32.
func opUI32ToI32(inputs []CXValue, outputs []CXValue) {
	outV0 := int32(inputs[0].Get_ui32())
    outputs[0].Set_i32(outV0)
}

// The built-in i64 function returns operand 1 casted from type ui32 to type i64.
func opUI32ToI64(inputs []CXValue, outputs []CXValue) {
	outV0 := int64(inputs[0].Get_ui32())
    outputs[0].Set_i64(outV0)
}

// The built-in ui8 function returns operand 1 casted from type ui32 to type ui8.
func opUI32ToUI8(inputs []CXValue, outputs []CXValue) {
	outV0 := uint8(inputs[0].Get_ui32())
    outputs[0].Set_ui8(outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui32 to type ui16.
func opUI32ToUI16(inputs []CXValue, outputs []CXValue) {
	outV0 := uint16(inputs[0].Get_ui32())
    outputs[0].Set_ui16(outV0)
}

// The built-in ui64 function returns the operand 1 casted from type ui32 to type ui64.
func opUI32ToUI64(inputs []CXValue, outputs []CXValue) {
	outV0 := uint64(inputs[0].Get_ui32())
    outputs[0].Set_ui64(outV0)
}

// The built-in f32 function returns operand 1 casted from type ui32 to type f32.
func opUI32ToF32(inputs []CXValue, outputs []CXValue) {
	outV0 := float32(inputs[0].Get_ui32())
    outputs[0].Set_f32(outV0)
}

// The built-in f64 function returns operand 1 casted from type ui32 to type f64.
func opUI32ToF64(inputs []CXValue, outputs []CXValue) {
	outV0 := float64(inputs[0].Get_ui32())
    outputs[0].Set_f64(outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI32Print(inputs []CXValue, outputs []CXValue) {
	fmt.Println(inputs[0].Get_ui32())
}

// The built-in add function returns the sum of two ui32 numbers.
func opUI32Add(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() + inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in sub function returns the difference of two ui32 numbers.
func opUI32Sub(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() - inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in mul function returns the product of two ui32 numbers.
func opUI32Mul(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() * inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in div function returns the quotient of two ui32 numbers.
func opUI32Div(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() / inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI32Gt(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() > inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI32Gteq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() >= inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI32Lt(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() < inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI32Lteq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() <= inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI32Eq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() == inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI32Uneq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() != inputs[1].Get_ui32()
	outputs[0].Set_bool(outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI32Mod(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() % inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI32Rand(inputs []CXValue, outputs []CXValue) {
	outV0 := rand.Uint32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI32Bitand(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() & inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI32Bitor(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() | inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI32Bitxor(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() ^ inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI32Bitclear(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() &^ inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI32Bitshl(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() << inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI32Bitshr(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_ui32() >> inputs[1].Get_ui32()
	outputs[0].Set_ui32(outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI32Max(inputs []CXValue, outputs []CXValue) {
	inpV0 := inputs[0].Get_ui32()
	inpV1 := inputs[1].Get_ui32()
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_ui32(inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI32Min(inputs []CXValue, outputs []CXValue) {
	inpV0 := inputs[0].Get_ui32()
	inpV1 := inputs[1].Get_ui32()
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_ui32(inpV0)
}
