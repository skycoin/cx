// +build cxfx,mobile

package cxfx

import (
	"github.com/skycoin/cx/cx"
	/*"bufio"
	  "github.com/amherag/skycoin/src/cipher/encoder"
	  "github.com/mjibson/go-dsp/wav"
	  "golang.org/x/mobile/exp/audio/al"*/)

func opAlCloseDevice(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteBuffers(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeleteSources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlDeviceError(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_i32(0)
}

func opAlError(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_i32(0)
}

func opAlExtensions(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_str("CX_RUNTIME_NOT_IMPLEMENTED")
}

func opAlOpenDevice(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPauseSources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlPlaySources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlRenderer(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_str("CX_RUNTIME_NOT_IMPLEMENTED")
}

func opAlRewindSources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlStopSources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlVendor(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_str("CX_RUNTIME_NOT_IMPLEMENTED")
}

func opAlVersion(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_str("CX_RUNTIME_NOT_IMPLEMENTED")
}

func opAlGenBuffers(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := ast.GetPointerOffset(int32(outputSlicePointer))
    outputs[0].SetSlice(outputSliceOffset)
}

func opAlBufferData(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlGenSources(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := ast.GetPointerOffset(int32(outputSlicePointer))
    outputs[0].SetSlice(outputSliceOffset)
}

func opAlSourceBuffersProcessed(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceBuffersQueued(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceQueueBuffers(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceState(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opAlSourceUnqueueBuffers(inputs []ast.CXValue, outputs []ast.CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
