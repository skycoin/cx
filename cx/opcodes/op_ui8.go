package opcodes

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opUI8ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatUint(uint64(inputs[0].Get_ui8(prgrm)), 10)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i8 function returns operand 1 casted from type ui8 to type i8.
func opUI8ToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in i16 function returns operand 1 casted from type ui8 to type i16.
func opUI8ToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in i32 function return operand 1 casted from type ui8 to type i32.
func opUI8ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function returns operand 1 casted from type ui8 to type i64.
func opUI8ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui16 function returns the operand 1 casted from type ui8 to type ui16.
func opUI8ToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in ui32 function returns the operand 1 casted from type ui8 to type ui32.
func opUI8ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in ui64 function returns the operand 1 casted from type ui8 to type ui64.
func opUI8ToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in f32 function returns operand 1 casted from type ui8 to type f32.
func opUI8ToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type ui8 to type f64.
func opUI8ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_ui8(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opUI8Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_ui8(prgrm))
}

// The built-in add function returns the sum of two ui8 numbers.
func opUI8Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) + inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in sub function returns the difference of two ui8 numbers.
func opUI8Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) - inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in mul function returns the product of two ui8 numbers.
func opUI8Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) * inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in div function returns the quotient of two ui8 numbers.
func opUI8Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) / inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opUI8Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) > inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opUI8Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) >= inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opUI8Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) < inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opUI8Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) <= inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opUI8Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) == inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opUI8Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) != inputs[1].Get_ui8(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in mod function returns the remainder of operand 1 / operand 2.
func opUI8Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) % inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in rand function returns a pseudo-random number.
func opUI8Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(rand.Int31n(int32(math.MaxUint8)))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitand function returns the bitwise AND of 2 operands.
func opUI8Bitand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) & inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitor function returns the bitwise OR of 2 operands.
func opUI8Bitor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) | inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitxor function returns the bitwise XOR of 2 operands.
func opUI8Bitxor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) ^ inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitclear function returns the bitwise AND NOT of 2 operands.
func opUI8Bitclear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) &^ inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitshl function returns bits of operand 1 shifted to the left
// by number of positions specified in operand 2.
func opUI8Bitshl(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) << inputs[1].Get_ui8(prgrm)
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in bitshr function returns bits of operand 1 shifted to the right
// by number of positions specified in operand 2.
func opUI8Bitshr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_ui8(prgrm) >> uint32(inputs[1].Get_ui8(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in max function returns the biggest of the two operands.
func opUI8Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui8(prgrm)
	inpV1 := inputs[1].Get_ui8(prgrm)
	if inpV1 > inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui8(prgrm, inpV0)
}

// The built-in min function returns the smallest of the two operands.
func opUI8Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_ui8(prgrm)
	inpV1 := inputs[1].Get_ui8(prgrm)
	if inpV1 < inpV0 {
		inpV0 = inpV1
	}
	outputs[0].Set_ui8(prgrm, inpV0)
}
