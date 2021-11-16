package opcodes

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI16ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatInt(int64(inputs[0].Get_i16(prgrm)), 10)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i8 function returns operand 1 casted from type i16 to type i8.
func opI16ToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_i16(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in i32 function returns operand 1 casted from type i16 to type i32.
func opI16ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_i16(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function returns operand 1 casted from type i16 to type i64.
func opI16ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_i16(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui8 function returns operand 1 casted from type i16 to type ui8.
func opI16ToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_i16(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i16 to type ui16.
func opI16ToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_i16(prgrm))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in ui16 function returns the operand 1 casted from type i16 to type ui32.
func opI16ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_i16(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in ui64 function returns the operand 1 casted from type i16 to type ui64.
func opI16ToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_i16(prgrm))
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in f32 function returns operand 1 casted from type i16 to type f32.
func opI16ToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_i16(prgrm))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type i16 to type f64.
func opI16ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_i16(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opI16Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_i16(prgrm))
}

// The built-in add function returns the sum of two i16 numbers.
func opI16Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) + inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in sub function returns the difference of two i16 numbers.
func opI16Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) - inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opI16Neg(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := -inputs[0].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in mul function returns the product of two i16 numbers.
func opI16Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) * inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in div function returns the quotient of two i16 numbers.
func opI16Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) / inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in abs function returns the absolute number of the number.
func opI16Abs(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	V0 := inputs[0].Get_i16(prgrm)
	sign := V0 >> 15
	outV0 := (V0 ^ sign) - sign
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI16Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) > inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI16Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) >= inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI16Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) < inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI16Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) <= inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI16Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) == inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI16Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) != inputs[1].Get_i16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opI16Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) % inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI16Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	minimum := inputs[0].Get_i16(prgrm)
	maximum := inputs[1].Get_i16(prgrm)

	r := int(maximum - minimum)
	outV0 := int16(0)
	if r > 0 {
		outV0 = int16(rand.Intn(r) + int(minimum))
	}

	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI16Bitand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) & inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI16Bitor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) | inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI16Bitxor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) ^ inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI16Bitclear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i16(prgrm) &^ inputs[1].Get_i16(prgrm)
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI16Bitshl(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i16(prgrm) << uint16(inputs[1].Get_i16(prgrm)))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI16Bitshr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i16(prgrm) >> uint16(inputs[1].Get_i16(prgrm)))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in max function returns the biggest of the two operands.
func opI16Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i16(prgrm)
	inpV1 := inputs[1].Get_i16(prgrm)
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_i16(prgrm, inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opI16Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i16(prgrm)
	inpV1 := inputs[1].Get_i16(prgrm)
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_i16(prgrm, inpV0)
}
