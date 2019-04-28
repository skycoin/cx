package cxcore

import (
	"math/rand"
	"strconv"
)

func opUI64ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(ReadUI64(fp, expr.Inputs[0]), 10))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadUI64(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Add(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) + ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Sub(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) - ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Mul(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) * ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Div(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) / ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Gt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) > ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Gteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) >= ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Lt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) < ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Lteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) <= ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Eq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) == ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Uneq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI64(fp, expr.Inputs[0]) != ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Mod(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) % ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Rand(expr *CXExpression, fp int) {
	outB0 := FromUI64(rand.Uint64())
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitand(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) & ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitor(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) | ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitxor(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) ^ ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitclear(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) &^ ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitshl(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) << ReadUI64(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Bitshr(expr *CXExpression, fp int) {
	outB0 := FromUI64(ReadUI64(fp, expr.Inputs[0]) >> ReadUI64(fp, expr.Outputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI64Max(expr *CXExpression, fp int) {
	max := ReadUI64(fp, expr.Inputs[0])
	next := ReadUI64(fp, expr.Inputs[1])
	if next > max {
		max = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI64(max))
}

func opUI64Min(expr *CXExpression, fp int) {
	min := ReadUI64(fp, expr.Inputs[0])
	next := ReadUI64(fp, expr.Inputs[1])
	if next > min {
		min = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI64(min))
}
