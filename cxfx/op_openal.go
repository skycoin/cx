// +build cxfx

package cxfx

import (
	"bufio"
	. "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/mjibson/go-dsp/wav"
	//"golang.org/x/mobile/exp/audio/al"
)

func opAlLoadWav(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(wav.Header.AudioFormat))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(wav.Header.NumChannels))
	WriteI32(GetFinalOffset(fp, expr.Outputs[2]), int32(wav.Header.SampleRate))
	WriteI32(GetFinalOffset(fp, expr.Outputs[3]), int32(wav.Header.ByteRate))
	WriteI32(GetFinalOffset(fp, expr.Outputs[4]), int32(wav.Header.BlockAlign))
	WriteI32(GetFinalOffset(fp, expr.Outputs[5]), int32(wav.Header.BitsPerSample))
	WriteI32(GetFinalOffset(fp, expr.Outputs[6]), int32(wav.Samples))
	WriteI64(GetFinalOffset(fp, expr.Outputs[7]), int64(wav.Duration))

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[8])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(data)), 1))
	copy(GetSliceData(outputSliceOffset, 1), data)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func toBytes(in interface{}) []byte { // REFACTOR : ??
	if in != nil {
		return in.([]byte)
	}
	return nil
}

/*func toBuffers(in interface{}) []al.Buffer { // REFACTOR : ??
	var out []al.Buffer
	var buffers []int32 = in.([]int32)
	for _, b := range buffers {
		out = append(out, al.Buffer(b))
	}
	return out
}*/

/*func toSources(in interface{}) []al.Source { // REFACTOR : ??
	var out []al.Source
	var sources []int32 = in.([]int32)
	for _, s := range sources {
		out = append(out, al.Source(s))
	}
	return out
}*/

/*func opAlCloseDevice(_ *CXProgram) {
	al.CloseDevice()
}

func opAlDeleteBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	buffers := toBuffers(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.DeleteBuffers(buffers...)
}

func opAlDeleteSources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := toSources(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.DeleteSources(sources...)
}

func opAlDeviceError(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	err := al.DeviceError()
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(err))
}

func opAlError(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	err := al.Error()
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(err))
}

func opAlExtensions(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	extensions := al.Extensions()
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(extensions))
}

func opAlOpenDevice(_ *CXProgram) {
	if err := al.OpenDevice(); err != nil {
		panic(err)
	}
}

func opAlPauseSources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := toSources(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.PauseSources(sources...)
}

func opAlPlaySources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := toSources(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.PlaySources(sources...)
}

func opAlRenderer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	renderer := al.Renderer()
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(renderer))
}

func opAlRewindSources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := toSources(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.RewindSources(sources...)
}

func opAlStopSources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := toSources(ReadData(fp, expr.Inputs[0], TYPE_I32))
	al.StopSources(sources...)
}

func opAlVendor(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	vendor := al.Vendor()
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(vendor))
}

func opAlVersion(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	version := al.Version()
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(version))
}

func opAlGenBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	buffers := al.GenBuffers(int(ReadI32(fp, expr.Inputs[0])))
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, b := range buffers { // REFACTOR append with copy ?
		obj := FromI32(int32(b))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
	copy(PROGRAM.Memory[outputSlicePointer:], FromI32(outputSliceOffset))
}

func opAlBufferData(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	buffer := al.Buffer(ReadI32(fp, expr.Inputs[0]))
	format := ReadI32(fp, expr.Inputs[1])
	data := toBytes(ReadData(fp, expr.Inputs[2], TYPE_UI8))
	frequency := ReadI32(fp, expr.Inputs[3])
	buffer.BufferData(uint32(format), data, frequency)
}

func opAlGenSources(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	sources := al.GenSources(int(ReadI32(fp, expr.Inputs[0])))
	outputSlice := expr.Outputs[0]
	outputSlicePointer := GetFinalOffset(fp, outputSlice)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	for _, s := range sources { // REFACTOR append with copy ?
		obj := FromI32(int32(s))
		outputSliceOffset = int32(WriteToSlice(int(outputSliceOffset), obj))
	}
	copy(PROGRAM.Memory[outputSlicePointer:], FromI32(outputSliceOffset))
}

func opAlSourceBuffersProcessed(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.BuffersProcessed()))
}

func opAlSourceBuffersQueued(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.BuffersQueued()))
}

func opAlSourceQueueBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	buffers := toBuffers(ReadData(fp, expr.Inputs[1], TYPE_I32))
	source.QueueBuffers(buffers...)
}

func opAlSourceState(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromI32(source.State()))
}

func opAlSourceUnqueueBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	source := al.Source(ReadI32(fp, expr.Inputs[0]))
	buffers := toBuffers(ReadData(fp, expr.Inputs[1], TYPE_I32))
	source.UnqueueBuffers(buffers...)
}*/
