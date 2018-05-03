package compiled

import (
	"fmt"
	// "strconv"
	// "math"
	// "math/rand"
)

func byte_byte (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)
	switch out1.Type {
	case TYPE_BYTE: WriteMemory(stack, out1Offset, out1, FromByte(byte(ReadI32(stack, fp, inp1))))
	// case TYPE_STR: WriteMemory(stack, out1Offset, out1, FromStr(strconv.Itoa(ReadI32(stack, fp, inp1))))
	case TYPE_I32: WriteMemory(stack, out1Offset, out1, FromI32(ReadI32(stack, fp, inp1)))
	case TYPE_I64: WriteMemory(stack, out1Offset, out1, FromI64(int64(ReadI32(stack, fp, inp1))))
	case TYPE_F32: WriteMemory(stack, out1Offset, out1, FromF32(float32(ReadI32(stack, fp, inp1))))
	case TYPE_F64: WriteMemory(stack, out1Offset, out1, FromF64(float64(ReadI32(stack, fp, inp1))))
	}
}

func byte_print (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadByte(stack, fp, inp1))
}
