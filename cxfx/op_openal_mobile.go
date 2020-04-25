// +build cxfx,mobile

package cxfx

import (
	. "github.com/SkycoinProject/cx/cx"
	/*"bufio"
	  "github.com/amherag/skycoin/src/cipher/encoder"
	  "github.com/mjibson/go-dsp/wav"
	  "golang.org/x/mobile/exp/audio/al"*/)

func opAlCloseDevice(_ *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteBuffers(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteSources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeviceError(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opAlError(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opAlExtensions(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr("CX_RUNTIME_NOT_IMPLEMENTED"))
}

func opAlOpenDevice(_ *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPauseSources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPlaySources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlRenderer(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr("CX_RUNTIME_NOT_IMPLEMENTED"))
}

func opAlRewindSources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlStopSources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlVendor(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr("CX_RUNTIME_NOT_IMPLEMENTED"))
}

func opAlVersion(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr("CX_RUNTIME_NOT_IMPLEMENTED"))
}

func opAlGenBuffers(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opAlBufferData(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlGenSources(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opAlSourceBuffersProcessed(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceBuffersQueued(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceQueueBuffers(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceState(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceUnqueueBuffers(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
