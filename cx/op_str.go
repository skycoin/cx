package base

import (
	"fmt"
)

// func op_str_str(expr *CXExpression, fp int) {
// 	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
// 	out1Offset := GetFinalOffset(fp, out1)

// 	switch out1.Type {
// 	case TYPE_BYTE:
// 		WriteMemory(out1Offset, out1, FromByte(byte(ReadI32(fp, inp1))))
// 		// case TYPE_STR: WriteMemory(out1Offset, out1, FromStr(strconv.Itoa(ReadI32(fp, inp1))))
// 	case TYPE_I32:
// 		WriteMemory(out1Offset, out1, FromI32(ReadI32(fp, inp1)))
// 	case TYPE_I64:
// 		WriteMemory(out1Offset, out1, FromI64(int64(ReadI32(fp, inp1))))
// 	case TYPE_F32:
// 		WriteMemory(out1Offset, out1, FromF32(float32(ReadI32(fp, inp1))))
// 	case TYPE_F64:
// 		WriteMemory(out1Offset, out1, FromF64(float64(ReadI32(fp, inp1))))
// 	}
// }

func op_str_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(fp, inp1))
}

func op_str_eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(fp, inp1) == ReadStr(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
