package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

//first section

//second section

// ReadBool ...
func ReadBool(fp int, inp *CXArgument) bool {
	offset := GetOffset_bool(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.BOOL_SIZE]
	return helper.Deserialize_bool(readMemory)
	// return helper.DeserializeBool(ReadMemory(GetOffset_bool(fp, inp), inp))
}

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) int8 {
	offset := GetOffset_i8(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I8_SIZE]
	return helper.Deserialize_i8(readMemory)
	// return helper.Deserialize_i8(ReadMemory(GetOffset_i8(fp, inp), inp))
}

// ReadI16 ...
func ReadI16(fp int, inp *CXArgument) int16 {
	offset := GetOffset_i16(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I16_SIZE]
	return helper.Deserialize_i16(readMemory)
	// return helper.Deserialize_i16(ReadMemory(GetOffset_i16(fp, inp), inp))
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) int32 {
	offset := GetOffset_i32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I32_SIZE]
	return helper.Deserialize_i32(readMemory)
	// return helper.Deserialize_i32(ReadMemory(GetOffset_i32(fp, inp), inp))
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) int64 {
	offset := GetOffset_i64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I64_SIZE]
	return helper.Deserialize_i64(readMemory)
	// return helper.Deserialize_i64(ReadMemory(GetOffset_i64(fp, inp), inp))
}

// ReadUI8 ...
func ReadUI8(fp int, inp *CXArgument) uint8 {
	offset := GetOffset_ui8(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I8_SIZE]
	return helper.Deserialize_ui8(readMemory)
	// return helper.Deserialize_ui8(ReadMemory(GetOffset_ui8(fp, inp), inp))
}

// ReadUI16 ...
func ReadUI16(fp int, inp *CXArgument) uint16 {
	offset := GetOffset_ui16(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I16_SIZE]
	return helper.Deserialize_ui16(readMemory)
	// return helper.Deserialize_ui16(ReadMemory(GetOffset_ui16(fp, inp), inp))
}

// ReadUI32 ...
func ReadUI32(fp int, inp *CXArgument) uint32 {
	offset := GetOffset_ui32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I32_SIZE]
	return helper.Deserialize_ui32(readMemory)
	// return helper.Deserialize_ui32(ReadMemory(GetOffset_ui32(fp, inp), inp))
}

// ReadUI64 ...
func ReadUI64(fp int, inp *CXArgument) uint64 {
	offset := GetOffset_ui64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.I64_SIZE]
	return helper.Deserialize_ui64(readMemory)
	// return helper.Deserialize_ui64(ReadMemory(GetOffset_ui64(fp, inp), inp))
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) float32 {
	offset := GetOffset_f32(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.F32_SIZE]
	return helper.Deserialize_f32(readMemory)
	// return helper.Deserialize_f32(ReadMemory(GetOffset_f32(fp, inp), inp))
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) float64 {
	offset := GetOffset_f64(fp, inp)
	readMemory := PROGRAM.Memory[offset : offset+constants.F64_SIZE]
	return helper.Deserialize_f64(readMemory)
	// return helper.Deserialize_f64(ReadMemory(GetOffset_f64(fp, inp), inp))
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
