package base

import (
	"fmt"
)

// func op_str_str(expr *CXExpression, mem []byte, fp int) {
// 	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
// 	out1Offset := GetFinalOffset(mem, fp, out1, MEM_WRITE)

// 	switch out1.Type {
// 	case TYPE_BYTE:
// 		WriteMemory(mem, out1Offset, out1, FromByte(byte(ReadI32(mem, fp, inp1))))
// 		// case TYPE_STR: WriteMemory(mem, out1Offset, out1, FromStr(strconv.Itoa(ReadI32(mem, fp, inp1))))
// 	case TYPE_I32:
// 		WriteMemory(mem, out1Offset, out1, FromI32(ReadI32(mem, fp, inp1)))
// 	case TYPE_I64:
// 		WriteMemory(mem, out1Offset, out1, FromI64(int64(ReadI32(mem, fp, inp1))))
// 	case TYPE_F32:
// 		WriteMemory(mem, out1Offset, out1, FromF32(float32(ReadI32(mem, fp, inp1))))
// 	case TYPE_F64:
// 		WriteMemory(mem, out1Offset, out1, FromF64(float64(ReadI32(mem, fp, inp1))))
// 	}
// }

func op_str_print(expr *CXExpression, mem []byte, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(mem, fp, inp1))
}

func op_str_eq(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(mem, fp, inp1) == ReadStr(mem, fp, inp2))
	WriteMemory(mem, GetFinalOffset(mem, fp, out1, MEM_WRITE), outB1)
}
