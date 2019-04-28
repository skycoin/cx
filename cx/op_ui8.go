package cxcore

import (
	"math"
	"math/rand"
	"strconv"
)

func opUI8ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(uint64(ReadUI8(fp, expr.Inputs[0])), 10))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToUI16(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadUI8(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Add(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) + ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Sub(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) - ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Mul(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) * ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Div(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) / ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Gt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) > ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Gteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) >= ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Lt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) < ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Lteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) <= ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Eq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) == ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Uneq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI8(fp, expr.Inputs[0]) != ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Mod(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) % ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Rand(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(rand.Int31n(int32(math.MaxUint8))))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitand(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) & ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitor(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) | ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitxor(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) ^ ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitclear(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) &^ ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitshl(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) << ReadUI8(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Bitshr(expr *CXExpression, fp int) {
	outB0 := FromUI8(ReadUI8(fp, expr.Inputs[0]) >> uint32(ReadUI8(fp, expr.Outputs[1])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI8Max(expr *CXExpression, fp int) {
	max := ReadUI8(fp, expr.Inputs[0])
	next := ReadUI8(fp, expr.Inputs[1])
	if next > max {
		max = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI8(max))
}

func opUI8Min(expr *CXExpression, fp int) {
	min := ReadUI8(fp, expr.Inputs[0])
	next := ReadUI8(fp, expr.Inputs[1])
	if next > min {
		min = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI8(min))
}
