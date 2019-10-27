package cxcore

// import "fmt"

// func opI8Print(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1 := expr.Inputs[0]
// 	fmt.Println(ReadI8(fp, inp1))
// }

// func opI8Add(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) + ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Sub(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) - ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Mul(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) * ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Div(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) / ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Gt(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) > ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Gteq(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) >= ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Lt(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) < ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Lteq(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) <= ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Eq(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) == ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Uneq(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromBool(ReadI8(fp, inp1) != ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Bitand(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) & ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Bitor(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) | ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Bitxor(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) ^ ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }

// func opI8Bitclear(prgrm *CXProgram) {
//	expr := prgrm.GetExpr()
//	fp := prgrm.GetFramePointer()

// 	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
// 	outB1 := FromI8(ReadI8(fp, inp1) &^ ReadI8(fp, inp2))
// 	WriteMemory(GetFinalOffset(fp, out1), outB1)
// }
