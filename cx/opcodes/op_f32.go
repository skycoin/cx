package opcodes

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
)

// The built-in str function returns the base 10 string representation of operand 1.
func opF32ToStr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := strconv.FormatFloat(float64(inputs[0].Get_f32(prgrm)), 'f', -1, 32)
	outputs[0].Set_str(prgrm, outV0)
}

// The built-in i8 function returns operand 1 casted from type f32 to type i8.
func opF32ToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int8(inputs[0].Get_f32(prgrm))
	outputs[0].Set_i8(prgrm, outV0)
}

// The built-in i16 function returns operand 1 casted from type f32 to type i16.
func opF32ToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int16(inputs[0].Get_f32(prgrm))
	outputs[0].Set_i16(prgrm, outV0)
}

// The built-in i32 function return operand 1 casted from type f32 to type i32.
func opF32ToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(inputs[0].Get_f32(prgrm))
	outputs[0].Set_i32(prgrm, outV0)
}

// The built-in i64 function returns operand 1 casted from type f32 to type i64.
func opF32ToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int64(inputs[0].Get_f32(prgrm))
	outputs[0].Set_i64(prgrm, outV0)
}

// The built-in ui8 function returns operand 1 casted from type f32 to type ui8.
func opF32ToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint8(inputs[0].Get_f32(prgrm))
	outputs[0].Set_ui8(prgrm, outV0)
}

// The built-in ui16 function returns the operand 1 casted from type f32 to type ui16.
func opF32ToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint16(inputs[0].Get_f32(prgrm))
	outputs[0].Set_ui16(prgrm, outV0)
}

// The built-in ui32 function returns the operand 1 casted from type f32 to type ui32.
func opF32ToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint32(inputs[0].Get_f32(prgrm))
	outputs[0].Set_ui32(prgrm, outV0)
}

// The built-in ui64 function returns the operand 1 casted from type f32 to type ui64.
func opF32ToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := uint64(inputs[0].Get_f32(prgrm))
	outputs[0].Set_ui64(prgrm, outV0)
}

// The built-in f64 function returns operand 1 casted from type f32 to type f64.
func opF32ToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float64(inputs[0].Get_f32(prgrm))
	outputs[0].Set_f64(prgrm, outV0)
}

// The built-in isnan function returns true if operand is nan value.
func opF32Isnan(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := math.IsNaN(float64(inputs[0].Get_f32(prgrm)))
	outputs[0].Set_bool(prgrm, outV0)
}

// The print built-in function formats its arguments and prints them.
func opF32Print(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_f32(prgrm))
}

// The built-in add function returns the sum of the two operands.
func opF32Add(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) + inputs[1].Get_f32(prgrm)
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in sub function returns the difference between the two operands.
func opF32Sub(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) - inputs[1].Get_f32(prgrm)
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in neg function returns the opposite of operand 1.
func opF32Neg(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := -inputs[0].Get_f32(prgrm)
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in mul function returns the product of the two operands.
func opF32Mul(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) * inputs[1].Get_f32(prgrm)
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in div function returns the quotient between the two operands.
func opF32Div(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) / inputs[1].Get_f32(prgrm)
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in mod function return the floating-point remainder of operand 1 divided by operand 2.
func opF32Mod(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Mod(float64(inputs[0].Get_f32(prgrm)), float64(inputs[1].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in abs function returns the absolute value of the operand.
func opF32Abs(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Abs(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in pow function returns x**n for n>0 otherwise 1.
func opF32Pow(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Pow(float64(inputs[0].Get_f32(prgrm)), float64(inputs[1].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opF32Gt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) > inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in gteq function returns true if the operand 1 is greater than or
// equal to operand 2.
func opF32Gteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) >= inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opF32Lt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) < inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in lteq function returns true if operand 1 is less than or
// equal to operand 2.
func opF32Lteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) <= inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opF32Eq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) == inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in uneq function returns true operand1 is different from operand 2.
func opF32Uneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_f32(prgrm) != inputs[1].Get_f32(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

// The built-in rand function returns a pseudo-random number in [0.0,1.0) from the default Source
func opF32Rand(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_f32(prgrm, rand.Float32())
}

// The built-in acos function returns the arc cosine of the operand.
func opF32Acos(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Acos(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in cos function returns the cosine of the operand.
func opF32Cos(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Cos(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in asin function returns the arc sine of the operand.
func opF32Asin(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Asin(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in sin function returns the sine of the operand.
func opF32Sin(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Sin(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in sqrt function returns the square root of the operand.
func opF32Sqrt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Sqrt(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in log function returns the natural logarithm of the operand.
func opF32Log(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Log(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in log2 function returns the 2-logarithm of the operand.
func opF32Log2(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Log2(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in log10 function returns the 10-logarithm of the operand.
func opF32Log10(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Log10(float64(inputs[0].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in max function returns the largest value of the two operands.
func opF32Max(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Max(float64(inputs[0].Get_f32(prgrm)), float64(inputs[1].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}

// The built-in min function returns the smallest value of the two operands.
func opF32Min(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := float32(math.Min(float64(inputs[0].Get_f32(prgrm)), float64(inputs[1].Get_f32(prgrm))))
	outputs[0].Set_f32(prgrm, outV0)
}
