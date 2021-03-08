package cxcore

import()

//first section

//second section

// ReadI8 ...
func ReadI8(fp int, inp *CXArgument) int8 {
	return Deserialize_i8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI16 ...
func ReadI16(fp int, inp *CXArgument) int16 {
	return Deserialize_i16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI32 ...
func ReadI32(fp int, inp *CXArgument) int32 {
	return Deserialize_i32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadI64 ...
func ReadI64(fp int, inp *CXArgument) int64 {
	return Deserialize_i64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI8 ...
func ReadUI8(fp int, inp *CXArgument) uint8 {
	return Deserialize_ui8(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI16 ...
func ReadUI16(fp int, inp *CXArgument) uint16 {
	return Deserialize_ui16(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI32 ...
func ReadUI32(fp int, inp *CXArgument) uint32 {
	return Deserialize_ui32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadUI64 ...
func ReadUI64(fp int, inp *CXArgument) uint64 {
	return Deserialize_ui64(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF32 ...
func ReadF32(fp int, inp *CXArgument) float32 {
	return Deserialize_f32(ReadMemory(GetFinalOffset(fp, inp), inp))
}

// ReadF64 ...
func ReadF64(fp int, inp *CXArgument) float64 {
	return Deserialize_f64(ReadMemory(GetFinalOffset(fp, inp), inp))
}