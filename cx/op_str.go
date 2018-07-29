package base

import (
	"fmt"
)

// func op_str_str(expr *CXExpression, stack *CXStack, fp int) {
// 	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
// 	out1Offset := GetFinalOffset(stack, fp, out1, MEM_WRITE)

// 	switch out1.Type {
// 	case TYPE_BYTE:
// 		WriteMemory(stack, out1Offset, out1, FromByte(byte(ReadI32(stack, fp, inp1))))
// 		// case TYPE_STR: WriteMemory(stack, out1Offset, out1, FromStr(strconv.Itoa(ReadI32(stack, fp, inp1))))
// 	case TYPE_I32:
// 		WriteMemory(stack, out1Offset, out1, FromI32(ReadI32(stack, fp, inp1)))
// 	case TYPE_I64:
// 		WriteMemory(stack, out1Offset, out1, FromI64(int64(ReadI32(stack, fp, inp1))))
// 	case TYPE_F32:
// 		WriteMemory(stack, out1Offset, out1, FromF32(float32(ReadI32(stack, fp, inp1))))
// 	case TYPE_F64:
// 		WriteMemory(stack, out1Offset, out1, FromF64(float64(ReadI32(stack, fp, inp1))))
// 	}
// }

func op_str_print(expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(stack, fp, inp1))
}

func op_str_eq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(stack, fp, inp1) == ReadStr(stack, fp, inp2))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
