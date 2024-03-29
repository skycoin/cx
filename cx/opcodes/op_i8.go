package opcodes

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opI8ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatInt(int64(inputs[0].Get_i8(prgrm)), 10)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i16 function returns operand 1 casted from type i8 to type i16.
func opI8ToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_i8(prgrm))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in i32 function returns operand 1 casted from type i8 to type i32.
func opI8ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_i8(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function returns operand 1 casted from type i8 to type i64.
func opI8ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_i8(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui8 function returns operand 1 casted from type i8 to type ui8.
func opI8ToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_i8(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in ui16 function returns operand 1 casted from type i8 to type ui16.
func opI8ToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_i8(prgrm))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in ui32 function returns operand 1 casted from type i8 to type ui32.
func opI8ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_i8(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in ui64 function returns operand 1 casted from type i8 to type ui64.
func opI8ToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_i8(prgrm))
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in f32 function returns operand 1 casted from type i8 to type f32.
func opI8ToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_i8(prgrm))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type i8 to type uf64.
func opI8ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_i8(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opI8Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_i8(prgrm))
}

// The built-in add function returns the sum of two i8 numbers.
func opI8Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) + inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in sub function returns the difference of two i8 numbers.
func opI8Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) - inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opI8Neg(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := -inputs[0].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in mul function returns the product of two i8 numbers.
func opI8Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) * inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in div function returns the quotient of two i8 numbers.
func opI8Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) / inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in abs function returns the absolute number of the number.
func opI8Abs(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i8(prgrm)
	sign := inpV0 >> 7
	outV0 := (inpV0 ^ sign) - sign
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI8Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) > inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI8Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) >= inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI8Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) < inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI8Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) <= inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI8Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) == inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI8Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) != inputs[1].Get_i8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opI8Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) % inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in rand function returns a pseudo random number in [operand 1, operand 2).
func opI8Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	minimum := inputs[0].Get_i8(prgrm)
	maximum := inputs[1].Get_i8(prgrm)

	r := int(maximum - minimum)
	outV0 := int8(0)
	if r > 0 {
		outV0 = int8(rand.Intn(r) + int(minimum))
	}

	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opI8Bitand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) & inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opI8Bitor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) | inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opI8Bitxor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) ^ inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opI8Bitclear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) &^ inputs[1].Get_i8(prgrm)
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opI8Bitshl(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) << uint8(inputs[1].Get_i8(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opI8Bitshr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_i8(prgrm) >> uint8(inputs[1].Get_i8(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in max function returns the biggest of the two operands.
func opI8Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i8(prgrm)
	inpV1 := inputs[1].Get_i8(prgrm)
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_i8(prgrm, inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opI8Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i8(prgrm)
	inpV1 := inputs[1].Get_i8(prgrm)
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_i8(prgrm, inpV0)
}
