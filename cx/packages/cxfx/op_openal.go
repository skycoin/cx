// +build cxfx

package cxfx

import (
	"bufio"
	"github.com/mjibson/go-dsp/wav"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/util"
	"github.com/skycoin/skycoin/src/cipher/encoder"

	//"golang.org/x/mobile/exp/audio/al"
)

func opAlLoadWav(inputs []ast.CXValue, outputs []ast.CXValue) {
	file, err := util.CXOpenFile(inputs[0].Get_str())
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

	outputs[0].Set_i32(int32(wav.Header.AudioFormat))
	outputs[1].Set_i32(int32(wav.Header.NumChannels))
	outputs[2].Set_i32(int32(wav.Header.SampleRate))
	outputs[3].Set_i32(int32(wav.Header.ByteRate))
	outputs[4].Set_i32(int32(wav.Header.BlockAlign))
	outputs[5].Set_i32(int32(wav.Header.BitsPerSample))
	outputs[6].Set_i32(int32(wav.Samples))
	outputs[7].Set_i64(int64(wav.Duration))

	outputSlicePointer := outputs[8].Offset
	outputSliceOffset := ast.GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(ast.SliceResizeEx(outputSliceOffset, int32(len(data)), 1))
	copy(ast.GetSliceData(outputSliceOffset, 1), data)
	outputs[8].SetSlice(outputSliceOffset)
}

func toBytes(in interface{}) []byte { // REFACTOR : ??
	if in != nil {
		return in.([]byte)
	}
	return nil
}

