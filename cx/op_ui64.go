package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI64ToStr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatUint(inputs[0].Get_ui64(), 10)
	outputs[0].Set_str(outV0)
}

// The built-in i8 function returns operand 1 casted from type ui64 to type i8.
func opUI64ToI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_ui64())
    outputs[0].Set_i8(outV0)
}

// The built-in i16 function retruns operand 1 casted from type ui64 to i16.
func opUI64ToI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_ui64())
    outputs[0].Set_i16(outV0)
}

// The built-in i32 function return operand 1 casted from type ui64 to type i32.
func opUI64ToI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_ui64())
    outputs[0].Set_i32(outV0)
}

// The built-in i64 function return operand 1 casted from type ui64 to type i64.
func opUI64ToI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_ui64())
    outputs[0].Set_i64(outV0)
}

// The built-in ui8 function returns operand 1 casted from type ui64 to type ui8.
func opUI64ToUI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_ui64())
    outputs[0].Set_ui8(outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui64 to type ui16.
func opUI64ToUI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_ui64())
    outputs[0].Set_ui16(outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui64 to type ui32.
func opUI64ToUI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_ui64())
    outputs[0].Set_ui32(outV0)
}

// The built-in f32 function returns operand 1 casted from type ui64 to type f32.
func opUI64ToF32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_ui64())
    outputs[0].Set_f32(outV0)
}

// The built-in f64 function returns operand 1 casted from type ui64 to type f64.
func opUI64ToF64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_ui64())
    outputs[0].Set_f64(outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI64Print(inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_ui64())
}

// The built-in add function returns the sum of two ui64 numbers.
func opUI64Add(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() + inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in sub function returns the difference of two ui64 numbers.
func opUI64Sub(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() - inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in mul function returns the product of two ui64 numbers.
func opUI64Mul(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() * inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in div function returns the quotient of two ui64 numbers.
func opUI64Div(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() / inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI64Gt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() > inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI64Gteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() >= inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI64Lt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() < inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI64Lteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() <= inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI64Eq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() == inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI64Uneq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() != inputs[1].Get_ui64()
	outputs[0].Set_bool(outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI64Mod(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() % inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI64Rand(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := rand.Uint64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI64Bitand(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() & inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI64Bitor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() | inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI64Bitxor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() ^ inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI64Bitclear(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() &^ inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI64Bitshl(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() << inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI64Bitshr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64() >> inputs[1].Get_ui64()
	outputs[0].Set_ui64(outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI64Max(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui64()
	inpV1 := inputs[1].Get_ui64()
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_ui64(inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI64Min(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui64()
	inpV1 := inputs[1].Get_ui64()
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_ui64(inpV0)
}
