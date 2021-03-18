// +build cxfx,!mobile

package cxfx

import (
	. "github.com/skycoin/cx/cx"
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

func opAlCloseDevice(inputs []CXValue, outputs []CXValue) {
	al.CloseDevice()
}

func opAlDeleteBuffers(inputs []CXValue, outputs []CXValue) {
	buffers := toBuffers(inputs[0].GetSlice_i32())
	al.DeleteBuffers(buffers...)
}

func opAlDeleteSources(inputs []CXValue, outputs []CXValue) {
	sources := toSources(inputs[0].GetSlice_i32())
	al.DeleteSources(sources...)
}

func opAlDeviceError(inputs []CXValue, outputs []CXValue) {
	err := al.DeviceError()
	outputs[0].Set_i32(err)
}

func opAlError(inputs []CXValue, outputs []CXValue) {
	err := al.Error()
	outputs[0].Set_i32(err)
}

func opAlExtensions(inputs []CXValue, outputs []CXValue) {
	extensions := al.Extensions()
	outputs[0].Set_str(extensions)
}

func opAlOpenDevice(inputs []CXValue, outputs []CXValue) {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

func opAlPauseSources(inputs []CXValue, outputs []CXValue) {
	sources := toSources(inputs[0].GetSlice_i32())
	al.PauseSources(sources...)
}

func opAlPlaySources(inputs []CXValue, outputs []CXValue) {
	sources := toSources(inputs[0].GetSlice_i32())
	al.PlaySources(sources...)
}

func opAlRenderer(inputs []CXValue, outputs []CXValue) {
	renderer := al.Renderer()
    outputs[0].Set_str(renderer)
}

func opAlRewindSources(inputs []CXValue, outputs []CXValue) {
	sources := toSources(inputs[0].GetSlice_i32())
	al.RewindSources(sources...)
}

func opAlStopSources(inputs []CXValue, outputs []CXValue) {
	sources := toSources(inputs[0].GetSlice_i32())
	al.StopSources(sources...)
}

func opAlVendor(inputs []CXValue, outputs []CXValue) {
	vendor := al.Vendor()
    outputs[0].Set_str(vendor)
}

func opAlVersion(inputs []CXValue, outputs []CXValue) {
	version := al.Version()
    outputs[0].Set_str(version)
}

func opAlGenBuffers(inputs []CXValue, outputs []CXValue) {
	buffers := al.GenBuffers(int(inputs[0].Get_i32()))
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, b := range buffers { // REFACTOR append with copy ?
		obj := FromI32(int32(b))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
    outputs[0].SetSlice(outputSliceOffset)
}

func opAlBufferData(inputs []CXValue, outputs []CXValue) {
	buffer := al.Buffer(inputs[0].Get_i32())
	format := inputs[1].Get_i32()
	data := toBytes(inputs[2].GetSlice_ui8())
	frequency := inputs[3].Get_i32()
	buffer.BufferData(uint32(format), data, frequency)
}

func opAlGenSources(inputs []CXValue, outputs []CXValue) {
	sources := al.GenSources(int(inputs[0].Get_i32()))
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, s := range sources { // REFACTOR append with copy ?
		obj := FromI32(int32(s))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
    outputs[0].SetSlice(outputSliceOffset)
}

func opAlSourceBuffersProcessed(inputs []CXValue, outputs []CXValue) {
	source := al.Source(inputs[0].Get_i32())
    outputs[0].Set_i32(source.BuffersProcessed())
}

func opAlSourceBuffersQueued(inputs []CXValue, outputs []CXValue) {
	source := al.Source(inputs[0].Get_i32())
	outputs[0].Set_i32(source.BuffersQueued())
}

func opAlSourceQueueBuffers(inputs []CXValue, outputs []CXValue) {
	source := al.Source(inputs[0].Get_i32())
	buffers := toBuffers(inputs[1].GetSlice_i32())
	source.QueueBuffers(buffers...)
}

func opAlSourceState(inputs []CXValue, outputs []CXValue) {
	source := al.Source(inputs[0].Get_i32())
    outputs[0].Set_i32(source.State())
}

func opAlSourceUnqueueBuffers(inputs []CXValue, outputs []CXValue) {
	source := al.Source(inputs[0].Get_i32())
	buffers := toBuffers(inputs[1].GetSlice_i32())
	source.UnqueueBuffers(buffers...)
}
