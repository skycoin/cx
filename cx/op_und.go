package base

import (
	// "fmt"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_lt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) < ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) < ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) < ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) < ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_gt(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) > ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) > ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) > ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) > ReadF64(stack, fp, inp2))
	}
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_lteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) <= ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) <= ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) <= ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) <= ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_gteq(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) >= ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) >= ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) >= ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) >= ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_equal(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) == ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) == ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) == ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) == ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_unequal(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromBool(ReadI32(stack, fp, inp1) != ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromBool(ReadI64(stack, fp, inp1) != ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromBool(ReadF32(stack, fp, inp1) != ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromBool(ReadF64(stack, fp, inp1) != ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitand(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) & ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) & ReadI64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) | ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) | ReadI64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitxor(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) ^ ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) ^ ReadI64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_mul(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) * ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) * ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) * ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) * ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_div(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) / ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) / ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) / ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) / ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_mod(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) % ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) % ReadI64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_add(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) + ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) + ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) + ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) + ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_sub(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(ReadI32(stack, fp, inp1) - ReadI32(stack, fp, inp2))
	case TYPE_I64:
		outB1 = FromI64(ReadI64(stack, fp, inp1) - ReadI64(stack, fp, inp2))
	case TYPE_F32:
		outB1 = FromF32(ReadF32(stack, fp, inp1) - ReadF32(stack, fp, inp2))
	case TYPE_F64:
		outB1 = FromF64(ReadF64(stack, fp, inp1) - ReadF64(stack, fp, inp2))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitshl(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(stack, fp, inp1)) << uint32(ReadI32(stack, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint64(ReadI64(stack, fp, inp1)) << uint64(ReadI64(stack, fp, inp2))))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_bitshr(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	var outB1 []byte
	switch inp1.Type {
	case TYPE_I32:
		outB1 = FromI32(int32(uint32(ReadI32(stack, fp, inp1)) >> uint32(ReadI32(stack, fp, inp2))))
	case TYPE_I64:
		outB1 = FromI64(int64(uint32(ReadI64(stack, fp, inp1)) >> uint32(ReadI64(stack, fp, inp2))))
	}
	
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}

func op_len(expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(inp1.Lengths[len(inp1.Lengths) - 1]))
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, outB1)
}
