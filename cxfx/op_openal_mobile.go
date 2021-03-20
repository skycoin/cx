// +build cxfx,mobile

package cxfx

import (
	"github.com/skycoin/cx/cx"
	/*"bufio"
	  "github.com/amherag/skycoin/src/cipher/encoder"
	  "github.com/mjibson/go-dsp/wav"
	  "golang.org/x/mobile/exp/audio/al"*/)

func opAlCloseDevice(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteBuffers(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteSources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeviceError(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opAlError(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opAlExtensions(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteString(fp, "CX_RUNTIME_NOT_IMPLEMENTED", expr.Outputs[0])
}

func opAlOpenDevice(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPauseSources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPlaySources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlRenderer(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteString(fp, "CX_RUNTIME_NOT_IMPLEMENTED", expr.Outputs[0])
}

func opAlRewindSources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlStopSources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlVendor(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteString(fp, "CX_RUNTIME_NOT_IMPLEMENTED", expr.Outputs[0])
}

func opAlVersion(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteString(fp, "CX_RUNTIME_NOT_IMPLEMENTED", expr.Outputs[0])
}

func opAlGenBuffers(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opAlBufferData(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlGenSources(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func opAlSourceBuffersProcessed(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceBuffersQueued(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceQueueBuffers(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceState(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceUnqueueBuffers(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
