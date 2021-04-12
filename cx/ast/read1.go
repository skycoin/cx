package ast

import (
	"github.com/skycoin/cx/cx/helper"
)

//first section

//second section

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) bool {
	return helper.DeserializeBool(ReadMemory(GetOffset_bool(fp, inp), inp))
}

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) int8 {
	offset := GetOffset_i8(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+1]
	return helper.Deserialize_i8(readMemory)
}

// ReadI16 ...
func ReadI16(fp int, inp *CXArgument) int16 {
	offset := GetOffset_i16(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+2]
	return helper.Deserialize_i16(readMemory)
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) int32 {
	offset := GetOffset_i32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+4]
	return helper.Deserialize_i32(readMemory)
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) int64 {
	offset := GetOffset_i64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+8]
	return helper.Deserialize_i64(readMemory)
}

// ReadUI8 ...
func ReadUI8(fp int, inp *CXArgument) uint8 {
	offset := GetOffset_ui8(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+1]
	return helper.Deserialize_ui8(readMemory)
}

// ReadUI16 ...
func ReadUI16(fp int, inp *CXArgument) uint16 {
	offset := GetOffset_ui16(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+2]
	return helper.Deserialize_ui16(readMemory)
}

// ReadUI32 ...
func ReadUI32(fp int, inp *CXArgument) uint32 {
	offset := GetOffset_ui32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+4]
	return helper.Deserialize_ui32(readMemory)
}

// ReadUI64 ...
func ReadUI64(fp int, inp *CXArgument) uint64 {
	offset := GetOffset_ui64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+8]
	return helper.Deserialize_ui64(readMemory)
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) float32 {
	offset := GetOffset_f32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+4]
	return helper.Deserialize_f32(readMemory)
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) float64 {
	offset := GetOffset_f64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+8]
	return helper.Deserialize_f64(readMemory)
}

// ReadSlice ...
func ReadSlice(fp int, inp *CXArgument) int32 {
	return helper.Deserialize_i32(ReadMemory(GetOffset_slice(fp, inp), inp))
}

// ReadArray ...
func ReadArray(fp int, inp *CXArgument) int32 {
	return helper.Deserialize_i32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadPtr ...
func ReadPtr(fp int, inp *CXArgument) int32 {
	return helper.Deserialize_i32(ReadMemory(GetFinalOffset(fp, inp), inp))
}
