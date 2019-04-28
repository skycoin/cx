package cxcore

import (
	"math"
	"math/rand"
	"strconv"
)

func opUI16ToStr(expr *CXExpression, fp int) {
	outB0 := FromStr(strconv.FormatUint(uint64(ReadUI16(fp, expr.Inputs[0])), 10))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToI32(expr *CXExpression, fp int) {
	outB0 := FromI32(int32(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToI64(expr *CXExpression, fp int) {
	outB0 := FromI64(int64(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToUI8(expr *CXExpression, fp int) {
	outB0 := FromUI8(uint8(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToUI32(expr *CXExpression, fp int) {
	outB0 := FromUI32(uint32(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToUI64(expr *CXExpression, fp int) {
	outB0 := FromUI64(uint64(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToF32(expr *CXExpression, fp int) {
	outB0 := FromF32(float32(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16ToF64(expr *CXExpression, fp int) {
	outB0 := FromF64(float64(ReadUI16(fp, expr.Inputs[0])))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Add(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) + ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Sub(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) - ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Mul(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) * ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Div(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) / ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Gt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) > ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Gteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) >= ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Lt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) < ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Lteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) <= ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Eq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) == ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Uneq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadUI16(fp, expr.Inputs[0]) != ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Mod(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) % ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Rand(expr *CXExpression, fp int) {
	outB0 := FromUI16(uint16(rand.Int31n(int32(math.MaxUint16))))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitand(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) & ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitor(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) | ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitxor(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) ^ ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitclear(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) &^ ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitshl(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) << ReadUI16(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Bitshr(expr *CXExpression, fp int) {
	outB0 := FromUI16(ReadUI16(fp, expr.Inputs[0]) >> ReadUI16(fp, expr.Outputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opUI16Max(expr *CXExpression, fp int) {
	max := ReadUI16(fp, expr.Inputs[0])
	next := ReadUI16(fp, expr.Inputs[1])
	if next > max {
		max = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI16(max))
}

func opUI16Min(expr *CXExpression, fp int) {
	min := ReadUI16(fp, expr.Inputs[0])
	next := ReadUI16(fp, expr.Inputs[1])
	if next > min {
		min = next
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI16(min))
}
