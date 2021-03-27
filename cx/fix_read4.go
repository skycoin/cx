package cxcore

import "github.com/skycoin/cx/cx/helper"

// ReadData ...

func readDataI8(bytes []byte) (out []int8) {
	count := len(bytes)
	if count > 0 {
		out = make([]int8, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_i8(bytes[i:])
		}
	}
	return
}

func readDataUI8(bytes []byte) (out []uint8) {
	count := len(bytes)
	if count > 0 {
		out = make([]uint8, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_ui8(bytes[i:])
		}
	}
	return
}

func readDataI16(bytes []byte) (out []int16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]int16, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_i16(bytes[i*2:])
		}
	}
	return
}

func readDataUI16(bytes []byte) (out []uint16) {
	count := len(bytes) / 2
	if count > 0 {
		out = make([]uint16, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_ui16(bytes[i*2:])
		}
	}
	return
}

func readDataI32(bytes []byte) (out []int32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]int32, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_i32(bytes[i*4:])
		}
	}
	return
}

func readDataUI32(bytes []byte) (out []uint32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]uint32, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_ui32(bytes[i*4:])
		}
	}
	return
}

func readDataI64(bytes []byte) (out []int64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]int64, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_i64(bytes[i*8:])
		}
	}
	return
}

func readDataUI64(bytes []byte) (out []uint64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]uint64, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_ui64(bytes[i*8:])
		}
	}
	return
}

func readDataF32(bytes []byte) (out []float32) {
	count := len(bytes) / 4
	if count > 0 {
		out = make([]float32, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_f32(bytes[i*4:])
		}
	}
	return
}

func readDataF64(bytes []byte) (out []float64) {
	count := len(bytes) / 8
	if count > 0 {
		out = make([]float64, count)
		for i := 0; i < count; i++ {
			out[i] = helper.Deserialize_f64(bytes[i*8:])
		}
	}
	return
}