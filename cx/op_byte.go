package base

import (
	"fmt"
	// "strconv"
	// "math"
	// "math/rand"
)

func op_byte_byte(expr *CXExpression, mem []byte, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(mem, fp, out1, MEM_WRITE)
	switch out1.Type {
	case TYPE_BYTE:
		WriteMemory(mem, out1Offset, FromByte(byte(ReadI32(mem, fp, inp1))))
	// case TYPE_STR: WriteMemory(mem, out1Offset, FromStr(strconv.Itoa(ReadI32(mem, fp, inp1))))
	case TYPE_I32:
		WriteMemory(mem, out1Offset, FromI32(ReadI32(mem, fp, inp1)))
	case TYPE_I64:
		WriteMemory(mem, out1Offset, FromI64(int64(ReadI32(mem, fp, inp1))))
	case TYPE_F32:
		WriteMemory(mem, out1Offset, FromF32(float32(ReadI32(mem, fp, inp1))))
	case TYPE_F64:
		WriteMemory(mem, out1Offset, FromF64(float64(ReadI32(mem, fp, inp1))))
	}
}

func op_byte_print(expr *CXExpression, mem []byte, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadByte(mem, fp, inp1))
}
