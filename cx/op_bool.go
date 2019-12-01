package cxcore

import (
	"fmt"
)

func opBoolPrint(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(fp, inp1))
}

func opBoolEqual(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) == ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolUnequal(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) != ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolNot(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(!ReadBool(fp, inp1))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolAnd(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) && ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolOr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) || ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
