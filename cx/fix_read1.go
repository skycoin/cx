package cxcore

import(
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/helper"
)

//first section

//second section

// ReadBool ...
func ReadBool(fp int, inp *ast.CXArgument) bool {
	return helper.DeserializeBool(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadI8 ...
func ReadI8(fp int, inp *ast.CXArgument) int8 {
	return helper.Deserialize_i8(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadI16 ...
func ReadI16(fp int, inp *ast.CXArgument) int16 {
	return helper.Deserialize_i16(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadI32 ...
func ReadI32(fp int, inp *ast.CXArgument) int32 {
	return helper.Deserialize_i32(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadI64 ...
func ReadI64(fp int, inp *ast.CXArgument) int64 {
	return helper.Deserialize_i64(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadUI8 ...
func ReadUI8(fp int, inp *ast.CXArgument) uint8 {
	return helper.Deserialize_ui8(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadUI16 ...
func ReadUI16(fp int, inp *ast.CXArgument) uint16 {
	return helper.Deserialize_ui16(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadUI32 ...
func ReadUI32(fp int, inp *ast.CXArgument) uint32 {
	return helper.Deserialize_ui32(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadUI64 ...
func ReadUI64(fp int, inp *ast.CXArgument) uint64 {
	return helper.Deserialize_ui64(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadF32 ...
func ReadF32(fp int, inp *ast.CXArgument) float32 {
	return helper.Deserialize_f32(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}

// ReadF64 ...
func ReadF64(fp int, inp *ast.CXArgument) float64 {
	return helper.Deserialize_f64(ast.ReadMemory(ast.GetFinalOffset(fp, inp), inp))
}