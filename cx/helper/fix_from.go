package helper

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	//"fmt"

	"math"
)

/*
NOTE:
- "FromStr" is serializing str to byte array
- should be EncodeString or SerializeString or PackString
*/

// FromI8 ...
func FromI8(in int8) []byte {
	//Serialize Atomic uses switch! Use serialize int8 directly
	//copy code over from encoder
	var b [1]byte
	b[0] = byte(in)
	return b[:1]
}

// FromI16 ...
func FromI16(in int16) []byte {
	//Serialize Atomic uses switch! Use serialize int16 directly
	//copy code over
	var b [2]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	return b[:2]
}

// FromI32 ...
func FromI32(in int32) []byte {
	//Serialize Atomic uses switch! Use serialize int32 directly
	//copy code over
	var b [4]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	b[2] = byte(in >> 16)
	b[3] = byte(in >> 24)
	return b[:4]
}

// FromI64 ...
func FromI64(in int64) []byte {
	//Serialize Atomic uses switch! Use serialize int64 directly
	//copy code over
	var b [8]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	b[2] = byte(in >> 16)
	b[3] = byte(in >> 24)
	b[4] = byte(in >> 32)
	b[5] = byte(in >> 40)
	b[6] = byte(in >> 48)
	b[7] = byte(in >> 56)
	return b[:8]
}

// FromUI8 ...
func FromUI8(in uint8) []byte {
	var b [1]byte
	b[0] = in
	return b[:1]
}

// FromUI16 ...
func FromUI16(in uint16) []byte {
	var b [2]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	return b[:2]
}

// FromUI32 ...
func FromUI32(in uint32) []byte {
	var b [4]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	b[2] = byte(in >> 16)
	b[3] = byte(in >> 24)
	return b[:4]
}

// FromUI64 ...
func FromUI64(in uint64) []byte {
	var b [8]byte
	b[0] = byte(in)
	b[1] = byte(in >> 8)
	b[2] = byte(in >> 16)
	b[3] = byte(in >> 24)
	b[4] = byte(in >> 32)
	b[5] = byte(in >> 40)
	b[6] = byte(in >> 48)
	b[7] = byte(in >> 56)
	return b[:8]
}

// FromF32 ...
func FromF32(in float32) []byte {
	return FromUI32(math.Float32bits(in))
}

// FromF64 ...
func FromF64(in float64) []byte {
	return FromUI64(math.Float64bits(in))
}

// FromStr ...
func FromStr(in string) []byte {
	return encoder.Serialize(in)
}

// FromBool ...
func FromBool(in bool) []byte {
	if in {
		return []byte{1}
	}
	return []byte{0}

}