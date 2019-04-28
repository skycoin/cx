package cxcore

import (
	"math/rand"
	"strconv"
)

func opUI32ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(uint64(ReadUI32(fp, expr.Inputs[0])), 10))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadUI32(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Add(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) + ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Sub(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) - ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Mul(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) * ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Div(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) / ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Gt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) > ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Gteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) >= ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Lt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) < ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Lteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) <= ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Eq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) == ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Uneq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI32(fp, expr.Inputs[0]) != ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Mod(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) % ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Rand(expr *CXExpression, fp int) {
	outB0 := FromUI32((rand.Uint32()))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitand(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) & ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitor(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) | ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitxor(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) ^ ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitclear(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) &^ ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitshl(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) << ReadUI32(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Bitshr(expr *CXExpression, fp int) {
	outB0 := FromUI32(ReadUI32(fp, expr.Inputs[0]) >> ReadUI32(fp, expr.Outputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI32Max(expr *CXExpression, fp int) {
	max := ReadUI32(fp, expr.Inputs[0])
	next := ReadUI32(fp, expr.Inputs[1])
	if next > max {
		max = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI32(max))
}

func opUI32Min(expr *CXExpression, fp int) {
	min := ReadUI32(fp, expr.Inputs[0])
	next := ReadUI32(fp, expr.Inputs[1])
	if next > min {
		min = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI32(min))
}
