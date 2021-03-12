// +build cxfx,!mobile

package cxfx

import (
	//	"bufio"
	. "github.com/skycoin/cx/cx"
	//"github.com/amherag/skycoin/src/cipher/encoder"
	//	"github.com/mjibson/go-dsp/wav"
	"golang.org/x/mobile/exp/audio/al"
)

/*func opAlLoadWav(expr *CXExpression, fp int) {
	file, err := CXOpenFile(ReadStr(fp, expr.Inputs[0]))
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)

	wav, err := wav.New(reader)
	if err != nil {
		panic(err)
	}

	samples, err := wav.ReadSamples(wav.Samples)
	if err != nil {
		panic(err)
	}

	data := encoder.Serialize(samples)

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(wav.Header.AudioFormat)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromI32(int32(wav.Header.NumChannels)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[2]), FromI32(int32(wav.Header.SampleRate)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[3]), FromI32(int32(wav.Header.ByteRate)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[4]), FromI32(int32(wav.Header.BlockAlign)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[5]), FromI32(int32(wav.Header.BitsPerSample)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[6]), FromI32(int32(wav.Samples)))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[7]), FromI64(int64(wav.Duration)))

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[8])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(data)), 1))
	copy(GetSliceData(outputSliceOffset, 1), data)
	copy(PROGRAM.Memory[outputSlicePointer:], FromI32(outputSliceOffset))
}*/

/*func toBytes(in interface{}) []byte { // REFACTOR : ??
	if in != nil {
		return in.([]byte)
	}
	return nil
}
*/
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

func opAlCloseDevice(expr *CXExpression, fp int) {
	al.CloseDevice()
}

func opAlDeleteBuffers(expr *CXExpression, fp int) {
	buffers := toBuffers(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.DeleteBuffers(buffers...)
}

func opAlDeleteSources(expr *CXExpression, fp int) {
	sources := toSources(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.DeleteSources(sources...)
}

func opAlDeviceError(expr *CXExpression, fp int) {
	err := al.DeviceError()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), err)
}

func opAlError(expr *CXExpression, fp int) {
	err := al.Error()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), err)
}

func opAlExtensions(expr *CXExpression, fp int) {
	extensions := al.Extensions()
	WriteString(fp, extensions, expr.Outputs[0])
}

func opAlOpenDevice(expr *CXExpression, fp int) {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

func opAlPauseSources(expr *CXExpression, fp int) {
	sources := toSources(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.PauseSources(sources...)
}

func opAlPlaySources(expr *CXExpression, fp int) {
	sources := toSources(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.PlaySources(sources...)
}

func opAlRenderer(expr *CXExpression, fp int) {
	renderer := al.Renderer()
	WriteString(fp, renderer, expr.Outputs[0])
}

func opAlRewindSources(expr *CXExpression, fp int) {
	sources := toSources(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.RewindSources(sources...)
}

func opAlStopSources(expr *CXExpression, fp int) {
	sources := toSources(ReadData_i32(fp, expr.Inputs[0], TYPE_I32))
	al.StopSources(sources...)
}

func opAlVendor(expr *CXExpression, fp int) {
	vendor := al.Vendor()
	WriteString(fp, vendor, expr.Outputs[0])
}

func opAlVersion(expr *CXExpression, fp int) {
	version := al.Version()
	WriteString(fp, version, expr.Outputs[0])
}

func opAlGenBuffers(expr *CXExpression, fp int) {
	buffers := al.GenBuffers(int(ReadI32(fp, expr.Inputs[0])))
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, b := range buffers { // REFACTOR append with copy ?
		obj := FromI32(int32(b))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
	WriteI32(outputSlicePointer, outputSliceOffset)
	//copy(PROGRAM.Memory[outputSlicePointer:], FromI32(outputSliceOffset))
}

func opAlBufferData(expr *CXExpression, fp int) {
	buffer := al.Buffer(ReadI32(fp, expr.Inputs[0]))
	format := ReadI32(fp, expr.Inputs[1])
	data := toBytes(ReadData_ui8(fp, expr.Inputs[2], TYPE_UI8))
	frequency := ReadI32(fp, expr.Inputs[3])
	buffer.BufferData(uint32(format), data, frequency)
}

func opAlGenSources(expr *CXExpression, fp int) {
	sources := al.GenSources(int(ReadI32(fp, expr.Inputs[0])))
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, s := range sources { // REFACTOR append with copy ?
		obj := FromI32(int32(s))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
	WriteI32(outputSlicePointer, outputSliceOffset)
	//copy(PROGRAM.Memory[outputSlicePointer:], FromI32(outputSliceOffset))
}

func opAlSourceBuffersProcessed(expr *CXExpression, fp int) {
	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.BuffersProcessed()))
}

func opAlSourceBuffersQueued(expr *CXExpression, fp int) {
	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.BuffersQueued()))
}

func opAlSourceQueueBuffers(expr *CXExpression, fp int) {
	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	buffers := toBuffers(ReadData_i32(fp, expr.Inputs[1], TYPE_I32))
	source.QueueBuffers(buffers...)
}

func opAlSourceState(expr *CXExpression, fp int) {
	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.State()))
}

func opAlSourceUnqueueBuffers(expr *CXExpression, fp int) {
	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	buffers := toBuffers(ReadData_i32(fp, expr.Inputs[1], TYPE_I32))
	source.UnqueueBuffers(buffers...)
}
