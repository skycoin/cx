package helper

import (
	"math"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const DEBUG_READ_INPUT_LENGTH bool = false

//TODO: replace with unsafe

// func DeserializeBool(b []byte) bool {
// 	switch b[0] {
// 	case 0:
// 		return false
// 	case 1:
// 		return true
// 	default:
// 		panic(encoder.ErrInvalidBool)
// 		//return false [2020 Jun 07 (ReewassSquared)]: does nothing
// 	}
// }

func Deserialize_bool(b []byte) bool {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 1 {
		panic("byte invalid length")
	}
	return b[0] != 0
}

func Deserialize_i8(b []byte) int8 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 1 {
		panic("byte invalid length")
	}
	return int8(b[0])
}

func Deserialize_i16(b []byte) int16 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 2 {
		panic("byte invalid length")
	}
	return int16(b[0]) | int16(b[1])<<8
}

func Deserialize_i32(b []byte) int32 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 4 {
		panic("byte invalid length")
	}
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func Deserialize_i64(b []byte) int64 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 8 {
		panic("byte invalid length")
	}
	return int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 |
		int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
}

func Deserialize_ui8(b []byte) uint8 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 1 {
		panic("byte invalid length")
	}
	return uint8(b[0])
}

func Deserialize_ui16(b []byte) uint16 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 2 {
		panic("byte invalid length")
	}
	return uint16(b[0]) | uint16(b[1])<<8
}

func Deserialize_ui32(b []byte) uint32 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 4 {
		panic("byte invalid length")
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func Deserialize_ui64(b []byte) uint64 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 8 {
		panic("byte invalid length")
	}
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func Deserialize_f32(b []byte) float32 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 4 {
		panic("byte invalid length")
	}
	return math.Float32frombits(Deserialize_ui32(b))
}

func Deserialize_f64(b []byte) float64 {
	if DEBUG_READ_INPUT_LENGTH && len(b) != 8 {
		panic("byte invalid length")
	}
	return math.Float64frombits(Deserialize_ui64(b))
}

func DeserializeRaw(byts []byte, item interface{}) {
	_, err := encoder.DeserializeRaw(byts, item)
	if err != nil {
		panic(err)
	}
}
