package opcodes

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"math/rand"
	"strconv"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI64ToStr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatInt(inputs[0].Get_i64(), 10)
	outputs[0].Set_str(outV0)
}

// The built-in i8 function returns operand 1 casted from type i64 to type i8.
func opI64ToI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_i64())
    outputs[0].Set_i8(outV0)
}

// The built-in i16 function returns operand 1 casted from type i64 to type i16.
func opI64ToI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i64())
    outputs[0].Set_i16(outV0)
}

// The built-in i32 function returns operand 1 casted from type i64 to type i32.
func opI64ToI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_i64())
    outputs[0].Set_i32(outV0)
}

// The built-in ui8 function returns operand 1 casted from type i64 to type ui8.
func opI64ToUI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_i64())
    outputs[0].Set_ui8(outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i64 to type ui16.
func opI64ToUI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_i64())
    outputs[0].Set_ui16(outV0)
}

// The built-in ui32 function returns the operand 1 casted from type i64 to type ui32.
func opI64ToUI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_i64())
    outputs[0].Set_ui32(outV0)
}

// The built-in ui64 function returns the operand 1 casted from type i64 to type ui64.
func opI64ToUI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_i64())
    outputs[0].Set_ui64(outV0)
}

// The built-in f32 function returns operand 1 casted from type i64 to type f32.
func opI64ToF32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_i64())
    outputs[0].Set_f32(outV0)
}

// The built-in f64 function returns operand 1 casted from type i64 to type f64.
func opI64ToF64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_i64())
    outputs[0].Set_f64(outV0)
}

// The print built-in function formats its arguments and prints them.
func opI64Print(inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_i64())
}

// The built-in add function returns the sum of the two operands.
func opI64Add(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() + inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in sub function returns the difference between the two operands.
func opI64Sub(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() - inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in neg function returns the opposit of operand 1.
func opI64Neg(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := -inputs[0].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in mul function returns the product of the two operands.
func opI64Mul(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() * inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in div function returns the quotient of the two operands.
func opI64Div(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() / inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in abs function returns the absolute value of the operand.
func opI64Abs(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i64()
	sign := inpV0 >> 63
	outV0 := (inpV0 ^ sign) - sign
	outputs[0].Set_i64(outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI64Gt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() > inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI64Gteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() >= inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in lt function returns true if operand 1 is less than oeprand 2.
func opI64Lt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() < inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in lteq function returns true if operand 1 is less than or
// equal to operand 2.
func opI64Lteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() <= inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI64Eq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() == inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI64Uneq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() != inputs[1].Get_i64()
	outputs[0].Set_bool(outV0)
}

// The built-in mod function returns the remainder of operand 1 divided by operand 2.
func opI64Mod(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() % inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI64Rand(inputs []ast.CXValue, outputs []ast.CXValue) {
	minimum := inputs[0].Get_i64()
	maximum := inputs[1].Get_i64()

	r := int(maximum - minimum)
	outV0 := int64(0)
	if r > 0 {
		outV0 = int64(rand.Intn(r) + int(minimum))
	}

	outputs[0].Set_i64(outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI64Bitand(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() & inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI64Bitor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() | inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI64Bitxor(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() ^ inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI64Bitclear(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i64() &^ inputs[1].Get_i64()
	outputs[0].Set_i64(outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI64Bitshl(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_i64() << uint64(inputs[1].Get_i64()))
	outputs[0].Set_i64(outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI64Bitshr(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_i64() >> uint64(inputs[1].Get_i64()))
	outputs[0].Set_i64(outV0)
}

// The built-in max function returns the greatest value of the two operands.
func opI64Max(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i64()
	inpV1 := inputs[1].Get_i64()
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_i64(inpV0)
}

// The built-in min function returns the smallest value of the two operands.
func opI64Min(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i64()
	inpV1 := inputs[1].Get_i64()
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
    outputs[0].Set_i64(inpV0)
}
