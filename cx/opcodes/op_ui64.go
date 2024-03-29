package opcodes

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI64ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatUint(inputs[0].Get_ui64(prgrm), 10)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i8 function returns operand 1 casted from type ui64 to type i8.
func opUI64ToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in i16 function retruns operand 1 casted from type ui64 to i16.
func opUI64ToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in i32 function return operand 1 casted from type ui64 to type i32.
func opUI64ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function return operand 1 casted from type ui64 to type i64.
func opUI64ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui8 function returns operand 1 casted from type ui64 to type ui8.
func opUI64ToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui64 to type ui16.
func opUI64ToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui64 to type ui32.
func opUI64ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in f32 function returns operand 1 casted from type ui64 to type f32.
func opUI64ToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type ui64 to type f64.
func opUI64ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_ui64(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI64Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_ui64(prgrm))
}

// The built-in add function returns the sum of two ui64 numbers.
func opUI64Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) + inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in sub function returns the difference of two ui64 numbers.
func opUI64Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) - inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in mul function returns the product of two ui64 numbers.
func opUI64Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) * inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in div function returns the quotient of two ui64 numbers.
func opUI64Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) / inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI64Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) > inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI64Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) >= inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI64Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) < inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI64Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) <= inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI64Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) == inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI64Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) != inputs[1].Get_ui64(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI64Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) % inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI64Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := rand.Uint64()
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI64Bitand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) & inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI64Bitor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) | inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI64Bitxor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) ^ inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI64Bitclear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) &^ inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI64Bitshl(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) << inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI64Bitshr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui64(prgrm) >> inputs[1].Get_ui64(prgrm)
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI64Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui64(prgrm)
	inpV1 := inputs[1].Get_ui64(prgrm)
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui64(prgrm, inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI64Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui64(prgrm)
	inpV1 := inputs[1].Get_ui64(prgrm)
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui64(prgrm, inpV0)
}
