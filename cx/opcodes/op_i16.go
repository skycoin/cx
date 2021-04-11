package opcodes

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI16ToStr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatInt(int64(inputs[0].Get_i16()), 10)
	outputs[0].Set_str(outV0)
}

// The built-in i8 function returns operand 1 casted from type i16 to type i8.
func opI16ToI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_i16())
    outputs[0].Set_i8(outV0)
}

// The built-in i32 function returns operand 1 casted from type i16 to type i32.
func opI16ToI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_i16())
    outputs[0].Set_i32(outV0)
}

// The built-in i64 function returns operand 1 casted from type i16 to type i64.
func opI16ToI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_i16())
    outputs[0].Set_i64(outV0)
}

// The built-in ui8 function returns operand 1 casted from type i16 to type ui8.
func opI16ToUI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_i16())
    outputs[0].Set_ui8(outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i16 to type ui16.
func opI16ToUI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_i16())
    outputs[0].Set_ui16(outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i16 to type ui32.
func opI16ToUI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_i16())
    outputs[0].Set_ui32(outV0)
}

// The built-in ui64 function returns the operand 1 casted from type i16 to type ui64.
func opI16ToUI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_i16())
    outputs[0].Set_ui64(outV0)
}

// The built-in f32 function returns operand 1 casted from type i16 to type f32.
func opI16ToF32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_i16())
    outputs[0].Set_f32(outV0)
}

// The built-in f64 function returns operand 1 casted from type i16 to type f64.
func opI16ToF64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_i16())
    outputs[0].Set_f64(outV0)
}

// The print built-in function formats its arguments and prints them.
func opI16Print(inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_i16())
}

// The built-in add function returns the sum of two i16 numbers.
func opI16Add(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() + inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in sub function returns the difference of two i16 numbers.
func opI16Sub(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() - inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opI16Neg(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := -inputs[0].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in mul function returns the product of two i16 numbers.
func opI16Mul(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() * inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in div function returns the quotient of two i16 numbers.
func opI16Div(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() / inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in abs function returns the absolute number of the number.
func opI16Abs(inputs []ast.CXValue, outputs []ast.CXValue) {
	V0 := inputs[0].Get_i16()
	sign := V0 >> 15
	outV0 := (V0 ^ sign) - sign
	outputs[0].Set_i16(outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI16Gt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() > inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI16Gteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() >= inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI16Lt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() < inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI16Lteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() <= inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI16Eq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() == inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI16Uneq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() != inputs[1].Get_i16()
	outputs[0].Set_bool(outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opI16Mod(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() % inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI16Rand(inputs []ast.CXValue, outputs []ast.CXValue) {
	minimum := inputs[0].Get_i16()
	maximum := inputs[1].Get_i16()

	r := int(maximum - minimum)
	outV0 := int16(0)
	if r > 0 {
		outV0 = int16(rand.Intn(r) + int(minimum))
	}

	outputs[0].Set_i16(outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI16Bitand(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() & inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI16Bitor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() | inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI16Bitxor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() ^ inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI16Bitclear(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16() &^ inputs[1].Get_i16()
	outputs[0].Set_i16(outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI16Bitshl(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i16() << uint16(inputs[1].Get_i16()))
	outputs[0].Set_i16(outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI16Bitshr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i16() >> uint16(inputs[1].Get_i16()))
	outputs[0].Set_i16(outV0)
}

// The built-in max function returns the biggest of the two operands.
func opI16Max(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i16()
	inpV1 := inputs[1].Get_i16()
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_i16(inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opI16Min(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i16()
	inpV1 := inputs[1].Get_i16()
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_i16(inpV0)
}
