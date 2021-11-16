// +build cxfx,!mobile

package cxfx

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
	"golang.org/x/mobile/exp/audio/al"
)

func toBuffers(in interface{}) []al.Buffer { // REFACTOR : ??
	var out []al.Buffer
	var buffers []int32 = in.([]int32)
	for _, b := range buffers {
		out = append(out, al.Buffer(b))
	}
	return out
}

func toSources(in interface{}) []al.Source { // REFACTOR : ??
	var out []al.Source
	var sources []int32 = in.([]int32)
	for _, s := range sources {
		out = append(out, al.Source(s))
	}
	return out
}

func opAlCloseDevice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	al.CloseDevice()
}

func opAlDeleteBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	buffers := toBuffers(inputs[0].GetSlice_i32(prgrm))
	al.DeleteBuffers(buffers...)
}

func opAlDeleteSources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := toSources(inputs[0].GetSlice_i32(prgrm))
	al.DeleteSources(sources...)
}

func opAlDeviceError(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	err := al.DeviceError()
	outputs[0].Set_i32(prgrm, err)
}

func opAlError(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	err := al.Error()
	outputs[0].Set_i32(prgrm, err)
}

func opAlExtensions(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	extensions := al.Extensions()
	outputs[0].Set_str(prgrm, extensions)
}

func opAlOpenDevice(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

func opAlPauseSources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := toSources(inputs[0].GetSlice_i32(prgrm))
	al.PauseSources(sources...)
}

func opAlPlaySources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := toSources(inputs[0].GetSlice_i32(prgrm))
	al.PlaySources(sources...)
}

func opAlRenderer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	renderer := al.Renderer()
	outputs[0].Set_str(prgrm, renderer)
}

func opAlRewindSources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := toSources(inputs[0].GetSlice_i32(prgrm))
	al.RewindSources(sources...)
}

func opAlStopSources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := toSources(inputs[0].GetSlice_i32(prgrm))
	al.StopSources(sources...)
}

func opAlVendor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	vendor := al.Vendor()
	outputs[0].Set_str(prgrm, vendor)
}

func opAlVersion(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	version := al.Version()
	outputs[0].Set_str(prgrm, version)
}

func opAlGenBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	buffers := al.GenBuffers(int(inputs[0].Get_i32(prgrm)))
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := types.Read_ptr(prgrm.Memory, outputSlicePointer)
	for _, b := range buffers { // REFACTOR append with copy ?
		var obj [4]byte
		types.Write_i32(obj[:], 0, int32(b))
		outputSliceOffset = ast.WriteToSlice(prgrm, outputSliceOffset, obj[:])
	}
	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

func opAlBufferData(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	buffer := al.Buffer(inputs[0].Get_i32(prgrm))
	format := inputs[1].Get_i32(prgrm)
	data := toBytes(inputs[2].GetSlice_ui8(prgrm))
	frequency := inputs[3].Get_i32(prgrm)
	buffer.BufferData(uint32(format), data, frequency)
}

func opAlGenSources(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	sources := al.GenSources(int(inputs[0].Get_i32(prgrm)))
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := types.Read_ptr(prgrm.Memory, outputSlicePointer)
	for _, s := range sources { // REFACTOR append with copy ?
		var obj [4]byte
		types.Write_i32(obj[:], 0, int32(s))
		outputSliceOffset = ast.WriteToSlice(prgrm, outputSliceOffset, obj[:])
	}
	outputs[0].Set_ptr(prgrm, outputSliceOffset)
}

func opAlSourceBuffersProcessed(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	source := al.Source(inputs[0].Get_i32(prgrm))
	outputs[0].Set_i32(prgrm, source.BuffersProcessed())
}

func opAlSourceBuffersQueued(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	source := al.Source(inputs[0].Get_i32(prgrm))
	outputs[0].Set_i32(prgrm, source.BuffersQueued())
}

func opAlSourceQueueBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	source := al.Source(inputs[0].Get_i32(prgrm))
	buffers := toBuffers(inputs[1].GetSlice_i32(prgrm))
	source.QueueBuffers(buffers...)
}

func opAlSourceState(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	source := al.Source(inputs[0].Get_i32(prgrm))
	outputs[0].Set_i32(prgrm, source.State())
}

func opAlSourceUnqueueBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	source := al.Source(inputs[0].Get_i32(prgrm))
	buffers := toBuffers(inputs[1].GetSlice_i32(prgrm))
	source.UnqueueBuffers(buffers...)
}
