package opcodes

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI16ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatUint(uint64(inputs[0].Get_ui16(prgrm)), 10)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i8 function returns operand 1 casted from type ui16 to type i8.
func opUI16ToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in i16 function returns operand 1 casted from type ui16 to type i16.
func opUI16ToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in i32 function return operand 1 casted from type ui16 to type i32.
func opUI16ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function returns operand 1 casted from type ui16 to type i64.
func opUI16ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui8 function returns operand 1 casted from type ui16 to type ui8.
func opUI16ToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui16 to type ui32.
func opUI16ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in ui64 function returns the operand 1 casted from type ui16 to type ui64.
func opUI16ToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in f32 function returns operand 1 casted from type ui16 to type f32.
func opUI16ToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type ui16 to type f64.
func opUI16ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_ui16(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI16Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_ui16(prgrm))
}

// The built-in add function returns the sum of two ui16 numbers.
func opUI16Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) + inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in sub function returns the difference of two ui16 numbers.
func opUI16Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) - inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in mul function returns the product of two ui16 numbers.
func opUI16Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) * inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in div function returns the quotient of two ui16 numbers.
func opUI16Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) / inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI16Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) > inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI16Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) >= inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI16Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) < inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI16Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) <= inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI16Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) == inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI16Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) != inputs[1].Get_ui16(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI16Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) % inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI16Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(rand.Int31n(int32(math.MaxUint16)))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI16Bitand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) & inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI16Bitor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) | inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI16Bitxor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) ^ inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI16Bitclear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) &^ inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI16Bitshl(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) << inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI16Bitshr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui16(prgrm) >> inputs[1].Get_ui16(prgrm)
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in max function returns the biggest of the two operands..
func opUI16Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui16(prgrm)
	inpV1 := inputs[1].Get_ui16(prgrm)
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui16(prgrm, inpV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI16Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui16(prgrm)
	inpV1 := inputs[1].Get_ui16(prgrm)
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui16(prgrm, inpV0)
}
